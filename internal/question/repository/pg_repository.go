package repository

import (
	"context"
	"github.com/Kichiyaki/gopgutil/v10"
	"github.com/pkg/errors"
	"github.com/zdam-egzamin-zawodowy/backend/internal"
	"strings"
	"time"

	"github.com/go-pg/pg/v10"

	"github.com/zdam-egzamin-zawodowy/backend/fstorage"
	"github.com/zdam-egzamin-zawodowy/backend/internal/question"
	"github.com/zdam-egzamin-zawodowy/backend/util/errorutil"
)

type PGRepositoryConfig struct {
	DB          *pg.DB
	FileStorage fstorage.FileStorage
}

type PGRepository struct {
	*pg.DB
	*repository
}

var _ question.Repository = &PGRepository{}

func NewPGRepository(cfg *PGRepositoryConfig) (*PGRepository, error) {
	if cfg == nil || cfg.DB == nil {
		return nil, errors.New("cfg.DB is required")
	}
	if cfg.FileStorage == nil {
		return nil, errors.New("cfg.FileStorage is required")
	}
	return &PGRepository{
		cfg.DB,
		&repository{
			fileStorage: cfg.FileStorage,
		},
	}, nil
}

func (repo *PGRepository) Store(ctx context.Context, input *internal.QuestionInput) (*internal.Question, error) {
	item := input.ToQuestion()
	baseQuery := repo.
		Model(item).
		Context(ctx).
		Returning("*")
	if _, err := baseQuery.
		Clone().
		Insert(); err != nil {
		return nil, handleInsertAndUpdateError(err)
	}

	repo.saveImages(item, input)
	if _, err := baseQuery.
		Clone().
		WherePK().
		Set("image = ?", item.Image).
		Set("answer_a_image = ?", item.AnswerAImage).
		Set("answer_b_image = ?", item.AnswerBImage).
		Set("answer_c_image = ?", item.AnswerCImage).
		Set("answer_d_image = ?", item.AnswerDImage).
		Update(); err != nil && err != pg.ErrNoRows {
		return nil, errorutil.Wrap(err, messageFailedToSaveModel)
	}

	return item, nil
}

func (repo *PGRepository) UpdateOneByID(ctx context.Context, id int, input *internal.QuestionInput) (*internal.Question, error) {
	item := &internal.Question{}
	baseQuery := repo.
		Model(item).
		Context(ctx).
		Returning("*").
		Where(gopgutil.BuildConditionEquals("?"), gopgutil.AddAliasToColumnName("id", "question"), id).
		Set("updated_at = ?", time.Now())

	if _, err := baseQuery.
		Clone().
		Apply(input.ApplyUpdate).
		Update(); err != nil && err != pg.ErrNoRows {
		return nil, handleInsertAndUpdateError(err)
	}

	repo.saveImages(item, input)
	repo.deleteImagesBasedOnInput(item, input)

	if _, err := baseQuery.
		Clone().
		Set("image = ?", item.Image).
		Set("answer_a_image = ?", item.AnswerAImage).
		Set("answer_b_image = ?", item.AnswerBImage).
		Set("answer_c_image = ?", item.AnswerCImage).
		Set("answer_d_image = ?", item.AnswerDImage).
		Update(); err != nil && err != pg.ErrNoRows {
		return nil, handleInsertAndUpdateError(err)
	}

	return item, nil
}

func (repo *PGRepository) Delete(ctx context.Context, f *internal.QuestionFilter) ([]*internal.Question, error) {
	items := make([]*internal.Question, 0)
	if _, err := repo.
		Model(&items).
		Context(ctx).
		Returning("*").
		Apply(f.Where).
		Delete(); err != nil && err != pg.ErrNoRows {
		return nil, errorutil.Wrap(err, messageFailedToDeleteModel)
	}

	go repo.getAllImagesAndDelete(items)

	return items, nil
}

func (repo *PGRepository) Fetch(ctx context.Context, cfg *question.FetchConfig) ([]*internal.Question, int, error) {
	var err error
	items := make([]*internal.Question, 0)
	total := 0
	query := repo.
		Model(&items).
		Context(ctx).
		Limit(cfg.Limit).
		Offset(cfg.Offset).
		Apply(cfg.Filter.Where).
		Apply(gopgutil.OrderAppender{
			Orders: cfg.Sort,
		}.Apply)

	if cfg.Count {
		total, err = query.SelectAndCount()
	} else {
		err = query.Select()
	}
	if err != nil && err != pg.ErrNoRows {
		return nil, 0, errorutil.Wrap(err, messageFailedToFetchModel)
	}
	return items, total, nil
}

func (repo *PGRepository) GenerateTest(ctx context.Context, cfg *question.GenerateTestConfig) ([]*internal.Question, error) {
	subquery := repo.
		Model(&internal.Question{}).
		Column("id").
		Where(gopgutil.BuildConditionArray("qualification_id"), pg.Array(cfg.Qualifications)).
		OrderExpr("random()").
		Limit(cfg.Limit)
	items := make([]*internal.Question, 0)
	if err := repo.
		Model(&items).
		Context(ctx).
		Where(gopgutil.BuildConditionIn("id"), subquery).
		Select(); err != nil && err != pg.ErrNoRows {
		return nil, errorutil.Wrap(err, messageFailedToFetchModel)
	}
	return items, nil
}

func handleInsertAndUpdateError(err error) error {
	if strings.Contains(err.Error(), "questions_from_content_correct_answer_qualification_id_key") {
		return errorutil.Wrap(err, messageSimilarRecordExists)
	}
	return errorutil.Wrap(err, messageFailedToSaveModel)
}
