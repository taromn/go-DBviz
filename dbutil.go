package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

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

func getTable(db *sql.DB, dbs *[]DB, i int, j int, query string, sname string) error {
	table_rows, err := db.Query(query, sname)
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
