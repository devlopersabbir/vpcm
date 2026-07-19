package inventory

import (
	"context"
	"log/slog"
	"time"

	"github.com/devlopersabbir/vpcm/internal/events"
)

type serverService struct {
	repo ServerRepository
}

func NewService(repo ServerRepository) ServerService {
	return &serverService{repo: repo}
}

func (s *serverService) AddServer(ctx context.Context, server *Server) error {
	slog.Debug("Adding server", "name", server.Name, "host", server.Host)
	if err := s.repo.Create(ctx, server); err != nil {
		return err
	}

	events.Publish(events.Event{
		Type:    "ServerAdded",
		Payload: server,
	})
	return nil
}

func (s *serverService) GetServer(ctx context.Context, id uint) (*Server, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *serverService) ListServers(ctx context.Context) ([]Server, error) {
	return s.repo.List(ctx)
}

func (s *serverService) UpdateServer(ctx context.Context, server *Server) error {
	slog.Debug("Updating server", "id", server.ID)

	existing, err := s.repo.GetByID(ctx, server.ID)
	if err == nil && existing != nil {
		// Prevent updates to critical connection/credentials settings
		server.Host = existing.Host
		server.Username = existing.Username
		server.Port = existing.Port
	}

	if err := s.repo.Update(ctx, server); err != nil {
		return err
	}

	events.Publish(events.Event{
		Type:    "ServerUpdated",
		Payload: server,
	})
	return nil
}

func (s *serverService) RemoveServer(ctx context.Context, id uint) error {
	slog.Debug("Removing server", "id", id)
	server, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	events.Publish(events.Event{
		Type:    "ServerDeleted",
		Payload: server,
	})
	return nil
}

func (s *serverService) RenameServer(ctx context.Context, id uint, newName string) error {
	server, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	server.Name = newName
	return s.UpdateServer(ctx, server)
}

func (s *serverService) ScanInventory(ctx context.Context, id uint) error {
	slog.Info("Scanning inventory for server", "id", id)
	// v0.0.1 scans are mocked/skeletons. We will implement basic trigger.
	events.Publish(events.Event{
		Type:    "InventoryScanned",
		Payload: id,
	})
	return nil
}

func (s *serverService) FlushServers(ctx context.Context) error {
	slog.Warn("Flushing all servers from database")
	return s.repo.Flush(ctx)
}

func (s *serverService) LogConnectionStart(ctx context.Context, server *Server) (*ConnectionLog, error) {
	log := &ConnectionLog{
		ServerID:   server.ID,
		ServerName: server.Name,
		Username:   server.Username,
		Host:       server.Host,
		LoggedInAt: time.Now(),
		Status:     "active",
	}
	if err := s.repo.CreateConnectionLog(ctx, log); err != nil {
		return nil, err
	}
	return log, nil
}

func (s *serverService) LogConnectionEnd(ctx context.Context, log *ConnectionLog, err error) error {
	now := time.Now()
	log.LoggedOutAt = &now
	log.Duration = now.Sub(log.LoggedInAt).Round(time.Second).String()
	if err != nil {
		log.Status = "failed"
		log.ErrorMessage = err.Error()
	} else {
		log.Status = "success"
	}
	return s.repo.UpdateConnectionLog(ctx, log)
}

func (s *serverService) GetConnectionHistory(ctx context.Context, serverID uint) ([]ConnectionLog, error) {
	return s.repo.GetConnectionLogs(ctx, serverID)
}
