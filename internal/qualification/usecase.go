package qualification

import (
	"context"

	"github.com/zdam-egzamin-zawodowy/backend/internal/model"
)

type Usecase interface {
	Store(ctx context.Context, input *model.QualificationInput) (*model.Qualification, error)
	UpdateOneByID(ctx context.Context, id int, input *model.QualificationInput) (*model.Qualification, error)
	Delete(ctx context.Context, f *model.QualificationFilter) ([]*model.Qualification, error)
	Fetch(ctx context.Context, cfg *FetchConfig) ([]*model.Qualification, int, error)
	GetByID(ctx context.Context, id int) (*model.Qualification, error)
	GetBySlug(ctx context.Context, slug string) (*model.Qualification, error)
	GetSimilar(ctx context.Context, cfg *GetSimilarConfig) ([]*model.Qualification, int, error)
}
