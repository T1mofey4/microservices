package pg

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/pressly/goose/v3"

	def "github.com/T1mofey4/microservices/platform/pkg/migrator"
)

var (
	_    def.Migrator = (*migrator)(nil)
	once sync.Once
	err  error
)

type migrator struct {
	db            *sql.DB
	migrationsDir string
}

func NewMigrator(db *sql.DB, migrationsDir string) *migrator {
	return &migrator{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

func (m *migrator) Up(ctx context.Context) error {
	dialectErr := ensureDialect("postgres")
	if dialectErr != nil {
		return fmt.Errorf("migrator failed set dialect: %w", dialectErr)
	}

	upErr := goose.UpContext(ctx, m.db, m.migrationsDir)
	if upErr != nil {
		return fmt.Errorf("apply migrations: %w", upErr)
	}

	return nil
}

func ensureDialect(dialect string) error {
	once.Do(func() {
		err = goose.SetDialect(dialect)
	})
	return err
}
