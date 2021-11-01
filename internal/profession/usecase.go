package profession

import (
	"context"
	"github.com/zdam-egzamin-zawodowy/backend/internal"
)

type Usecase interface {
	Store(ctx context.Context, input *internal.ProfessionInput) (*internal.Profession, error)
	UpdateOneByID(ctx context.Context, id int, input *internal.ProfessionInput) (*internal.Profession, error)
	Delete(ctx context.Context, f *internal.ProfessionFilter) ([]*internal.Profession, error)
	Fetch(ctx context.Context, cfg *FetchConfig) ([]*internal.Profession, int, error)
	GetByID(ctx context.Context, id int) (*internal.Profession, error)
	GetBySlug(ctx context.Context, slug string) (*internal.Profession, error)
}
