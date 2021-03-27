package usecase

import (
	"context"
	"github.com/pkg/errors"

	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
	"github.com/zdam-egzamin-zawodowy/backend/internal/user"
	"github.com/zdam-egzamin-zawodowy/backend/pkg/utils"
	sqlutils "github.com/zdam-egzamin-zawodowy/backend/pkg/utils/sql"
)

type usecase struct {
	userRepository user.Repository
}

type Config struct {
	UserRepository user.Repository
}

func New(cfg *Config) (user.Usecase, error) {
	if cfg == nil || cfg.UserRepository == nil {
		return nil, errors.New("user/usecase: UserRepository is required")
	}
	return &usecase{
		cfg.UserRepository,
	}, nil
}

func (ucase *usecase) Store(ctx context.Context, input *models.UserInput) (*models.User, error) {
	if err := validateInput(input.Sanitize(), validateOptions{false}); err != nil {
		return nil, err
	}
	return ucase.userRepository.Store(ctx, input)
}

func (ucase *usecase) UpdateOneByID(ctx context.Context, id int, input *models.UserInput) (*models.User, error) {
	if id <= 0 {
		return nil, errors.New(messageInvalidID)
	}
	items, err := ucase.UpdateMany(
		ctx,
		&models.UserFilter{
			ID: []int{id},
		},
		input.Sanitize(),
	)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, errors.New(messageItemNotFound)
	}
	return items[0], nil
}

func (ucase *usecase) UpdateMany(ctx context.Context, f *models.UserFilter, input *models.UserInput) ([]*models.User, error) {
	if f == nil {
		return []*models.User{}, nil
	}
	if err := validateInput(input.Sanitize(), validateOptions{true}); err != nil {
		return nil, err
	}
	items, err := ucase.userRepository.UpdateMany(ctx, f, input)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (ucase *usecase) Delete(ctx context.Context, f *models.UserFilter) ([]*models.User, error) {
	return ucase.userRepository.Delete(ctx, f)
}

func (ucase *usecase) Fetch(ctx context.Context, cfg *user.FetchConfig) ([]*models.User, int, error) {
	if cfg == nil {
		cfg = &user.FetchConfig{
			Limit: user.FetchMaxLimit,
			Count: true,
		}
	}
	if cfg.Limit > user.FetchMaxLimit || cfg.Limit <= 0 {
		cfg.Limit = user.FetchMaxLimit
	}
	cfg.Sort = sqlutils.SanitizeSorts(cfg.Sort)
	return ucase.userRepository.Fetch(ctx, cfg)
}

func (ucase *usecase) GetByID(ctx context.Context, id int) (*models.User, error) {
	items, _, err := ucase.Fetch(ctx, &user.FetchConfig{
		Limit: 1,
		Count: false,
		Filter: &models.UserFilter{
			ID: []int{id},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, errors.New(messageItemNotFound)
	}
	return items[0], nil
}

func (ucase *usecase) GetByCredentials(ctx context.Context, email, password string) (*models.User, error) {
	items, _, err := ucase.Fetch(ctx, &user.FetchConfig{
		Limit: 1,
		Count: false,
		Filter: &models.UserFilter{
			Email: []string{email},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, errors.New(messageInvalidCredentials)
	}
	if err := items[0].CompareHashAndPassword(password); err != nil {
		return nil, errors.New(messageInvalidCredentials)
	}
	return items[0], nil
}

type validateOptions struct {
	acceptNilValues bool
}

func validateInput(input *models.UserInput, opts validateOptions) error {
	if input.IsEmpty() {
		return errors.New(messageEmptyPayload)
	}

	if input.DisplayName != nil {
		displayNameLength := len(*input.DisplayName)
		if displayNameLength < user.MinDisplayNameLength {
			return errors.New(messageDisplayNameIsRequired)
		} else if displayNameLength > user.MaxDisplayNameLength {
			return errors.Errorf(messageDisplayNameIsTooLong, user.MaxDisplayNameLength)
		}
	} else if !opts.acceptNilValues {
		return errors.New(messageDisplayNameIsRequired)
	}

	if input.Email != nil {
		if !utils.IsEmailValid(*input.Email) {
			return errors.New(messageEmailIsInvalid)
		}
	} else if !opts.acceptNilValues {
		return errors.New(messageEmailIsRequired)
	}

	if input.Password != nil {
		passwordLength := len(*input.Password)
		if passwordLength > user.MaxPasswordLength || passwordLength < user.MinPasswordLength {
			return errors.Errorf(messagePasswordInvalidLength, user.MinPasswordLength, user.MaxPasswordLength)
		}
	} else if !opts.acceptNilValues {
		return errors.New(messagePasswordIsRequired)
	}

	if input.Role != nil {
		if !input.Role.IsValid() {
			return errors.New(messageInvalidRole)
		}
	} else if !opts.acceptNilValues {
		return errors.New(messageInvalidRole)
	}

	return nil
}
