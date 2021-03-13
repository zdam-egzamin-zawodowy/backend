package usecase

import (
	"context"
	"fmt"

	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
	"github.com/zdam-egzamin-zawodowy/backend/internal/question"
	sqlutils "github.com/zdam-egzamin-zawodowy/backend/pkg/utils/sql"
)

var (
	imageValidMIMETypes = map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
	}
)

type usecase struct {
	questionRepository question.Repository
}

type Config struct {
	QuestionRepository question.Repository
}

func New(cfg *Config) (question.Usecase, error) {
	if cfg == nil || cfg.QuestionRepository == nil {
		return nil, fmt.Errorf("question/usecase: cfg.QuestionRepository is required")
	}
	return &usecase{
		cfg.QuestionRepository,
	}, nil
}

func (ucase *usecase) Store(ctx context.Context, input *models.QuestionInput) (*models.Question, error) {
	if err := validateInput(input.Sanitize(), validateOptions{false}); err != nil {
		return nil, err
	}
	return ucase.questionRepository.Store(ctx, input)
}

func (ucase *usecase) UpdateOneByID(ctx context.Context, id int, input *models.QuestionInput) (*models.Question, error) {
	if id <= 0 {
		return nil, fmt.Errorf(messageInvalidID)
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
		return nil, fmt.Errorf(messageItemNotFound)
	}
	return item, nil
}

func (ucase *usecase) Delete(ctx context.Context, f *models.QuestionFilter) ([]*models.Question, error) {
	return ucase.questionRepository.Delete(ctx, f)
}

func (ucase *usecase) Fetch(ctx context.Context, cfg *question.FetchConfig) ([]*models.Question, int, error) {
	if cfg == nil {
		cfg = &question.FetchConfig{
			Limit: question.FetchMaxLimit,
			Count: true,
		}
	}
	if cfg.Limit > question.FetchMaxLimit {
		cfg.Limit = question.FetchMaxLimit
	}
	cfg.Sort = sqlutils.SanitizeSorts(cfg.Sort)
	return ucase.questionRepository.Fetch(ctx, cfg)
}

func (ucase *usecase) GetByID(ctx context.Context, id int) (*models.Question, error) {
	items, _, err := ucase.Fetch(ctx, &question.FetchConfig{
		Limit: 1,
		Count: false,
		Filter: &models.QuestionFilter{
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

func (ucase *usecase) GenerateTest(ctx context.Context, cfg *question.GenerateTestConfig) ([]*models.Question, error) {
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

func validateInput(input *models.QuestionInput, opts validateOptions) error {
	if input.IsEmpty() {
		return fmt.Errorf(messageEmptyPayload)
	}

	if input.Content != nil {
		if *input.Content == "" {
			return fmt.Errorf(messageContentIsRequired)
		}
	} else if !opts.allowNilValues {
		return fmt.Errorf(messageContentIsRequired)
	}

	if input.CorrectAnswer != nil {
		if !input.CorrectAnswer.IsValid() {
			return fmt.Errorf(messageCorrectAnswerIsInvalid)
		}
	} else if !opts.allowNilValues {
		return fmt.Errorf(messageCorrectAnswerIsInvalid)
	}

	if input.QualificationID != nil {
		if *input.QualificationID <= 0 {
			return fmt.Errorf(messageQualificationIDIsRequired)
		}
	} else if !opts.allowNilValues {
		return fmt.Errorf(messageQualificationIDIsRequired)
	}

	if input.AnswerA != nil {
		if !input.AnswerA.IsValid() {
			return fmt.Errorf(messageAnswerIsInvalid, "A")
		}
	}

	if input.AnswerB != nil {
		if !input.AnswerB.IsValid() {
			return fmt.Errorf(messageAnswerIsInvalid, "B")
		}
	}

	if input.AnswerC != nil {
		if !input.AnswerC.IsValid() {
			return fmt.Errorf(messageAnswerIsInvalid, "C")
		}
	}

	if input.AnswerD != nil {
		if !input.AnswerD.IsValid() {
			return fmt.Errorf(messageAnswerIsInvalid, "D")
		}
	}

	if input.Image != nil {
		if validateMimeType(input.Image.ContentType) {
			return fmt.Errorf(messageImageNotAcceptableMIMEType, "Obrazek pytanie")
		}
	}

	if input.AnswerAImage != nil {
		if !validateMimeType(input.AnswerAImage.ContentType) {
			return fmt.Errorf(messageAnswerIsInvalid, "Obrazek odpowiedź A")
		}
	}

	if input.AnswerBImage != nil {
		if !validateMimeType(input.AnswerBImage.ContentType) {
			return fmt.Errorf(messageAnswerIsInvalid, "Obrazek odpowiedź B")
		}
	}

	if input.AnswerCImage != nil {
		if !validateMimeType(input.AnswerCImage.ContentType) {
			return fmt.Errorf(messageAnswerIsInvalid, "Obrazek odpowiedź C")
		}
	}

	if input.AnswerDImage != nil {
		if !validateMimeType(input.AnswerDImage.ContentType) {
			return fmt.Errorf(messageAnswerIsInvalid, "Obrazek odpowiedź D")
		}
	}

	if input.DeleteAnswerAImage != nil && input.AnswerA == nil && input.AnswerAImage == nil {
		return fmt.Errorf(messageCannotDeleteImageWithoutNewAnswer, "Obrazek odpowiedź A")
	}

	if input.DeleteAnswerBImage != nil && input.AnswerB == nil && input.AnswerBImage == nil {
		return fmt.Errorf(messageCannotDeleteImageWithoutNewAnswer, "Obrazek odpowiedź B")
	}

	if input.DeleteAnswerCImage != nil && input.AnswerC == nil && input.AnswerCImage == nil {
		return fmt.Errorf(messageCannotDeleteImageWithoutNewAnswer, "Obrazek odpowiedź C")
	}

	if input.DeleteAnswerDImage != nil && input.AnswerD == nil && input.AnswerDImage == nil {
		return fmt.Errorf(messageCannotDeleteImageWithoutNewAnswer, "Obrazek odpowiedź D")

	}

	if !opts.allowNilValues {
		if input.AnswerA == nil && input.AnswerAImage == nil {
			return fmt.Errorf(messageAnswerIsRequired, "A")
		}

		if input.AnswerB == nil && input.AnswerBImage == nil {
			return fmt.Errorf(messageAnswerIsRequired, "B")
		}

		if input.AnswerC == nil && input.AnswerCImage == nil {
			return fmt.Errorf(messageAnswerIsRequired, "C")
		}

		if input.AnswerD == nil && input.AnswerDImage == nil {
			return fmt.Errorf(messageAnswerIsRequired, "D")
		}
	}

	return nil
}

func validateMimeType(contentType string) bool {
	return imageValidMIMETypes[contentType]
}
