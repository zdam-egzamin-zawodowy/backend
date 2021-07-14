package usecase

import (
	"context"
	"github.com/pkg/errors"

	"github.com/zdam-egzamin-zawodowy/backend/internal/model"
	"github.com/zdam-egzamin-zawodowy/backend/internal/profession"
)

type usecase struct {
	professionRepository profession.Repository
}

type Config struct {
	ProfessionRepository profession.Repository
}

func New(cfg *Config) (profession.Usecase, error) {
	if cfg == nil || cfg.ProfessionRepository == nil {
		return nil, errors.New("cfg.ProfessionRepository is required")
	}
	return &usecase{
		cfg.ProfessionRepository,
	}, nil
}

func (ucase *usecase) Store(ctx context.Context, input *model.ProfessionInput) (*model.Profession, error) {
	if err := validateInput(input.Sanitize(), validateOptions{false}); err != nil {
		return nil, err
	}
	return ucase.professionRepository.Store(ctx, input)
}

func (ucase *usecase) UpdateOneByID(ctx context.Context, id int, input *model.ProfessionInput) (*model.Profession, error) {
	if id <= 0 {
		return nil, errors.New(messageInvalidID)
	}
	if err := validateInput(input.Sanitize(), validateOptions{true}); err != nil {
		return nil, err
	}
	items, err := ucase.professionRepository.UpdateMany(ctx,
		&model.ProfessionFilter{
			ID: []int{id},
		},
		input)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, errors.New(messageItemNotFound)
	}
	return items[0], nil
}

func (ucase *usecase) Delete(ctx context.Context, f *model.ProfessionFilter) ([]*model.Profession, error) {
	return ucase.professionRepository.Delete(ctx, f)
}

func (ucase *usecase) Fetch(ctx context.Context, cfg *profession.FetchConfig) ([]*model.Profession, int, error) {
	if cfg == nil {
		cfg = &profession.FetchConfig{
			Limit: profession.FetchDefaultLimit,
			Count: true,
		}
	}
	if len(cfg.Sort) > profession.MaxOrders {
		cfg.Sort = cfg.Sort[0:profession.MaxOrders]
	}

	return ucase.professionRepository.Fetch(ctx, cfg)
}

func (ucase *usecase) GetByID(ctx context.Context, id int) (*model.Profession, error) {
	items, _, err := ucase.Fetch(ctx, &profession.FetchConfig{
		Limit: 1,
		Count: false,
		Filter: &model.ProfessionFilter{
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

func (ucase *usecase) GetBySlug(ctx context.Context, slug string) (*model.Profession, error) {
	items, _, err := ucase.Fetch(ctx, &profession.FetchConfig{
		Limit: 1,
		Count: false,
		Filter: &model.ProfessionFilter{
			Slug: []string{slug},
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

type validateOptions struct {
	allowNilValues bool
}

func validateInput(input *model.ProfessionInput, opts validateOptions) error {
	if input.IsEmpty() {
		return errors.New(messageEmptyPayload)
	}

	if input.Name != nil {
		if *input.Name == "" {
			return errors.New(messageNameIsRequired)
		} else if len(*input.Name) > profession.MaxNameLength {
			return errors.Errorf(messageNameIsTooLong, profession.MaxNameLength)
		}
	} else if !opts.allowNilValues {
		return errors.New(messageNameIsRequired)
	}

	return nil
}
