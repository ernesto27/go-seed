package goseed

import (
	"reflect"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/brianvoe/gofakeit/v7/source"
)

type seeder struct {
	tableName string
	count     int
	engine    myDB
}

type Options struct {
	Engine   string
	Host     string
	Port     string
	Database string
	User     string
	Password string
	File     string
	Table    string
}

func NewSeeder(options Options) *seeder {
	var engine myDB

	switch options.Engine {
	case "mysql":
		engine = &Mysql{Options: options}
		err := engine.New()
		if err != nil {
			panic(err)
		}
	case "postgres":
		engine = &Postgres{Options: options}
		err := engine.New()
		if err != nil {
			panic(err)
		}
	case "sqlite":
		engine = &Sqlite{Options: options}
		err := engine.New()
		if err != nil {
			panic(err)
		}
	case "cassandra":
		engine = &Cassandra{Options: options}
		err := engine.New()
		if err != nil {
			panic(err)
		}
	case "oracle":
		engine = &Oracle{Options: options}
		err := engine.New()
		if err != nil {
			panic(err)
		}
	}

	return &seeder{
		tableName: options.Table,
		engine:    engine,
	}
}

func (s *seeder) doInsert(data map[string][]any) {
	dataInsert := s.getData(data)

	err := s.engine.Query(dataInsert)
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
