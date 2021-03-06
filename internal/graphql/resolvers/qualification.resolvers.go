package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/zdam-egzamin-zawodowy/backend/internal/graphql/generated"
	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
)

func (r *mutationResolver) CreateQualification(ctx context.Context, input models.QualificationInput) (*models.Qualification, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateQualification(ctx context.Context, id int, input models.QualificationInput) (*models.Qualification, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteQualifications(ctx context.Context, ids []int) ([]*models.Qualification, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Qualifications(ctx context.Context, filter *models.QualificationFilter, limit *int, offset *int, sort []string) (*generated.QualificationList, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Qualification(ctx context.Context, id *int, slug *string) (*models.Qualification, error) {
	panic(fmt.Errorf("not implemented"))
}
