package question

import (
	"context"
	"github.com/zdam-egzamin-zawodowy/backend/internal"
)

type FetchConfig struct {
	Filter *internal.QuestionFilter
	Offset int
	Limit  int
	Sort   []string
	Count  bool
}

type GenerateTestConfig struct {
	Qualifications []int
	Limit          int
}

type Repository interface {
	Store(ctx context.Context, input *internal.QuestionInput) (*internal.Question, error)
	UpdateOneByID(ctx context.Context, id int, input *internal.QuestionInput) (*internal.Question, error)
	Delete(ctx context.Context, f *internal.QuestionFilter) ([]*internal.Question, error)
	Fetch(ctx context.Context, cfg *FetchConfig) ([]*internal.Question, int, error)
	GenerateTest(ctx context.Context, cfg *GenerateTestConfig) ([]*internal.Question, error)
}
