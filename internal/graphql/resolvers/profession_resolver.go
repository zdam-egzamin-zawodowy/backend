package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/Kichiyaki/goutil/safeptr"

	"github.com/zdam-egzamin-zawodowy/backend/internal/chi/middleware"

	"github.com/zdam-egzamin-zawodowy/backend/internal/graphql/generated"
	"github.com/zdam-egzamin-zawodowy/backend/internal/model"
	"github.com/zdam-egzamin-zawodowy/backend/internal/profession"
)

func (r *mutationResolver) CreateProfession(ctx context.Context, input model.ProfessionInput) (*model.Profession, error) {
	return r.ProfessionUsecase.Store(ctx, &input)
}

func (r *mutationResolver) UpdateProfession(ctx context.Context, id int, input model.ProfessionInput) (*model.Profession, error) {
	return r.ProfessionUsecase.UpdateOneByID(ctx, id, &input)
}

func (r *mutationResolver) DeleteProfessions(ctx context.Context, ids []int) ([]*model.Profession, error) {
	return r.ProfessionUsecase.Delete(ctx, &model.ProfessionFilter{
		ID: ids,
	})
}

func (r *queryResolver) Professions(
	ctx context.Context,
	filter *model.ProfessionFilter,
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
			Limit:  safeptr.SafeIntPointer(limit, profession.FetchDefaultLimit),
			Offset: safeptr.SafeIntPointer(offset, 0),
			Sort:   sort,
		},
	)
	return list, err
}

func (r *queryResolver) Profession(ctx context.Context, id *int, slug *string) (*model.Profession, error) {
	if id != nil {
		return r.ProfessionUsecase.GetByID(ctx, *id)
	} else if slug != nil {
		return r.ProfessionUsecase.GetBySlug(ctx, *slug)
	}

	return nil, nil
}

func (r *professionResolver) Qualifications(
	ctx context.Context,
	obj *model.Profession,
) ([]*model.Qualification, error) {
	if obj != nil {
		if dataloader, err := middleware.DataLoaderFromContext(ctx); err == nil && dataloader != nil {
			return dataloader.QualificationsByProfessionID.Load(obj.ID)
		}
	}
	return []*model.Qualification{}, nil
}
