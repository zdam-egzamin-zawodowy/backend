package profession

import (
	"context"

	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
)

type FetchConfig struct {
	Filter *models.ProfessionFilter
	Offset int
	Limit  int
	Sort   []string
	Count  bool
}

type Repository interface {
	Store(ctx context.Context, input *models.ProfessionInput) (*models.Profession, error)
	UpdateMany(ctx context.Context, f *models.ProfessionFilter, input *models.ProfessionInput) ([]*models.Profession, error)
	Delete(ctx context.Context, f *models.ProfessionFilter) ([]*models.Profession, error)
	Fetch(ctx context.Context, cfg *FetchConfig) ([]*models.Profession, int, error)
	GetAssociatedQualifications(ctx context.Context, ids ...int) (map[int][]*models.Qualification, error)
}
