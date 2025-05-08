package db

import (
	"database/sql"
	"mymodule/gin/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitMySQL(config config.DatabaseConfig) error {
	var err error
	db, err = sql.Open("mysql", config.DataSourceName)
	if err != nil {
		return err
	}

	// Set the maximum number of open connections
	db.SetMaxOpenConns(config.MaxOpenConns)

	// Set the maximum number of idle connections
	db.SetMaxIdleConns(config.MaxIdleConns)

	// Set the maximum lifetime of a connection
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	return db.Ping()
}

func GetDB() *sql.DB {
	return db
}

func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
