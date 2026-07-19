package inventory

import "context"

type ServerRepository interface {
	Create(ctx context.Context, server *Server) error
	GetByID(ctx context.Context, id uint) (*Server, error)
	GetByUUID(ctx context.Context, uuid string) (*Server, error)
	List(ctx context.Context) ([]Server, error)
	Update(ctx context.Context, server *Server) error
	Delete(ctx context.Context, id uint) error
	AddTag(ctx context.Context, serverID uint, tagName string) error
	RemoveTag(ctx context.Context, serverID uint, tagName string) error
	Flush(ctx context.Context) error
	CreateConnectionLog(ctx context.Context, log *ConnectionLog) error
	UpdateConnectionLog(ctx context.Context, log *ConnectionLog) error
	GetConnectionLogs(ctx context.Context, serverID uint) ([]ConnectionLog, error)
}

type ServerService interface {
	AddServer(ctx context.Context, server *Server) error
	GetServer(ctx context.Context, id uint) (*Server, error)
	ListServers(ctx context.Context) ([]Server, error)
	UpdateServer(ctx context.Context, server *Server) error
	RemoveServer(ctx context.Context, id uint) error
	RenameServer(ctx context.Context, id uint, newName string) error
	ScanInventory(ctx context.Context, id uint) error
	FlushServers(ctx context.Context) error
	LogConnectionStart(ctx context.Context, server *Server) (*ConnectionLog, error)
	LogConnectionEnd(ctx context.Context, log *ConnectionLog, err error) error
	GetConnectionHistory(ctx context.Context, serverID uint) ([]ConnectionLog, error)
}
