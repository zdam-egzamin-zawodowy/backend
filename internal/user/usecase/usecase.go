package usecase

import (
	"context"
	"github.com/Kichiyaki/goutil/strutil"
	"github.com/pkg/errors"

	"github.com/zdam-egzamin-zawodowy/backend/internal/model"
	"github.com/zdam-egzamin-zawodowy/backend/internal/user"
)

type usecase struct {
	userRepository user.Repository
}

type Config struct {
	UserRepository user.Repository
}

func New(cfg *Config) (user.Usecase, error) {
	if cfg == nil || cfg.UserRepository == nil {
		return nil, errors.New("cfg.UserRepository is required")
	}
	return &usecase{
		cfg.UserRepository,
	}, nil
}

func (ucase *usecase) Store(ctx context.Context, input *model.UserInput) (*model.User, error) {
	if err := validateInput(input.Sanitize(), validateOptions{false}); err != nil {
		return nil, err
	}
	return ucase.userRepository.Store(ctx, input)
}

func (ucase *usecase) UpdateOneByID(ctx context.Context, id int, input *model.UserInput) (*model.User, error) {
	if id <= 0 {
		return nil, errors.New(messageInvalidID)
	}
	items, err := ucase.UpdateMany(
		ctx,
		&model.UserFilter{
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

func (ucase *usecase) UpdateMany(ctx context.Context, f *model.UserFilter, input *model.UserInput) ([]*model.User, error) {
	if f == nil {
		return []*model.User{}, nil
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

func (ucase *usecase) Delete(ctx context.Context, f *model.UserFilter) ([]*model.User, error) {
	return ucase.userRepository.Delete(ctx, f)
}

func (ucase *usecase) Fetch(ctx context.Context, cfg *user.FetchConfig) ([]*model.User, int, error) {
	if cfg == nil {
		cfg = &user.FetchConfig{
			Limit: user.FetchMaxLimit,
			Count: true,
		}
	}
	if cfg.Limit > user.FetchMaxLimit || cfg.Limit <= 0 {
		cfg.Limit = user.FetchMaxLimit
	}
	if len(cfg.Sort) > user.MaxOrders {
		cfg.Sort = cfg.Sort[0:user.MaxOrders]
	}
	return ucase.userRepository.Fetch(ctx, cfg)
}

func (ucase *usecase) GetByID(ctx context.Context, id int) (*model.User, error) {
	items, _, err := ucase.Fetch(ctx, &user.FetchConfig{
		Limit: 1,
		Count: false,
		Filter: &model.UserFilter{
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

func (ucase *usecase) GetByCredentials(ctx context.Context, email, password string) (*model.User, error) {
	items, _, err := ucase.Fetch(ctx, &user.FetchConfig{
		Limit: 1,
		Count: false,
		Filter: &model.UserFilter{
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

func validateInput(input *model.UserInput, opts validateOptions) error {
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
		if !strutil.IsEmail(*input.Email) {
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
