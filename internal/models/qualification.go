package models

import (
	"context"
	"time"

	"github.com/gosimple/slug"
)

type Qualification struct {
	ID          int       `json:"id" xml:"id" gqlgen:"id"`
	Slug        string    `json:"slug" pg:",unique" xml:"slug" gqlgen:"slug"`
	Name        string    `json:"name" pg:",unique:group_1" xml:"name" gqlgen:"name"`
	Code        string    `json:"code" pg:",unique:group_1" xml:"code" gqlgen:"code"`
	Formula     string    `json:"formula" xml:"formula" gqlgen:"formula"`
	Description string    `json:"description" xml:"description" gqlgen:"description"`
	CreatedAt   time.Time `json:"createdAt,omitempty" pg:"default:now()" xml:"createdAt" gqlgen:"createdAt"`
}

func (q *Qualification) BeforeInsert(ctx context.Context) (context.Context, error) {
	q.CreatedAt = time.Now()
	q.Slug = slug.Make(q.Code)

	return ctx, nil
}

func (q *Qualification) BeforeUpdate(ctx context.Context) (context.Context, error) {
	if q.Code != "" {
		q.Slug = slug.Make(q.Code)
	}

	return ctx, nil
}
