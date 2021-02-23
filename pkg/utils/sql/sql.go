package sqlutils

import "strings"

func AddAliasToColumnName(column, prefix string) string {
	if prefix != "" && !strings.HasPrefix(column, prefix+".") {
		column = WrapStringInDoubleQuotes(prefix) + "." + WrapStringInDoubleQuotes(column)
	} else {
		column = WrapStringInDoubleQuotes(column)
	}
	return column
}

func WrapStringInDoubleQuotes(str string) string {
	return `"` + str + `"`
}

func BuildConditionEquals(column string) string {
	return column + " = ?"
}

func BuildConditionLT(column string) string {
	return column + " < ?"
}

func BuildConditionLTE(column string) string {
	return column + " <= ?"
}

func BuildConditionGT(column string) string {
	return column + " > ?"
}

func BuildConditionGTE(column string) string {
	return column + " >= ?"
}

func BuildConditionMatch(column string) string {
	return column + " LIKE ?"
}

func BuildConditionIEQ(column string) string {
	return column + " ILIKE ?"
}

func BuildConditionArray(column string) string {
	return column + " = ANY(?)"
}

func BuildConditionNotInArray(column string) string {
	return "NOT (" + BuildConditionArray(column) + ")"
}
