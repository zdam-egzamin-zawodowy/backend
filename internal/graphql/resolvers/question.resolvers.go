package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/zdam-egzamin-zawodowy/backend/pkg/util/safepointer"

	"github.com/zdam-egzamin-zawodowy/backend/internal/gin/middleware"
	"github.com/zdam-egzamin-zawodowy/backend/internal/graphql/generated"
	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
	"github.com/zdam-egzamin-zawodowy/backend/internal/question"
)

func (r *mutationResolver) CreateQuestion(ctx context.Context, input models.QuestionInput) (*models.Question, error) {
	return r.QuestionUsecase.Store(ctx, &input)
}

func (r *mutationResolver) UpdateQuestion(ctx context.Context, id int, input models.QuestionInput) (*models.Question, error) {
	return r.QuestionUsecase.UpdateOneByID(ctx, id, &input)
}

func (r *mutationResolver) DeleteQuestions(ctx context.Context, ids []int) ([]*models.Question, error) {
	return r.QuestionUsecase.Delete(ctx, &models.QuestionFilter{
		ID: ids,
	})
}

func (r *queryResolver) GenerateTest(ctx context.Context, qualificationIDs []int, limit *int) ([]*models.Question, error) {
	return r.QuestionUsecase.GenerateTest(ctx, &question.GenerateTestConfig{
		Qualifications: qualificationIDs,
		Limit:          safepointer.SafeIntPointer(limit, question.TestMaxLimit),
	})
}

func (r *queryResolver) Questions(
	ctx context.Context,
	filter *models.QuestionFilter,
	limit *int,
	offset *int,
	sort []string,
) (*generated.QuestionList, error) {
	var err error
	list := &generated.QuestionList{}
	list.Items, list.Total, err = r.QuestionUsecase.Fetch(
		ctx,
		&question.FetchConfig{
			Count:  shouldCount(ctx),
			Filter: filter,
			Limit:  safepointer.SafeIntPointer(limit, question.FetchDefaultLimit),
			Offset: safepointer.SafeIntPointer(offset, 0),
			Sort:   sort,
		},
	)
	return list, err
}

func (r *questionResolver) Qualification(ctx context.Context, obj *models.Question) (*models.Qualification, error) {
	if obj != nil && obj.Qualification != nil {
		return obj.Qualification, nil
	}

	if obj != nil && obj.QualificationID > 0 {
		if dataloader, err := middleware.DataLoaderFromContext(ctx); err == nil && dataloader != nil {
			return dataloader.QualificationByID.Load(obj.QualificationID)
		}
	}

	return nil, nil
}
