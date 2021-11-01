package internal

import (
	"context"
	"github.com/Kichiyaki/gopgutil/v10"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

var _ pg.BeforeInsertHook = (*Question)(nil)

type Question struct {
	tableName struct{} `pg:"alias:question"`

	ID              int            `json:"id" xml:"id" gqlgen:"id"`
	From            string         `pg:",unique:group_1" json:"from" xml:"from" gqlgen:"from"`
	Content         string         `pg:",unique:group_1,notnull" json:"content" xml:"content" gqlgen:"content"`
	Explanation     string         `json:"explanation" xml:"explanation" gqlgen:"explanation"`
	CorrectAnswer   Answer         `pg:",unique:group_1,notnull" json:"correctAnswer" xml:"correctAnswer" gqlgen:"correctAnswer"`
	Image           string         `json:"image" xml:"image" gqlgen:"image"`
	AnswerA         string         `pg:"answer_a" json:"answerA" xml:"answerA" gqlgen:"answerA"`
	AnswerAImage    string         `pg:"answer_a_image" json:"answerAImage" xml:"answerAImage" gqlgen:"answerAImage"`
	AnswerB         string         `pg:"answer_b" json:"answerB" xml:"answerB" gqlgen:"answerB"`
	AnswerBImage    string         `pg:"answer_b_image" json:"answerBImage" xml:"answerBImage" gqlgen:"answerBImage"`
	AnswerC         string         `pg:"answer_c" json:"answerC" xml:"answerC" gqlgen:"answerC"`
	AnswerCImage    string         `pg:"answer_c_image" json:"answerCImage" xml:"answerCImage" gqlgen:"answerCImage"`
	AnswerD         string         `pg:"answer_d" json:"answerD" xml:"answerD" gqlgen:"answerD"`
	AnswerDImage    string         `pg:"answer_d_image" json:"answerDImage" xml:"answerDImage" gqlgen:"answerDImage"`
	QualificationID int            `pg:",unique:group_1,on_delete:CASCADE" json:"qualificationID" xml:"qualificationID" gqlgen:"qualificationID"`
	Qualification   *Qualification `pg:"rel:has-one" json:"qualification" xml:"qualification" gqlgen:"qualification"`
	CreatedAt       time.Time      `json:"createdAt,omitempty" pg:"default:now()" xml:"createdAt" gqlgen:"createdAt"`
	UpdatedAt       time.Time      `pg:"default:now()" json:"updatedAt" xml:"updatedAt" gqlgen:"updatedAt"`
}

func (q *Question) BeforeInsert(ctx context.Context) (context.Context, error) {
	q.CreatedAt = time.Now()
	q.UpdatedAt = time.Now()

	return ctx, nil
}

type QuestionInput struct {
	Content            *string         `json:"content" xml:"content" gqlgen:"content"`
	From               *string         `json:"from" xml:"from" gqlgen:"from"`
	Explanation        *string         `json:"explanation" xml:"explanation" gqlgen:"explanation"`
	CorrectAnswer      *Answer         `json:"correctAnswer" xml:"correctAnswer" gqlgen:"correctAnswer"`
	AnswerA            *string         `gqlgen:"answerA" json:"answerA" xml:"answerA"`
	AnswerB            *string         `gqlgen:"answerB" json:"answerB" xml:"answerB"`
	AnswerC            *string         `gqlgen:"answerC" json:"answerC" xml:"answerC"`
	AnswerD            *string         `gqlgen:"answerD" json:"answerD" xml:"answerD"`
	QualificationID    *int            `gqlgen:"qualificationID" json:"qualificationID" xml:"qualificationID"`
	Image              *graphql.Upload `json:"image" xml:"image" gqlgen:"image"`
	DeleteImage        *bool           `json:"deleteImage" xml:"deleteImage" gqlgen:"deleteImage"`
	AnswerAImage       *graphql.Upload `json:"answerAImage" gqlgen:"answerAImage" xml:"answerAImage"`
	DeleteAnswerAImage *bool           `json:"deleteAnswerAImage" xml:"deleteAnswerAImage" gqlgen:"deleteAnswerAImage"`
	AnswerBImage       *graphql.Upload `json:"answerBImage" gqlgen:"answerBImage" xml:"answerBImage"`
	DeleteAnswerBImage *bool           `json:"deleteAnswerBImage" xml:"deleteAnswerBImage" gqlgen:"deleteAnswerBImage"`
	AnswerCImage       *graphql.Upload `json:"answerCImage" gqlgen:"answerCImage" xml:"answerCImage"`
	DeleteAnswerCImage *bool           `json:"deleteAnswerCImage" xml:"deleteAnswerCImage" gqlgen:"deleteAnswerCImage"`
	AnswerDImage       *graphql.Upload `json:"answerDImage" gqlgen:"answerDImage" xml:"answerDImage"`
	DeleteAnswerDImage *bool           `json:"deleteAnswerDImage" xml:"deleteAnswerDImage" gqlgen:"deleteAnswerDImage"`
}

func (input *QuestionInput) IsEmpty() bool {
	return input == nil &&
		input.Content == nil &&
		input.From == nil &&
		input.Explanation == nil &&
		input.CorrectAnswer == nil &&
		input.AnswerA == nil &&
		input.AnswerAImage == nil &&
		input.AnswerB == nil &&
		input.AnswerBImage == nil &&
		input.AnswerC == nil &&
		input.AnswerCImage == nil &&
		input.AnswerD == nil &&
		input.AnswerDImage == nil &&
		input.Image == nil &&
		input.QualificationID == nil
}

func (input *QuestionInput) HasBasicDataToUpdate() bool {
	return input != nil &&
		(input.Content != nil ||
			input.Explanation != nil ||
			input.CorrectAnswer != nil ||
			input.AnswerA != nil ||
			input.AnswerB != nil ||
			input.AnswerC != nil ||
			input.AnswerD != nil ||
			input.QualificationID != nil)
}

func (input *QuestionInput) Sanitize() *QuestionInput {
	if input.Content != nil {
		*input.Content = strings.TrimSpace(*input.Content)
	}
	if input.From != nil {
		*input.From = strings.TrimSpace(*input.From)
	}
	if input.Explanation != nil {
		*input.Explanation = strings.TrimSpace(*input.Explanation)
	}
	if input.AnswerA != nil {
		*input.AnswerA = strings.TrimSpace(*input.AnswerA)
	}
	if input.AnswerB != nil {
		*input.AnswerB = strings.TrimSpace(*input.AnswerB)
	}
	if input.AnswerC != nil {
		*input.AnswerC = strings.TrimSpace(*input.AnswerC)
	}
	if input.AnswerD != nil {
		*input.AnswerD = strings.TrimSpace(*input.AnswerD)
	}

	return input
}

func (input *QuestionInput) ToQuestion() *Question {
	q := &Question{}
	if input.Content != nil {
		q.Content = *input.Content
	}
	if input.From != nil {
		q.From = *input.From
	}
	if input.Explanation != nil {
		q.Explanation = *input.Explanation
	}
	if input.CorrectAnswer != nil {
		q.CorrectAnswer = *input.CorrectAnswer
	}
	if input.AnswerA != nil {
		q.AnswerA = *input.AnswerA
	}
	if input.AnswerB != nil {
		q.AnswerB = *input.AnswerB
	}
	if input.AnswerC != nil {
		q.AnswerC = *input.AnswerC
	}
	if input.AnswerD != nil {
		q.AnswerD = *input.AnswerD
	}
	if input.QualificationID != nil {
		q.QualificationID = *input.QualificationID
	}
	return q
}

func (input *QuestionInput) ApplyUpdate(q *orm.Query) (*orm.Query, error) {
	if !input.IsEmpty() {
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
	}

	return q, nil
}

type QuestionFilter struct {
	ID    []int `gqlgen:"id" json:"id" xml:"id"`
	IDNEQ []int `gqlgen:"idNEQ" json:"idNEQ" xml:"idNEQ"`

	From []string `json:"from" xml:"from" gqlgen:"from"`

	ContentMATCH string `json:"contentMATCH" xml:"contentMATCH" gqlgen:"contentMATCH"`
	ContentIEQ   string `json:"contentIEQ" xml:"contentIEQ" gqlgen:"contentIEQ"`

	QualificationID     []int               `json:"qualificationID" xml:"qualificationID" gqlgen:"qualificationID"`
	QualificationIDNEQ  []int               `json:"qualificationIDNEQ" xml:"qualificationIDNEQ" gqlgen:"qualificationIDNEQ"`
	QualificationFilter QualificationFilter `json:"qualificationFilter" xml:"qualificationFilter" gqlgen:"qualificationFilter"`

	CreatedAt    time.Time `gqlgen:"createdAt" json:"createdAt" xml:"createdAt"`
	CreatedAtGT  time.Time `gqlgen:"createdAtGT" json:"createdAtGT" xml:"createdAtGT"`
	CreatedAtGTE time.Time `json:"createdAtGTE" xml:"createdAtGTE" gqlgen:"createdAtGTE"`
	CreatedAtLT  time.Time `gqlgen:"createdAtLT" json:"createdAtLT" xml:"createdAtLT"`
	CreatedAtLTE time.Time `json:"createdAtLTE" xml:"createdAtLTE" gqlgen:"createdAtLTE"`
}

func (f *QuestionFilter) WhereWithAlias(q *orm.Query, alias string) (*orm.Query, error) {
	if f == nil {
		return q, nil
	}

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

	var err error
	q, err = f.QualificationFilter.WhereWithAlias(q.Relation("Qualification._"), "qualification")
	if err != nil {
		return q, err
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

func (f *QuestionFilter) Where(q *orm.Query) (*orm.Query, error) {
	return f.WhereWithAlias(q, "question")
}
