package usecase

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zdam-egzamin-zawodowy/backend/internal"

	"github.com/zdam-egzamin-zawodowy/backend/internal/question"
)

var (
	imageValidMIMETypes = map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
	}
)

type Config struct {
	QuestionRepository question.Repository
}

type Usecase struct {
	questionRepository question.Repository
}

var _ question.Usecase = &Usecase{}

func New(cfg *Config) (*Usecase, error) {
	if cfg == nil || cfg.QuestionRepository == nil {
		return nil, errors.New("cfg.QuestionRepository is required")
	}
	return &Usecase{
		cfg.QuestionRepository,
	}, nil
}

func (ucase *Usecase) Store(ctx context.Context, input *internal.QuestionInput) (*internal.Question, error) {
	if err := validateInput(input.Sanitize(), validateOptions{false}); err != nil {
		return nil, err
	}
	return ucase.questionRepository.Store(ctx, input)
}

func (ucase *Usecase) UpdateOneByID(ctx context.Context, id int, input *internal.QuestionInput) (*internal.Question, error) {
	if id <= 0 {
		return nil, errors.New(messageInvalidID)
	}
	if err := validateInput(input.Sanitize(), validateOptions{true}); err != nil {
		return nil, err
	}
	item, err := ucase.questionRepository.UpdateOneByID(ctx,
		id,
		input)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.New(messageItemNotFound)
	}
	return item, nil
}

func (ucase *Usecase) Delete(ctx context.Context, f *internal.QuestionFilter) ([]*internal.Question, error) {
	return ucase.questionRepository.Delete(ctx, f)
}

func (ucase *Usecase) Fetch(ctx context.Context, cfg *question.FetchConfig) ([]*internal.Question, int, error) {
	if cfg == nil {
		cfg = &question.FetchConfig{
			Limit: question.FetchMaxLimit,
			Count: true,
		}
	}
	if cfg.Limit > question.FetchMaxLimit {
		cfg.Limit = question.FetchMaxLimit
	}
	if len(cfg.Sort) > question.MaxOrders {
		cfg.Sort = cfg.Sort[0:question.MaxOrders]
	}
	return ucase.questionRepository.Fetch(ctx, cfg)
}

func (ucase *Usecase) GetByID(ctx context.Context, id int) (*internal.Question, error) {
	items, _, err := ucase.Fetch(ctx, &question.FetchConfig{
		Limit: 1,
		Count: false,
		Filter: &internal.QuestionFilter{
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

func (ucase *Usecase) GenerateTest(ctx context.Context, cfg *question.GenerateTestConfig) ([]*internal.Question, error) {
	if cfg == nil {
		cfg = &question.GenerateTestConfig{
			Limit: question.TestMaxLimit,
		}
	}
	if cfg.Limit > question.TestMaxLimit {
		cfg.Limit = question.TestMaxLimit
	}
	return ucase.questionRepository.GenerateTest(ctx, cfg)
}

type validateOptions struct {
	allowNilValues bool
}

func validateInput(input *internal.QuestionInput, opts validateOptions) error {
	if input.IsEmpty() {
		return errors.New(messageEmptyPayload)
	}

	if input.Content != nil {
		if *input.Content == "" {
			return errors.New(messageContentIsRequired)
		}
	} else if !opts.allowNilValues {
		return errors.New(messageContentIsRequired)
	}

	if input.CorrectAnswer != nil {
		if !input.CorrectAnswer.IsValid() {
			return errors.New(messageCorrectAnswerIsInvalid)
		}
	} else if !opts.allowNilValues {
		return errors.New(messageCorrectAnswerIsInvalid)
	}

	if input.QualificationID != nil {
		if *input.QualificationID <= 0 {
			return errors.New(messageQualificationIDIsRequired)
		}
	} else if !opts.allowNilValues {
		return errors.New(messageQualificationIDIsRequired)
	}

	if input.AnswerA != nil {
		if *input.AnswerA == "" {
			return errors.Errorf(messageAnswerIsRequired, "A")
		}
	}

	if input.AnswerB != nil {
		if *input.AnswerA == "" {
			return errors.Errorf(messageAnswerIsRequired, "B")
		}
	}

	if input.AnswerC != nil {
		if *input.AnswerC == "" {
			return errors.Errorf(messageAnswerIsRequired, "C")
		}
	}

	if input.AnswerD != nil {
		if *input.AnswerD == "" {
			return errors.Errorf(messageAnswerIsRequired, "D")
		}
	}

	if input.Image != nil {
		if !isValidMIMEType(input.Image.ContentType) {
			return errors.Errorf(messageImageNotAcceptableMIMEType, "Obrazek pytanie")
		}
	}

	if input.AnswerAImage != nil {
		if !isValidMIMEType(input.AnswerAImage.ContentType) {
			return errors.Errorf(messageImageNotAcceptableMIMEType, "Obrazek odpowiedź A")
		}
	}

	if input.AnswerBImage != nil {
		if !isValidMIMEType(input.AnswerBImage.ContentType) {
			return errors.Errorf(messageImageNotAcceptableMIMEType, "Obrazek odpowiedź B")
		}
	}

	if input.AnswerCImage != nil {
		if !isValidMIMEType(input.AnswerCImage.ContentType) {
			return errors.Errorf(messageImageNotAcceptableMIMEType, "Obrazek odpowiedź C")
		}
	}

	if input.AnswerDImage != nil {
		if !isValidMIMEType(input.AnswerDImage.ContentType) {
			return errors.Errorf(messageImageNotAcceptableMIMEType, "Obrazek odpowiedź D")
		}
	}

	if input.DeleteAnswerAImage != nil && input.AnswerA == nil && input.AnswerAImage == nil {
		return errors.Errorf(messageCannotDeleteImageWithoutNewAnswer, "Obrazek odpowiedź A")
	}

	if input.DeleteAnswerBImage != nil && input.AnswerB == nil && input.AnswerBImage == nil {
		return errors.Errorf(messageCannotDeleteImageWithoutNewAnswer, "Obrazek odpowiedź B")
	}

	if input.DeleteAnswerCImage != nil && input.AnswerC == nil && input.AnswerCImage == nil {
		return errors.Errorf(messageCannotDeleteImageWithoutNewAnswer, "Obrazek odpowiedź C")
	}

	if input.DeleteAnswerDImage != nil && input.AnswerD == nil && input.AnswerDImage == nil {
		return errors.Errorf(messageCannotDeleteImageWithoutNewAnswer, "Obrazek odpowiedź D")

	}

	if !opts.allowNilValues {
		if input.AnswerA == nil && input.AnswerAImage == nil {
			return errors.Errorf(messageAnswerIsRequired, "A")
		}

		if input.AnswerB == nil && input.AnswerBImage == nil {
			return errors.Errorf(messageAnswerIsRequired, "B")
		}

		if input.AnswerC == nil && input.AnswerCImage == nil {
			return errors.Errorf(messageAnswerIsRequired, "C")
		}

		if input.AnswerD == nil && input.AnswerDImage == nil {
			return errors.Errorf(messageAnswerIsRequired, "D")
		}
	}

	return nil
}

func isValidMIMEType(contentType string) bool {
	return imageValidMIMETypes[contentType]
}
