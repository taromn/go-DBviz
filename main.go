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

	var dbs []Db
	var datname string

	// show rows
	for db_rows.Next() {
		err = db_rows.Scan(&datname) // *database/sql.Rows
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(datname)
		dbs = append(dbs, Db{DBname: datname})
	}
	fmt.Println(dbs)

	// connect to each DB and get schemas
	for i := range dbs {
		each_dbname := dbs[i].DBname
		fmt.Println(each_dbname)
		each_dsn := dsn + fmt.Sprintf(" dbname=%s", each_dbname)
		db_conn, err := sql.Open("postgres", each_dsn)
		if err != nil {
			log.Print(err)
			continue
		}
		defer db_conn.Close()

		schema_rows, err := db_conn.Query("select nspname from pg_namespace;")
		if err != nil {
			log.Print(err)
			continue
		}
		defer schema_rows.Close()

		var schemas []Schema
		var schema_name string

		// show schema rows
		for schema_rows.Next() {
			err = schema_rows.Scan(&schema_name)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(schema_name)
			s := Schema{Sname: schema_name}
			schemas = append(schemas, s)
		}
		fmt.Println(schemas)

		dbs[i].SchemaList = schemas

		fmt.Println("final dbs are", dbs)

		// after writing all processes,
		// normalize the func creating connections
		// normalize the func of query results -> slices
	}

}
