package qualification

import (
	"context"

	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
)

type FetchConfig struct {
	Filter *models.QualificationFilter
	Offset int
	Limit  int
	Sort   []string
	Count  bool
}

type Repository interface {
	Store(ctx context.Context, input *models.QualificationInput) (*models.Qualification, error)
	UpdateMany(ctx context.Context, f *models.QualificationFilter, input *models.QualificationInput) ([]*models.Qualification, error)
	Delete(ctx context.Context, f *models.QualificationFilter) ([]*models.Qualification, error)
	Fetch(ctx context.Context, cfg *FetchConfig) ([]*models.Qualification, int, error)
}
