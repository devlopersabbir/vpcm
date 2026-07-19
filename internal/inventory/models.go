package inventory

import (
	"time"
)

type Server struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UUID      string     `gorm:"type:uuid;uniqueIndex" json:"uuid"`
	Name      string     `gorm:"size:255;not null" json:"name"`
	Host      string     `gorm:"size:255;not null" json:"host"`
	Port      int        `gorm:"default:22" json:"port"`
	Username  string     `gorm:"size:255;not null" json:"username"`
	AuthType   string     `gorm:"size:50" json:"auth_type"` // e.g. key, password
	AuthSecret string     `gorm:"type:text" json:"auth_secret,omitempty"`
	Provider   string     `gorm:"size:100" json:"provider"` // e.g. aws, digitalocean
	Region     string     `gorm:"size:100" json:"region"`
	OSFamily   string     `gorm:"size:100" json:"os_family"`
	OSVersion  string     `gorm:"size:100" json:"os_version"`
	CPUModel   string     `gorm:"size:255" json:"cpu_model"`
	CPUCores   int        `gorm:"default:0" json:"cpu_cores"`
	RAMTotal   string     `gorm:"size:50" json:"ram_total"`
	DiskTotal  string     `gorm:"size:50" json:"disk_total"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	LastSeen   *time.Time `json:"last_seen,omitempty"`
	Tags       []Tag      `gorm:"many2many:server_tags;" json:"tags,omitempty"`
	Software   []Software `json:"software,omitempty"`
}

type Tag struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:100;uniqueIndex;not null" json:"name"`
}

type Software struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	ServerID uint   `gorm:"index" json:"server_id"`
	Name     string `gorm:"size:255;not null" json:"name"`
	Version  string `gorm:"size:100" json:"version"`
}

type ConnectionLog struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	ServerID     uint       `gorm:"index" json:"server_id"`
	ServerName   string     `gorm:"size:255;not null" json:"server_name"`
	Username     string     `gorm:"size:255;not null" json:"username"`
	Host         string     `gorm:"size:255;not null" json:"host"`
	LoggedInAt   time.Time  `json:"logged_in_at"`
	LoggedOutAt  *time.Time `json:"logged_out_at,omitempty"`
	Duration     string     `gorm:"size:50" json:"duration,omitempty"` // e.g. "5m32s"
	Status       string     `gorm:"size:50;not null" json:"status"`    // e.g. success, failed, active
	ErrorMessage string     `gorm:"type:text" json:"error_message,omitempty"`
}
