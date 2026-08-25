package cli

import (
	"testing"

	"github.com/devlopersabbir/vpcm/internal/inventory"
)

func TestResolveImportFormat(t *testing.T) {
	cases := []struct {
		name      string
		requested string
		path      string
		data      string
		want      string
		wantErr   bool
	}{
		{name: "explicit json", requested: "json", path: "servers.csv", data: "", want: "json"},
		{name: "explicit yml alias", requested: "yml", path: "", data: "", want: "yaml"},
		{name: "unsupported explicit", requested: "toml", path: "", data: "", wantErr: true},
		{name: "json extension", requested: "auto", path: "backups/servers.json", data: "", want: "json"},
		{name: "yaml extension", requested: "", path: "backups/servers.yml", data: "", want: "yaml"},
		{name: "csv extension", requested: "auto", path: "servers.CSV", data: "", want: "csv"},
		{name: "ssh config filename", requested: "auto", path: "/home/me/.ssh/config", data: "", want: "ssh"},
		{name: "sniff json array", requested: "auto", path: "dump", data: "  [{\"name\":\"a\"}]", want: "json"},
		{name: "sniff ssh config", requested: "auto", path: "dump", data: "Host web\n  HostName 1.1.1.1", want: "ssh"},
		{name: "sniff csv header", requested: "auto", path: "dump", data: "ID,UUID,Name,Host,Port\n1,x,a,1.1.1.1,22", want: "csv"},
		{name: "sniff yaml list", requested: "auto", path: "dump", data: "- name: web\n  host: 1.1.1.1", want: "yaml"},
		{name: "undetectable", requested: "auto", path: "dump", data: "nothing structured here", wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := resolveImportFormat(tc.requested, tc.path, []byte(tc.data))
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got format %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("got format %q, want %q", got, tc.want)
			}
		})
	}
}

func TestParseCSVServers(t *testing.T) {
	data := `ID,UUID,Name,Host,Port,Username,AuthType,Provider
1,abc-123,web-prod,192.168.1.10,2222,deploy,password,DigitalOcean
2,def-456,db-prod,10.0.0.5,,root,,Generic VPS
`
	servers, err := parseCSVServers([]byte(data))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(servers) != 2 {
		t.Fatalf("got %d servers, want 2", len(servers))
	}

	first := servers[0]
	if first.UUID != "abc-123" || first.Name != "web-prod" || first.Host != "192.168.1.10" {
		t.Fatalf("unexpected identity fields: %+v", first)
	}
	if first.Port != 2222 || first.Username != "deploy" || first.AuthType != "password" {
		t.Fatalf("unexpected connection fields: %+v", first)
	}
	if first.Provider != "DigitalOcean" {
		t.Fatalf("got provider %q, want DigitalOcean", first.Provider)
	}

	// An empty port cell stays zero here and is defaulted during normalization.
	if servers[1].Port != 0 {
		t.Fatalf("got port %d for blank cell, want 0", servers[1].Port)
	}
}

func TestParseCSVServersRejectsMissingHeaders(t *testing.T) {
	if _, err := parseCSVServers([]byte("ID,UUID,Port\n1,abc,22\n")); err == nil {
		t.Fatal("expected an error for CSV input without Name and Host columns")
	}
}

func TestParseCSVServersRejectsBadPort(t *testing.T) {
	if _, err := parseCSVServers([]byte("Name,Host,Port\nweb,1.1.1.1,http\n")); err == nil {
		t.Fatal("expected an error for a non-numeric port")
	}
}

func TestParseSSHConfigServers(t *testing.T) {
	data := `# global defaults
Host *
    ServerAliveInterval 60

Host bastion jump
    HostName bastion.example.com
    User admin
    Port 2200
    IdentityFile ~/.ssh/id_ed25519

Host web-*
    HostName wildcard.example.com

Host db
    HostName=10.0.0.5
`
	servers, err := parseSSHConfigServers([]byte(data))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(servers) != 2 {
		t.Fatalf("got %d servers, want 2 (wildcard patterns must be ignored)", len(servers))
	}

	bastion := servers[0]
	if bastion.Name != "bastion" {
		t.Fatalf("got name %q, want the first alias 'bastion'", bastion.Name)
	}
	if bastion.Host != "bastion.example.com" || bastion.Username != "admin" || bastion.Port != 2200 {
		t.Fatalf("unexpected connection fields: %+v", bastion)
	}
	if bastion.AuthType != "key" || bastion.AuthSecret != "~/.ssh/id_ed25519" {
		t.Fatalf("expected IdentityFile to become key auth, got %+v", bastion)
	}

	if servers[1].Host != "10.0.0.5" {
		t.Fatalf("got host %q, want '=' separated value to be parsed", servers[1].Host)
	}
}

func TestParseSSHConfigServersRejectsBadPort(t *testing.T) {
	if _, err := parseSSHConfigServers([]byte("Host web\n    Port ssh\n")); err == nil {
		t.Fatal("expected an error for a non-numeric Port")
	}
}

func TestNormalizeImportedServer(t *testing.T) {
	t.Run("applies defaults and drops the source id", func(t *testing.T) {
		got, err := normalizeImportedServer(inventory.Server{ID: 42, Name: "  web  ", Host: " 1.1.1.1 "})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.ID != 0 {
			t.Fatalf("got id %d, want the imported id to be discarded", got.ID)
		}
		if got.Name != "web" || got.Host != "1.1.1.1" {
			t.Fatalf("expected surrounding whitespace to be trimmed, got %+v", got)
		}
		if got.Port != 22 || got.Username != "root" {
			t.Fatalf("got port %d and user %q, want 22 and root", got.Port, got.Username)
		}
	})

	t.Run("falls back to the host as name", func(t *testing.T) {
		got, err := normalizeImportedServer(inventory.Server{Host: "1.1.1.1"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Name != "1.1.1.1" {
			t.Fatalf("got name %q, want the host as fallback", got.Name)
		}
	})

	for _, tc := range []struct {
		name   string
		server inventory.Server
	}{
		{name: "missing host", server: inventory.Server{Name: "web"}},
		{name: "port too high", server: inventory.Server{Name: "web", Host: "1.1.1.1", Port: 70000}},
		{name: "negative port", server: inventory.Server{Name: "web", Host: "1.1.1.1", Port: -1}},
		{name: "unknown auth type", server: inventory.Server{Name: "web", Host: "1.1.1.1", AuthType: "magic"}},
	} {
		t.Run("rejects "+tc.name, func(t *testing.T) {
			if _, err := normalizeImportedServer(tc.server); err == nil {
				t.Fatalf("expected an error for %s", tc.name)
			}
		})
	}
}

func TestFindImportConflict(t *testing.T) {
	stored := inventory.Server{ID: 7, UUID: "abc-123", Name: "web-prod", Host: "1.1.1.1"}
	byUUID := map[string]inventory.Server{stored.UUID: stored}
	byName := map[string]inventory.Server{"web-prod": stored}

	if match, on := findImportConflict(inventory.Server{UUID: "abc-123", Name: "other"}, byUUID, byName); match == nil || on != "uuid" {
		t.Fatalf("expected a uuid match, got %v on %q", match, on)
	}
	if match, on := findImportConflict(inventory.Server{Name: "WEB-PROD"}, byUUID, byName); match == nil || on != "name" {
		t.Fatalf("expected a case-insensitive name match, got %v on %q", match, on)
	}
	if match, _ := findImportConflict(inventory.Server{UUID: "zzz", Name: "fresh"}, byUUID, byName); match != nil {
		t.Fatalf("expected no match, got %+v", match)
	}
}

func TestMergeImportedServer(t *testing.T) {
	existing := inventory.Server{
		ID:         7,
		UUID:       "abc-123",
		Name:       "web-prod",
		Host:       "1.1.1.1",
		Port:       22,
		Username:   "root",
		AuthType:   "password",
		AuthSecret: "stored-secret",
		Provider:   "DigitalOcean",
		Tags:       []inventory.Tag{{Name: "prod"}},
	}

	// A CSV import carries no credentials, so they must survive the merge.
	merged := mergeImportedServer(existing, inventory.Server{
		Name:     "web-prod",
		Host:     "2.2.2.2",
		Port:     2222,
		Username: "deploy",
	})

	if merged.ID != 7 || merged.UUID != "abc-123" {
		t.Fatalf("expected the stored identity to be kept, got %+v", merged)
	}
	if merged.Host != "2.2.2.2" || merged.Port != 2222 || merged.Username != "deploy" {
		t.Fatalf("expected connection fields to be updated, got %+v", merged)
	}
	if merged.AuthSecret != "stored-secret" || merged.AuthType != "password" {
		t.Fatalf("expected credentials to be preserved, got %+v", merged)
	}
	if merged.Provider != "DigitalOcean" || len(merged.Tags) != 1 {
		t.Fatalf("expected provider and tags to be preserved, got %+v", merged)
	}

	overridden := mergeImportedServer(existing, inventory.Server{
		Name:       "web-prod",
		Host:       "1.1.1.1",
		Port:       22,
		Username:   "root",
		AuthType:   "key",
		AuthSecret: "new-secret",
		Provider:   "AWS",
		Tags:       []inventory.Tag{{Name: "staging"}, {Name: "eu"}},
	})
	if overridden.AuthType != "key" || overridden.AuthSecret != "new-secret" {
		t.Fatalf("expected supplied credentials to win, got %+v", overridden)
	}
	if overridden.Provider != "AWS" || len(overridden.Tags) != 2 {
		t.Fatalf("expected supplied provider and tags to win, got %+v", overridden)
	}
}

func TestUniqueImportName(t *testing.T) {
	taken := map[string]bool{"web": true, "web-1": true}

	if got := uniqueImportName("db", taken); got != "db" {
		t.Fatalf("got %q, want an unused name to be returned as is", got)
	}
	if got := uniqueImportName("web", taken); got != "web-2" {
		t.Fatalf("got %q, want web-2", got)
	}
	if got := uniqueImportName("WEB", taken); got != "WEB-2" {
		t.Fatalf("got %q, want the collision check to be case-insensitive", got)
	}
}

func TestParseImportDataRejectsUnknownFormat(t *testing.T) {
	if _, err := parseImportData("toml", []byte("a = 1")); err == nil {
		t.Fatal("expected an error for an unsupported format")
	}
}

func TestParseJSONServersAcceptsSingleObject(t *testing.T) {
	servers, err := parseJSONServers([]byte(`{"name":"web","host":"1.1.1.1"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(servers) != 1 || servers[0].Name != "web" {
		t.Fatalf("got %+v, want a single server named web", servers)
	}
}
