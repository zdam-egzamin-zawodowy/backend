package db

import (
	"strings"

	"github.com/go-pg/pg/v10/orm"
)

type Sort struct {
	Relationships map[string]string
	Orders        []string
}

func (s Sort) Apply(q *orm.Query) (*orm.Query, error) {
	for _, order := range s.Orders {
		if alias := s.extractAlias(order); alias != "" && s.Relationships[alias] != "" {
			q = q.Relation(s.Relationships[alias])
		}
	}
	return q.Order(s.Orders...), nil
}

func (s Sort) extractAlias(order string) string {
	if strings.Contains(order, ".") {
		return strings.Split(order, ".")[0]
	}
	return ""
}
