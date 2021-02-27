package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
	"github.com/zdam-egzamin-zawodowy/backend/internal/qualification"
	sqlutils "github.com/zdam-egzamin-zawodowy/backend/pkg/utils/sql"
)

type usecase struct {
	qualificationRepository qualification.Repository
}

type Config struct {
	QualificationRepository qualification.Repository
}

func New(cfg *Config) (qualification.Usecase, error) {
	if cfg == nil || cfg.QualificationRepository == nil {
		return nil, fmt.Errorf("qualification/usecase: QualificationRepository is required")
	}
	return &usecase{
		cfg.QualificationRepository,
	}, nil
}

func (ucase *usecase) Store(ctx context.Context, input *models.QualificationInput) (*models.Qualification, error) {
	if err := ucase.validateInput(input, validateOptions{false}); err != nil {
		return nil, err
	}
	return ucase.qualificationRepository.Store(ctx, input)
}

func (ucase *usecase) UpdateOne(ctx context.Context, id int, input *models.QualificationInput) (*models.Qualification, error) {
	if id <= 0 {
		return nil, fmt.Errorf(messageInvalidID)
	}
	if err := ucase.validateInput(input, validateOptions{true}); err != nil {
		return nil, err
	}
	items, err := ucase.qualificationRepository.UpdateMany(ctx,
		&models.QualificationFilter{
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

func (ucase *usecase) Delete(ctx context.Context, f *models.QualificationFilter) ([]*models.Qualification, error) {
	return ucase.qualificationRepository.Delete(ctx, f)
}

func (ucase *usecase) Fetch(ctx context.Context, cfg *qualification.FetchConfig) ([]*models.Qualification, int, error) {
	if cfg == nil {
		cfg = &qualification.FetchConfig{
			Limit: qualification.DefaultLimit,
			Count: true,
		}
	}
	cfg.Sort = sqlutils.SanitizeSortExpressions(cfg.Sort)
	return ucase.qualificationRepository.Fetch(ctx, cfg)
}

func (ucase *usecase) GetByID(ctx context.Context, id int) (*models.Qualification, error) {
	items, _, err := ucase.Fetch(ctx, &qualification.FetchConfig{
		Limit: 1,
		Count: false,
		Filter: &models.QualificationFilter{
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

func (ucase *usecase) GetBySlug(ctx context.Context, slug string) (*models.Qualification, error) {
	items, _, err := ucase.Fetch(ctx, &qualification.FetchConfig{
		Limit: 1,
		Count: false,
		Filter: &models.QualificationFilter{
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

func (ucase *usecase) validateInput(input *models.QualificationInput, opts validateOptions) error {
	if input.IsEmpty() {
		return fmt.Errorf(messageEmptyPayload)
	}

	if input.Name != nil {
		trimmedName := strings.TrimSpace(*input.Name)
		input.Name = &trimmedName
		if trimmedName == "" {
			return fmt.Errorf(messageNameIsRequired)
		} else if len(trimmedName) > qualification.MaxNameLength {
			return fmt.Errorf(messageNameIsTooLong, qualification.MaxNameLength)
		}
	} else if !opts.allowNilValues {
		return fmt.Errorf(messageNameIsRequired)
	}

	if input.Code != nil {
		trimmedCode := strings.TrimSpace(*input.Code)
		input.Code = &trimmedCode
		if trimmedCode == "" {
			return fmt.Errorf(messageCodeIsRequired)
		}
	} else if !opts.allowNilValues {
		return fmt.Errorf(messageCodeIsRequired)
	}

	return nil
}
