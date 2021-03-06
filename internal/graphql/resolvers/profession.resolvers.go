package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/zdam-egzamin-zawodowy/backend/internal/graphql/generated"
	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
)

func (r *mutationResolver) CreateProfession(ctx context.Context, input models.ProfessionInput) (*models.Profession, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateProfession(ctx context.Context, id int, input models.ProfessionInput) (*models.Profession, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteProfessions(ctx context.Context, ids []int) ([]*models.Profession, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Professions(ctx context.Context, filter *models.ProfessionFilter, limit *int, offset *int, sort []string) (*generated.ProfessionList, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Profession(ctx context.Context, id *int, slug *string) (*models.Profession, error) {
	panic(fmt.Errorf("not implemented"))
}
