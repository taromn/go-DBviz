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

type Db struct {
	DBname     string
	SchemaList []Schema
}

func OpenRun(d_str string, q_str string) (*sql.Rows, error) {
	db, err := sql.Open("postgres", d_str)
	if err != nil {
		return nil, fmt.Errorf("failed to connect DB: %w", err)
	}
	defer db.Close()

	rows, err := db.Query(q_str)
	if err != nil {
		return nil, fmt.Errorf("failed to run query: %s %w", q_str, err)
	}

	return rows, nil // do not close rows here
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
	dsn := os.Getenv("DSN")

	db_rows, err := OpenRun(dsn, "SELECT datname FROM pg_database WHERE datistemplate = false;")

	if err != nil {
		log.Println("failed to get DB info:", err)
		os.Exit(1)
	}

	var dbs []Db
	var datname string

	// show rows
	for db_rows.Next() {
		err := db_rows.Scan(&datname) // *database/sql.Rows
		if err != nil {
			log.Println(err)
		}
		fmt.Println(datname)
		dbs = append(dbs, Db{DBname: datname})
	}
	db_rows.Close()

	fmt.Println(dbs)

	// connect to each DB and get schemas
	for i := range dbs {
		each_dbname := dbs[i].DBname
		fmt.Println(each_dbname)
		each_dsn := dsn + fmt.Sprintf(" dbname=%s", each_dbname)

		schema_rows, err := OpenRun(each_dsn, "select nspname from pg_namespace;")

		if err != nil {
			log.Println(err)
			continue // prevent termination
		}

		var schemas []Schema
		var schema_name string

		// show schema rows
		for schema_rows.Next() {
			err := schema_rows.Scan(&schema_name)
			if err != nil {
				log.Println(err)
			}
			fmt.Println(schema_name)
			s := Schema{Sname: schema_name}
			schemas = append(schemas, s)
		}
		schema_rows.Close()

		fmt.Println(schemas)

		dbs[i].SchemaList = schemas

		for j := range dbs[i].SchemaList {
			each_sname := dbs[i].SchemaList[j].Sname
			fmt.Printf("sname is: %s", each_sname)

			table_query := fmt.Sprintf("SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname = '%s';", each_sname)

			table_rows, err := OpenRun(each_dsn, table_query)

			if err != nil {
				log.Println(err)
				continue // prevent termination
			}

			var tables []Table
			var table_name string

			for table_rows.Next() {
				err := table_rows.Scan(&table_name)
				if err != nil {
					log.Println(err)
				}
				fmt.Println(table_name)
				t := Table{Tname: table_name}
				tables = append(tables, t)
			}
			table_rows.Close()

			fmt.Println(tables)

			dbs[i].SchemaList[j].TableList = tables

		}

	}
	PrintS(dbs)

}
