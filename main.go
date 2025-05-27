package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	// example: export DSN="host=example-instance.rds.amazonaws.com port=5432 user=testuser password=xxx sslmode=disable"
	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Println("DSN is not set")
		os.Exit(1)
	}

	var dbs []DB

	err := getDB(&dbs, dsn, "SELECT datname FROM pg_database WHERE datistemplate = false;")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// connect to each DB and get schemas
	for i := range dbs {
		each_dsn := dsn + fmt.Sprintf(" dbname=%s", dbs[i].DBname)

		db, err := getSchema(&dbs, i, each_dsn, "select nspname from pg_namespace;")
		if err != nil {
			log.Println(err)
			continue
		}
		// tables for each schema
		for j := range dbs[i].SchemaList {
			each_sname := dbs[i].SchemaList[j].Sname
			table_query := "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname = $1;"

			err := getTable(db, &dbs, i, j, table_query, each_sname)
			if err != nil {
				log.Println(err)
				continue
			}
		}
		db.Close()

	}
	PrintS(dbs)

}
