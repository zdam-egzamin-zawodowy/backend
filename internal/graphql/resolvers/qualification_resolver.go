package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/Kichiyaki/goutil/safeptr"

	"github.com/zdam-egzamin-zawodowy/backend/internal/graphql/generated"
	"github.com/zdam-egzamin-zawodowy/backend/internal/model"
	"github.com/zdam-egzamin-zawodowy/backend/internal/qualification"
)

func (r *mutationResolver) CreateQualification(ctx context.Context, input model.QualificationInput) (*model.Qualification, error) {
	return r.QualificationUsecase.Store(ctx, &input)
}

func (r *mutationResolver) UpdateQualification(ctx context.Context, id int, input model.QualificationInput) (*model.Qualification, error) {
	return r.QualificationUsecase.UpdateOneByID(ctx, id, &input)
}

func (r *mutationResolver) DeleteQualifications(ctx context.Context, ids []int) ([]*model.Qualification, error) {
	return r.QualificationUsecase.Delete(ctx, &model.QualificationFilter{
		ID: ids,
	})
}

func (r *queryResolver) Qualifications(
	ctx context.Context,
	filter *model.QualificationFilter,
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
			Limit:  safeptr.SafeIntPointer(limit, qualification.FetchDefaultLimit),
			Offset: safeptr.SafeIntPointer(offset, 0),
			Sort:   sort,
		},
	)
	return list, err
}

func (r *queryResolver) SimilarQualifications(
	ctx context.Context,
	qualificationID int,
	limit *int,
	offset *int,
	sort []string,
) (*generated.QualificationList, error) {
	var err error
	list := &generated.QualificationList{}
	list.Items, list.Total, err = r.QualificationUsecase.GetSimilar(
		ctx,
		&qualification.GetSimilarConfig{
			Count:           shouldCount(ctx),
			QualificationID: qualificationID,
			Limit:           safeptr.SafeIntPointer(limit, qualification.FetchDefaultLimit),
			Offset:          safeptr.SafeIntPointer(offset, 0),
			Sort:            sort,
		},
	)
	return list, err
}

func (r *queryResolver) Qualification(ctx context.Context, id *int, slug *string) (*model.Qualification, error) {
	if id != nil {
		return r.QualificationUsecase.GetByID(ctx, *id)
	} else if slug != nil {
		return r.QualificationUsecase.GetBySlug(ctx, *slug)
	}

	return nil, nil
}
