package orm

import (
	"fmt"
	"strings"
)

const (
	QueryDbParamsType = iota
	QueryDbNamedType  = iota
)

func getQuerySql(queryType int, table, selectElem string, offset, limit int, columns ...string) string {
	sql := fmt.Sprintf("select %s from %s", selectElem, table)
	if len(columns) > 0 {
		sql += " where " + strings.Join(processColumns1(queryType, columns...), " and ")
	}

	sql = fmt.Sprintf("%s offset %d", sql, offset)
	if limit < 0 {
		return sql
	}
	return fmt.Sprintf("%s limit %d", sql, limit)
}

func getInsertSql(queryType int, table string, columns ...string) string {
	colNames := getSnakeColumns(columns...)
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
	colNames := getSnakeColumns(columns...)
	columnsPhrase := make([]string, 0, len(columns))
	if queryType == QueryDbParamsType {
		for _, col := range colNames {
			columnsPhrase = append(columnsPhrase, fmt.Sprintf("%s = ?", col))
		}
	}
	if queryType == QueryDbNamedType {
		for _, col := range colNames {
			columnsPhrase = append(columnsPhrase, fmt.Sprintf("%s = :%s", col, col))
		}
	}
	return columnsPhrase
}

func processColumns2(queryType int, columns ...string) string {
	colNames := getSnakeColumns(columns...)
	var columnsPhrase string
	if queryType == QueryDbParamsType {
		newColNames := make([]string, 0, len(colNames))
		for _ := range colNames {
			newColNames = append(newColNames, "?")
		}
		columnsPhrase = strings.Join(newColNames, ", ")
	}
	if queryType == QueryDbNamedType {
		columnsPhrase = ":" + strings.Join(colNames, ", :")
	}
	return columnsPhrase
}

func getSnakeColumns(columns ...string) []string {
	colNames := make([]string, 0, len(columns))
	for _, col := range columns {
		colNames = append(colNames, mapper(col))
	}
	return colNames
}
