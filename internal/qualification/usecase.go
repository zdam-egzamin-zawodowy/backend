package qualification

import (
	"context"
	"github.com/zdam-egzamin-zawodowy/backend/internal"
)

type Usecase interface {
	Store(ctx context.Context, input *internal.QualificationInput) (*internal.Qualification, error)
	UpdateOneByID(ctx context.Context, id int, input *internal.QualificationInput) (*internal.Qualification, error)
	Delete(ctx context.Context, f *internal.QualificationFilter) ([]*internal.Qualification, error)
	Fetch(ctx context.Context, cfg *FetchConfig) ([]*internal.Qualification, int, error)
	GetByID(ctx context.Context, id int) (*internal.Qualification, error)
	GetBySlug(ctx context.Context, slug string) (*internal.Qualification, error)
	GetSimilar(ctx context.Context, cfg *GetSimilarConfig) ([]*internal.Qualification, int, error)
}
