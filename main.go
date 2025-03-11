package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Schema struct {
	Sname string
}

type Db struct {
	DBname     string
	SchemaList []Schema
}

func OpenRun(d_str string, q_str string) (*sql.Rows, error) {
	db, err := sql.Open("postgres", d_str)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(q_str)
	if err != nil {
		return nil, err
	}

	return rows, nil // do not close rows here
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

	}
	fmt.Println("final dbs are", dbs)

}
