package inventory

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"github.com/devlopersabbir/vpcm/internal/database"
)

// The repository is exercised through the real initializer so that these tests
// see the same connection pool and pragmas as the shipped binaries.
func newTestRepo(t *testing.T) (ServerRepository, *sql.DB) {
	t.Helper()

	db, err := database.InitSQLite(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	repo, err := NewSQLiteRepository(db)
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return repo, db
}

func mustCreate(t *testing.T, repo ServerRepository, s *Server) *Server {
	t.Helper()
	if err := repo.Create(context.Background(), s); err != nil {
		t.Fatalf("failed to create server %q: %v", s.Name, err)
	}
	return s
}

func TestSQLiteCreateAndList(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	mustCreate(t, repo, &Server{UUID: "u-1", Name: "web", Host: "10.0.0.1", Port: 22, Username: "root"})
	mustCreate(t, repo, &Server{UUID: "u-2", Name: "db", Host: "10.0.0.2", Port: 2222, Username: "admin"})

	servers, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(servers) != 2 {
		t.Fatalf("got %d servers, want 2", len(servers))
	}
	if servers[0].ID == 0 || servers[1].ID == 0 {
		t.Fatalf("expected autoincrement ids to be assigned, got %+v", servers)
	}
}

func TestSQLiteCreateRejectsDuplicateUUID(t *testing.T) {
	repo, _ := newTestRepo(t)

	mustCreate(t, repo, &Server{UUID: "same", Name: "a", Host: "10.0.0.1", Port: 22, Username: "root"})
	err := repo.Create(context.Background(), &Server{UUID: "same", Name: "b", Host: "10.0.0.2", Port: 22, Username: "root"})
	if err == nil {
		t.Fatal("expected the UNIQUE constraint on uuid to reject a duplicate")
	}
}

// Update must persist the connection fields; server import relies on it to apply
// an --on-conflict overwrite, and every other caller reads the row back first.
func TestSQLiteUpdatePersistsConnectionFields(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	s := mustCreate(t, repo, &Server{
		UUID: "u-1", Name: "web", Host: "10.0.0.1", Port: 22,
		Username: "root", AuthType: "password", AuthSecret: "pw", Provider: "Generic VPS",
	})

	s.Name = "web-renamed"
	s.Host = "172.16.0.9"
	s.Port = 2222
	s.Username = "deploy"
	s.Provider = "AWS"
	s.IsFavorite = true
	s.Tags = []Tag{{Name: "prod"}}
	if err := repo.Update(ctx, s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := repo.GetByID(ctx, s.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "web-renamed" || got.Host != "172.16.0.9" || got.Port != 2222 || got.Username != "deploy" {
		t.Fatalf("connection fields were not persisted: %+v", got)
	}
	if got.Provider != "AWS" || !got.IsFavorite {
		t.Fatalf("provider or favorite flag not persisted: %+v", got)
	}
	if len(got.Tags) != 1 || got.Tags[0].Name != "prod" {
		t.Fatalf("tags not persisted: %+v", got.Tags)
	}
	if got.AuthSecret != "pw" {
		t.Fatalf("got auth secret %q, want it preserved", got.AuthSecret)
	}
}

// A favorite toggle or rename must not disturb anything else on the row.
func TestSQLiteRepeatedUpdatesDoNotDrift(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	s := mustCreate(t, repo, &Server{
		UUID: "u-1", Name: "web", Host: "10.0.0.1", Port: 2201,
		Username: "alice", AuthType: "key", AuthSecret: "KEY", Provider: "Hetzner",
	})

	for i := 0; i < 5; i++ {
		cur, err := repo.GetByID(ctx, s.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		cur.IsFavorite = !cur.IsFavorite
		if err := repo.Update(ctx, cur); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	got, err := repo.GetByID(ctx, s.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Host != "10.0.0.1" || got.Port != 2201 || got.Username != "alice" {
		t.Fatalf("connection fields drifted after repeated toggles: %+v", got)
	}
	if got.AuthSecret != "KEY" || got.AuthType != "key" {
		t.Fatalf("credentials drifted after repeated toggles: %+v", got)
	}
}

// ListServerViews queries the software table for each row while the outer result
// set is still open, so the pool must allow more than one connection. Capping it
// at one deadlocks this call and hangs GET /servers.
func TestSQLiteListServerViewsDoesNotDeadlock(t *testing.T) {
	repo, db := newTestRepo(t)
	ctx := context.Background()

	for i := 1; i <= 5; i++ {
		s := mustCreate(t, repo, &Server{
			UUID: "u-" + string(rune('0'+i)), Name: "srv-" + string(rune('0'+i)),
			Host: "10.0.0." + string(rune('0'+i)), Port: 22, Username: "root",
		})
		if err := repo.ReplaceSoftware(ctx, s.ID, []Software{{Name: "nginx", Version: "1.25"}}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	if got := db.Stats().MaxOpenConnections; got == 1 {
		t.Fatal("connection pool is capped at 1, which deadlocks ListServerViews")
	}

	done := make(chan error, 1)
	go func() {
		views, err := repo.ListServerViews(ctx)
		if err == nil && len(views) != 5 {
			t.Errorf("got %d views, want 5", len(views))
		}
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("ListServerViews deadlocked")
	}
}

func TestSQLiteDeleteAndFlush(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	a := mustCreate(t, repo, &Server{UUID: "u-1", Name: "a", Host: "10.0.0.1", Port: 22, Username: "root"})
	mustCreate(t, repo, &Server{UUID: "u-2", Name: "b", Host: "10.0.0.2", Port: 22, Username: "root"})

	if err := repo.Delete(ctx, a.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := repo.Delete(ctx, a.ID); err == nil {
		t.Fatal("expected deleting a missing row to report ErrServerNotFound")
	}

	servers, _ := repo.List(ctx)
	if len(servers) != 1 {
		t.Fatalf("got %d servers after delete, want 1", len(servers))
	}

	if err := repo.Flush(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	servers, _ = repo.List(ctx)
	if len(servers) != 0 {
		t.Fatalf("got %d servers after flush, want 0", len(servers))
	}
}

func TestSQLiteTagRoundTrip(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	s := mustCreate(t, repo, &Server{UUID: "u-1", Name: "web", Host: "10.0.0.1", Port: 22, Username: "root"})

	if err := repo.AddTag(ctx, s.ID, "prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := repo.AddTag(ctx, s.ID, "prod"); err != nil {
		t.Fatalf("adding the same tag twice should be a no-op: %v", err)
	}
	if err := repo.AddTag(ctx, s.ID, "eu"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, _ := repo.GetByID(ctx, s.ID)
	if len(got.Tags) != 2 {
		t.Fatalf("got tags %+v, want exactly prod and eu", got.Tags)
	}

	if err := repo.RemoveTag(ctx, s.ID, "prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ = repo.GetByID(ctx, s.ID)
	if len(got.Tags) != 1 || got.Tags[0].Name != "eu" {
		t.Fatalf("got tags %+v, want only eu", got.Tags)
	}

	// Tags must survive an unrelated update.
	got.IsFavorite = true
	if err := repo.Update(ctx, got); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ = repo.GetByID(ctx, s.ID)
	if len(got.Tags) != 1 {
		t.Fatalf("tags lost on update: %+v", got.Tags)
	}
}
