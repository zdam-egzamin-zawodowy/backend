package qualification

import (
	"context"

	"github.com/zdam-egzamin-zawodowy/backend/internal/model"
)

type FetchConfig struct {
	Filter *model.QualificationFilter
	Offset int
	Limit  int
	Sort   []string
	Count  bool
}

type GetSimilarConfig struct {
	Limit           int
	Offset          int
	QualificationID int
	Sort            []string
	Count           bool
}

type Repository interface {
	Store(ctx context.Context, input *model.QualificationInput) (*model.Qualification, error)
	UpdateMany(ctx context.Context, f *model.QualificationFilter, input *model.QualificationInput) ([]*model.Qualification, error)
	Delete(ctx context.Context, f *model.QualificationFilter) ([]*model.Qualification, error)
	Fetch(ctx context.Context, cfg *FetchConfig) ([]*model.Qualification, int, error)
	GetSimilar(ctx context.Context, cfg *GetSimilarConfig) ([]*model.Qualification, int, error)
}
