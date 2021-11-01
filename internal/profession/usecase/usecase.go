package usecase

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zdam-egzamin-zawodowy/backend/internal"

	"github.com/zdam-egzamin-zawodowy/backend/internal/profession"
)

type Config struct {
	ProfessionRepository profession.Repository
}

type Usecase struct {
	professionRepository profession.Repository
}

var _ profession.Usecase = &Usecase{}

func New(cfg *Config) (*Usecase, error) {
	if cfg == nil || cfg.ProfessionRepository == nil {
		return nil, errors.New("cfg.ProfessionRepository is required")
	}
	return &Usecase{
		cfg.ProfessionRepository,
	}, nil
}

func (ucase *Usecase) Store(ctx context.Context, input *internal.ProfessionInput) (*internal.Profession, error) {
	if err := validateInput(input.Sanitize(), validateOptions{false}); err != nil {
		return nil, err
	}
	return ucase.professionRepository.Store(ctx, input)
}

func (ucase *Usecase) UpdateOneByID(ctx context.Context, id int, input *internal.ProfessionInput) (*internal.Profession, error) {
	if id <= 0 {
		return nil, errors.New(messageInvalidID)
	}
	if err := validateInput(input.Sanitize(), validateOptions{true}); err != nil {
		return nil, err
	}
	items, err := ucase.professionRepository.UpdateMany(ctx,
		&internal.ProfessionFilter{
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

func (ucase *Usecase) Delete(ctx context.Context, f *internal.ProfessionFilter) ([]*internal.Profession, error) {
	return ucase.professionRepository.Delete(ctx, f)
}

func (ucase *Usecase) Fetch(ctx context.Context, cfg *profession.FetchConfig) ([]*internal.Profession, int, error) {
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

func (ucase *Usecase) GetByID(ctx context.Context, id int) (*internal.Profession, error) {
	items, _, err := ucase.Fetch(ctx, &profession.FetchConfig{
		Limit: 1,
		Count: false,
		Filter: &internal.ProfessionFilter{
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

func (ucase *Usecase) GetBySlug(ctx context.Context, slug string) (*internal.Profession, error) {
	items, _, err := ucase.Fetch(ctx, &profession.FetchConfig{
		Limit: 1,
		Count: false,
		Filter: &internal.ProfessionFilter{
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

func validateInput(input *internal.ProfessionInput, opts validateOptions) error {
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
