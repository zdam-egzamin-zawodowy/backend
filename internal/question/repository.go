package question

import (
	"context"

	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
)

type FetchConfig struct {
	Filter *models.QuestionFilter
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
	Store(ctx context.Context, input *models.QuestionInput) (*models.Question, error)
	FindByIDAndUpdate(ctx context.Context, id int, input *models.QuestionInput) (*models.Question, error)
	Delete(ctx context.Context, f *models.QuestionFilter) ([]*models.Question, error)
	Fetch(ctx context.Context, cfg *FetchConfig) ([]*models.Question, int, error)
	GenerateTest(ctx context.Context, cfg *GenerateTestConfig) ([]*models.Question, error)
}
