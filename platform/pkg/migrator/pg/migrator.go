package pg

import (
	"context"
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

type Migrator struct {
	db  *sql.DB
	fs  embed.FS
	dir string
}

func New(db *sql.DB, fs embed.FS, dir string) *Migrator {
	return &Migrator{
		db:  db,
		fs:  fs,
		dir: dir,
	}
}

func (m *Migrator) Up(ctx context.Context) error {
	goose.SetBaseFS(m.fs)
	return goose.UpContext(ctx, m.db, m.dir)
}

func (m *Migrator) Down(ctx context.Context) error {
	goose.SetBaseFS(m.fs)
	return goose.DownContext(ctx, m.db, m.dir)
}

func (m *Migrator) Version(ctx context.Context) (int64, error) {
	return goose.GetDBVersionContext(ctx, m.db)
}