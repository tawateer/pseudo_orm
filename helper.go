package orm

import (
	"fmt"
	"strings"
)

const (
	QueryDbParamsType = iota
	QueryDbNamedType

	ColumnSplitSymbol = ","
)

func getQuerySql(queryType int, table, selectColumns string, offset, limit int, columns ...string) string {
	sql := fmt.Sprintf("select %s from %s", selectColumns, table)
	if len(columns) > 0 {
		sql += " where " + strings.Join(processColumns1(queryType, columns...), " and ")
	}

	if offset >= 0 {
		sql = fmt.Sprintf("%s offset %d", sql, offset)
	}
	if limit >= 0 {
		sql = fmt.Sprintf("%s limit %d", sql, limit)
	}
	return sql
}

func getInsertSql(queryType int, table string, columns ...string) string {
	colNames, _ := getColumns(columns...)
	var secureColNames []string
	for _, colName := range colNames {
		secureColNames = append(secureColNames, "`"+colName+"`")
	}
	valStr := processColumns2(queryType, columns...)

	return fmt.Sprintf("insert into %s (%s) values (%s)",
		table,
		strings.Join(secureColNames, ", "),
		valStr)
}

func getUpdateSql(queryType int, table string, primaries []string, columns ...string) string {
	setPhrase := processColumns1(queryType, columns...)
	wherePhrase := processColumns1(queryType, primaries...)

	return fmt.Sprintf("update %s set %s where %s",
		table,
		strings.Join(setPhrase, ", "),
		strings.Join(wherePhrase, " and "))
}

func getDeleteSql(queryType int, table string, columns ...string) string {
	sql := fmt.Sprintf("delete from %s", table)
	if len(columns) > 0 {
		sql += " where " + strings.Join(processColumns1(queryType, columns...), " and ")
	}
	return sql
}

//

func processColumns1(queryType int, columns ...string) []string {
	colNames, colSymbols := getColumns(columns...)
	columnsPhrase := make([]string, 0, len(columns))
	if queryType == QueryDbParamsType {
		for i, col := range colNames {
			columnsPhrase = append(columnsPhrase, fmt.Sprintf("%s %s ?", col, colSymbols[i]))
		}
	}
	if queryType == QueryDbNamedType {
		for i, col := range colNames {
			columnsPhrase = append(columnsPhrase, fmt.Sprintf("%s %s :%s", col, col, colSymbols[i]))
		}
	}
	return columnsPhrase
}

func processColumns2(queryType int, columns ...string) string {
	colNames, _ := getColumns(columns...)
	var columnsPhrase string
	if queryType == QueryDbParamsType {
		newColNames := make([]string, 0, len(colNames))
		for range colNames {
			newColNames = append(newColNames, "?")
		}
		columnsPhrase = strings.Join(newColNames, ", ")
	}
	if queryType == QueryDbNamedType {
		columnsPhrase = ":" + strings.Join(colNames, ", :")
	}
	return columnsPhrase
}

func getColumns(columns ...string) ([]string, []string) {
	colNames := make([]string, 0, len(columns))
	colSymbols := make([]string, 0, len(columns))

	for _, col := range columns {
		n, s := splitColumn(col)
		colNames = append(colNames, mapper(n))
		colSymbols = append(colSymbols, s)
	}
	return colNames, colSymbols
}

func splitColumn(column string) (string, string) {
	i := strings.Index(column, ColumnSplitSymbol)
	if i == -1 {
		return column, "="
	}
	return column[:i], column[i+1:]
}
