package sql

import "strings"

func AddAliasToColumnName(column, alias string) string {
	if alias != "" && !strings.HasPrefix(column, alias+".") {
		column = wrapStringInDoubleQuotes(alias) + "." + wrapStringInDoubleQuotes(column)
	} else {
		column = wrapStringInDoubleQuotes(column)
	}
	return column
}

func wrapStringInDoubleQuotes(str string) string {
	return `"` + str + `"`
}
