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

var _ pg.BeforeInsertHook = (*Qualification)(nil)
var _ pg.BeforeUpdateHook = (*Qualification)(nil)

type Qualification struct {
	tableName struct{} `pg:"alias:qualification"`

	ID          int
	Slug        string `pg:",unique"`
	Name        string `pg:",unique:group_1"`
	Code        string `pg:",unique:group_1"`
	Formula     string
	Description string
	CreatedAt   time.Time `pg:"default:now()"`
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

type QualificationToProfession struct {
	ID              int
	QualificationID int            `pg:"on_delete:CASCADE,unique:group_1"`
	Qualification   *Qualification `pg:"rel:has-one"`
	ProfessionID    int            `pg:"on_delete:CASCADE,unique:group_1"`
	Profession      *Profession    `pg:"rel:has-one"`
}

func applyQualificationInputUpdates(input internal.QualificationInput) func(q *orm.Query) (*orm.Query, error) {
	return func(q *orm.Query) (*orm.Query, error) {
		if input.Name != nil {
			q = q.Set(gopgutil.BuildConditionEquals("name"), *input.Name)
		}
		if input.Code != nil {
			q = q.Set(gopgutil.BuildConditionEquals("code"), *input.Code)
		}
		if input.Formula != nil {
			q = q.Set(gopgutil.BuildConditionEquals("formula"), *input.Formula)
		}
		if input.Description != nil {
			q = q.Set(gopgutil.BuildConditionEquals("description"), *input.Description)
		}

		return q, nil
	}
}

func applyQualificationFilterOr(f internal.QualificationFilterOr, alias string) func(q *orm.Query) (*orm.Query, error) {
	return func(q *orm.Query) (*orm.Query, error) {
		q = q.WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			if !isZero(f.NameMATCH) {
				q = q.WhereOr(gopgutil.BuildConditionMatch("?"), gopgutil.AddAliasToColumnName("name", alias), f.NameMATCH)
			}
			if !isZero(f.NameIEQ) {
				q = q.WhereOr(gopgutil.BuildConditionIEQ("?"), gopgutil.AddAliasToColumnName("name", alias), f.NameIEQ)
			}

			if !isZero(f.CodeMATCH) {
				q = q.WhereOr(gopgutil.BuildConditionMatch("?"), gopgutil.AddAliasToColumnName("code", alias), f.CodeMATCH)
			}
			if !isZero(f.CodeIEQ) {
				q = q.WhereOr(gopgutil.BuildConditionIEQ("?"), gopgutil.AddAliasToColumnName("code", alias), f.CodeIEQ)
			}

			return q, nil
		})

		return q, nil
	}
}

func applyQualificationFilter(f internal.QualificationFilter, alias string) func(q *orm.Query) (*orm.Query, error) {
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

		if !isZero(f.Code) {
			q = q.Where(gopgutil.BuildConditionArray("?"), gopgutil.AddAliasToColumnName("code", alias), pg.Array(f.Code))
		}
		if !isZero(f.CodeNEQ) {
			q = q.Where(gopgutil.BuildConditionNotInArray("?"), gopgutil.AddAliasToColumnName("code", alias), pg.Array(f.CodeNEQ))
		}
		if !isZero(f.CodeMATCH) {
			q = q.Where(gopgutil.BuildConditionMatch("?"), gopgutil.AddAliasToColumnName("code", alias), f.CodeMATCH)
		}
		if !isZero(f.CodeIEQ) {
			q = q.Where(gopgutil.BuildConditionIEQ("?"), gopgutil.AddAliasToColumnName("code", alias), f.CodeIEQ)
		}

		if !isZero(f.Formula) {
			q = q.Where(gopgutil.BuildConditionArray("?"), gopgutil.AddAliasToColumnName("formula", alias), pg.Array(f.Formula))
		}
		if !isZero(f.FormulaNEQ) {
			q = q.Where(gopgutil.BuildConditionNotInArray("?"), gopgutil.AddAliasToColumnName("formula", alias), pg.Array(f.FormulaNEQ))
		}

		if !isZero(f.DescriptionMATCH) {
			q = q.Where(gopgutil.BuildConditionMatch("?"), gopgutil.AddAliasToColumnName("description", alias), f.DescriptionMATCH)
		}
		if !isZero(f.DescriptionIEQ) {
			q = q.Where(gopgutil.BuildConditionIEQ("?"), gopgutil.AddAliasToColumnName("description", alias), f.DescriptionIEQ)
		}

		if !isZero(f.ProfessionID) {
			sliceLength := len(f.ProfessionID)
			subquery := q.
				New().
				Model(&QualificationToProfession{}).
				ColumnExpr(gopgutil.BuildCountColumnExpr("profession_id", "count")).
				Column("qualification_id").
				Where(gopgutil.BuildConditionArray("profession_id"), pg.Array(f.ProfessionID)).
				Group("qualification_id")

			q = q.
				Join(`INNER JOIN (?) AS qualification_to_professions ON qualification_to_professions.qualification_id = qualification.id`, subquery).
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

		q = q.Apply(applyQualificationFilterOr(f.Or, alias))

		return q, nil
	}
}
