package postgres

import (
	"context"
	"github.com/Kichiyaki/gopgutil/v10"
	"github.com/zdam-egzamin-zawodowy/backend/internal"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/gosimple/slug"
)

var _ pg.BeforeInsertHook = (*Profession)(nil)
var _ pg.BeforeUpdateHook = (*Profession)(nil)

type Profession struct {
	tableName struct{} `pg:"alias:profession"`

	ID          int
	Slug        string `pg:",unique"`
	Name        string `pg:",unique"`
	Description string
	CreatedAt   time.Time `pg:"default:now()"`
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

func applyProfessionInputUpdates(input internal.ProfessionInput) func(q *orm.Query) (*orm.Query, error) {
	return func(q *orm.Query) (*orm.Query, error) {
		if input.Name != nil {
			q = q.Set(gopgutil.BuildConditionEquals("name"), *input.Name)
		}

		if input.Description != nil {
			q = q.Set(gopgutil.BuildConditionEquals("description"), *input.Description)
		}

		return q, nil
	}
}

func applyProfessionFilter(f internal.ProfessionFilter, alias string) func(q *orm.Query) (*orm.Query, error) {
	return func(q *orm.Query) (*orm.Query, error) {
		if !isZero(f.ID) {
			q = q.Where(gopgutil.BuildConditionArray("?"), gopgutil.AddAliasToColumnName("id", alias), pg.Array(f.ID))
		}
		if !isZero(f.IDNEQ) {
			q = q.Where(gopgutil.BuildConditionNotInArray("?"), gopgutil.AddAliasToColumnName("id", alias), pg.Array(f.IDNEQ))
		}

		if !isZero(f.Slug) {
			q = q.Where(gopgutil.BuildConditionArray("?"), gopgutil.AddAliasToColumnName("slug", alias), pg.Array(f.Slug))
		}
		if !isZero(f.SlugNEQ) {
			q = q.Where(gopgutil.BuildConditionNotInArray("?"), gopgutil.AddAliasToColumnName("slug", alias), pg.Array(f.SlugNEQ))
		}

		if !isZero(f.Name) {
			q = q.Where(gopgutil.BuildConditionArray("?"), gopgutil.AddAliasToColumnName("name", alias), pg.Array(f.Name))
		}
		if !isZero(f.NameNEQ) {
			q = q.Where(gopgutil.BuildConditionNotInArray("?"), gopgutil.AddAliasToColumnName("name", alias), pg.Array(f.NameNEQ))
		}
		if !isZero(f.NameMATCH) {
			q = q.Where(gopgutil.BuildConditionMatch("?"), gopgutil.AddAliasToColumnName("name", alias), f.NameMATCH)
		}
		if !isZero(f.NameIEQ) {
			q = q.Where(gopgutil.BuildConditionIEQ("?"), gopgutil.AddAliasToColumnName("name", alias), f.NameIEQ)
		}

		if !isZero(f.DescriptionMATCH) {
			q = q.Where(gopgutil.BuildConditionMatch("?"), gopgutil.AddAliasToColumnName("description", alias), f.DescriptionMATCH)
		}
		if !isZero(f.DescriptionIEQ) {
			q = q.Where(gopgutil.BuildConditionIEQ("?"), gopgutil.AddAliasToColumnName("description", alias), f.DescriptionIEQ)
		}

		if !isZero(f.QualificationID) {
			sliceLength := len(f.QualificationID)
			subquery := q.
				New().
				Model(&QualificationToProfession{}).
				ColumnExpr(gopgutil.BuildCountColumnExpr("qualification_id", "count")).
				Column("profession_id").
				Where(gopgutil.BuildConditionArray("qualification_id"), pg.Array(f.QualificationID)).
				Group("profession_id")

			q = q.
				Join(`INNER JOIN (?) AS qualification_to_professions ON qualification_to_professions.profession_id = profession.id`, subquery).
				Where(gopgutil.BuildConditionGTE("count"), sliceLength)
		}

		if !isZero(f.CreatedAt) {
			q = q.Where(gopgutil.BuildConditionEquals("?"), gopgutil.AddAliasToColumnName("created_at", alias), f.CreatedAt)
		}
		if !isZero(f.CreatedAtGT) {
			q = q.Where(gopgutil.BuildConditionGT("?"), gopgutil.AddAliasToColumnName("created_at", alias), f.CreatedAtGT)
		}
		if !isZero(f.CreatedAtGTE) {
			q = q.Where(gopgutil.BuildConditionGTE("?"), gopgutil.AddAliasToColumnName("created_at", alias), f.CreatedAtGTE)
		}
		if !isZero(f.CreatedAtLT) {
			q = q.Where(gopgutil.BuildConditionLT("?"), gopgutil.AddAliasToColumnName("created_at", alias), f.CreatedAtLT)
		}
		if !isZero(f.CreatedAtLTE) {
			q = q.Where(gopgutil.BuildConditionLTE("?"), gopgutil.AddAliasToColumnName("created_at", alias), f.CreatedAtLTE)
		}

		return q, nil
	}
}
