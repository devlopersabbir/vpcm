package inventory

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

type sqliteServerRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) (ServerRepository, error) {
	r := &sqliteServerRepository{db: db}
	if err := r.migrate(); err != nil {
		return nil, err
	}
	return r, nil
}

// ─── Schema ───────────────────────────────────────────────────────────────────

func (r *sqliteServerRepository) migrate() error {
	_, err := r.db.Exec(`
	PRAGMA foreign_keys = ON;

	CREATE TABLE IF NOT EXISTS servers (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid        TEXT UNIQUE NOT NULL,
		name        TEXT NOT NULL,
		host        TEXT NOT NULL,
		port        INTEGER DEFAULT 22,
		username    TEXT NOT NULL,
		auth_type   TEXT,
		auth_secret TEXT,
		provider    TEXT,
		created_at  DATETIME,
		updated_at  DATETIME,
		last_seen   DATETIME,
		is_favorite INTEGER DEFAULT 0,
		tags        TEXT
	);

	CREATE TABLE IF NOT EXISTS server_network (
		id                INTEGER PRIMARY KEY AUTOINCREMENT,
		server_id         INTEGER UNIQUE NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
		hostname          TEXT,
		public_ip         TEXT,
		private_ip        TEXT,
		mac_address       TEXT,
		region            TEXT,
		availability_zone TEXT
	);

	CREATE TABLE IF NOT EXISTS server_hardware (
		id              INTEGER PRIMARY KEY AUTOINCREMENT,
		server_id       INTEGER UNIQUE NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
		cpu_model       TEXT,
		cpu_cores       INTEGER DEFAULT 0,
		ram_total       TEXT,
		swap_total      TEXT,
		disk_total      TEXT,
		virtualization  TEXT,
		instance_type   TEXT,
		serial_number   TEXT,
		bios_version    TEXT,
		uptime          TEXT
	);

	CREATE TABLE IF NOT EXISTS server_os (
		id              INTEGER PRIMARY KEY AUTOINCREMENT,
		server_id       INTEGER UNIQUE NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
		os_family       TEXT,
		os_version      TEXT,
		kernel_version  TEXT,
		architecture    TEXT,
		init_system     TEXT,
		timezone        TEXT,
		locale          TEXT,
		package_manager TEXT
	);

	CREATE TABLE IF NOT EXISTS software (
		id        INTEGER PRIMARY KEY AUTOINCREMENT,
		server_id INTEGER NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
		name      TEXT NOT NULL,
		version   TEXT
	);

	CREATE TABLE IF NOT EXISTS connection_logs (
		id            INTEGER PRIMARY KEY AUTOINCREMENT,
		server_id     INTEGER,
		server_name   TEXT NOT NULL,
		username      TEXT NOT NULL,
		host          TEXT NOT NULL,
		logged_in_at  DATETIME NOT NULL,
		logged_out_at DATETIME,
		duration      TEXT,
		status        TEXT NOT NULL,
		error_message TEXT
	);`)
	if err != nil {
		return err
	}
	_, _ = r.db.Exec("ALTER TABLE servers ADD COLUMN is_favorite INTEGER DEFAULT 0")
	return nil
}

// ─── Core CRUD ────────────────────────────────────────────────────────────────

func (r *sqliteServerRepository) Create(ctx context.Context, s *Server) error {
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	tagsJSON, _ := json.Marshal(s.Tags)

	var lastSeen interface{}
	if s.LastSeen != nil {
		lastSeen = *s.LastSeen
	}

	isFavorite := 0
	if s.IsFavorite {
		isFavorite = 1
	}

	res, err := r.db.ExecContext(ctx, `
	INSERT INTO servers (uuid, name, host, port, username, auth_type, auth_secret, provider, created_at, updated_at, last_seen, is_favorite, tags)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		s.UUID, s.Name, s.Host, s.Port, s.Username, s.AuthType, s.AuthSecret,
		s.Provider, s.CreatedAt, s.UpdatedAt, lastSeen, isFavorite, string(tagsJSON))
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		s.ID = uint(id)
	}
	return nil
}

func (r *sqliteServerRepository) GetByID(ctx context.Context, id uint) (*Server, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, uuid, name, host, port, username, auth_type, auth_secret, provider, created_at, updated_at, last_seen, is_favorite, tags
		 FROM servers WHERE id = ?`, id)
	return r.scanServer(row)
}

func (r *sqliteServerRepository) GetByUUID(ctx context.Context, uuid string) (*Server, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, uuid, name, host, port, username, auth_type, auth_secret, provider, created_at, updated_at, last_seen, is_favorite, tags
		 FROM servers WHERE uuid = ?`, uuid)
	return r.scanServer(row)
}

func (r *sqliteServerRepository) List(ctx context.Context) ([]Server, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, uuid, name, host, port, username, auth_type, auth_secret, provider, created_at, updated_at, last_seen, is_favorite, tags
		 FROM servers`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servers []Server
	for rows.Next() {
		s, err := r.scanServerRow(rows)
		if err != nil {
			return nil, err
		}
		servers = append(servers, *s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if servers == nil {
		servers = []Server{}
	}
	return servers, nil
}

func (r *sqliteServerRepository) Update(ctx context.Context, s *Server) error {
	s.UpdatedAt = time.Now()
	tagsJSON, _ := json.Marshal(s.Tags)

	var lastSeen interface{}
	if s.LastSeen != nil {
		lastSeen = *s.LastSeen
	}

	isFavorite := 0
	if s.IsFavorite {
		isFavorite = 1
	}

	_, err := r.db.ExecContext(ctx, `
	UPDATE servers SET name = ?, auth_type = ?, auth_secret = ?, provider = ?, updated_at = ?, last_seen = ?, is_favorite = ?, tags = ?
	WHERE id = ?`,
		s.Name, s.AuthType, s.AuthSecret, s.Provider, s.UpdatedAt, lastSeen, isFavorite, string(tagsJSON), s.ID)
	return err
}

func (r *sqliteServerRepository) Delete(ctx context.Context, id uint) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM servers WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err == nil && affected == 0 {
		return ErrServerNotFound
	}
	return nil
}

// ─── Tags ─────────────────────────────────────────────────────────────────────

func (r *sqliteServerRepository) AddTag(ctx context.Context, serverID uint, tagName string) error {
	s, err := r.GetByID(ctx, serverID)
	if err != nil {
		return err
	}
	for _, t := range s.Tags {
		if t.Name == tagName {
			return nil
		}
	}
	s.Tags = append(s.Tags, Tag{Name: tagName})
	return r.Update(ctx, s)
}

func (r *sqliteServerRepository) RemoveTag(ctx context.Context, serverID uint, tagName string) error {
	s, err := r.GetByID(ctx, serverID)
	if err != nil {
		return err
	}
	var next []Tag
	for _, t := range s.Tags {
		if t.Name != tagName {
			next = append(next, t)
		}
	}
	s.Tags = next
	return r.Update(ctx, s)
}

// ─── Metadata Upserts ─────────────────────────────────────────────────────────

func (r *sqliteServerRepository) UpsertNetwork(ctx context.Context, n *ServerNetwork) error {
	_, err := r.db.ExecContext(ctx, `
	INSERT INTO server_network (server_id, hostname, public_ip, private_ip, mac_address, region, availability_zone)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(server_id) DO UPDATE SET
		hostname          = excluded.hostname,
		public_ip         = excluded.public_ip,
		private_ip        = excluded.private_ip,
		mac_address       = excluded.mac_address,
		region            = excluded.region,
		availability_zone = excluded.availability_zone`,
		n.ServerID, n.Hostname, n.PublicIP, n.PrivateIP, n.MACAddress, n.Region, n.AvailabilityZone)
	return err
}

func (r *sqliteServerRepository) UpsertHardware(ctx context.Context, h *ServerHardware) error {
	_, err := r.db.ExecContext(ctx, `
	INSERT INTO server_hardware (server_id, cpu_model, cpu_cores, ram_total, swap_total, disk_total, virtualization, instance_type, serial_number, bios_version, uptime)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(server_id) DO UPDATE SET
		cpu_model      = excluded.cpu_model,
		cpu_cores      = excluded.cpu_cores,
		ram_total      = excluded.ram_total,
		swap_total     = excluded.swap_total,
		disk_total     = excluded.disk_total,
		virtualization = excluded.virtualization,
		instance_type  = excluded.instance_type,
		serial_number  = excluded.serial_number,
		bios_version   = excluded.bios_version,
		uptime         = excluded.uptime`,
		h.ServerID, h.CPUModel, h.CPUCores, h.RAMTotal, h.SwapTotal, h.DiskTotal,
		h.Virtualization, h.InstanceType, h.SerialNumber, h.BIOSVersion, h.Uptime)
	return err
}

func (r *sqliteServerRepository) UpsertOS(ctx context.Context, o *ServerOS) error {
	_, err := r.db.ExecContext(ctx, `
	INSERT INTO server_os (server_id, os_family, os_version, kernel_version, architecture, init_system, timezone, locale, package_manager)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(server_id) DO UPDATE SET
		os_family       = excluded.os_family,
		os_version      = excluded.os_version,
		kernel_version  = excluded.kernel_version,
		architecture    = excluded.architecture,
		init_system     = excluded.init_system,
		timezone        = excluded.timezone,
		locale          = excluded.locale,
		package_manager = excluded.package_manager`,
		o.ServerID, o.OSFamily, o.OSVersion, o.KernelVersion, o.Architecture,
		o.InitSystem, o.Timezone, o.Locale, o.PackageManager)
	return err
}

// ─── Software ─────────────────────────────────────────────────────────────────

func (r *sqliteServerRepository) ReplaceSoftware(ctx context.Context, serverID uint, software []Software) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck

	if _, err := tx.ExecContext(ctx, `DELETE FROM software WHERE server_id = ?`, serverID); err != nil {
		return err
	}
	for _, sw := range software {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO software (server_id, name, version) VALUES (?, ?, ?)`,
			serverID, sw.Name, sw.Version); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *sqliteServerRepository) GetSoftware(ctx context.Context, serverID uint) ([]Software, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, server_id, name, version FROM software WHERE server_id = ?`, serverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Software
	for rows.Next() {
		var sw Software
		if err := rows.Scan(&sw.ID, &sw.ServerID, &sw.Name, &sw.Version); err != nil {
			return nil, err
		}
		list = append(list, sw)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if list == nil {
		list = []Software{}
	}
	return list, nil
}

// ─── Joined Views ─────────────────────────────────────────────────────────────

const serverViewQuery = `
SELECT
	s.id, s.uuid, s.name, s.host, s.port, s.username, s.auth_type, s.auth_secret, s.provider,
	s.created_at, s.updated_at, s.last_seen, s.is_favorite, s.tags,
	-- network
	n.id, n.hostname, n.public_ip, n.private_ip, n.mac_address, n.region, n.availability_zone,
	-- hardware
	h.id, h.cpu_model, h.cpu_cores, h.ram_total, h.swap_total, h.disk_total,
	       h.virtualization, h.instance_type, h.serial_number, h.bios_version, h.uptime,
	-- os
	o.id, o.os_family, o.os_version, o.kernel_version, o.architecture,
	      o.init_system, o.timezone, o.locale, o.package_manager
FROM servers s
LEFT JOIN server_network  n ON n.server_id = s.id
LEFT JOIN server_hardware h ON h.server_id = s.id
LEFT JOIN server_os       o ON o.server_id = s.id`

func (r *sqliteServerRepository) GetServerView(ctx context.Context, id uint) (*ServerView, error) {
	row := r.db.QueryRowContext(ctx, serverViewQuery+` WHERE s.id = ?`, id)
	v, err := r.scanView(row)
	if err != nil {
		return nil, err
	}
	sw, _ := r.GetSoftware(ctx, v.ID)
	v.Software = sw
	return v, nil
}

func (r *sqliteServerRepository) GetServerViewByUUID(ctx context.Context, uuid string) (*ServerView, error) {
	row := r.db.QueryRowContext(ctx, serverViewQuery+` WHERE s.uuid = ?`, uuid)
	v, err := r.scanView(row)
	if err != nil {
		return nil, err
	}
	sw, _ := r.GetSoftware(ctx, v.ID)
	v.Software = sw
	return v, nil
}

func (r *sqliteServerRepository) ListServerViews(ctx context.Context) ([]ServerView, error) {
	rows, err := r.db.QueryContext(ctx, serverViewQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var views []ServerView
	for rows.Next() {
		v, err := r.scanViewRow(rows)
		if err != nil {
			return nil, err
		}
		sw, _ := r.GetSoftware(ctx, v.ID)
		v.Software = sw
		views = append(views, *v)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if views == nil {
		views = []ServerView{}
	}
	return views, nil
}

// ─── Misc ─────────────────────────────────────────────────────────────────────

func (r *sqliteServerRepository) Flush(ctx context.Context) error {
	for _, tbl := range []string{"software", "server_os", "server_hardware", "server_network", "connection_logs", "servers"} {
		if _, err := r.db.ExecContext(ctx, `DELETE FROM `+tbl); err != nil {
			return err
		}
	}
	_, _ = r.db.ExecContext(ctx, `DELETE FROM sqlite_sequence`)
	return nil
}

// ─── Connection Logs ──────────────────────────────────────────────────────────

func (r *sqliteServerRepository) CreateConnectionLog(ctx context.Context, s *ConnectionLog) error {
	var loggedOut interface{}
	if s.LoggedOutAt != nil {
		loggedOut = *s.LoggedOutAt
	}
	res, err := r.db.ExecContext(ctx, `
	INSERT INTO connection_logs (server_id, server_name, username, host, logged_in_at, logged_out_at, duration, status, error_message)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		s.ServerID, s.ServerName, s.Username, s.Host, s.LoggedInAt, loggedOut, s.Duration, s.Status, s.ErrorMessage)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		s.ID = uint(id)
	}
	return nil
}

func (r *sqliteServerRepository) UpdateConnectionLog(ctx context.Context, s *ConnectionLog) error {
	var loggedOut interface{}
	if s.LoggedOutAt != nil {
		loggedOut = *s.LoggedOutAt
	}
	_, err := r.db.ExecContext(ctx, `
	UPDATE connection_logs SET server_id = ?, server_name = ?, username = ?, host = ?,
		logged_in_at = ?, logged_out_at = ?, duration = ?, status = ?, error_message = ?
	WHERE id = ?`,
		s.ServerID, s.ServerName, s.Username, s.Host, s.LoggedInAt,
		loggedOut, s.Duration, s.Status, s.ErrorMessage, s.ID)
	return err
}

func (r *sqliteServerRepository) GetConnectionLogs(ctx context.Context, serverID uint) ([]ConnectionLog, error) {
	var rows *sql.Rows
	var err error
	if serverID > 0 {
		rows, err = r.db.QueryContext(ctx, `
		SELECT id, server_id, server_name, username, host, logged_in_at, logged_out_at, duration, status, error_message
		FROM connection_logs WHERE server_id = ? ORDER BY logged_in_at DESC`, serverID)
	} else {
		rows, err = r.db.QueryContext(ctx, `
		SELECT id, server_id, server_name, username, host, logged_in_at, logged_out_at, duration, status, error_message
		FROM connection_logs ORDER BY logged_in_at DESC`)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []ConnectionLog
	for rows.Next() {
		var l ConnectionLog
		var loggedOut sql.NullTime
		var duration, errMsg sql.NullString
		if err := rows.Scan(&l.ID, &l.ServerID, &l.ServerName, &l.Username, &l.Host,
			&l.LoggedInAt, &loggedOut, &duration, &l.Status, &errMsg); err != nil {
			return nil, err
		}
		if loggedOut.Valid {
			l.LoggedOutAt = &loggedOut.Time
		}
		l.Duration = duration.String
		l.ErrorMessage = errMsg.String
		logs = append(logs, l)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if logs == nil {
		logs = []ConnectionLog{}
	}
	return logs, nil
}

// ─── Scan Helpers ─────────────────────────────────────────────────────────────

type rowScanner interface {
	Scan(dest ...any) error
}

func (r *sqliteServerRepository) scanServer(row *sql.Row) (*Server, error) {
	var s Server
	var tagsStr sql.NullString
	var lastSeen sql.NullTime
	var isFavorite int
	err := row.Scan(&s.ID, &s.UUID, &s.Name, &s.Host, &s.Port, &s.Username,
		&s.AuthType, &s.AuthSecret, &s.Provider, &s.CreatedAt, &s.UpdatedAt, &lastSeen, &isFavorite, &tagsStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrServerNotFound
		}
		return nil, err
	}
	if lastSeen.Valid {
		s.LastSeen = &lastSeen.Time
	}
	s.IsFavorite = (isFavorite != 0)
	if tagsStr.Valid && tagsStr.String != "" {
		_ = json.Unmarshal([]byte(tagsStr.String), &s.Tags)
	}
	if s.Tags == nil {
		s.Tags = []Tag{}
	}
	return &s, nil
}

func (r *sqliteServerRepository) scanServerRow(rows *sql.Rows) (*Server, error) {
	var s Server
	var tagsStr sql.NullString
	var lastSeen sql.NullTime
	var isFavorite int
	err := rows.Scan(&s.ID, &s.UUID, &s.Name, &s.Host, &s.Port, &s.Username,
		&s.AuthType, &s.AuthSecret, &s.Provider, &s.CreatedAt, &s.UpdatedAt, &lastSeen, &isFavorite, &tagsStr)
	if err != nil {
		return nil, err
	}
	if lastSeen.Valid {
		s.LastSeen = &lastSeen.Time
	}
	s.IsFavorite = (isFavorite != 0)
	if tagsStr.Valid && tagsStr.String != "" {
		_ = json.Unmarshal([]byte(tagsStr.String), &s.Tags)
	}
	if s.Tags == nil {
		s.Tags = []Tag{}
	}
	return &s, nil
}

func scanViewColumns(scanner rowScanner) (*ServerView, error) {
	var v ServerView
	var tagsStr sql.NullString
	var lastSeen sql.NullTime

	// nullable network columns
	var nID sql.NullInt64
	var nHostname, nPublicIP, nPrivateIP, nMAC, nRegion, nAZ sql.NullString

	// nullable hardware columns
	var hID sql.NullInt64
	var hCPUModel, hRAM, hSwap, hDisk, hVirt, hInstType, hSerial, hBIOS, hUptime sql.NullString
	var hCPUCores sql.NullInt64

	// nullable os columns
	var oID sql.NullInt64
	var oFamily, oVersion, oKernel, oArch, oInit, oTZ, oLocale, oPkgMgr sql.NullString

	var isFavorite int

	err := scanner.Scan(
		&v.ID, &v.UUID, &v.Name, &v.Host, &v.Port, &v.Username, &v.AuthType, &v.AuthSecret, &v.Provider,
		&v.CreatedAt, &v.UpdatedAt, &lastSeen, &isFavorite, &tagsStr,
		// network
		&nID, &nHostname, &nPublicIP, &nPrivateIP, &nMAC, &nRegion, &nAZ,
		// hardware
		&hID, &hCPUModel, &hCPUCores, &hRAM, &hSwap, &hDisk, &hVirt, &hInstType, &hSerial, &hBIOS, &hUptime,
		// os
		&oID, &oFamily, &oVersion, &oKernel, &oArch, &oInit, &oTZ, &oLocale, &oPkgMgr,
	)
	if err != nil {
		return nil, err
	}

	if lastSeen.Valid {
		v.LastSeen = &lastSeen.Time
	}
	v.IsFavorite = (isFavorite != 0)
	if tagsStr.Valid && tagsStr.String != "" {
		_ = json.Unmarshal([]byte(tagsStr.String), &v.Tags)
	}
	if v.Tags == nil {
		v.Tags = []Tag{}
	}

	if nID.Valid {
		v.Network = &ServerNetwork{
			ServerID:         v.ID,
			Hostname:         nHostname.String,
			PublicIP:         nPublicIP.String,
			PrivateIP:        nPrivateIP.String,
			MACAddress:       nMAC.String,
			Region:           nRegion.String,
			AvailabilityZone: nAZ.String,
		}
	}
	if hID.Valid {
		v.Hardware = &ServerHardware{
			ServerID:       v.ID,
			CPUModel:       hCPUModel.String,
			CPUCores:       int(hCPUCores.Int64),
			RAMTotal:       hRAM.String,
			SwapTotal:      hSwap.String,
			DiskTotal:      hDisk.String,
			Virtualization: hVirt.String,
			InstanceType:   hInstType.String,
			SerialNumber:   hSerial.String,
			BIOSVersion:    hBIOS.String,
			Uptime:         hUptime.String,
		}
	}
	if oID.Valid {
		v.OS = &ServerOS{
			ServerID:       v.ID,
			OSFamily:       oFamily.String,
			OSVersion:      oVersion.String,
			KernelVersion:  oKernel.String,
			Architecture:   oArch.String,
			InitSystem:     oInit.String,
			Timezone:       oTZ.String,
			Locale:         oLocale.String,
			PackageManager: oPkgMgr.String,
		}
	}

	v.Software = []Software{}
	return &v, nil
}

func (r *sqliteServerRepository) scanView(row *sql.Row) (*ServerView, error) {
	v, err := scanViewColumns(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrServerNotFound
		}
		return nil, err
	}
	return v, nil
}

func (r *sqliteServerRepository) scanViewRow(rows *sql.Rows) (*ServerView, error) {
	return scanViewColumns(rows)
}
