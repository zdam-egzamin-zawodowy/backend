package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/zdam-egzamin-zawodowy/backend/internal/graphql/generated"
	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input models.UserInput) (*models.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id int, input models.UserInput) (*models.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateManyUsers(ctx context.Context, ids []int, input models.UserInput) ([]*models.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteUsers(ctx context.Context, ids []int) ([]*models.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SignIn(ctx context.Context, email string, password string, staySignedIn *bool) (*models.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Users(ctx context.Context, filter *models.UserFilter, limit *int, offset *int, sort []string) (*generated.UserList, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) User(ctx context.Context, id int) (*models.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Me(ctx context.Context) (*models.User, error) {
	panic(fmt.Errorf("not implemented"))
}
