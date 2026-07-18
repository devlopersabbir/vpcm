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

func NewSQLiteRepository(db *sql.DB) ServerRepository {
	repo := &sqliteServerRepository{db: db}
	if err := repo.migrate(); err != nil {
		panic("failed to migrate SQLite database: " + err.Error())
	}
	return repo
}

func (r *sqliteServerRepository) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS servers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid TEXT UNIQUE NOT NULL,
		name TEXT NOT NULL,
		host TEXT NOT NULL,
		port INTEGER DEFAULT 22,
		username TEXT NOT NULL,
		auth_type TEXT,
		auth_secret TEXT,
		provider TEXT,
		region TEXT,
		created_at DATETIME,
		updated_at DATETIME,
		last_seen DATETIME,
		tags TEXT,
		software TEXT
	);`
	_, err := r.db.Exec(query)
	return err
}

func (r *sqliteServerRepository) Create(ctx context.Context, s *Server) error {
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()

	tagsJSON, _ := json.Marshal(s.Tags)
	softwareJSON, _ := json.Marshal(s.Software)

	var lastSeen interface{}
	if s.LastSeen != nil {
		lastSeen = *s.LastSeen
	}

	query := `
	INSERT INTO servers (uuid, name, host, port, username, auth_type, auth_secret, provider, region, created_at, updated_at, last_seen, tags, software)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, s.UUID, s.Name, s.Host, s.Port, s.Username, s.AuthType, s.AuthSecret, s.Provider, s.Region, s.CreatedAt, s.UpdatedAt, lastSeen, string(tagsJSON), string(softwareJSON))
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
	query := `SELECT id, uuid, name, host, port, username, auth_type, auth_secret, provider, region, created_at, updated_at, last_seen, tags, software FROM servers WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	var s Server
	var tagsStr, softwareStr sql.NullString
	var lastSeenTime sql.NullTime

	err := row.Scan(&s.ID, &s.UUID, &s.Name, &s.Host, &s.Port, &s.Username, &s.AuthType, &s.AuthSecret, &s.Provider, &s.Region, &s.CreatedAt, &s.UpdatedAt, &lastSeenTime, &tagsStr, &softwareStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrServerNotFound
		}
		return nil, err
	}

	if lastSeenTime.Valid {
		s.LastSeen = &lastSeenTime.Time
	}

	if tagsStr.Valid && tagsStr.String != "" {
		_ = json.Unmarshal([]byte(tagsStr.String), &s.Tags)
	}
	if s.Tags == nil {
		s.Tags = []Tag{}
	}

	if softwareStr.Valid && softwareStr.String != "" {
		_ = json.Unmarshal([]byte(softwareStr.String), &s.Software)
	}
	if s.Software == nil {
		s.Software = []Software{}
	}

	return &s, nil
}

func (r *sqliteServerRepository) GetByUUID(ctx context.Context, uuid string) (*Server, error) {
	query := `SELECT id, uuid, name, host, port, username, auth_type, auth_secret, provider, region, created_at, updated_at, last_seen, tags, software FROM servers WHERE uuid = ?`
	row := r.db.QueryRowContext(ctx, query, uuid)

	var s Server
	var tagsStr, softwareStr sql.NullString
	var lastSeenTime sql.NullTime

	err := row.Scan(&s.ID, &s.UUID, &s.Name, &s.Host, &s.Port, &s.Username, &s.AuthType, &s.AuthSecret, &s.Provider, &s.Region, &s.CreatedAt, &s.UpdatedAt, &lastSeenTime, &tagsStr, &softwareStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrServerNotFound
		}
		return nil, err
	}

	if lastSeenTime.Valid {
		s.LastSeen = &lastSeenTime.Time
	}

	if tagsStr.Valid && tagsStr.String != "" {
		_ = json.Unmarshal([]byte(tagsStr.String), &s.Tags)
	}
	if s.Tags == nil {
		s.Tags = []Tag{}
	}

	if softwareStr.Valid && softwareStr.String != "" {
		_ = json.Unmarshal([]byte(softwareStr.String), &s.Software)
	}
	if s.Software == nil {
		s.Software = []Software{}
	}

	return &s, nil
}

func (r *sqliteServerRepository) List(ctx context.Context) ([]Server, error) {
	query := `SELECT id, uuid, name, host, port, username, auth_type, auth_secret, provider, region, created_at, updated_at, last_seen, tags, software FROM servers`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servers []Server
	for rows.Next() {
		var s Server
		var tagsStr, softwareStr sql.NullString
		var lastSeenTime sql.NullTime

		err := rows.Scan(&s.ID, &s.UUID, &s.Name, &s.Host, &s.Port, &s.Username, &s.AuthType, &s.AuthSecret, &s.Provider, &s.Region, &s.CreatedAt, &s.UpdatedAt, &lastSeenTime, &tagsStr, &softwareStr)
		if err != nil {
			return nil, err
		}

		if lastSeenTime.Valid {
			s.LastSeen = &lastSeenTime.Time
		}

		if tagsStr.Valid && tagsStr.String != "" {
			_ = json.Unmarshal([]byte(tagsStr.String), &s.Tags)
		}
		if s.Tags == nil {
			s.Tags = []Tag{}
		}

		if softwareStr.Valid && softwareStr.String != "" {
			_ = json.Unmarshal([]byte(softwareStr.String), &s.Software)
		}
		if s.Software == nil {
			s.Software = []Software{}
		}

		servers = append(servers, s)
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
	softwareJSON, _ := json.Marshal(s.Software)

	var lastSeen interface{}
	if s.LastSeen != nil {
		lastSeen = *s.LastSeen
	}

	query := `
	UPDATE servers SET uuid = ?, name = ?, host = ?, port = ?, username = ?, auth_type = ?, auth_secret = ?, provider = ?, region = ?, updated_at = ?, last_seen = ?, tags = ?, software = ?
	WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, s.UUID, s.Name, s.Host, s.Port, s.Username, s.AuthType, s.AuthSecret, s.Provider, s.Region, s.UpdatedAt, lastSeen, string(tagsJSON), string(softwareJSON), s.ID)
	return err
}

func (r *sqliteServerRepository) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM servers WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err == nil && affected == 0 {
		return ErrServerNotFound
	}
	return nil
}

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

	var newTags []Tag
	for _, t := range s.Tags {
		if t.Name != tagName {
			newTags = append(newTags, t)
		}
	}
	s.Tags = newTags
	return r.Update(ctx, s)
}

func (r *sqliteServerRepository) Flush(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM servers`)
	if err != nil {
		return err
	}
	_, _ = r.db.ExecContext(ctx, `DELETE FROM sqlite_sequence WHERE name = 'servers'`)
	return nil
}
