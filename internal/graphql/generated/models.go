// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package generated

import (
	"github.com/zdam-egzamin-zawodowy/backend/internal/model"
)

type ProfessionList struct {
	Total int                 `json:"total"`
	Items []*model.Profession `json:"items"`
}

type QualificationList struct {
	Total int                    `json:"total"`
	Items []*model.Qualification `json:"items"`
}

type QuestionList struct {
	Total int               `json:"total"`
	Items []*model.Question `json:"items"`
}

type UserList struct {
	Total int           `json:"total"`
	Items []*model.User `json:"items"`
}

type UserWithToken struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}
