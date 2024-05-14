package goseed

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocql/gocql"
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

func (mysql *Mysql) New() error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", mysql.User, mysql.Password, mysql.Host, mysql.Port, mysql.Database))
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	mysql.db = db
	return nil
}

func (mysql *Mysql) Query(data map[string][]any) error {
	defer mysql.db.Close()

	query, params := getQueryParams(data, mysql.Table, "?", "`")
	_, err := mysql.db.Exec(query, params...)
	if err != nil {
		return err
	}

	return nil
}

type Postgres struct {
	db *sql.DB
	Options
}

func (postgres *Postgres) New() error {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", postgres.Host, postgres.Port, postgres.User, postgres.Password, postgres.Database))
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	postgres.db = db
	return nil
}

func (postgres *Postgres) Query(data map[string][]any) error {
	defer postgres.db.Close()

	query, params := getQueryParams(data, postgres.Table, "$", "")
	_, err := postgres.db.Exec(query, params...)
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
	defer sqlite.db.Close()
	query, params := getQueryParams(data, sqlite.Table, "?", "")
	_, err := sqlite.db.Exec(query, params...)
	if err != nil {
		return err
	}
	return nil
}

func getQueryParams(data map[string][]any, table string, placeholder string, delimeter string) (string, []interface{}) {
	query := fmt.Sprintf("INSERT INTO "+delimeter+"%s"+delimeter+" (", table)
	values := "VALUES ("
	params := []interface{}{}

	for key, value := range data {
		query += fmt.Sprintf(delimeter+"%s"+delimeter+", ", key)
		values += placeholder + ", "
		params = append(params, value[0])
	}

	query = query[:len(query)-2] + ") "
	values = values[:len(values)-2] + ");"

	fmt.Println(query + values)

	return query + values, params
}

func queryCassandra(data map[string][]any, table string) (string, []interface{}) {
	query := fmt.Sprintf("INSERT INTO %s (", table)
	values := "VALUES ("
	params := []interface{}{}

	for key, value := range data {
		query += fmt.Sprintf("%s, ", key)
		values += "?, "
		params = append(params, value[0])
	}

	query = query[:len(query)-2] + ") "
	values = values[:len(values)-2] + ");"

	return query + values, params
}

type Cassandra struct {
	session *gocql.Session
	Options
}

func (cassandra *Cassandra) New() error {
	cluster := gocql.NewCluster(cassandra.Host)
	cluster.Keyspace = cassandra.Database
	cluster.Consistency = gocql.One

	session, err := cluster.CreateSession()
	if err != nil {
		return err
	}

	cassandra.session = session
	return nil
}

func (cassandra *Cassandra) Query(data map[string][]any) error {
	query, params := getQueryParams(data, cassandra.Table, "?", "")

	for key, param := range params {
		if param == "UUID" {
			params[key] = gocql.TimeUUID()
		}
	}

	err := cassandra.session.Query(query, params...).Exec()
	if err != nil {
		return err
	}

	return nil
}
