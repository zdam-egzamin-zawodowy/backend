package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/zdam-egzamin-zawodowy/backend/internal/gin/middleware"
	"github.com/zdam-egzamin-zawodowy/backend/internal/graphql/generated"
	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
	"github.com/zdam-egzamin-zawodowy/backend/internal/user"
	"github.com/zdam-egzamin-zawodowy/backend/pkg/utils"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input models.UserInput) (*models.User, error) {
	return r.UserUsecase.Store(ctx, &input)
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id int, input models.UserInput) (*models.User, error) {
	return r.UserUsecase.UpdateOneByID(ctx, id, &input)
}

func (r *mutationResolver) UpdateManyUsers(ctx context.Context, ids []int, input models.UserInput) ([]*models.User, error) {
	return r.UserUsecase.UpdateMany(
		ctx,
		&models.UserFilter{
			ID: ids,
		},
		&input,
	)
}

func (r *mutationResolver) DeleteUsers(ctx context.Context, ids []int) ([]*models.User, error) {
	return r.UserUsecase.Delete(ctx, &models.UserFilter{
		ID: ids,
	})
}

func (r *mutationResolver) SignIn(
	ctx context.Context,
	email string,
	password string,
	staySignedIn *bool,
) (*generated.UserWithToken, error) {
	var err error
	userWithToken := &generated.UserWithToken{}
	userWithToken.User, userWithToken.Token, err = r.AuthUsecase.SignIn(
		ctx,
		email,
		password,
		utils.SafeBoolPointer(staySignedIn, false),
	)
	if err != nil {
		return nil, err
	}
	return userWithToken, nil
}

func (r *queryResolver) Users(
	ctx context.Context,
	filter *models.UserFilter,
	limit *int,
	offset *int,
	sort []string,
) (*generated.UserList, error) {
	var err error
	userList := &generated.UserList{}
	userList.Items, userList.Total, err = r.UserUsecase.Fetch(
		ctx,
		&user.FetchConfig{
			Count:  shouldCount(ctx),
			Filter: filter,
			Limit:  utils.SafeIntPointer(limit, user.FetchMaxLimit),
			Offset: utils.SafeIntPointer(offset, 0),
			Sort:   sort,
		},
	)
	return userList, err
}

func (r *queryResolver) User(ctx context.Context, id int) (*models.User, error) {
	return r.UserUsecase.GetByID(ctx, id)
}

func (r *queryResolver) Me(ctx context.Context) (*models.User, error) {
	user, _ := middleware.UserFromContext(ctx)
	return user, nil
}
