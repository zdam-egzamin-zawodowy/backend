package question

import (
	"context"

	"github.com/zdam-egzamin-zawodowy/backend/internal/model"
)

type FetchConfig struct {
	Filter *model.QuestionFilter
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
	Store(ctx context.Context, input *model.QuestionInput) (*model.Question, error)
	UpdateOneByID(ctx context.Context, id int, input *model.QuestionInput) (*model.Question, error)
	Delete(ctx context.Context, f *model.QuestionFilter) ([]*model.Question, error)
	Fetch(ctx context.Context, cfg *FetchConfig) ([]*model.Question, int, error)
	GenerateTest(ctx context.Context, cfg *GenerateTestConfig) ([]*model.Question, error)
}
