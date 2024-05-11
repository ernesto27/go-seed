package goseed

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/brianvoe/gofakeit/v7/source"
	_ "github.com/go-sql-driver/mysql"
)

type seeder struct {
	tableName string
	db        *sql.DB
	count     int
}

type Options struct {
	Engine   string
	Host     string
	Port     string
	Database string
	User     string
	Password string
	Table    string
}

func NewSeeder(options Options) *seeder {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", options.User, options.Password, options.Host, options.Port, options.Database))
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return &seeder{
		tableName: options.Table,
		db:        db,
	}
}

func (s *seeder) doInsert(data map[string][]any) {
	dataInsert := s.getData(data)

	query := fmt.Sprintf("INSERT INTO `%s` (", s.tableName)
	values := "VALUES ("
	params := []interface{}{}

	for key, value := range dataInsert {
		query += fmt.Sprintf("`%s`, ", key)
		values += "?, "
		params = append(params, value[0])
	}

	query = query[:len(query)-2] + ") "
	values = values[:len(values)-2] + ");"

	_, err := s.db.Exec(query+values, params...)
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}
}

func (*seeder) getData(data map[string][]any) map[string][]any {
	faker := gofakeit.NewFaker(source.NewCrypto(), true)

	dataInsert := map[string][]any{}

	for key, value := range data {
		method := value[0]
		meth := reflect.ValueOf(faker).MethodByName(method.(string))

		if meth.Kind() == reflect.Invalid {
			dataInsert[key] = []any{method}
			continue
		}

		var resp []reflect.Value
		if len(value) > 1 {
			resp = meth.Call([]reflect.Value{reflect.ValueOf(value[1])})
		} else {
			resp = meth.Call(nil)
		}
		dataInsert[key] = []any{resp[0].String()}
	}
	return dataInsert
}

func (s *seeder) Insert(data map[string][]any) {
	if s.count == 0 {
		s.count = 1
	}

	for i := 0; i < s.count; i++ {
		s.doInsert(data)
	}
}

func (s *seeder) WithCount(count int) *seeder {
	s.count = count
	return s
}
