package sql

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	sortRegex = regexp.MustCompile(`^[\p{L}\_\.]+$`)
)

func SanitizeOrder(order string) string {
	parts := strings.Split(strings.TrimSpace(order), " ")
	length := len(parts)

	if length != 2 || !sortRegex.Match([]byte(parts[0])) {
		return ""
	}

	table := ""
	column := parts[0]
	if strings.Contains(parts[0], ".") {
		columnAndTable := strings.Split(parts[0], ".")
		table = underscore(columnAndTable[0]) + "."
		column = columnAndTable[1]
	}

	direction := "ASC"
	if strings.ToUpper(parts[1]) == "DESC" {
		direction = "DESC"
	}

	return strings.ToLower(table+underscore(column)) + " " + direction
}

func SanitizeOrders(orders []string) []string {
	var sanitizedOrders []string
	for _, sort := range orders {
		sanitized := SanitizeOrder(sort)
		if sanitized != "" {
			sanitizedOrders = append(sanitizedOrders, sanitized)
		}
	}
	return sanitizedOrders
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
