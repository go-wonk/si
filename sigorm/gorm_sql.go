package sigorm

import (
	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenPostgres(db *sql.DB) (*gorm.DB, error) {
	d := NewPostgresDialector(postgres.Config{Conn: db})
	return Open(d, &gorm.Config{})
}

func OpenPostgresWithConfig(db *sql.DB, config *gorm.Config) (*gorm.DB, error) {
	d := NewPostgresDialector(postgres.Config{Conn: db})
	return Open(d, config)
}

func Open(gormDialector gorm.Dialector, config *gorm.Config) (*gorm.DB, error) {
	gormDB, err := gorm.Open(
		gormDialector,
		config)

	if err != nil {
		return nil, err
	}

	return gormDB, nil
}

func NewPostgresDialector(config postgres.Config) gorm.Dialector {
	return postgres.New(config)
}
