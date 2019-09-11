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


func GetTable(data interface{}, table string, column []string, args interface{}) error {
	if vi := reflect.ValueOf(data); vi.Kind() != reflect.Ptr {
		panic("must be reflect.Ptr")
	}
	query := getQuerySql(QueryDbParamsType, table, "*", 0, -1, column...)
	return getSqlByTx(data, query, args)
}

func NamedGetTable(data interface{}, table string, column []string, arg interface{}) error {
	if vi := reflect.ValueOf(data); vi.Kind() != reflect.Ptr {
		panic("must be reflect.Ptr")
	}
	query := getQuerySql(QueryDbNamedType, table, "*", 0, -1, column...)
	return namedGetSqlByTx(data, query, arg)
}

func SelectTable(data interface{}, table string, column []string, args interface{}) error {
	if vi := reflect.ValueOf(data); vi.Kind() != reflect.Ptr {
		panic("must be reflect.Ptr")
	}
	query := getQuerySql(QueryDbParamsType, table, "*", 0, -1, column...)
	return selectSqlByTx(data, query, args)
}

func NamedSelectTable(data interface{}, table string, column []string, arg interface{}) error {
	if vi := reflect.ValueOf(data); vi.Kind() != reflect.Ptr {
		panic("must be reflect.Ptr")
	}
	query := getQuerySql(QueryDbNamedType, table, "*", 0, -1, column...)
	return namedSelectSqlByTx(data, query, arg)
}

func UpdateTable(table string, primiries, columns []string, args interface{}) error {
	query := getUpdateSql(QueryDbParamsType, table, primiries, columns...)
	return execSqlByTx(query, args)
}

func NamedUpdateTable(table string, primiries, columns []string, arg interface{}) error {
	query := getUpdateSql(QueryDbNamedType, table, primiries, columns...)
	return namedExecSqlByTx(query, arg)
}

func InsertTable(table string, columns []string, args interface{}) error {
	query := getInsertSql(QueryDbParamsType, table, columns...)
	return execSqlByTx(query, args)
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
	return namedExecSqlByTx(query, arg)
}

func DeleteTable(table string, columns []string, args interface{}) error {
	query := getDeleteSql(QueryDbParamsType, table, columns...)
	return execSqlByTx(query, args)
}

func NamedDeleteTable(table string, columns []string, arg interface{}) error {
	query := getDeleteSql(QueryDbNamedType, table, columns...)
	return namedExecSqlByTx(query, arg)
}
