package user

import (
	"context"

	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
)

type FetchConfig struct {
	Filter *models.UserFilter
	Offset int
	Limit  int
	Sort   []string
	Count  bool
}

type Repository interface {
	Store(ctx context.Context, input *models.UserInput) (*models.User, error)
	UpdateMany(ctx context.Context, f *models.UserFilter, input *models.UserInput) ([]*models.User, error)
	Delete(ctx context.Context, f *models.UserFilter) ([]*models.User, error)
	Fetch(ctx context.Context, cfg *FetchConfig) ([]*models.User, int, error)
}
