package user

import (
	"context"

	"github.com/zdam-egzamin-zawodowy/backend/internal/model"
)

type Usecase interface {
	Store(ctx context.Context, input *model.UserInput) (*model.User, error)
	UpdateOneByID(ctx context.Context, id int, input *model.UserInput) (*model.User, error)
	UpdateMany(ctx context.Context, f *model.UserFilter, input *model.UserInput) ([]*model.User, error)
	Delete(ctx context.Context, f *model.UserFilter) ([]*model.User, error)
	Fetch(ctx context.Context, cfg *FetchConfig) ([]*model.User, int, error)
	GetByID(ctx context.Context, id int) (*model.User, error)
	GetByCredentials(ctx context.Context, email, password string) (*model.User, error)
}
