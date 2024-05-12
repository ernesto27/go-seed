# Go seeder database

## Description
Seed a database in a local environment.

## Installation
```bash
$ go get github.com/ernesto27/go-seed
```

## Example

```go
package main

import goseed "github.com/ernesto27/go-seed"

func main() {

	// mysql example
	dataMysql := map[string][]any{
		"username":      {"Name"},
		"email":         {"Email"},
		"full_name":     {"Name"},
		"password_hash": {"Country"},
	}

	goseed.NewSeeder(goseed.Options{
		Engine:   "mysql",
		Host:     "localhost",
		Port:     "3306",
		Database: "yourdb",
		User:     "root",
		Password: "1111",
		Table:    "users",
	}).
		WithCount(10).
		Insert(dataMysql)

	// postgres example
	dataPostgres := map[string][]any{
		"name":    {"Company"},
		"address": {"Street"},
		"phone":   {"Phone"},
		"website": {"DomainName"},
		"email":   {"Email"},
	}

	goseed.NewSeeder(goseed.Options{
		Engine:   "postgres",
		Host:     "localhost",
		Port:     "5433",
		Database: "yourdb",
		User:     "postgres",
		Password: "1111",
		Table:    "providers",
	}).
		WithCount(20).
		Insert(dataPostgres)

	// sqlite example
	dataSqlite := map[string][]any{
		"name":     {"FirstName"},
		"email":    {"Email"},
		"password": {"13456"},
	}

	goseed.NewSeeder(goseed.Options{
		Engine: "sqlite",
		File:   "yourdb.db",
		Table:  "users",
	}).
		WithCount(10).
		Insert(dataSqlite)

}


```

You can use random data on the value of the columns of your tables using the name of the function of the gofakeit library.

https://github.com/brianvoe/gofakeit

https://github.com/brianvoe/gofakeit?tab=readme-ov-file#functions

Pass the name of the function, and the parameter required in a slice.

```go   
data := map[string][]any{
    "username":      {"Name"},
    "email":         {"Email"},
    "text":          {"Sentence", 3},
}
```
If you pass a method that does not exists in the gofakeit library, the value will be use in the insert query.

