package inventory

import "context"

// ServerRepository is the storage contract. SQLite and MongoDB both implement this.
type ServerRepository interface {
	// Core CRUD
	Create(ctx context.Context, server *Server) error
	GetByID(ctx context.Context, id uint) (*Server, error)
	GetByUUID(ctx context.Context, uuid string) (*Server, error)
	List(ctx context.Context) ([]Server, error)
	Update(ctx context.Context, server *Server) error
	Delete(ctx context.Context, id uint) error

	// Tags
	AddTag(ctx context.Context, serverID uint, tagName string) error
	RemoveTag(ctx context.Context, serverID uint, tagName string) error

	// Metadata sub-tables (INSERT OR REPLACE keyed by server_id)
	UpsertNetwork(ctx context.Context, n *ServerNetwork) error
	UpsertHardware(ctx context.Context, h *ServerHardware) error
	UpsertOS(ctx context.Context, o *ServerOS) error

	// Software
	ReplaceSoftware(ctx context.Context, serverID uint, software []Software) error
	GetSoftware(ctx context.Context, serverID uint) ([]Software, error)

	// Joined read (all child tables)
	GetServerView(ctx context.Context, id uint) (*ServerView, error)
	GetServerViewByUUID(ctx context.Context, uuid string) (*ServerView, error)
	ListServerViews(ctx context.Context) ([]ServerView, error)

	// Misc
	Flush(ctx context.Context) error

	// Connection logs
	CreateConnectionLog(ctx context.Context, log *ConnectionLog) error
	UpdateConnectionLog(ctx context.Context, log *ConnectionLog) error
	GetConnectionLogs(ctx context.Context, serverID uint) ([]ConnectionLog, error)

	// Terminal Preferences
	GetTerminalPreference(ctx context.Context) (*TerminalPreference, error)
	SaveTerminalPreference(ctx context.Context, pref *TerminalPreference) error
}

// ServerService is the application-level contract consumed by CLI and API.
type ServerService interface {
	AddServer(ctx context.Context, server *Server) error
	GetServer(ctx context.Context, id uint) (*ServerView, error)
	ListServers(ctx context.Context) ([]ServerView, error)
	UpdateServer(ctx context.Context, server *Server) error
	RemoveServer(ctx context.Context, id uint) error
	RenameServer(ctx context.Context, id uint, newName string) error
	ScanInventory(ctx context.Context, id uint) error
	FlushServers(ctx context.Context) error

	// Metadata upserts
	UpsertServerNetwork(ctx context.Context, n *ServerNetwork) error
	UpsertServerHardware(ctx context.Context, h *ServerHardware) error
	UpsertServerOS(ctx context.Context, o *ServerOS) error
	ReplaceSoftware(ctx context.Context, serverID uint, software []Software) error

	// Terminal Preferences
	GetTerminalPreference(ctx context.Context) (*TerminalPreference, error)
	SaveTerminalPreference(ctx context.Context, pref *TerminalPreference) error

	// Connection tracking
	LogConnectionStart(ctx context.Context, server *Server) (*ConnectionLog, error)
	LogConnectionEnd(ctx context.Context, log *ConnectionLog, err error) error
	GetConnectionHistory(ctx context.Context, serverID uint) ([]ConnectionLog, error)
	ToggleFavorite(ctx context.Context, id uint) (bool, error)
}
