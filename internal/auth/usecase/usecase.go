package usecase

import (
	"context"
	"fmt"

	"github.com/zdam-egzamin-zawodowy/backend/internal/auth"
	"github.com/zdam-egzamin-zawodowy/backend/internal/auth/jwt"
	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
	"github.com/zdam-egzamin-zawodowy/backend/internal/user"
	errorutils "github.com/zdam-egzamin-zawodowy/backend/pkg/utils/error"
)

type usecase struct {
	userRepository user.Repository
	tokenGenerator jwt.TokenGenerator
}

type Config struct {
	UserRepository user.Repository
	TokenGenerator jwt.TokenGenerator
}

func New(cfg *Config) (auth.Usecase, error) {
	if cfg == nil || cfg.UserRepository == nil {
		return nil, fmt.Errorf("user/usecase: UserRepository is required")
	}
	return &usecase{
		cfg.UserRepository,
		cfg.TokenGenerator,
	}, nil
}

func (ucase *usecase) SignIn(ctx context.Context, email, password string, staySignedIn bool) (*models.User, string, error) {
	user, err := ucase.GetUserByCredentials(ctx, email, password)
	if err != nil {
		return nil, "", err
	}

	token, err := ucase.tokenGenerator.Generate(jwt.Metadata{
		StaySignedIn: staySignedIn,
		Credentials: jwt.Credentials{
			Email:    user.Email,
			Password: user.Password,
		},
	})
	if err != nil {
		return nil, "", errorutils.Wrap(err, messageInvalidCredentials)
	}

	return user, token, nil
}

func (ucase *usecase) ExtractAccessTokenMetadata(ctx context.Context, accessToken string) (*models.User, error) {
	metadata, err := ucase.tokenGenerator.ExtractAccessTokenMetadata(accessToken)
	if err != nil {
		return nil, errorutils.Wrap(err, messageInvalidAccessToken)
	}

	return ucase.GetUserByCredentials(ctx, metadata.Credentials.Email, metadata.Credentials.Password)
}

func (ucase *usecase) GetUserByCredentials(ctx context.Context, email, password string) (*models.User, error) {
	users, _, err := ucase.userRepository.Fetch(ctx, &user.FetchConfig{
		Limit: 1,
		Count: false,
		Filter: &models.UserFilter{
			Email: []string{email},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(users) <= 0 {
		return nil, fmt.Errorf(messageInvalidCredentials)
	}

	user := users[0]
	if err := user.CompareHashAndPassword(password); err != nil {
		return nil, fmt.Errorf(messageInvalidCredentials)
	}

	return user, nil
}
