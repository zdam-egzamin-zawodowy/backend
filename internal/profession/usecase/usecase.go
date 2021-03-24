package usecase

import (
	"context"
	"fmt"

	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
	"github.com/zdam-egzamin-zawodowy/backend/internal/profession"
	sqlutils "github.com/zdam-egzamin-zawodowy/backend/pkg/utils/sql"
)

type usecase struct {
	professionRepository profession.Repository
}

type Config struct {
	ProfessionRepository profession.Repository
}

func New(cfg *Config) (profession.Usecase, error) {
	if cfg == nil || cfg.ProfessionRepository == nil {
		return nil, fmt.Errorf("profession/usecase: ProfessionRepository is required")
	}
	return &usecase{
		cfg.ProfessionRepository,
	}, nil
}

func (ucase *usecase) Store(ctx context.Context, input *models.ProfessionInput) (*models.Profession, error) {
	if err := validateInput(input.Sanitize(), validateOptions{false}); err != nil {
		return nil, err
	}
	return ucase.professionRepository.Store(ctx, input)
}

func (ucase *usecase) UpdateOneByID(ctx context.Context, id int, input *models.ProfessionInput) (*models.Profession, error) {
	if id <= 0 {
		return nil, fmt.Errorf(messageInvalidID)
	}
	if err := validateInput(input.Sanitize(), validateOptions{true}); err != nil {
		return nil, err
	}
	items, err := ucase.professionRepository.UpdateMany(ctx,
		&models.ProfessionFilter{
			ID: []int{id},
		},
		input)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf(messageItemNotFound)
	}
	return items[0], nil
}

func (ucase *usecase) Delete(ctx context.Context, f *models.ProfessionFilter) ([]*models.Profession, error) {
	return ucase.professionRepository.Delete(ctx, f)
}

func (ucase *usecase) Fetch(ctx context.Context, cfg *profession.FetchConfig) ([]*models.Profession, int, error) {
	if cfg == nil {
		cfg = &profession.FetchConfig{
			Limit: profession.FetchDefaultLimit,
			Count: true,
		}
	}
	cfg.Sort = sqlutils.SanitizeSorts(cfg.Sort)
	return ucase.professionRepository.Fetch(ctx, cfg)
}

func (ucase *usecase) GetByID(ctx context.Context, id int) (*models.Profession, error) {
	items, _, err := ucase.Fetch(ctx, &profession.FetchConfig{
		Limit: 1,
		Count: false,
		Filter: &models.ProfessionFilter{
			ID: []int{id},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf(messageItemNotFound)
	}
	return items[0], nil
}

func (ucase *usecase) GetBySlug(ctx context.Context, slug string) (*models.Profession, error) {
	items, _, err := ucase.Fetch(ctx, &profession.FetchConfig{
		Limit: 1,
		Count: false,
		Filter: &models.ProfessionFilter{
			Slug: []string{slug},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf(messageItemNotFound)
	}
	return items[0], nil
}

type validateOptions struct {
	allowNilValues bool
}

func validateInput(input *models.ProfessionInput, opts validateOptions) error {
	if input.IsEmpty() {
		return fmt.Errorf(messageEmptyPayload)
	}

	if input.Name != nil {
		if *input.Name == "" {
			return fmt.Errorf(messageNameIsRequired)
		} else if len(*input.Name) > profession.MaxNameLength {
			return fmt.Errorf(messageNameIsTooLong, profession.MaxNameLength)
		}
	} else if !opts.allowNilValues {
		return fmt.Errorf(messageNameIsRequired)
	}

	return nil
}
