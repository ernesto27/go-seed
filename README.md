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

import dbseed "github.com/ernesto27/go-seed"

func main() {

	data := map[string][]any{
		"username":      {"Name"},
		"email":         {"Email"},
	}

	goseed.NewSeeder(goseed.Options{
		Engine:   "mysql",
		Host:     "localhost",
		Port:     "3388",
		Database: "mydb",
		User:     "root",
		Password: "1111",
		Table:    "users",
	}).
		WithCount(10).
		Insert(data)
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

