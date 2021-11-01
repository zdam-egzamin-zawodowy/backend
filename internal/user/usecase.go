package user

import (
	"context"
	"github.com/zdam-egzamin-zawodowy/backend/internal"
)

type Usecase interface {
	Store(ctx context.Context, input *internal.UserInput) (*internal.User, error)
	UpdateOneByID(ctx context.Context, id int, input *internal.UserInput) (*internal.User, error)
	UpdateMany(ctx context.Context, f *internal.UserFilter, input *internal.UserInput) ([]*internal.User, error)
	Delete(ctx context.Context, f *internal.UserFilter) ([]*internal.User, error)
	Fetch(ctx context.Context, cfg *FetchConfig) ([]*internal.User, int, error)
	GetByID(ctx context.Context, id int) (*internal.User, error)
	GetByCredentials(ctx context.Context, email, password string) (*internal.User, error)
}
