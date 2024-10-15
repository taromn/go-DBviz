package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Getenv("DSN")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// select database names
	db_rows, err := db.Query("SELECT datname FROM pg_database WHERE datistemplate = false;")
	if err != nil {
		log.Fatal(err)
	}
	defer db_rows.Close()

	var dbs []string

	// show rows
	for db_rows.Next() {
		var datname string
		err = db_rows.Scan(&datname)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(datname)
		dbs = append(dbs, datname)
	}
	fmt.Println(dbs)

	// connect to each DB and get schemas
	for _, each_dbname := range dbs {
		fmt.Println(each_dbname)
		each_dsn := dsn + fmt.Sprintf("dbname=%s", each_dbname)
		each_db, err := sql.Open("postgres", each_dsn)
		if err != nil {
			log.Fatal(err)
		}
		defer each_db.Close()

		schema_rows, err := each_db.Query("select nspname from pg_namespace;")
		if err != nil {
			log.Fatal(err)
		}
		defer schema_rows.Close()

		// after writing all processes,
		// normalize the func creating connections
		// normalize the func of query results -> slices

	}
}
