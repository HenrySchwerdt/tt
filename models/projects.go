package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Project struct {
	bun.BaseModel `bun:"table:projects,alias:p"`

	ID        int64     `bun:",pk,autoincrement"`
	ParentID  int64     `bun:"parent_id,nullzero"`
	Path      string    `bun:"path,notnull,unique"`
	Name      string    `bun:"name,notnull"`
	CreatedAt time.Time `bun:"created_at,notnull,default:CURRENT_TIMESTAMP"`
	Finished  bool      `bun:"finished,notnull,default:false"`
	TimeSpend int64     `bun:"-"`

	Children []*Project `bun:"rel:has-many,join:id=parent_id"`
}
