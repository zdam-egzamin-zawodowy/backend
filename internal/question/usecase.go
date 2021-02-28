package question

import (
	"context"

	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
)

type Usecase interface {
	Store(ctx context.Context, input *models.QuestionInput) (*models.Question, error)
	UpdateOne(ctx context.Context, id int, input *models.QuestionInput) (*models.Question, error)
	Delete(ctx context.Context, f *models.QuestionFilter) ([]*models.Question, error)
	Fetch(ctx context.Context, cfg *FetchConfig) ([]*models.Question, int, error)
	GetByID(ctx context.Context, id int) (*models.Question, error)
}
