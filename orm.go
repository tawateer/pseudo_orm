package orm

import (
	"reflect"

	"github.com/jmoiron/sqlx"
)


type TxGetter func() (*sqlx.Tx, error)
type Mapper func(string) string
type TagFilter func(reflect.StructField) bool

var (
	txGetter TxGetter
	mapper Mapper
	tagFilter TagFilter
)

func Init(t TxGetter, m Mapper, f TagFilter) {
	txGetter = t
	mapper = m
	tagFilter = f
}


func GetTable(data interface{}, table, selectColumns string, column []string, offset, limit int, args interface{}) error {
	if vi := reflect.ValueOf(data); vi.Kind() != reflect.Ptr {
		panic("must be reflect.Ptr")
	}
	query := getQuerySql(QueryDbParamsType, table, selectColumns, offset, limit, column...)
	return getByTx(data, query, args)
}

func NamedGetTable(data interface{}, table, selectColumns string, column []string, offset, limit int, arg interface{}) error {
	if vi := reflect.ValueOf(data); vi.Kind() != reflect.Ptr {
		panic("must be reflect.Ptr")
	}
	query := getQuerySql(QueryDbNamedType, table, selectColumns, offset, limit, column...)
	return namedGetByTx(data, query, arg)
}

func SelectTable(data interface{}, table string, column []string, args interface{}) error {
	if vi := reflect.ValueOf(data); vi.Kind() != reflect.Ptr {
		panic("must be reflect.Ptr")
	}
	query := getQuerySql(QueryDbParamsType, table, "*", 0, -1, column...)
	return selectByTx(data, query, args)
}

func NamedSelectTable(data interface{}, table string, column []string, arg interface{}) error {
	if vi := reflect.ValueOf(data); vi.Kind() != reflect.Ptr {
		panic("must be reflect.Ptr")
	}
	query := getQuerySql(QueryDbNamedType, table, "*", 0, -1, column...)
	return namedSelectByTx(data, query, arg)
}

func UpdateTable(table string, primiries, columns []string, args interface{}) error {
	query := getUpdateSql(QueryDbParamsType, table, primiries, columns...)
	return execByTx(query, args)
}

func NamedUpdateTable(table string, primiries, columns []string, arg interface{}) error {
	query := getUpdateSql(QueryDbNamedType, table, primiries, columns...)
	return namedExecByTx(query, arg)
}

func InsertTable(table string, columns []string, args interface{}) error {
	query := getInsertSql(QueryDbParamsType, table, columns...)
	return execByTx(query, args)
}

func NamedInsertTable(table string, arg interface{}) error {
	// TODO: Support Anonymous.
	var columns []string
	v := reflect.ValueOf(arg)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	if t.Kind() != reflect.Struct {
		panic("must be reflect.Struct")
	}
	for fieldPos := 0; fieldPos < t.NumField(); fieldPos++ {
		f := t.Field(fieldPos)
		if tagFilter(f) {
			columns = append(columns, f.Name)
		}
	}

	query := getInsertSql(QueryDbNamedType, table, columns...)
	return namedExecByTx(query, arg)
}

func DeleteTable(table string, columns []string, args interface{}) error {
	query := getDeleteSql(QueryDbParamsType, table, columns...)
	return execByTx(query, args)
}

func NamedDeleteTable(table string, columns []string, arg interface{}) error {
	query := getDeleteSql(QueryDbNamedType, table, columns...)
	return namedExecByTx(query, arg)
}
