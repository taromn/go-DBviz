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
	rows, err := db.Query("SELECT datname FROM pg_database WHERE datistemplate = false;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// show rows
	for rows.Next() {
		var datname string
		err = rows.Scan(&datname)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(datname)
	}
}
