package models

import (
	"context"
	"github.com/Kichiyaki/gopgutil/v10"
	"strings"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/gosimple/slug"
)

var _ pg.BeforeInsertHook = (*Qualification)(nil)
var _ pg.BeforeUpdateHook = (*Qualification)(nil)

type Qualification struct {
	tableName struct{} `pg:"alias:qualification"`

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

type QualificationToProfession struct {
	ID              int            `json:"id" xml:"id" gqlgen:"id"`
	QualificationID int            `pg:"on_delete:CASCADE,unique:group_1" json:"qualificationID" xml:"qualificationID" gqlgen:"qualificationID"`
	Qualification   *Qualification `pg:"rel:has-one" json:"qualification" xml:"qualification" gqlgen:"qualification"`
	ProfessionID    int            `pg:"on_delete:CASCADE,unique:group_1" json:"professionID" xml:"professionID" gqlgen:"professionID"`
	Profession      *Profession    `pg:"rel:has-one" json:"profession" xml:"profession" gqlgen:"profession"`
}

type QualificationInput struct {
	Name                 *string `json:"name" xml:"name" gqlgen:"name"`
	Description          *string `json:"description" xml:"description" gqlgen:"description"`
	Code                 *string `json:"code" xml:"code" gqlgen:"code"`
	Formula              *string `json:"formula" xml:"formula" gqlgen:"formula"`
	AssociateProfession  []int   `json:"associateProfession" xml:"associateProfession" gqlgen:"associateProfession"`
	DissociateProfession []int   `json:"dissociateProfession" xml:"dissociateProfession" gqlgen:"dissociateProfession"`
}

func (input *QualificationInput) IsEmpty() bool {
	return input == nil &&
		input.Name == nil &&
		input.Code == nil &&
		input.Formula == nil &&
		input.Description == nil &&
		len(input.AssociateProfession) == 0 &&
		len(input.DissociateProfession) == 0
}

func (input *QualificationInput) HasBasicDataToUpdate() bool {
	return input != nil &&
		(input.Name != nil ||
			input.Code != nil ||
			input.Formula != nil ||
			input.Description != nil)
}

func (input *QualificationInput) Sanitize() *QualificationInput {
	if input.Name != nil {
		*input.Name = strings.TrimSpace(*input.Name)
	}
	if input.Description != nil {
		*input.Description = strings.TrimSpace(*input.Description)
	}
	if input.Code != nil {
		*input.Code = strings.TrimSpace(*input.Code)
	}
	if input.Formula != nil {
		*input.Formula = strings.TrimSpace(*input.Formula)
	}

	return input
}

func (input *QualificationInput) ToQualification() *Qualification {
	q := &Qualification{}
	if input.Name != nil {
		q.Name = *input.Name
	}
	if input.Description != nil {
		q.Description = *input.Description
	}
	if input.Code != nil {
		q.Code = *input.Code
	}
	if input.Formula != nil {
		q.Formula = *input.Formula
	}
	return q
}

func (input *QualificationInput) ApplyUpdate(q *orm.Query) (*orm.Query, error) {
	if !input.IsEmpty() {
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
	}

	return q, nil
}

type QualificationFilterOr struct {
	NameMATCH string `json:"nameMATCH" xml:"nameMATCH" gqlgen:"nameMATCH"`
	NameIEQ   string `gqlgen:"nameIEQ" json:"nameIEQ" xml:"nameIEQ"`

	CodeMATCH string `json:"codeMATCH" xml:"codeMATCH" gqlgen:"codeMATCH"`
	CodeIEQ   string `gqlgen:"codeIEQ" json:"codeIEQ" xml:"codeIEQ"`
}

func (f *QualificationFilterOr) WhereWithAlias(q *orm.Query, alias string) *orm.Query {
	if f == nil {
		return q
	}

	q = q.WhereGroup(func(q *orm.Query) (*orm.Query, error) {
		if !isZero(f.NameMATCH) {
			q = q.WhereOr(gopgutil.BuildConditionMatch(gopgutil.AddAliasToColumnName("name", alias)), f.NameMATCH)
		}
		if !isZero(f.NameIEQ) {
			q = q.WhereOr(gopgutil.BuildConditionIEQ(gopgutil.AddAliasToColumnName("name", alias)), f.NameIEQ)
		}

		if !isZero(f.CodeMATCH) {
			q = q.WhereOr(gopgutil.BuildConditionMatch(gopgutil.AddAliasToColumnName("code", alias)), f.CodeMATCH)
		}
		if !isZero(f.CodeIEQ) {
			q = q.WhereOr(gopgutil.BuildConditionIEQ(gopgutil.AddAliasToColumnName("code", alias)), f.CodeIEQ)
		}

		return q, nil
	})
	return q
}

type QualificationFilter struct {
	ID    []int `gqlgen:"id" json:"id" xml:"id"`
	IDNEQ []int `gqlgen:"idNEQ" json:"idNEQ" xml:"idNEQ"`

	Slug    []string `json:"slug" xml:"slug" gqlgen:"slug"`
	SlugNEQ []string `json:"slugNEQ" xml:"slugNEQ" gqlgen:"slugNEQ"`

	Formula    []string `gqlgen:"formula" json:"formula" xml:"formula"`
	FormulaNEQ []string `json:"formulaNEQ" xml:"formulaNEQ" gqlgen:"formulaNEQ"`

	Name      []string `gqlgen:"name" json:"name" xml:"name"`
	NameNEQ   []string `json:"nameNEQ" xml:"nameNEQ" gqlgen:"nameNEQ"`
	NameMATCH string   `json:"nameMATCH" xml:"nameMATCH" gqlgen:"nameMATCH"`
	NameIEQ   string   `gqlgen:"nameIEQ" json:"nameIEQ" xml:"nameIEQ"`

	Code      []string `gqlgen:"code" json:"code" xml:"code"`
	CodeNEQ   []string `json:"codeNEQ" xml:"codeNEQ" gqlgen:"codeNEQ"`
	CodeMATCH string   `json:"codeMATCH" xml:"codeMATCH" gqlgen:"codeMATCH"`
	CodeIEQ   string   `gqlgen:"codeIEQ" json:"codeIEQ" xml:"codeIEQ"`

	DescriptionMATCH string `gqlgen:"descriptionMATCH" json:"descriptionMATCH" xml:"descriptionMATCH"`
	DescriptionIEQ   string `json:"descriptionIEQ" xml:"descriptionIEQ" gqlgen:"descriptionIEQ"`

	ProfessionID []int `json:"professionID" xml:"professionID" gqlgen:"professionID"`

	CreatedAt    time.Time `gqlgen:"createdAt" json:"createdAt" xml:"createdAt"`
	CreatedAtGT  time.Time `gqlgen:"createdAtGT" json:"createdAtGT" xml:"createdAtGT"`
	CreatedAtGTE time.Time `json:"createdAtGTE" xml:"createdAtGTE" gqlgen:"createdAtGTE"`
	CreatedAtLT  time.Time `gqlgen:"createdAtLT" json:"createdAtLT" xml:"createdAtLT"`
	CreatedAtLTE time.Time `json:"createdAtLTE" xml:"createdAtLTE" gqlgen:"createdAtLTE"`

	Or *QualificationFilterOr `json:"or" xml:"or" gqlgen:"or"`
}

func (f *QualificationFilter) WhereWithAlias(q *orm.Query, alias string) (*orm.Query, error) {
	if f == nil {
		return q, nil
	}

	if !isZero(f.ID) {
		q = q.Where(gopgutil.BuildConditionArray(gopgutil.AddAliasToColumnName("id", alias)), pg.Array(f.ID))
	}
	if !isZero(f.IDNEQ) {
		q = q.Where(gopgutil.BuildConditionNotInArray(gopgutil.AddAliasToColumnName("id", alias)), pg.Array(f.IDNEQ))
	}

	if !isZero(f.Slug) {
		q = q.Where(gopgutil.BuildConditionArray(gopgutil.AddAliasToColumnName("slug", alias)), pg.Array(f.Slug))
	}
	if !isZero(f.SlugNEQ) {
		q = q.Where(gopgutil.BuildConditionNotInArray(gopgutil.AddAliasToColumnName("slug", alias)), pg.Array(f.SlugNEQ))
	}

	if !isZero(f.Name) {
		q = q.Where(gopgutil.BuildConditionArray(gopgutil.AddAliasToColumnName("name", alias)), pg.Array(f.Name))
	}
	if !isZero(f.NameNEQ) {
		q = q.Where(gopgutil.BuildConditionNotInArray(gopgutil.AddAliasToColumnName("name", alias)), pg.Array(f.NameNEQ))
	}
	if !isZero(f.NameMATCH) {
		q = q.Where(gopgutil.BuildConditionMatch(gopgutil.AddAliasToColumnName("name", alias)), f.NameMATCH)
	}
	if !isZero(f.NameIEQ) {
		q = q.Where(gopgutil.BuildConditionIEQ(gopgutil.AddAliasToColumnName("name", alias)), f.NameIEQ)
	}

	if !isZero(f.Code) {
		q = q.Where(gopgutil.BuildConditionArray(gopgutil.AddAliasToColumnName("code", alias)), pg.Array(f.Code))
	}
	if !isZero(f.CodeNEQ) {
		q = q.Where(gopgutil.BuildConditionNotInArray(gopgutil.AddAliasToColumnName("code", alias)), pg.Array(f.CodeNEQ))
	}
	if !isZero(f.CodeMATCH) {
		q = q.Where(gopgutil.BuildConditionMatch(gopgutil.AddAliasToColumnName("code", alias)), f.CodeMATCH)
	}
	if !isZero(f.CodeIEQ) {
		q = q.Where(gopgutil.BuildConditionIEQ(gopgutil.AddAliasToColumnName("code", alias)), f.CodeIEQ)
	}

	if !isZero(f.Formula) {
		q = q.Where(gopgutil.BuildConditionArray(gopgutil.AddAliasToColumnName("formula", alias)), pg.Array(f.Formula))
	}
	if !isZero(f.FormulaNEQ) {
		q = q.Where(gopgutil.BuildConditionNotInArray(gopgutil.AddAliasToColumnName("formula", alias)), pg.Array(f.FormulaNEQ))
	}

	if !isZero(f.DescriptionMATCH) {
		q = q.Where(gopgutil.BuildConditionMatch(gopgutil.AddAliasToColumnName("description", alias)), f.DescriptionMATCH)
	}
	if !isZero(f.DescriptionIEQ) {
		q = q.Where(gopgutil.BuildConditionIEQ(gopgutil.AddAliasToColumnName("description", alias)), f.DescriptionIEQ)
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
		q = q.Where(gopgutil.BuildConditionEquals(gopgutil.AddAliasToColumnName("created_at", alias)), f.CreatedAt)
	}
	if !isZero(f.CreatedAtGT) {
		q = q.Where(gopgutil.BuildConditionGT(gopgutil.AddAliasToColumnName("created_at", alias)), f.CreatedAtGT)
	}
	if !isZero(f.CreatedAtGTE) {
		q = q.Where(gopgutil.BuildConditionGTE(gopgutil.AddAliasToColumnName("created_at", alias)), f.CreatedAtGTE)
	}
	if !isZero(f.CreatedAtLT) {
		q = q.Where(gopgutil.BuildConditionLT(gopgutil.AddAliasToColumnName("created_at", alias)), f.CreatedAtLT)
	}
	if !isZero(f.CreatedAtLTE) {
		q = q.Where(gopgutil.BuildConditionLTE(gopgutil.AddAliasToColumnName("created_at", alias)), f.CreatedAtLTE)
	}

	if f.Or != nil {
		q = f.Or.WhereWithAlias(q, alias)
	}

	return q, nil
}

func (f *QualificationFilter) Where(q *orm.Query) (*orm.Query, error) {
	return f.WhereWithAlias(q, "qualification")
}
