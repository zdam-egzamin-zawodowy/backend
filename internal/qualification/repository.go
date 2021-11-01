package qualification

import (
	"context"
	"github.com/zdam-egzamin-zawodowy/backend/internal"
)

type FetchConfig struct {
	Filter *internal.QualificationFilter
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
	Store(ctx context.Context, input *internal.QualificationInput) (*internal.Qualification, error)
	UpdateMany(ctx context.Context, f *internal.QualificationFilter, input *internal.QualificationInput) ([]*internal.Qualification, error)
	Delete(ctx context.Context, f *internal.QualificationFilter) ([]*internal.Qualification, error)
	Fetch(ctx context.Context, cfg *FetchConfig) ([]*internal.Qualification, int, error)
	GetSimilar(ctx context.Context, cfg *GetSimilarConfig) ([]*internal.Qualification, int, error)
}
