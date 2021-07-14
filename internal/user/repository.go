package user

import (
	"context"

	"github.com/zdam-egzamin-zawodowy/backend/internal/model"
)

type FetchConfig struct {
	Filter *model.UserFilter
	Offset int
	Limit  int
	Sort   []string
	Count  bool
}

type Repository interface {
	Store(ctx context.Context, input *model.UserInput) (*model.User, error)
	UpdateMany(ctx context.Context, f *model.UserFilter, input *model.UserInput) ([]*model.User, error)
	Delete(ctx context.Context, f *model.UserFilter) ([]*model.User, error)
	Fetch(ctx context.Context, cfg *FetchConfig) ([]*model.User, int, error)
}
