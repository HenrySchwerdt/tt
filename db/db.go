package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

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
	ErrProjectExists        = errors.New("project already exists")
	ErrProjectDoesNotExist  = errors.New("project does not exist")
	ErrTimeEntryNotClosed   = errors.New("close current time entry before starting a new one")
	ErrNoOpenTimeEntryFound = errors.New("no open time entry found")
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
	// Init Tables if not exist
	if _, err := client.NewCreateTable().IfNotExists().Model((*models.Project)(nil)).Exec(db.ctx); err != nil {
		return nil, err
	}
	if _, err := client.NewCreateTable().IfNotExists().Model((*models.TimeEntry)(nil)).Exec(db.ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func (d *Database) GetProjectByPath(path string) (*models.Project, error) {
	var p models.Project
	err := d.client.NewSelect().Model(&p).Where("path = ?", path).Relation("Children").Limit(1).Scan(d.ctx)
	fmt.Println(err)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}

func (d *Database) GetProjectByPathRecursive(path string) (*models.Project, error) {
	var p models.Project

	err := d.client.NewSelect().
		Model(&p).
		Where("path = ?", path).
		Relation("Children").
		Limit(1).
		Scan(d.ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	localTime, err := d.SumTimeForProject(p.ID)
	if err != nil {
		return nil, err
	}
	p.TimeSpend = localTime
	var totalChildrenTime int64 = 0

	for _, child := range p.Children {
		fullChild, err := d.GetProjectByPathRecursive(child.Path)
		if err != nil {
			return nil, err
		}

		*child = *fullChild
		totalChildrenTime += child.TimeSpend
	}
	p.TimeSpend += totalChildrenTime

	return &p, nil
}

func (d *Database) SumTimeForProject(projectID int64) (int64, error) {
	var sum *int64
	err := d.client.NewSelect().
		Table("time_entries").
		ColumnExpr("SUM(delta)").
		Where("project_id = ?", projectID).
		Where("end IS NOT NULL").
		Scan(d.ctx, &sum)
	if err != nil {
		return 0, err
	}

	if sum == nil {
		return 0, nil
	}
	return *sum, nil
}

func (d *Database) RemoveProject(path string) error {
	root, err := d.GetProjectByPathRecursive(path)
	if err != nil {
		return err
	}
	if root == nil {
		return ErrProjectDoesNotExist
	}

	var collectIDs func(p *models.Project, ids *[]int64)
	collectIDs = func(p *models.Project, ids *[]int64) {
		*ids = append(*ids, p.ID)
		for _, child := range p.Children {
			collectIDs(child, ids)
		}
	}

	var ids []int64
	collectIDs(root, &ids)

	if len(ids) == 0 {
		return nil
	}

	tx, err := d.client.BeginTx(d.ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.NewDelete().
		Model((*models.TimeEntry)(nil)).
		Where("project_id IN (?)", bun.In(ids)).
		Exec(d.ctx)
	if err != nil {
		return err
	}

	_, err = tx.NewDelete().
		Model((*models.Project)(nil)).
		Where("id IN (?)", bun.In(ids)).
		Exec(d.ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
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

func (d *Database) assertNoOpenTimeEntries() error {
	var exists bool

	// SELECT EXISTS(SELECT 1 FROM time_entries WHERE end IS NULL LIMIT 1)
	err := d.client.NewSelect().
		Model((*models.TimeEntry)(nil)).
		Where("end IS NULL").
		Limit(1).
		Scan(d.ctx, &exists)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}

	if exists {
		return ErrTimeEntryNotClosed
	}

	return nil
}

func (d *Database) StartTimeEntry(path string) error {

	if err := d.assertNoOpenTimeEntries(); err != nil {
		return err
	}

	project, err := d.GetProjectByPath(path)

	if err != nil || project == nil {
		return ErrProjectDoesNotExist
	}
	startTime := time.Now()
	timeEntry := &models.TimeEntry{
		ProjectID: project.ID,
		Start:     startTime,
		End:       nil,
		Delta:     nil,
		Message:   nil,
	}

	_, err = d.client.NewInsert().
		Model(timeEntry).
		Exec(d.ctx)

	if err != nil {
		return err
	}
	return nil
}

func (d *Database) EndTimeEntry(message string) error {
	var entry models.TimeEntry
	err := d.client.NewSelect().
		Model(&entry).
		Where("end IS NULL").
		Limit(1).
		Scan(d.ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoOpenTimeEntryFound
		}
		return err
	}

	endTime := time.Now()
	delta := int64(endTime.Sub(entry.Start).Seconds())

	entry.End = &endTime
	entry.Delta = &delta

	if message != "" {
		entry.Message = &message
	}

	_, err = d.client.NewUpdate().
		Model(&entry).
		WherePK().
		Exec(d.ctx)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) loadProject(path string) (*models.Project, error) {
	var p models.Project
	err := d.client.NewSelect().
		Model(&p).
		Where("path = ?", path).
		Relation("Children").
		Limit(1).
		Scan(d.ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}

func (d *Database) buildProjectTree(p *models.Project) error {
	// Sum local time
	localTime, err := d.SumTimeForProject(p.ID)
	if err != nil {
		return err
	}

	p.TimeSpend = localTime
	var totalChildren int64 = 0

	// Recursively process children
	for i, child := range p.Children {
		fullChild, err := d.loadProject(child.Path)
		if err != nil {
			return err
		}
		if fullChild == nil {
			return fmt.Errorf("child %s not found but listed in DB", child.Path)
		}

		// Recursively build the child
		if err := d.buildProjectTree(fullChild); err != nil {
			return err
		}

		p.Children[i] = fullChild
		totalChildren += fullChild.TimeSpend
	}

	p.TimeSpend += totalChildren
	return nil
}

func (d *Database) GetProjectByPathRecursive2(path string) (*models.Project, error) {
	project, err := d.loadProject(path)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, nil
	}

	if err := d.buildProjectTree(project); err != nil {
		return nil, err
	}

	return project, nil
}

func (d *Database) GetAllProjectsRecursive() ([]*models.Project, error) {
	// Load all root-level projects
	var roots []*models.Project
	err := d.client.
		NewSelect().
		Model(&roots).
		Where("parent_id IS NULL").
		Relation("Children").
		Order("name ASC").
		Scan(d.ctx)
	if err != nil {
		return nil, err
	}

	// Build full recursive trees
	for _, root := range roots {
		if err := d.buildProjectTree(root); err != nil {
			return nil, err
		}
	}
	return roots, nil
}

func (d *Database) GetProjectById(id int64) (*models.Project, error) {
	var p models.Project
	err := d.client.NewSelect().
		Model(&p).
		Where("id = ?", id).
		Limit(1).
		Scan(d.ctx)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (d *Database) GetTimelineForProject(path string) ([]models.FeedItem, error) {
	// Load project with full tree
	root, err := d.GetProjectByPathRecursive2(path)
	if err != nil {
		return nil, err
	}
	if root == nil {
		return nil, ErrProjectDoesNotExist
	}

	// Collect all project IDs
	var ids []int64
	var collect func(p *models.Project)
	collect = func(p *models.Project) {
		ids = append(ids, p.ID)
		for _, c := range p.Children {
			collect(c)
		}
	}
	collect(root)

	// Query time entries for all projects
	var entries []models.TimeEntry
	err = d.client.NewSelect().
		Model(&entries).
		Where("project_id IN (?)", bun.In(ids)).
		Where("end IS NOT NULL").
		Order("start ASC").
		Scan(d.ctx)
	if err != nil {
		return nil, err
	}

	// Convert to feed items
	feed := make([]models.FeedItem, 0, len(entries))

	for _, e := range entries {
		var dur time.Duration
		if e.Delta != nil {
			dur = time.Duration(*e.Delta) * time.Second
		} else if e.End != nil {
			dur = e.End.Sub(e.Start)
		}

		// get path for each ID
		p, err := d.GetProjectById(e.ProjectID)
		if err != nil {
			return nil, err
		}

		feed = append(feed, models.FeedItem{
			ProjectPath: p.Path,
			Start:       e.Start,
			End:         *e.End,
			Duration:    dur,
			Message:     "",
		})

		if e.Message != nil {
			feed[len(feed)-1].Message = *e.Message
		}
	}

	return feed, nil
}
