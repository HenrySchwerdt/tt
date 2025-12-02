package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/HenrySchwerdt/tt/models"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

type Database struct {
	client *bun.DB
	ctx    context.Context
}

var (
	ErrProjectExists = errors.New("project already exists")
)

func Init(path string) (*Database, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		file.Close()
	}
	sqldb, err := sql.Open(sqliteshim.ShimName, fmt.Sprintf("file:%s?cache=shared&mode=rwc", path))
	if err != nil {
		return nil, err
	}
	client := bun.NewDB(sqldb, sqlitedialect.New())

	db := &Database{
		client: client,
		ctx:    context.Background(),
	}

	_, err = client.NewCreateTable().IfNotExists().Model((*models.Project)(nil)).Exec(db.ctx)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func (d *Database) GetProjectByPath(path string) (*models.Project, error) {
	var p models.Project
	err := d.client.NewSelect().Model(&p).Where("path = ?", path).Limit(1).Scan(d.ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}

func (d *Database) createProject(parent *models.Project, name string) (*models.Project, error) {
	var path string
	if parent == nil {
		path = name
	} else {
		path = parent.Path + "/" + name
	}

	// Already exists?
	existing, err := d.GetProjectByPath(path)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrProjectExists
	}

	p := &models.Project{
		Name:     name,
		Path:     path,
		ParentID: 0,
	}

	if parent != nil {
		p.ParentID = parent.ID
	}

	_, err = d.client.NewInsert().
		Model(p).
		Exec(d.ctx)

	if err != nil {
		return nil, err
	}

	return p, nil
}

func (d *Database) CreateProject(path string) (*models.Project, error) {
	parts := strings.Split(path, "/")

	var parent *models.Project

	for i := range parts {
		sub := strings.Join(parts[:i+1], "/")
		name := parts[i]

		existing, err := d.GetProjectByPath(sub)
		if err != nil {
			return nil, err
		}

		if existing != nil {
			parent = existing
			continue
		}

		created, err := d.createProject(parent, name)
		if err != nil {
			if errors.Is(err, ErrProjectExists) {
				parent, _ = d.GetProjectByPath(sub)
				continue
			}
			return nil, err
		}

		parent = created
	}

	return parent, nil
}
