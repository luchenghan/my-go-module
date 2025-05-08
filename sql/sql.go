package sql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Config struct {
	dataSourceName  string
	maxOpenConns    int
	maxIdleConns    int
	connMaxLifetime time.Duration
}

type DB struct {
	*sql.DB
}

func NewDB(config Config) (*DB, error) {
	db, err := sql.Open("mysql", config.dataSourceName)
	if err != nil {
		return nil, err
	}

	// Set the maximum number of open connections
	db.SetMaxOpenConns(config.maxOpenConns)

	// Set the maximum number of idle connections
	db.SetMaxIdleConns(config.maxIdleConns)

	// Set the maximum lifetime of a connection
	db.SetConnMaxLifetime(config.connMaxLifetime)

	// Verify the connection
	if err = db.Ping(); err != nil {
		return nil, err
	}

	d := new(DB)
	d.DB = db

	return d, nil
}

// Implement Batch Insert method
func (d *DB) BatchInsert(table string, columns []string, values [][]interface{}, chunkSize int) error {
	if len(values) == 0 {
		return nil
	}

	tx, err := d.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	placeholder := "(" + strings.Repeat("?,", len(columns)-1) + "?)"
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", table, strings.Join(columns, ","), placeholder)

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i := 0; i < len(values); i += chunkSize {
		end := i + chunkSize
		if end > len(values) {
			end = len(values)
		}

		for _, value := range values[i:end] {
			if _, err := stmt.Exec(value...); err != nil {
				return err
			}
		}
	}

	return nil
}

// Implement Pagination method
func (d *DB) Paginate(table string, columns []string, page int, pageSize int) (*sql.Rows, error) {
	if page < 1 || pageSize < 1 {
		return nil, fmt.Errorf("invalid page or page size")
	}

	offset := (page - 1) * pageSize
	query := fmt.Sprintf("SELECT %s FROM %s LIMIT %d OFFSET %d", strings.Join(columns, ","), table, pageSize, offset)

	rows, err := d.Query(query)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
