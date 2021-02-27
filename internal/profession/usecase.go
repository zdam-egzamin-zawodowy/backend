package profession

import (
	"context"

	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
)

type Usecase interface {
	Store(ctx context.Context, input *models.ProfessionInput) (*models.Profession, error)
	UpdateOne(ctx context.Context, id int, input *models.ProfessionInput) (*models.Profession, error)
	UpdateMany(ctx context.Context, f *models.ProfessionFilter, input *models.ProfessionInput) ([]*models.Profession, error)
	Delete(ctx context.Context, f *models.ProfessionFilter) ([]*models.Profession, error)
	Fetch(ctx context.Context, cfg *FetchConfig) ([]*models.Profession, int, error)
	GetByID(ctx context.Context, id int) (*models.Profession, error)
	GetBySlug(ctx context.Context, slug string) (*models.Profession, error)
}
