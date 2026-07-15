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
}

type ServerService interface {
	AddServer(ctx context.Context, server *Server) error
	GetServer(ctx context.Context, id uint) (*Server, error)
	ListServers(ctx context.Context) ([]Server, error)
	UpdateServer(ctx context.Context, server *Server) error
	RemoveServer(ctx context.Context, id uint) error
	RenameServer(ctx context.Context, id uint, newName string) error
	ScanInventory(ctx context.Context, id uint) error
}
