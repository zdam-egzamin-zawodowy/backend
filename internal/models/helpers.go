package models

import (
	"strings"
	"time"
)

func addAliasToColumnName(column, prefix string) string {
	if prefix != "" && !strings.HasPrefix(column, prefix+".") {
		column = wrapStringInDoubleQuotes(prefix) + "." + wrapStringInDoubleQuotes(column)
	} else {
		column = wrapStringInDoubleQuotes(column)
	}
	return column
}

func wrapStringInDoubleQuotes(str string) string {
	return `"` + str + `"`
}

func buildConditionEquals(column string) string {
	return column + " = ?"
}

func buildConditionLT(column string) string {
	return column + " < ?"
}

func buildConditionLTE(column string) string {
	return column + " <= ?"
}

func buildConditionGT(column string) string {
	return column + " > ?"
}

func buildConditionGTE(column string) string {
	return column + " >= ?"
}

func buildConditionMatch(column string) string {
	return column + " LIKE ?"
}

func buildConditionIEQ(column string) string {
	return column + " ILIKE ?"
}

func buildConditionArray(column string) string {
	return column + " = ANY(?)"
}

func buildConditionNotInArray(column string) string {
	return "NOT (" + buildConditionArray(column) + ")"
}

func isZero(v interface{}) bool {
	switch c := v.(type) {
	case string:
		return c == ""
	case *string:
		return c == nil
	case []string:
		return c == nil || len(c) == 0
	case []Role:
		return c == nil || len(c) == 0
	case int:
		return c == 0
	case *int:
		return c == nil
	case []int:
		return c == nil || len(c) == 0
	case float64:
		return c == 0
	case *float64:
		return c == nil
	case float32:
		return c == 0
	case *float32:
		return c == nil
	case bool:
		return !c
	case *bool:
		return c == nil
	case time.Time:
		return c.IsZero()
	case *time.Time:
		return c == nil
	default:
		return false
	}
}
