package database

import (
	"path/filepath"
	"testing"
)

func TestInitSQLiteCreatesParentDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "dir", "vpsm.db")

	db, err := InitSQLite(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("database is not usable: %v", err)
	}
}

// A zero busy timeout makes concurrent writers fail immediately with
// "database is locked" instead of waiting for the current writer.
func TestInitSQLiteSetsBusyTimeout(t *testing.T) {
	db, err := InitSQLite(filepath.Join(t.TempDir(), "vpsm.db"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer db.Close()

	var timeout int
	if err := db.QueryRow("PRAGMA busy_timeout").Scan(&timeout); err != nil {
		t.Fatalf("failed to read busy_timeout: %v", err)
	}
	if timeout < 1000 {
		t.Fatalf("got busy_timeout %dms, want at least 1000ms", timeout)
	}
}

// The pragma is carried in the DSN so that every connection the pool opens gets
// it, not just the first one.
func TestInitSQLiteAppliesBusyTimeoutToNewConnections(t *testing.T) {
	db, err := InitSQLite(filepath.Join(t.TempDir(), "vpsm.db"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer db.Close()

	// Hold one connection open so the second query must open a fresh one.
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed to begin: %v", err)
	}
	defer func() { _ = tx.Rollback() }()

	var timeout int
	if err := db.QueryRow("PRAGMA busy_timeout").Scan(&timeout); err != nil {
		t.Fatalf("failed to read busy_timeout on a second connection: %v", err)
	}
	if timeout < 1000 {
		t.Fatalf("got busy_timeout %dms on a new connection, want at least 1000ms", timeout)
	}
}

// inventory.ListServerViews queries the software table for each row while the
// outer result set is still open, so a pool capped at one connection deadlocks it
// and hangs GET /servers.
func TestInitSQLiteDoesNotCapThePoolAtOneConnection(t *testing.T) {
	db, err := InitSQLite(filepath.Join(t.TempDir(), "vpsm.db"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer db.Close()

	if got := db.Stats().MaxOpenConnections; got == 1 {
		t.Fatal("the pool is capped at 1 connection, which deadlocks nested queries")
	}
}

func TestInitSQLiteRejectsUnusablePath(t *testing.T) {
	// A path whose parent is an existing file cannot be created.
	dir := t.TempDir()
	blocker := filepath.Join(dir, "afile")
	if _, err := InitSQLite(blocker); err != nil {
		t.Fatalf("expected a plain file path to work: %v", err)
	}
	if _, err := InitSQLite(filepath.Join(blocker, "child.db")); err == nil {
		t.Fatal("expected an error when the parent path is a file")
	}
}
