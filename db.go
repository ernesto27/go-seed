package goseed

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type myDB interface {
	New() error
	Query(data map[string][]any) error
}

type Mysql struct {
	db *sql.DB
	Options
}

func (m *Mysql) New() error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", m.User, m.Password, m.Host, m.Port, m.Database))
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	m.db = db
	return nil
}

func (m *Mysql) Query(data map[string][]any) error {
	return query(data, m.Table, m.db)
}

type Postgres struct {
	db *sql.DB
	Options
}

func (p *Postgres) New() error {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", p.Host, p.Port, p.User, p.Password, p.Database))
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	p.db = db
	return nil
}

func (p *Postgres) Query(data map[string][]any) error {
	query := fmt.Sprintf("INSERT INTO %s (", p.Table)
	values := "VALUES ("
	params := []interface{}{}

	index := 1
	for key, value := range data {
		query += fmt.Sprintf("%s, ", key)
		values += "$" + fmt.Sprint(index) + ", "
		params = append(params, value[0])
		index++
	}

	query = query[:len(query)-2] + ") "
	values = values[:len(values)-2] + ");"

	_, err := p.db.Exec(query+values, params...)
	if err != nil {
		return err
	}

	return nil
}

type Sqlite struct {
	db *sql.DB
	Options
}

func (sqlite *Sqlite) New() error {
	db, err := sql.Open("sqlite3", sqlite.File)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	sqlite.db = db
	return nil
}

func (sqlite *Sqlite) Query(data map[string][]any) error {
	return query(data, sqlite.Table, sqlite.db)
}

func query(data map[string][]any, table string, db *sql.DB) error {
	query := fmt.Sprintf("INSERT INTO `%s` (", table)
	values := "VALUES ("
	params := []interface{}{}

	for key, value := range data {
		query += fmt.Sprintf("`%s`, ", key)
		values += "?, "
		params = append(params, value[0])
	}

	query = query[:len(query)-2] + ") "
	values = values[:len(values)-2] + ");"

	_, err := db.Exec(query+values, params...)
	if err != nil {
		return err
	}

	return nil
}
