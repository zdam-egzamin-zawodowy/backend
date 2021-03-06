package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/zdam-egzamin-zawodowy/backend/internal/graphql/generated"
	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
)

func (r *mutationResolver) CreateQuestion(ctx context.Context, input models.QuestionInput) (*models.Question, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateQuestion(ctx context.Context, id int, input models.QuestionInput) (*models.Question, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteQuestions(ctx context.Context, ids []int) ([]*models.Question, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) GenerateTest(ctx context.Context, qualificationIDs []int, limit *int) ([]*models.Question, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Questions(ctx context.Context, filter *models.QuestionFilter, limit *int, offset *int, sort []string) (*generated.QuestionList, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *questionResolver) Qualification(ctx context.Context, obj *models.Question) (*models.Qualification, error) {
	panic(fmt.Errorf("not implemented"))
}
