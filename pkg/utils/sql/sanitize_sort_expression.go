package sqlutils

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	sortexprRegex = regexp.MustCompile(`^[\p{L}\_\.]+$`)
)

func SanitizeSortExpression(expr string) string {
	trimmed := strings.TrimSpace(expr)
	splitted := strings.Split(trimmed, " ")
	length := len(splitted)
	if length != 2 || !sortexprRegex.Match([]byte(splitted[0])) {
		return ""
	}
	table := ""
	column := splitted[0]
	if strings.Contains(splitted[0], ".") {
		columnAndTable := strings.Split(splitted[0], ".")
		table = underscore(columnAndTable[0]) + "."
		column = columnAndTable[1]
	}
	keyword := "ASC"
	if strings.ToUpper(splitted[1]) == "DESC" {
		keyword = "DESC"
	}
	return strings.ToLower(table+underscore(column)) + " " + keyword
}

func SanitizeSortExpressions(exprs []string) []string {
	sanitizedExprs := []string{}
	for _, expr := range exprs {
		sanitized := SanitizeSortExpression(expr)
		if sanitized != "" {
			sanitizedExprs = append(sanitizedExprs, sanitized)
		}
	}
	return sanitizedExprs
}

type buffer struct {
	r         []byte
	runeBytes [utf8.UTFMax]byte
}

func (b *buffer) write(r rune) {
	if r < utf8.RuneSelf {
		b.r = append(b.r, byte(r))
		return
	}
	n := utf8.EncodeRune(b.runeBytes[0:], r)
	b.r = append(b.r, b.runeBytes[0:n]...)
}

func (b *buffer) indent() {
	if len(b.r) > 0 {
		b.r = append(b.r, '_')
	}
}

func underscore(s string) string {
	b := buffer{
		r: make([]byte, 0, len(s)),
	}
	var m rune
	var w bool
	for _, ch := range s {
		if unicode.IsUpper(ch) {
			if m != 0 {
				if !w {
					b.indent()
					w = true
				}
				b.write(m)
			}
			m = unicode.ToLower(ch)
		} else {
			if m != 0 {
				b.indent()
				b.write(m)
				m = 0
				w = false
			}
			b.write(ch)
		}
	}
	if m != 0 {
		if !w {
			b.indent()
		}
		b.write(m)
	}
	return string(b.r)
}
