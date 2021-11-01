package postgres

import (
	"context"
	"github.com/Kichiyaki/gopgutil/v10"
	"github.com/zdam-egzamin-zawodowy/backend/internal"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

var _ pg.BeforeInsertHook = (*Question)(nil)

type Question struct {
	tableName struct{} `pg:"alias:question"`

	ID              int
	From            string `pg:",unique:group_1"`
	Content         string `pg:",unique:group_1,notnull"`
	Explanation     string
	CorrectAnswer   Answer `pg:",unique:group_1,notnull"`
	Image           string
	AnswerA         string         `pg:"answer_a"`
	AnswerAImage    string         `pg:"answer_a_image"`
	AnswerB         string         `pg:"answer_b" `
	AnswerBImage    string         `pg:"answer_b_image"`
	AnswerC         string         `pg:"answer_c" `
	AnswerCImage    string         `pg:"answer_c_image"`
	AnswerD         string         `pg:"answer_d" `
	AnswerDImage    string         `pg:"answer_d_image"`
	QualificationID int            `pg:",unique:group_1,on_delete:CASCADE"`
	Qualification   *Qualification `pg:"rel:has-one"`
	CreatedAt       time.Time      `pg:"default:now()"`
	UpdatedAt       time.Time      `pg:"default:now()"`
}

func (q *Question) BeforeInsert(ctx context.Context) (context.Context, error) {
	q.CreatedAt = time.Now()
	q.UpdatedAt = time.Now()

	return ctx, nil
}

func applyQuestionInputUpdates(input internal.QuestionInput) func(q *orm.Query) (*orm.Query, error) {
	return func(q *orm.Query) (*orm.Query, error) {
		if input.Content != nil {
			q = q.Set(
				gopgutil.BuildConditionEquals("content"),
				*input.Content,
			)
		}

		if input.From != nil {
			q = q.Set(
				gopgutil.BuildConditionEquals("?"),
				pg.Ident("from"),
				*input.From,
			)
		}

		if input.Explanation != nil {
			q = q.Set(gopgutil.BuildConditionEquals("explanation"), *input.Explanation)
		}

		if input.CorrectAnswer != nil {
			q = q.Set(gopgutil.BuildConditionEquals("correct_answer"), *input.CorrectAnswer)
		}

		if input.AnswerA != nil {
			q = q.Set(gopgutil.BuildConditionEquals("answer_a"), *input.AnswerA)
		}

		if input.AnswerB != nil {
			q = q.Set(gopgutil.BuildConditionEquals("answer_b"), *input.AnswerB)
		}

		if input.AnswerC != nil {
			q = q.Set(gopgutil.BuildConditionEquals("answer_c"), *input.AnswerC)
		}

		if input.AnswerD != nil {
			q = q.Set(gopgutil.BuildConditionEquals("answer_d"), *input.AnswerD)
		}

		if input.QualificationID != nil {
			q = q.Set(gopgutil.BuildConditionEquals("qualification_id"), *input.QualificationID)
		}

		return q, nil
	}
}

func applyQuestionFilter(f internal.QuestionFilter, alias string) func(q *orm.Query) (*orm.Query, error) {
	return func(q *orm.Query) (*orm.Query, error) {
		if !isZero(f.ID) {
			q = q.Where(gopgutil.BuildConditionArray("?"), gopgutil.AddAliasToColumnName("id", alias), pg.Array(f.ID))
		}
		if !isZero(f.IDNEQ) {
			q = q.Where(gopgutil.BuildConditionNotInArray("?"), gopgutil.AddAliasToColumnName("id", alias), pg.Array(f.IDNEQ))
		}

		if !isZero(f.From) {
			q = q.Where(gopgutil.BuildConditionArray("?"), gopgutil.AddAliasToColumnName("from", alias), pg.Array(f.From))
		}

		if !isZero(f.ContentMATCH) {
			q = q.Where(gopgutil.BuildConditionMatch("?"), gopgutil.AddAliasToColumnName("content", alias), f.ContentMATCH)
		}
		if !isZero(f.ContentIEQ) {
			q = q.Where(gopgutil.BuildConditionIEQ("?"), gopgutil.AddAliasToColumnName("content", alias), f.ContentIEQ)
		}

		if !isZero(f.QualificationID) {
			q = q.Where(gopgutil.BuildConditionArray("?"), gopgutil.AddAliasToColumnName("qualification_id", alias), pg.Array(f.QualificationID))
		}
		if !isZero(f.QualificationIDNEQ) {
			q = q.Where(gopgutil.BuildConditionNotInArray("?"), gopgutil.AddAliasToColumnName("qualification_id", alias), pg.Array(f.QualificationIDNEQ))
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

		q = q.Apply(applyQualificationFilter(f.QualificationFilter, "qualification"))

		return q, nil
	}
}
