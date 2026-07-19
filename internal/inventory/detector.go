package inventory

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/devlopersabbir/vpcm/internal/ssh"
)

// DetectProvider executes checks locally and remotely on the given SSH client to identify the VPS provider.
func DetectProvider(ctx context.Context, sshClient ssh.Client, host string) string {
	slog.Info("Starting VPS provider detection", "host", host)

	if sshClient != nil {
		// Step 1: Detect Cloud Metadata Service (Highest Confidence)
		if provider := detectMetadataService(ctx, sshClient); provider != "" {
			slog.Info("Provider detected via Metadata Service", "provider", provider)
			return provider
		}

		// Step 2: Read DMI Information
		if provider := detectDMIVendor(ctx, sshClient); provider != "" {
			slog.Info("Provider detected via DMI", "provider", provider)
			return provider
		}

		// Step 3: Installed Cloud Agents
		if provider := detectCloudAgents(ctx, sshClient); provider != "" {
			slog.Info("Provider detected via Cloud Agents", "provider", provider)
			return provider
		}
	}

	// Step 4: Reverse DNS Lookup (Local)
	if provider := detectReverseDNS(host); provider != "" {
		slog.Info("Provider detected via Reverse DNS", "provider", provider)
		return provider
	}

	// Step 5: Resolve Public IP + ASN (Local)
	if provider := detectASN(host); provider != "" {
		slog.Info("Provider detected via ASN Lookup", "provider", provider)
		return provider
	}

	slog.Info("VPS provider detection finished with no matches (Unknown/Generic)")
	return "Generic VPS"
}

func detectMetadataService(ctx context.Context, client ssh.Client) string {
	// AWS IMDSv2 Token Request or IMDSv1
	awsCmd := `curl -s -f -X PUT --connect-timeout 1 "http://169.254.169.254/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 60" || curl -s -f --connect-timeout 1 http://169.254.169.254/latest/meta-data/`
	if out, err := client.RunCommand(ctx, awsCmd); err == nil && len(out) > 0 {
		return "AWS"
	}

	// GCP Metadata Flavor Header
	gcpCmd := `curl -s -f -H "Metadata-Flavor: Google" --connect-timeout 1 http://169.254.169.254/computeMetadata/v1/`
	if _, err := client.RunCommand(ctx, gcpCmd); err == nil {
		return "GCP"
	}

	// Azure Instance Metadata Service (IMDS)
	azureCmd := `curl -s -f -H Metadata:true --connect-timeout 1 "http://169.254.169.254/metadata/instance?api-version=2021-02-01"`
	if _, err := client.RunCommand(ctx, azureCmd); err == nil {
		return "Azure"
	}

	// DigitalOcean Metadata Service
	doCmd := `curl -s -f --connect-timeout 1 http://169.254.169.254/metadata/v1.json`
	if _, err := client.RunCommand(ctx, doCmd); err == nil {
		return "DigitalOcean"
	}

	return ""
}

func detectDMIVendor(ctx context.Context, client ssh.Client) string {
	sysVendorCmd := `cat /sys/class/dmi/id/sys_vendor 2>/dev/null || cat /sys/devices/virtual/dmi/id/sys_vendor 2>/dev/null`
	out, err := client.RunCommand(ctx, sysVendorCmd)
	if err == nil {
		vendor := strings.ToLower(out)
		if strings.Contains(vendor, "amazon") {
			return "AWS"
		}
		if strings.Contains(vendor, "google") {
			return "GCP"
		}
		if strings.Contains(vendor, "microsoft") {
			return "Azure"
		}
		if strings.Contains(vendor, "digitalocean") {
			return "DigitalOcean"
		}
	}

	productNameCmd := `cat /sys/class/dmi/id/product_name 2>/dev/null || cat /sys/devices/virtual/dmi/id/product_name 2>/dev/null`
	out, err = client.RunCommand(ctx, productNameCmd)
	if err == nil {
		product := strings.ToLower(out)
		if strings.Contains(product, "amazon") || strings.Contains(product, "ec2") {
			return "AWS"
		}
		if strings.Contains(product, "google") {
			return "GCP"
		}
		if strings.Contains(product, "microsoft") || strings.Contains(product, "virtual machine") {
			// Note: strings.Contains(product, "microsoft") handles it, but "or" is invalid in Go. We'll use ||
			return "Azure"
		}
		if strings.Contains(product, "digitalocean") {
			return "DigitalOcean"
		}
	}

	return ""
}

func detectCloudAgents(ctx context.Context, client ssh.Client) string {
	// Look for AWS SSM agent
	if _, err := client.RunCommand(ctx, "ls /usr/bin/amazon-ssm-agent /var/lib/amazon 2>/dev/null"); err == nil {
		return "AWS"
	}
	// Look for Azure linux agent
	if _, err := client.RunCommand(ctx, "ls /usr/sbin/waagent /var/lib/waagent 2>/dev/null"); err == nil {
		return "Azure"
	}
	// Look for Google guest agent
	if _, err := client.RunCommand(ctx, "ls /usr/bin/google_guest_agent /var/lib/google 2>/dev/null"); err == nil {
		return "GCP"
	}
	// Look for DigitalOcean agent
	if _, err := client.RunCommand(ctx, "ls /usr/bin/droplet-agent /var/lib/droplet-agent 2>/dev/null"); err == nil {
		return "DigitalOcean"
	}
	return ""
}

func detectReverseDNS(host string) string {
	names, err := net.LookupAddr(host)
	if err != nil || len(names) == 0 {
		return ""
	}

	for _, name := range names {
		n := strings.ToLower(name)
		if strings.Contains(n, "amazonaws.com") {
			return "AWS"
		}
		if strings.Contains(n, "googleusercontent.com") {
			return "GCP"
		}
		if strings.Contains(n, "cloudapp.azure.com") || strings.Contains(n, "windows.net") {
			return "Azure"
		}
		if strings.Contains(n, "digitalocean.com") {
			return "DigitalOcean"
		}
	}
	return ""
}

func detectASN(host string) string {
	// Check if host is not a loopback or private IP
	ip := net.ParseIP(host)
	if ip != nil && (ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified()) {
		return ""
	}

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(fmt.Sprintf("http://ip-api.com/json/%s?fields=org,as", host))
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var result struct {
		Org string `json:"org"`
		AS  string `json:"as"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return ""
	}

	org := strings.ToLower(result.Org + " " + result.AS)
	if strings.Contains(org, "amazon") || strings.Contains(org, "aws") {
		return "AWS"
	}
	if strings.Contains(org, "google") {
		return "GCP"
	}
	if strings.Contains(org, "microsoft") || strings.Contains(org, "azure") {
		return "Azure"
	}
	if strings.Contains(org, "digitalocean") || strings.Contains(org, "digital ocean") {
		return "DigitalOcean"
	}
	if strings.Contains(org, "hetzner") {
		return "Hetzner"
	}
	if strings.Contains(org, "linode") || strings.Contains(org, "akamai") {
		return "Linode"
	}
	if strings.Contains(org, "vultr") {
		return "Vultr"
	}
	if strings.Contains(org, "ovh") {
		return "OVH"
	}

	return ""
}

// DetectOS detects the remote machine OS family and version via SSH.
func DetectOS(ctx context.Context, client ssh.Client) (string, string) {
	if client == nil {
		return "Unknown", ""
	}

	familyCmd := `cat /etc/os-release | grep -E '^ID=' | cut -d= -f2 | tr -d '"'`
	family, err := client.RunCommand(ctx, familyCmd)
	if err != nil {
		family, err = client.RunCommand(ctx, "uname -s")
		if err != nil {
			family = "Unknown"
		}
	}
	family = strings.TrimSpace(strings.ToLower(family))

	versionCmd := `cat /etc/os-release | grep -E '^VERSION_ID=' | cut -d= -f2 | tr -d '"'`
	version, err := client.RunCommand(ctx, versionCmd)
	if err != nil {
		version = ""
	}
	version = strings.TrimSpace(version)

	return family, version
}

// DetectSpecs detects remote CPU model, core counts, memory size, and disk space.
func DetectSpecs(ctx context.Context, client ssh.Client) (string, int, string, string) {
	if client == nil {
		return "Unknown", 0, "", ""
	}

	cpuModelCmd := `cat /proc/cpuinfo 2>/dev/null | grep -E '^model name' | head -n 1 | cut -d: -f2 | xargs || echo "Unknown"`
	cpuModel, err := client.RunCommand(ctx, cpuModelCmd)
	if err != nil {
		cpuModel = "Unknown"
	}
	cpuModel = strings.TrimSpace(cpuModel)

	cpuCoresCmd := `nproc 2>/dev/null || echo "1"`
	cpuCoresStr, err := client.RunCommand(ctx, cpuCoresCmd)
	cpuCores, _ := strconv.Atoi(strings.TrimSpace(cpuCoresStr))
	if cpuCores <= 0 {
		cpuCores = 1
	}

	ramCmd := `free -h 2>/dev/null | awk '/Mem:/ {print $2}' || echo ""`
	ramTotal, _ := client.RunCommand(ctx, ramCmd)
	ramTotal = strings.TrimSpace(ramTotal)

	diskCmd := `df -h / 2>/dev/null | awk 'NR==2 {print $2}' || echo ""`
	diskTotal, _ := client.RunCommand(ctx, diskCmd)
	diskTotal = strings.TrimSpace(diskTotal)

	return cpuModel, cpuCores, ramTotal, diskTotal
}
