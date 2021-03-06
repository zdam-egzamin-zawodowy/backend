package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/zdam-egzamin-zawodowy/backend/internal/graphql/generated"
	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
	"github.com/zdam-egzamin-zawodowy/backend/internal/qualification"
	"github.com/zdam-egzamin-zawodowy/backend/pkg/utils"
)

func (r *mutationResolver) CreateQualification(ctx context.Context, input models.QualificationInput) (*models.Qualification, error) {
	return r.QualificationUsecase.Store(ctx, &input)
}

func (r *mutationResolver) UpdateQualification(ctx context.Context, id int, input models.QualificationInput) (*models.Qualification, error) {
	return r.QualificationUsecase.UpdateOneByID(ctx, id, &input)
}

func (r *mutationResolver) DeleteQualifications(ctx context.Context, ids []int) ([]*models.Qualification, error) {
	return r.QualificationUsecase.Delete(ctx, &models.QualificationFilter{
		ID: ids,
	})
}

func (r *queryResolver) Qualifications(
	ctx context.Context,
	filter *models.QualificationFilter,
	limit *int,
	offset *int,
	sort []string,
) (*generated.QualificationList, error) {
	var err error
	list := &generated.QualificationList{}
	list.Items, list.Total, err = r.QualificationUsecase.Fetch(
		ctx,
		&qualification.FetchConfig{
			Count:  shouldCount(ctx),
			Filter: filter,
			Limit:  utils.SafeIntPointer(limit, qualification.DefaultLimit),
			Offset: utils.SafeIntPointer(offset, 0),
			Sort:   sort,
		},
	)
	return list, err
}

func (r *queryResolver) Qualification(ctx context.Context, id *int, slug *string) (*models.Qualification, error) {
	if id != nil {
		return r.QualificationUsecase.GetByID(ctx, *id)
	} else if slug != nil {
		return r.QualificationUsecase.GetBySlug(ctx, *slug)
	}

	return nil, nil
}
