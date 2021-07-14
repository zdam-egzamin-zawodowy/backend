package profession

import (
	"context"

	"github.com/zdam-egzamin-zawodowy/backend/internal/model"
)

type Usecase interface {
	Store(ctx context.Context, input *model.ProfessionInput) (*model.Profession, error)
	UpdateOneByID(ctx context.Context, id int, input *model.ProfessionInput) (*model.Profession, error)
	Delete(ctx context.Context, f *model.ProfessionFilter) ([]*model.Profession, error)
	Fetch(ctx context.Context, cfg *FetchConfig) ([]*model.Profession, int, error)
	GetByID(ctx context.Context, id int) (*model.Profession, error)
	GetBySlug(ctx context.Context, slug string) (*model.Profession, error)
}
