package qualification

import (
	"context"

	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
)

type Usecase interface {
	Store(ctx context.Context, input *models.QualificationInput) (*models.Qualification, error)
	UpdateOne(ctx context.Context, id int, input *models.QualificationInput) (*models.Qualification, error)
	Delete(ctx context.Context, f *models.QualificationFilter) ([]*models.Qualification, error)
	Fetch(ctx context.Context, cfg *FetchConfig) ([]*models.Qualification, int, error)
	GetByID(ctx context.Context, id int) (*models.Qualification, error)
	GetBySlug(ctx context.Context, slug string) (*models.Qualification, error)
}
