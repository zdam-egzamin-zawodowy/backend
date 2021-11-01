package user

import (
	"context"
	"github.com/zdam-egzamin-zawodowy/backend/internal"
)

type FetchConfig struct {
	Filter *internal.UserFilter
	Offset int
	Limit  int
	Sort   []string
	Count  bool
}

type Repository interface {
	Store(ctx context.Context, input *internal.UserInput) (*internal.User, error)
	UpdateMany(ctx context.Context, f *internal.UserFilter, input *internal.UserInput) ([]*internal.User, error)
	Delete(ctx context.Context, f *internal.UserFilter) ([]*internal.User, error)
	Fetch(ctx context.Context, cfg *FetchConfig) ([]*internal.User, int, error)
}
