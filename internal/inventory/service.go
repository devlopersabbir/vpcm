package inventory

import (
	"context"
	"log/slog"
	"time"

	"github.com/devlopersabbir/vpcm/internal/events"
	"github.com/devlopersabbir/vpcm/internal/ssh"
	"github.com/google/uuid"
)

type serverService struct {
	repo ServerRepository
}

func NewService(repo ServerRepository) ServerService {
	return &serverService{repo: repo}
}

// ─── Core ─────────────────────────────────────────────────────────────────────

func (s *serverService) AddServer(ctx context.Context, server *Server) error {
	if server.UUID == "" {
		server.UUID = uuid.NewString()
	}
	return s.repo.Create(ctx, server)
}

func (s *serverService) GetServer(ctx context.Context, id uint) (*ServerView, error) {
	return s.repo.GetServerView(ctx, id)
}

func (s *serverService) ListServers(ctx context.Context) ([]ServerView, error) {
	return s.repo.ListServerViews(ctx)
}

func (s *serverService) UpdateServer(ctx context.Context, server *Server) error {
	slog.Debug("Updating server", "id", server.ID)

	existing, err := s.repo.GetByID(ctx, server.ID)
	if err == nil && existing != nil {
		// Protect immutable connection fields
		server.Host = existing.Host
		server.Username = existing.Username
		server.Port = existing.Port
	}

	if err := s.repo.Update(ctx, server); err != nil {
		return err
	}
	events.Publish(events.Event{Type: "ServerUpdated", Payload: server})
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
	events.Publish(events.Event{Type: "ServerDeleted", Payload: server})
	return nil
}

func (s *serverService) RenameServer(ctx context.Context, id uint, newName string) error {
	server, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	server.Name = newName
	return s.repo.Update(ctx, server)
}

func (s *serverService) ScanInventory(ctx context.Context, id uint) error {
	slog.Info("Scanning inventory for server", "id", id)

	server, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	sshSvc := ssh.NewService(10 * time.Second)
	var client ssh.Client
	if server.AuthType == "key" {
		client, err = sshSvc.Connect(ctx, server.Host, server.Port, server.Username, server.AuthType, server.AuthSecret)
	} else {
		client, err = sshSvc.Connect(ctx, server.Host, server.Port, server.Username, "password", server.AuthSecret)
	}
	if err != nil {
		return err
	}
	defer client.Close()

	if server.Provider == "" || server.Provider == "Generic VPS" {
		server.Provider = DetectProvider(ctx, client, server.Host)
		_ = s.repo.Update(ctx, server)
	}

	osFamily, osVersion := DetectOS(ctx, client)
	cpuModel, cpuCores, ramTotal, diskTotal := DetectSpecs(ctx, client)
	si := DetectServerInfo(ctx, client)
	swList := DetectSoftware(ctx, client)
	_ = s.repo.ReplaceSoftware(ctx, server.ID, swList)

	_ = s.repo.UpsertNetwork(ctx, &ServerNetwork{
		ServerID:         server.ID,
		Hostname:         si.Hostname,
		PublicIP:         si.PublicIP,
		PrivateIP:        si.PrivateIP,
		MACAddress:       si.MACAddress,
		Region:           si.Region,
		AvailabilityZone: si.AvailabilityZone,
	})

	_ = s.repo.UpsertHardware(ctx, &ServerHardware{
		ServerID:       server.ID,
		CPUModel:       cpuModel,
		CPUCores:       cpuCores,
		RAMTotal:       ramTotal,
		SwapTotal:      si.SwapTotal,
		DiskTotal:      diskTotal,
		Virtualization: si.Virtualization,
		InstanceType:   si.InstanceType,
		SerialNumber:   si.SerialNumber,
		BIOSVersion:    si.BIOSVersion,
		Uptime:         si.Uptime,
	})

	_ = s.repo.UpsertOS(ctx, &ServerOS{
		ServerID:       server.ID,
		OSFamily:       osFamily,
		OSVersion:      osVersion,
		KernelVersion:  si.KernelVersion,
		Architecture:   si.Architecture,
		InitSystem:     si.InitSystem,
		Timezone:       si.Timezone,
		Locale:         si.Locale,
		PackageManager: si.PackageManager,
	})

	events.Publish(events.Event{Type: "InventoryScanned", Payload: id})
	return nil
}

func (s *serverService) FlushServers(ctx context.Context) error {
	slog.Warn("Flushing all servers from database")
	return s.repo.Flush(ctx)
}

// ─── Metadata Upserts ────────────────────────────────────────────────────────

func (s *serverService) UpsertServerNetwork(ctx context.Context, n *ServerNetwork) error {
	slog.Debug("Upserting server network info", "server_id", n.ServerID)
	return s.repo.UpsertNetwork(ctx, n)
}

func (s *serverService) UpsertServerHardware(ctx context.Context, h *ServerHardware) error {
	slog.Debug("Upserting server hardware info", "server_id", h.ServerID)
	return s.repo.UpsertHardware(ctx, h)
}

func (s *serverService) UpsertServerOS(ctx context.Context, o *ServerOS) error {
	slog.Debug("Upserting server OS info", "server_id", o.ServerID)
	return s.repo.UpsertOS(ctx, o)
}

func (s *serverService) ReplaceSoftware(ctx context.Context, serverID uint, software []Software) error {
	slog.Debug("Replacing software list", "server_id", serverID, "count", len(software))
	return s.repo.ReplaceSoftware(ctx, serverID, software)
}

// ─── Connection Logs ──────────────────────────────────────────────────────────

func (s *serverService) LogConnectionStart(ctx context.Context, server *Server) (*ConnectionLog, error) {
	now := time.Now()
	server.LastSeen = &now
	_ = s.repo.Update(ctx, server)

	log := &ConnectionLog{
		ServerID:   server.ID,
		ServerName: server.Name,
		Username:   server.Username,
		Host:       server.Host,
		LoggedInAt: now,
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
	if log.ID == 0 {
		return s.repo.CreateConnectionLog(ctx, log)
	}
	return s.repo.UpdateConnectionLog(ctx, log)
}

func (s *serverService) GetConnectionHistory(ctx context.Context, serverID uint) ([]ConnectionLog, error) {
	return s.repo.GetConnectionLogs(ctx, serverID)
}

func (s *serverService) ToggleFavorite(ctx context.Context, id uint) (bool, error) {
	server, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return false, err
	}
	server.IsFavorite = !server.IsFavorite
	if err := s.repo.Update(ctx, server); err != nil {
		return false, err
	}
	return server.IsFavorite, nil
}

func (s *serverService) GetTerminalPreference(ctx context.Context) (*TerminalPreference, error) {
	return s.repo.GetTerminalPreference(ctx)
}

func (s *serverService) SaveTerminalPreference(ctx context.Context, pref *TerminalPreference) error {
	return s.repo.SaveTerminalPreference(ctx, pref)
}
