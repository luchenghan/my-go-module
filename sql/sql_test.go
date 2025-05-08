package sql

import (
	"fmt"
	"log"
	"testing"
	"time"
)

var config = Config{
	dataSourceName:  "root:password@tcp(localhost:3306)/my_schema",
	maxOpenConns:    25,
	maxIdleConns:    25,
	connMaxLifetime: 5 * time.Minute,
}

var db *DB

func Test_NewDB(t *testing.T) {
	var err error
	db, err = NewDB(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
	t.Logf("Successfully connected to the database")
}

func Test_BatchInsert(t *testing.T) {
	var err error
	db, err = NewDB(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	start := time.Now()
	defer func() {
		t.Logf("Batch insert took %s", time.Since(start))
	}()

	table := "users"
	columns := []string{"username", "email"}
	var values [][]interface{}
	for i := 0; i < 2000; i++ {
		values = append(values, []interface{}{fmt.Sprintf("value%d", i), fmt.Sprintf("value%d", i+1)})
	}

	if err := db.BatchInsert(table, columns, values, 1000); err != nil {
		t.Fatalf("Failed to batch insert: %v", err)
	}
	t.Logf("Successfully batch inserted data into the database")
}

func Test_Paginate(t *testing.T) {
	var err error
	db, err = NewDB(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	start := time.Now()
	defer func() {
		t.Logf("Pagination took %s", time.Since(start))
	}()

	table := "users"
	columns := []string{"username", "email"}
	page := 1
	pageSize := 20

	rows, err := db.Paginate(table, columns, page, pageSize)
	if err != nil {
		t.Fatalf("Failed to paginate: %v", err)
	}

	for rows.Next() {
		var column1 string
		var column2 string
		err = rows.Scan(&column1, &column2)
		if err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}
		t.Logf("Row: %s, %s", column1, column2)
	}

	rows.Close()
}
