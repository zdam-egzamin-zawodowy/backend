package usecase

import (
	"context"
	"github.com/pkg/errors"

	"github.com/zdam-egzamin-zawodowy/backend/internal/auth"
	"github.com/zdam-egzamin-zawodowy/backend/internal/auth/jwt"
	"github.com/zdam-egzamin-zawodowy/backend/internal/model"
	"github.com/zdam-egzamin-zawodowy/backend/internal/user"
	"github.com/zdam-egzamin-zawodowy/backend/util/errorutil"
)

type Config struct {
	UserRepository user.Repository
	TokenGenerator *jwt.TokenGenerator
}

type Usecase struct {
	userRepository user.Repository
	tokenGenerator *jwt.TokenGenerator
}

var _ auth.Usecase = &Usecase{}

func New(cfg *Config) (*Usecase, error) {
	if cfg == nil || cfg.UserRepository == nil {
		return nil, errors.New("cfg.UserRepository is required")
	}
	return &Usecase{
		cfg.UserRepository,
		cfg.TokenGenerator,
	}, nil
}

func (ucase *Usecase) SignIn(ctx context.Context, email, password string, staySignedIn bool) (*model.User, string, error) {
	u, err := ucase.GetUserByCredentials(ctx, email, password)
	if err != nil {
		return nil, "", err
	}

	token, err := ucase.tokenGenerator.Generate(jwt.Metadata{
		StaySignedIn: staySignedIn,
		Credentials: jwt.Credentials{
			Email:    u.Email,
			Password: u.Password,
		},
	})
	if err != nil {
		return nil, "", errorutil.Wrap(err, messageInvalidCredentials)
	}

	return u, token, nil
}

func (ucase *Usecase) ExtractAccessTokenMetadata(ctx context.Context, accessToken string) (*model.User, error) {
	metadata, err := ucase.tokenGenerator.ExtractAccessTokenMetadata(accessToken)
	if err != nil {
		return nil, errorutil.Wrap(err, messageInvalidAccessToken)
	}

	return ucase.GetUserByCredentials(ctx, metadata.Credentials.Email, metadata.Credentials.Password)
}

func (ucase *Usecase) GetUserByCredentials(ctx context.Context, email, password string) (*model.User, error) {
	users, _, err := ucase.userRepository.Fetch(ctx, &user.FetchConfig{
		Limit: 1,
		Count: false,
		Filter: &model.UserFilter{
			Email: []string{email},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(users) <= 0 {
		return nil, errors.New(messageInvalidCredentials)
	}

	u := users[0]
	if err := u.CompareHashAndPassword(password); err != nil {
		return nil, errorutil.Wrap(err, messageInvalidCredentials)
	}

	return u, nil
}
