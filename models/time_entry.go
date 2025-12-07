package models

import (
	"time"

	"github.com/uptrace/bun"
)

type TimeEntry struct {
	bun.BaseModel `bun:"table:time_entries,alias:te"`
	ID            int64      `bun:",pk,autoincrement"`
	ProjectID     int64      `bun:"project_id"`
	Start         time.Time  `bun:"start,notnull"`
	End           *time.Time `bun:"end"`
	Delta         *int64     `bun:"delta"`
	Message       *string    `bun:"message"`
}
