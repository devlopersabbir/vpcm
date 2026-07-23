package inventory

import (
	"time"
)

// ─── Core Identity ────────────────────────────────────────────────────────────

type Server struct {
	ID         uint       `json:"id" bson:"id"`
	UUID       string     `json:"uuid" bson:"uuid"`
	Name       string     `json:"name" bson:"name"`
	Host       string     `json:"host" bson:"host"`
	Port       int        `json:"port" bson:"port"`
	Username   string     `json:"username" bson:"username"`
	AuthType   string     `json:"auth_type" bson:"auth_type"`
	AuthSecret string     `json:"auth_secret,omitempty" bson:"auth_secret,omitempty"`
	Provider   string     `json:"provider" bson:"provider"`
	CreatedAt  time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" bson:"updated_at"`
	LastSeen   *time.Time `json:"last_seen,omitempty" bson:"last_seen,omitempty"`
	IsFavorite bool       `json:"is_favorite" bson:"is_favorite"`
	Tags       []Tag      `json:"tags,omitempty" bson:"tags,omitempty"`
}

// ─── Network & Location ───────────────────────────────────────────────────────

type ServerNetwork struct {
	ID               uint   `json:"id" bson:"id"`
	ServerID         uint   `json:"server_id" bson:"server_id"`
	Hostname         string `json:"hostname" bson:"hostname"`
	PublicIP         string `json:"public_ip" bson:"public_ip"`
	PrivateIP        string `json:"private_ip" bson:"private_ip"`
	MACAddress       string `json:"mac_address" bson:"mac_address"`
	Region           string `json:"region" bson:"region"`
	AvailabilityZone string `json:"availability_zone" bson:"availability_zone"`
}

// ─── Hardware & Firmware ──────────────────────────────────────────────────────

type ServerHardware struct {
	ID             uint   `json:"id" bson:"id"`
	ServerID       uint   `json:"server_id" bson:"server_id"`
	CPUModel       string `json:"cpu_model" bson:"cpu_model"`
	CPUCores       int    `json:"cpu_cores" bson:"cpu_cores"`
	RAMTotal       string `json:"ram_total" bson:"ram_total"`
	SwapTotal      string `json:"swap_total" bson:"swap_total"`
	DiskTotal      string `json:"disk_total" bson:"disk_total"`
	Virtualization string `json:"virtualization" bson:"virtualization"`
	InstanceType   string `json:"instance_type" bson:"instance_type"`
	SerialNumber   string `json:"serial_number" bson:"serial_number"`
	BIOSVersion    string `json:"bios_version" bson:"bios_version"`
	Uptime         string `json:"uptime" bson:"uptime"`
}

// ─── Operating System ─────────────────────────────────────────────────────────

type ServerOS struct {
	ID             uint   `json:"id" bson:"id"`
	ServerID       uint   `json:"server_id" bson:"server_id"`
	OSFamily       string `json:"os_family" bson:"os_family"`
	OSVersion      string `json:"os_version" bson:"os_version"`
	KernelVersion  string `json:"kernel_version" bson:"kernel_version"`
	Architecture   string `json:"architecture" bson:"architecture"`
	InitSystem     string `json:"init_system" bson:"init_system"`
	Timezone       string `json:"timezone" bson:"timezone"`
	Locale         string `json:"locale" bson:"locale"`
	PackageManager string `json:"package_manager" bson:"package_manager"`
}

// ─── Installed Software ───────────────────────────────────────────────────────

type Software struct {
	ID       uint   `json:"id" bson:"id"`
	ServerID uint   `json:"server_id" bson:"server_id"`
	Name     string `json:"name" bson:"name"`
	Version  string `json:"version" bson:"version"`
}

// ─── Misc ─────────────────────────────────────────────────────────────────────

type Tag struct {
	ID   uint   `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

type ConnectionLog struct {
	ID           uint       `json:"id" bson:"id"`
	ServerID     uint       `json:"server_id" bson:"server_id"`
	ServerName   string     `json:"server_name" bson:"server_name"`
	Username     string     `json:"username" bson:"username"`
	Host         string     `json:"host" bson:"host"`
	LoggedInAt   time.Time  `json:"logged_in_at" bson:"logged_in_at"`
	LoggedOutAt  *time.Time `json:"logged_out_at,omitempty" bson:"logged_out_at,omitempty"`
	Duration     string     `json:"duration,omitempty" bson:"duration,omitempty"`
	Status       string     `json:"status" bson:"status"`
	ErrorMessage string     `json:"error_message,omitempty" bson:"error_message,omitempty"`
}

// ─── Joined API View ──────────────────────────────────────────────────────────

type ServerView struct {
	ID         uint       `json:"id" bson:"id"`
	UUID       string     `json:"uuid" bson:"uuid"`
	Name       string     `json:"name" bson:"name"`
	Host       string     `json:"host" bson:"host"`
	Port       int        `json:"port" bson:"port"`
	Username   string     `json:"username" bson:"username"`
	AuthType   string     `json:"auth_type" bson:"auth_type"`
	AuthSecret string     `json:"auth_secret,omitempty" bson:"auth_secret,omitempty"`
	Provider   string     `json:"provider" bson:"provider"`
	CreatedAt  time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" bson:"updated_at"`
	LastSeen   *time.Time `json:"last_seen,omitempty" bson:"last_seen,omitempty"`
	IsFavorite bool       `json:"is_favorite" bson:"is_favorite"`
	Tags       []Tag      `json:"tags" bson:"tags"`

	Network  *ServerNetwork  `json:"network" bson:"network"`
	Hardware *ServerHardware `json:"hardware" bson:"hardware"`
	OS       *ServerOS       `json:"os" bson:"os"`
	Software []Software      `json:"software" bson:"software"`
}

// ─── Terminal Preference ──────────────────────────────────────────────────────

type TerminalPreference struct {
	ID             uint    `json:"id" bson:"id"`
	FontSize       int     `json:"font_size" bson:"font_size"`
	FontFamily     string  `json:"font_family" bson:"font_family"`
	Background     string  `json:"background" bson:"background"`
	Foreground     string  `json:"foreground" bson:"foreground"`
	Opacity        float64 `json:"opacity" bson:"opacity"`
	Blur           int     `json:"blur" bson:"blur"`
	CursorStyle    string  `json:"cursor_style" bson:"cursor_style"`
	CursorBlink    bool    `json:"cursor_blink" bson:"cursor_blink"`
	UpdatedAt      time.Time `json:"updated_at" bson:"updated_at"`
}
