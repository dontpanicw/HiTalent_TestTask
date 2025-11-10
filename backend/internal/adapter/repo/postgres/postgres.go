package postgres

import (
	migrations "HiTalent_TestTask/backend/pkg/migration/postgres"
	"database/sql"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormDB(pgConnStr string) (*gorm.DB, error) {
	// Сначала применяем миграции через goose
	db, err := sql.Open("pgx", pgConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := migrations.Migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	if err = db.Close(); err != nil {
		return nil, fmt.Errorf("failed to close database: %w", err)
	}

	// Затем создаем GORM подключение
	gormDB, err := gorm.Open(postgres.Open(pgConnStr), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return gormDB, nil
}
