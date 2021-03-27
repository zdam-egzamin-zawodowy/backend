package sqlutils

import (
	"fmt"
	"strings"
)

func AddAliasToColumnName(column, alias string) string {
	if alias != "" && !strings.HasPrefix(column, alias+".") {
		column = WrapStringInDoubleQuotes(alias) + "." + WrapStringInDoubleQuotes(column)
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

func BuildConditionNEQ(column string) string {
	return column + " != ?"
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

func BuildConditionIn(column string) string {
	return column + " IN (?)"
}

func BuildConditionArray(column string) string {
	return column + " = ANY(?)"
}

func BuildConditionNotInArray(column string) string {
	return "NOT (" + BuildConditionArray(column) + ")"
}

func BuildCountColumnExpr(column, alias string) string {
	base := fmt.Sprintf("count(%s)", column)
	if alias != "" {
		return base + " as " + alias
	}
	return base
}
