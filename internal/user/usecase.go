package user

import (
	"context"

	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
)

type Usecase interface {
	Store(ctx context.Context, input *models.UserInput) (*models.User, error)
	UpdateOne(ctx context.Context, id int, input *models.UserInput) (*models.User, error)
	UpdateMany(ctx context.Context, f *models.UserFilter, input *models.UserInput) ([]*models.User, error)
	Delete(ctx context.Context, f *models.UserFilter) ([]*models.User, error)
	Fetch(ctx context.Context, cfg *FetchConfig) ([]*models.User, int, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByCredentials(ctx context.Context, email, password string) (*models.User, error)
}
