package si

import (
	"database/sql"
	"time"
)

func OpenSqlDB(driver, dsn string, maxIdleConns, maxOpenConns int, connMaxLifetime time.Duration) (*sql.DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime)

	return db, nil
}
