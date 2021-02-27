package models

import (
	"context"
	"time"
)

type Question struct {
	ID              int            `json:"id" xml:"id" gqlgen:"id"`
	From            string         `pg:",unique:group_1" json:"from" xml:"from" gqlgen:"from"`
	Content         string         `pg:",unique:group_1,notnull" json:"content" xml:"content" gqlgen:"content"`
	Explanation     string         `json:"explanation" xml:"explanation" gqlgen:"explanation"`
	CorrectAnswer   Answer         `pg:",unique:group_1,notnull" json:"correctAnswer" xml:"correctAnswer" gqlgen:"correctAnswer"`
	Image           string         `json:"image" xml:"image" gqlgen:"image"`
	AnswerA         Answer         `json:"answerA" xml:"answerA" gqlgen:"answerA"`
	AnswerAImage    string         `json:"answerAImage" xml:"answerAImage" gqlgen:"answerAImage"`
	AnswerB         Answer         `json:"answerB" xml:"answerB" gqlgen:"answerB"`
	AnswerBImage    string         `json:"answerBImage" xml:"answerBImage" gqlgen:"answerBImage"`
	AnswerC         Answer         `json:"answerC" xml:"answerC" gqlgen:"answerC"`
	AnswerCImage    string         `json:"answerCImage" xml:"answerCImage" gqlgen:"answerCImage"`
	AnswerD         Answer         `json:"answerD" xml:"answerD" gqlgen:"answerD"`
	AnswerDImage    string         `json:"answerDImage" xml:"answerDImage" gqlgen:"answerDImage"`
	QualificationID int            `pg:",unique:group_1,notnull" json:"qualificationID" xml:"qualificationID" gqlgen:"qualificationID"`
	Qualification   *Qualification `pg:"rel:has-one" json:"qualification" xml:"qualification" gqlgen:"qualification"`
	CreatedAt       time.Time      `json:"createdAt,omitempty" pg:"default:now()" xml:"createdAt" gqlgen:"createdAt"`
}

func (q *Question) BeforeInsert(ctx context.Context) (context.Context, error) {
	q.CreatedAt = time.Now()

	return ctx, nil
}
