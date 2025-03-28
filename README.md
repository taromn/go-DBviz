List database resources in Postgres

## How to use

0/ install Golang

https://go.dev/doc/install

1/ git clone
```
git clone https://github.com/taromn/go-DBviz.git
```

2/ set DSN
```
export DSN="host=example-instance.rds.amazonaws.com port=5432 user=testuser password=xxx sslmode=disable" 
```

3/ navigate to go-DBviz, run it
```
go run .
```

## Requirements
User specified in DSN should be able to SELECT the following views
- pg_database
- pg_namespace
- pg_catalog.pg_tables


## Disclaimer

The developer of this program assumes no responsibility for any problems that may occur as a result of executing this code.

Do not use in a production environment.
