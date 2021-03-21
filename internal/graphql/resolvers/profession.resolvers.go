package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/zdam-egzamin-zawodowy/backend/internal/gin/middleware"

	"github.com/zdam-egzamin-zawodowy/backend/internal/graphql/generated"
	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
	"github.com/zdam-egzamin-zawodowy/backend/internal/profession"
	"github.com/zdam-egzamin-zawodowy/backend/pkg/utils"
)

func (r *mutationResolver) CreateProfession(ctx context.Context, input models.ProfessionInput) (*models.Profession, error) {
	return r.ProfessionUsecase.Store(ctx, &input)
}

func (r *mutationResolver) UpdateProfession(ctx context.Context, id int, input models.ProfessionInput) (*models.Profession, error) {
	return r.ProfessionUsecase.UpdateOneByID(ctx, id, &input)
}

func (r *mutationResolver) DeleteProfessions(ctx context.Context, ids []int) ([]*models.Profession, error) {
	return r.ProfessionUsecase.Delete(ctx, &models.ProfessionFilter{
		ID: ids,
	})
}

func (r *queryResolver) Professions(
	ctx context.Context,
	filter *models.ProfessionFilter,
	limit *int,
	offset *int,
	sort []string,
) (*generated.ProfessionList, error) {
	var err error
	list := &generated.ProfessionList{}
	list.Items, list.Total, err = r.ProfessionUsecase.Fetch(
		ctx,
		&profession.FetchConfig{
			Count:  shouldCount(ctx),
			Filter: filter,
			Limit:  utils.SafeIntPointer(limit, profession.DefaultLimit),
			Offset: utils.SafeIntPointer(offset, 0),
			Sort:   sort,
		},
	)
	return list, err
}

func (r *queryResolver) Profession(ctx context.Context, id *int, slug *string) (*models.Profession, error) {
	if id != nil {
		return r.ProfessionUsecase.GetByID(ctx, *id)
	} else if slug != nil {
		return r.ProfessionUsecase.GetBySlug(ctx, *slug)
	}

	return nil, nil
}

func (r *professionResolver) Qualifications(
	ctx context.Context,
	obj *models.Profession,
) ([]*models.Qualification, error) {
	if obj != nil {
		if dataloader, err := middleware.DataLoaderFromContext(ctx); err == nil && dataloader != nil {
			return dataloader.QualificationsByProfessionID.Load(obj.ID)
		}
	}
	return []*models.Qualification{}, nil
}
