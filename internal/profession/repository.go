package profession

import (
	"context"
	"github.com/zdam-egzamin-zawodowy/backend/internal"
)

type FetchConfig struct {
	Filter *internal.ProfessionFilter
	Offset int
	Limit  int
	Sort   []string
	Count  bool
}

type Repository interface {
	Store(ctx context.Context, input *internal.ProfessionInput) (*internal.Profession, error)
	UpdateMany(ctx context.Context, f *internal.ProfessionFilter, input *internal.ProfessionInput) ([]*internal.Profession, error)
	Delete(ctx context.Context, f *internal.ProfessionFilter) ([]*internal.Profession, error)
	Fetch(ctx context.Context, cfg *FetchConfig) ([]*internal.Profession, int, error)
	GetAssociatedQualifications(ctx context.Context, ids ...int) (map[int][]*internal.Qualification, error)
}
