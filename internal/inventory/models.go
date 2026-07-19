package inventory

import (
	"time"
)

// ─── Core Identity ────────────────────────────────────────────────────────────
// Server holds only the minimal fields required to identify and connect to a
// remote machine. All metadata is stored in child tables keyed by server_id.

type Server struct {
	ID         uint       `json:"id"`
	UUID       string     `json:"uuid"`
	Name       string     `json:"name"`
	Host       string     `json:"host"`
	Port       int        `json:"port"`
	Username   string     `json:"username"`
	AuthType   string     `json:"auth_type"`   // "key" or "password"
	AuthSecret string     `json:"auth_secret,omitempty"`
	Provider   string     `json:"provider"`    // e.g. AWS, DigitalOcean
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	LastSeen   *time.Time `json:"last_seen,omitempty"`
	Tags       []Tag      `json:"tags,omitempty"`
}

// ─── Network & Location ───────────────────────────────────────────────────────

type ServerNetwork struct {
	ID               uint   `json:"id"`
	ServerID         uint   `json:"server_id"`
	Hostname         string `json:"hostname"`
	PublicIP         string `json:"public_ip"`
	PrivateIP        string `json:"private_ip"`
	MACAddress       string `json:"mac_address"`
	Region           string `json:"region"`            // e.g. us-east-1
	AvailabilityZone string `json:"availability_zone"` // e.g. us-east-1a
}

// ─── Hardware & Firmware ──────────────────────────────────────────────────────

type ServerHardware struct {
	ID             uint   `json:"id"`
	ServerID       uint   `json:"server_id"`
	CPUModel       string `json:"cpu_model"`
	CPUCores       int    `json:"cpu_cores"`
	RAMTotal       string `json:"ram_total"`
	SwapTotal      string `json:"swap_total"`
	DiskTotal      string `json:"disk_total"`
	Virtualization string `json:"virtualization"` // kvm, xen, none, …
	InstanceType   string `json:"instance_type"`  // t3.small, n1-standard-1, …
	SerialNumber   string `json:"serial_number"`
	BIOSVersion    string `json:"bios_version"`
	Uptime         string `json:"uptime"`
}

// ─── Operating System ─────────────────────────────────────────────────────────

type ServerOS struct {
	ID             uint   `json:"id"`
	ServerID       uint   `json:"server_id"`
	OSFamily       string `json:"os_family"`       // Ubuntu, CentOS, …
	OSVersion      string `json:"os_version"`      // 22.04, 9, …
	KernelVersion  string `json:"kernel_version"`
	Architecture   string `json:"architecture"`    // x86_64, aarch64, …
	InitSystem     string `json:"init_system"`     // systemd, openrc, …
	Timezone       string `json:"timezone"`        // UTC, America/New_York, …
	Locale         string `json:"locale"`          // en_US.UTF-8, …
	PackageManager string `json:"package_manager"` // apt, dnf, apk, …
}

// ─── Installed Software ───────────────────────────────────────────────────────

type Software struct {
	ID       uint   `json:"id"`
	ServerID uint   `json:"server_id"`
	Name     string `json:"name"`
	Version  string `json:"version"`
}

// ─── Misc ─────────────────────────────────────────────────────────────────────

type Tag struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ConnectionLog struct {
	ID           uint       `json:"id"`
	ServerID     uint       `json:"server_id"`
	ServerName   string     `json:"server_name"`
	Username     string     `json:"username"`
	Host         string     `json:"host"`
	LoggedInAt   time.Time  `json:"logged_in_at"`
	LoggedOutAt  *time.Time `json:"logged_out_at,omitempty"`
	Duration     string     `json:"duration,omitempty"`
	Status       string     `json:"status"`
	ErrorMessage string     `json:"error_message,omitempty"`
}

// ─── Joined API View ──────────────────────────────────────────────────────────
// ServerView is the rich response returned by the API. It is assembled by
// LEFT JOINing the four child tables onto the servers row.

type ServerView struct {
	// Core identity
	ID        uint       `json:"id"`
	UUID      string     `json:"uuid"`
	Name      string     `json:"name"`
	Host      string     `json:"host"`
	Port      int        `json:"port"`
	Username  string     `json:"username"`
	AuthType  string     `json:"auth_type"`
	Provider  string     `json:"provider"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	LastSeen  *time.Time `json:"last_seen,omitempty"`
	Tags      []Tag      `json:"tags"`

	// Child domains — nil if not yet collected
	Network  *ServerNetwork  `json:"network,omitempty"`
	Hardware *ServerHardware `json:"hardware,omitempty"`
	OS       *ServerOS       `json:"os,omitempty"`
	Software []Software      `json:"software"`
}
