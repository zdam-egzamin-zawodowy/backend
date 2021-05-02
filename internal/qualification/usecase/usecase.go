package usecase

import (
	"context"
	"github.com/Kichiyaki/gopgutil/v10"
	"github.com/pkg/errors"

	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
	"github.com/zdam-egzamin-zawodowy/backend/internal/qualification"
)

type usecase struct {
	qualificationRepository qualification.Repository
}

type Config struct {
	QualificationRepository qualification.Repository
}

func New(cfg *Config) (qualification.Usecase, error) {
	if cfg == nil || cfg.QualificationRepository == nil {
		return nil, errors.New("qualification/usecase: QualificationRepository is required")
	}
	return &usecase{
		cfg.QualificationRepository,
	}, nil
}

func (ucase *usecase) Store(ctx context.Context, input *models.QualificationInput) (*models.Qualification, error) {
	if err := validateInput(input.Sanitize(), validateOptions{false}); err != nil {
		return nil, err
	}
	return ucase.qualificationRepository.Store(ctx, input)
}

func (ucase *usecase) UpdateOneByID(ctx context.Context, id int, input *models.QualificationInput) (*models.Qualification, error) {
	if id <= 0 {
		return nil, errors.New(messageInvalidID)
	}
	if err := validateInput(input.Sanitize(), validateOptions{true}); err != nil {
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
		return nil, errors.New(messageItemNotFound)
	}
	return items[0], nil
}

func (ucase *usecase) Delete(ctx context.Context, f *models.QualificationFilter) ([]*models.Qualification, error) {
	return ucase.qualificationRepository.Delete(ctx, f)
}

func (ucase *usecase) Fetch(ctx context.Context, cfg *qualification.FetchConfig) ([]*models.Qualification, int, error) {
	if cfg == nil {
		cfg = &qualification.FetchConfig{
			Limit: qualification.FetchDefaultLimit,
			Count: true,
		}
	}
	cfg.Sort = gopgutil.SanitizeOrders(cfg.Sort)
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
		return nil, errors.New(messageItemNotFound)
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
		return nil, errors.New(messageItemNotFound)
	}
	return items[0], nil
}

func (ucase *usecase) GetSimilar(ctx context.Context, cfg *qualification.GetSimilarConfig) ([]*models.Qualification, int, error) {
	if cfg == nil || cfg.QualificationID <= 0 {
		return nil, 0, errors.New(messageQualificationIDIsRequired)
	}
	cfg.Sort = gopgutil.SanitizeOrders(cfg.Sort)
	return ucase.qualificationRepository.GetSimilar(ctx, cfg)
}

type validateOptions struct {
	allowNilValues bool
}

func validateInput(input *models.QualificationInput, opts validateOptions) error {
	if input.IsEmpty() {
		return errors.New(messageEmptyPayload)
	}

	if input.Name != nil {
		if *input.Name == "" {
			return errors.New(messageNameIsRequired)
		} else if len(*input.Name) > qualification.MaxNameLength {
			return errors.Errorf(messageNameIsTooLong, qualification.MaxNameLength)
		}
	} else if !opts.allowNilValues {
		return errors.New(messageNameIsRequired)
	}

	if input.Code != nil {
		if *input.Code == "" {
			return errors.New(messageCodeIsRequired)
		}
	} else if !opts.allowNilValues {
		return errors.New(messageCodeIsRequired)
	}

	return nil
}
