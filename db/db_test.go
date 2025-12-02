package db_test

import (
	"path/filepath"
	"testing"

	"github.com/HenrySchwerdt/tt/db"
)

func newTestDB(t *testing.T) *db.Database {
	t.Helper()

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	d, err := db.Init(dbPath)
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}

	return d
}

func TestCreateProject_Simple(t *testing.T) {
	d := newTestDB(t)

	p, err := d.CreateProject("foo")
	if err != nil {
		t.Fatalf("CreateProject failed: %v", err)
	}

	if p.Name != "foo" {
		t.Errorf("expected name foo, got %s", p.Name)
	}
	if p.Path != "foo" {
		t.Errorf("expected path foo, got %s", p.Path)
	}

	// Fetch again
	got, err := d.GetProjectByPath("foo")
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Fatalf("expected project foo to exist")
	}
}

func TestCreateProject_Nested(t *testing.T) {
	d := newTestDB(t)

	p, err := d.CreateProject("foo/bar/baz")
	if err != nil {
		t.Fatalf("CreateProject failed: %v", err)
	}

	if p.Path != "foo/bar/baz" {
		t.Fatalf("unexpected path %s", p.Path)
	}

	// Check intermediate paths
	check := []string{"foo", "foo/bar", "foo/bar/baz"}
	for _, path := range check {
		exists, err := d.GetProjectByPath(path)
		if err != nil {
			t.Fatal(err)
		}
		if exists == nil {
			t.Fatalf("expected %s to exist", path)
		}
	}
}

func TestCreateProject_Duplicate(t *testing.T) {
	d := newTestDB(t)

	_, err := d.CreateProject("foo/bar")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Creating again should return the same final project (not error)
	p2, err := d.CreateProject("foo/bar")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p2.Path != "foo/bar" {
		t.Fatalf("expected foo/bar, got %s", p2.Path)
	}

	// Only 1 project with path "foo/bar"
	p, err := d.GetProjectByPath("foo/bar")
	if err != nil {
		t.Fatal(err)
	}
	if p == nil {
		t.Fatalf("expected foo/bar to exist")
	}
}

func TestGetProjectByPath_NotFound(t *testing.T) {
	d := newTestDB(t)

	p, err := d.GetProjectByPath("does/not/exist")
	if err != nil {
		t.Fatal(err)
	}
	if p != nil {
		t.Fatalf("expected nil project for nonexistent path")
	}
}
