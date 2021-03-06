package models

import (
	"context"
	"strings"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/gosimple/slug"
	sqlutils "github.com/zdam-egzamin-zawodowy/backend/pkg/utils/sql"
)

type Profession struct {
	tableName struct{} `pg:"alias:profession"`

	ID          int       `json:"id,omitempty" xml:"id" gqlgen:"id"`
	Slug        string    `json:"slug" pg:",unique" xml:"slug" gqlgen:"slug"`
	Name        string    `json:"name,omitempty" pg:",unique" xml:"name" gqlgen:"name"`
	Description string    `json:"description,omitempty" xml:"description" gqlgen:"description"`
	CreatedAt   time.Time `json:"createdAt,omitempty" pg:"default:now()" xml:"createdAt" gqlgen:"createdAt"`
}

func (p *Profession) BeforeInsert(ctx context.Context) (context.Context, error) {
	p.CreatedAt = time.Now()
	p.Slug = slug.Make(p.Name)

	return ctx, nil
}

func (p *Profession) BeforeUpdate(ctx context.Context) (context.Context, error) {
	if p.Name != "" {
		p.Slug = slug.Make(p.Name)
	}

	return ctx, nil
}

type ProfessionInput struct {
	Name        *string `json:"name,omitempty" xml:"name" gqlgen:"name"`
	Description *string `json:"description,omitempty" xml:"description" gqlgen:"description"`
}

func (input *ProfessionInput) IsEmpty() bool {
	return input == nil && input.Name == nil && input.Description == nil
}

func (input *ProfessionInput) Sanitize() *ProfessionInput {
	if input.Name != nil {
		*input.Name = strings.TrimSpace(*input.Name)
	}
	if input.Description != nil {
		*input.Description = strings.TrimSpace(*input.Description)
	}

	return input
}

func (input *ProfessionInput) ToProfession() *Profession {
	p := &Profession{}
	if input.Name != nil {
		p.Name = *input.Name
	}
	if input.Description != nil {
		p.Description = *input.Description
	}
	return p
}

func (input *ProfessionInput) ApplyUpdate(q *orm.Query) (*orm.Query, error) {
	if !input.IsEmpty() {
		if input.Name != nil {
			q = q.Set("name = ?", *input.Name)
		}
		if input.Description != nil {
			q = q.Set("description = ?", *input.Description)
		}
	}

	return q, nil
}

type ProfessionFilter struct {
	ID    []int `gqlgen:"id" json:"id" xml:"id"`
	IDNEQ []int `gqlgen:"idNEQ" json:"idNEQ" xml:"idNEQ"`

	Slug    []string `json:"slug" xml:"slug" gqlgen:"slug"`
	SlugNEQ []string `json:"slugNEQ" xml:"slugNEQ" gqlgen:"slugNEQ"`

	Name      []string `gqlgen:"name" json:"name" xml:"name"`
	NameNEQ   []string `json:"nameNEQ" xml:"nameNEQ" gqlgen:"nameNEQ"`
	NameMATCH string   `json:"nameMATCH" xml:"nameMATCH" gqlgen:"nameMATCH"`
	NameIEQ   string   `gqlgen:"nameIEQ" json:"nameIEQ" xml:"nameIEQ"`

	DescriptionMATCH string `gqlgen:"descriptionMATCH" json:"descriptionMATCH" xml:"descriptionMATCH"`
	DescriptionIEQ   string `json:"descriptionIEQ" xml:"descriptionIEQ" gqlgen:"descriptionIEQ"`

	CreatedAt    time.Time `gqlgen:"createdAt" json:"createdAt" xml:"createdAt"`
	CreatedAtGT  time.Time `gqlgen:"createdAtGT" json:"createdAtGT" xml:"createdAtGT"`
	CreatedAtGTE time.Time `json:"createdAtGTE" xml:"createdAtGTE" gqlgen:"createdAtGTE"`
	CreatedAtLT  time.Time `gqlgen:"createdAtLT" json:"createdAtLT" xml:"createdAtLT"`
	CreatedAtLTE time.Time `json:"createdAtLTE" xml:"createdAtLTE" gqlgen:"createdAtLTE"`
}

func (f *ProfessionFilter) WhereWithAlias(q *orm.Query, alias string) (*orm.Query, error) {
	if f == nil {
		return q, nil
	}

	if !isZero(f.ID) {
		q = q.Where(sqlutils.BuildConditionArray(sqlutils.AddAliasToColumnName("id", alias)), pg.Array(f.ID))
	}
	if !isZero(f.IDNEQ) {
		q = q.Where(sqlutils.BuildConditionNotInArray(sqlutils.AddAliasToColumnName("id", alias)), pg.Array(f.IDNEQ))
	}

	if !isZero(f.Slug) {
		q = q.Where(sqlutils.BuildConditionArray(sqlutils.AddAliasToColumnName("slug", alias)), pg.Array(f.Slug))
	}
	if !isZero(f.SlugNEQ) {
		q = q.Where(sqlutils.BuildConditionNotInArray(sqlutils.AddAliasToColumnName("slug", alias)), pg.Array(f.SlugNEQ))
	}

	if !isZero(f.Name) {
		q = q.Where(sqlutils.BuildConditionArray(sqlutils.AddAliasToColumnName("name", alias)), pg.Array(f.Name))
	}
	if !isZero(f.NameNEQ) {
		q = q.Where(sqlutils.BuildConditionNotInArray(sqlutils.AddAliasToColumnName("name", alias)), pg.Array(f.NameNEQ))
	}
	if !isZero(f.NameMATCH) {
		q = q.Where(sqlutils.BuildConditionMatch(sqlutils.AddAliasToColumnName("name", alias)), f.NameMATCH)
	}
	if !isZero(f.NameIEQ) {
		q = q.Where(sqlutils.BuildConditionIEQ(sqlutils.AddAliasToColumnName("name", alias)), f.NameIEQ)
	}

	if !isZero(f.DescriptionMATCH) {
		q = q.Where(sqlutils.BuildConditionMatch(sqlutils.AddAliasToColumnName("description", alias)), f.DescriptionMATCH)
	}
	if !isZero(f.DescriptionIEQ) {
		q = q.Where(sqlutils.BuildConditionIEQ(sqlutils.AddAliasToColumnName("description", alias)), f.DescriptionIEQ)
	}

	if !isZero(f.CreatedAt) {
		q = q.Where(sqlutils.BuildConditionEquals(sqlutils.AddAliasToColumnName("created_at", alias)), f.CreatedAt)
	}
	if !isZero(f.CreatedAtGT) {
		q = q.Where(sqlutils.BuildConditionGT(sqlutils.AddAliasToColumnName("created_at", alias)), f.CreatedAtGT)
	}
	if !isZero(f.CreatedAtGTE) {
		q = q.Where(sqlutils.BuildConditionGTE(sqlutils.AddAliasToColumnName("created_at", alias)), f.CreatedAtGTE)
	}
	if !isZero(f.CreatedAtLT) {
		q = q.Where(sqlutils.BuildConditionLT(sqlutils.AddAliasToColumnName("created_at", alias)), f.CreatedAtLT)
	}
	if !isZero(f.CreatedAtLTE) {
		q = q.Where(sqlutils.BuildConditionLTE(sqlutils.AddAliasToColumnName("created_at", alias)), f.CreatedAtLTE)
	}

	return q, nil
}

func (f *ProfessionFilter) Where(q *orm.Query) (*orm.Query, error) {
	return f.WhereWithAlias(q, "profession")
}
