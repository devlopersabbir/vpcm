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
	AuthType  string     `gorm:"size:50" json:"auth_type"` // e.g. key, password
	Provider  string     `gorm:"size:100" json:"provider"` // e.g. aws, digitalocean
	Region    string     `gorm:"size:100" json:"region"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	LastSeen  *time.Time `json:"last_seen,omitempty"`
	Tags      []Tag      `gorm:"many2many:server_tags;" json:"tags,omitempty"`
	Software  []Software `json:"software,omitempty"`
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
