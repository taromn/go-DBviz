package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Table struct {
	Tname string `json:"table_name"`
}

type Schema struct {
	Sname     string `json:"schema_name"`
	TableList []Table
}

type DB struct {
	DBname     string
	SchemaList []Schema
}

func getDB(dbs *[]DB, dsn string, query string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect DB: %w", err)
	}
	defer db.Close()

	db_rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to run query: %s %w", query, err)
	}

	var datname string

	// show rows
	for db_rows.Next() {
		err := db_rows.Scan(&datname) // *database/sql.Rows
		if err != nil {
			return fmt.Errorf("failed to scan DB rows: %w", err)
		}
		*dbs = append(*dbs, DB{DBname: datname})
	}
	db_rows.Close()

	return nil
}

func getSchema(dbs *[]DB, i int, dsn string, query string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect DB: %w", err)
	}

	schema_rows, err := db.Query(query)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run query: %s %w", query, err)
	}

	var schemas []Schema
	var schema_name string

	// show schema rows
	for schema_rows.Next() {
		err := schema_rows.Scan(&schema_name)
		if err != nil {
			log.Println(err)
			continue
		}
		s := Schema{Sname: schema_name}
		schemas = append(schemas, s)
	}
	schema_rows.Close()

	(*dbs)[i].SchemaList = schemas

	return db, nil
}

func getTable(db *sql.DB, dbs *[]DB, i int, j int, query string) error {
	table_rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to run query: %s %w", query, err)
	}

	var tables []Table
	var table_name string

	for table_rows.Next() {
		err := table_rows.Scan(&table_name)
		if err != nil {
			log.Println(err)
			continue
		}
		t := Table{Tname: table_name}
		tables = append(tables, t)
	}
	table_rows.Close()

	(*dbs)[i].SchemaList[j].TableList = tables

	return nil
}

func PrintS(st any) {
	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(data))
}

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

			table_query := fmt.Sprintf("SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname = '%s';", each_sname)

			err := getTable(db, &dbs, i, j, table_query)
			if err != nil {
				log.Println(err)
				continue
			}
		}
		db.Close()

	}
	PrintS(dbs)

}
