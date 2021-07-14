package usecase

import (
	"context"
	"github.com/pkg/errors"

	"github.com/zdam-egzamin-zawodowy/backend/internal/model"
	"github.com/zdam-egzamin-zawodowy/backend/internal/qualification"
)

type Config struct {
	QualificationRepository qualification.Repository
}

type Usecase struct {
	qualificationRepository qualification.Repository
}

var _ qualification.Usecase = &Usecase{}

func New(cfg *Config) (*Usecase, error) {
	if cfg == nil || cfg.QualificationRepository == nil {
		return nil, errors.New("cfg.QualificationRepository is required")
	}
	return &Usecase{
		cfg.QualificationRepository,
	}, nil
}

func (ucase *Usecase) Store(ctx context.Context, input *model.QualificationInput) (*model.Qualification, error) {
	if err := validateInput(input.Sanitize(), validateOptions{false}); err != nil {
		return nil, err
	}
	return ucase.qualificationRepository.Store(ctx, input)
}

func (ucase *Usecase) UpdateOneByID(ctx context.Context, id int, input *model.QualificationInput) (*model.Qualification, error) {
	if id <= 0 {
		return nil, errors.New(messageInvalidID)
	}
	if err := validateInput(input.Sanitize(), validateOptions{true}); err != nil {
		return nil, err
	}
	items, err := ucase.qualificationRepository.UpdateMany(ctx,
		&model.QualificationFilter{
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

func (ucase *Usecase) Delete(ctx context.Context, f *model.QualificationFilter) ([]*model.Qualification, error) {
	return ucase.qualificationRepository.Delete(ctx, f)
}

func (ucase *Usecase) Fetch(ctx context.Context, cfg *qualification.FetchConfig) ([]*model.Qualification, int, error) {
	if cfg == nil {
		cfg = &qualification.FetchConfig{
			Limit: qualification.FetchDefaultLimit,
			Count: true,
		}
	}
	if len(cfg.Sort) > qualification.MaxOrders {
		cfg.Sort = cfg.Sort[0:qualification.MaxOrders]
	}
	return ucase.qualificationRepository.Fetch(ctx, cfg)
}

func (ucase *Usecase) GetByID(ctx context.Context, id int) (*model.Qualification, error) {
	items, _, err := ucase.Fetch(ctx, &qualification.FetchConfig{
		Limit: 1,
		Count: false,
		Filter: &model.QualificationFilter{
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

func (ucase *Usecase) GetBySlug(ctx context.Context, slug string) (*model.Qualification, error) {
	items, _, err := ucase.Fetch(ctx, &qualification.FetchConfig{
		Limit: 1,
		Count: false,
		Filter: &model.QualificationFilter{
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

func (ucase *Usecase) GetSimilar(ctx context.Context, cfg *qualification.GetSimilarConfig) ([]*model.Qualification, int, error) {
	if cfg == nil || cfg.QualificationID <= 0 {
		return nil, 0, errors.New(messageQualificationIDIsRequired)
	}
	return ucase.qualificationRepository.GetSimilar(ctx, cfg)
}

type validateOptions struct {
	allowNilValues bool
}

func validateInput(input *model.QualificationInput, opts validateOptions) error {
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
