package profession

import (
	"context"

	"github.com/zdam-egzamin-zawodowy/backend/internal/model"
)

type FetchConfig struct {
	Filter *model.ProfessionFilter
	Offset int
	Limit  int
	Sort   []string
	Count  bool
}

type Repository interface {
	Store(ctx context.Context, input *model.ProfessionInput) (*model.Profession, error)
	UpdateMany(ctx context.Context, f *model.ProfessionFilter, input *model.ProfessionInput) ([]*model.Profession, error)
	Delete(ctx context.Context, f *model.ProfessionFilter) ([]*model.Profession, error)
	Fetch(ctx context.Context, cfg *FetchConfig) ([]*model.Profession, int, error)
	GetAssociatedQualifications(ctx context.Context, ids ...int) (map[int][]*model.Qualification, error)
}
