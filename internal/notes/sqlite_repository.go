package notes

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type sqliteNoteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) NoteRepository {
	repo := &sqliteNoteRepository{db: db}
	if err := repo.migrate(); err != nil {
		panic("failed to migrate SQLite notes: " + err.Error())
	}
	return repo
}

func (r *sqliteNoteRepository) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		server_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		content TEXT,
		created_at DATETIME,
		updated_at DATETIME
	);`
	_, err := r.db.Exec(query)
	return err
}

func (r *sqliteNoteRepository) Create(ctx context.Context, n *Note) error {
	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()

	query := `
	INSERT INTO notes (server_id, title, content, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, n.ServerID, n.Title, n.Content, n.CreatedAt, n.UpdatedAt)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		n.ID = uint(id)
	}
	return nil
}

func (r *sqliteNoteRepository) GetByID(ctx context.Context, id uint) (*Note, error) {
	query := `SELECT id, server_id, title, content, created_at, updated_at FROM notes WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	var n Note
	err := row.Scan(&n.ID, &n.ServerID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoteNotFound
		}
		return nil, err
	}
	return &n, nil
}

func (r *sqliteNoteRepository) ListByServer(ctx context.Context, serverID uint) ([]Note, error) {
	query := `SELECT id, server_id, title, content, created_at, updated_at FROM notes WHERE server_id = ?`
	rows, err := r.db.QueryContext(ctx, query, serverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Note
	for rows.Next() {
		var n Note
		err := rows.Scan(&n.ID, &n.ServerID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if list == nil {
		list = []Note{}
	}
	return list, nil
}

func (r *sqliteNoteRepository) Update(ctx context.Context, n *Note) error {
	n.UpdatedAt = time.Now()
	query := `UPDATE notes SET title = ?, content = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, n.Title, n.Content, n.UpdatedAt, n.ID)
	return err
}

func (r *sqliteNoteRepository) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM notes WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err == nil && affected == 0 {
		return ErrNoteNotFound
	}
	return nil
}
