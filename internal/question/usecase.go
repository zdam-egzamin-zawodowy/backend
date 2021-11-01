package question

import (
	"context"
	"github.com/zdam-egzamin-zawodowy/backend/internal"
)

type Usecase interface {
	Store(ctx context.Context, input *internal.QuestionInput) (*internal.Question, error)
	UpdateOneByID(ctx context.Context, id int, input *internal.QuestionInput) (*internal.Question, error)
	Delete(ctx context.Context, f *internal.QuestionFilter) ([]*internal.Question, error)
	Fetch(ctx context.Context, cfg *FetchConfig) ([]*internal.Question, int, error)
	GetByID(ctx context.Context, id int) (*internal.Question, error)
	GenerateTest(ctx context.Context, cfg *GenerateTestConfig) ([]*internal.Question, error)
}
