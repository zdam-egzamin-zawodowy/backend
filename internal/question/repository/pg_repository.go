package repository

import (
	"context"
	"github.com/Kichiyaki/gopgutil/v10"
	"github.com/pkg/errors"
	"strings"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
	"github.com/zdam-egzamin-zawodowy/backend/internal/question"
	"github.com/zdam-egzamin-zawodowy/backend/pkg/fstorage"
	"github.com/zdam-egzamin-zawodowy/backend/pkg/util/errorutil"
)

type pgRepository struct {
	*pg.DB
	*repository
}

type PGRepositoryConfig struct {
	DB          *pg.DB
	FileStorage fstorage.FileStorage
}

func NewPGRepository(cfg *PGRepositoryConfig) (question.Repository, error) {
	if cfg == nil || cfg.DB == nil || cfg.FileStorage == nil {
		return nil, errors.New("question/pg_repository: *pg.DB and filestorage.FileStorage are required")
	}
	return &pgRepository{
		cfg.DB,
		&repository{
			fileStorage: cfg.FileStorage,
		},
	}, nil
}

func (repo *pgRepository) Store(ctx context.Context, input *models.QuestionInput) (*models.Question, error) {
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

func (repo *pgRepository) UpdateOneByID(ctx context.Context, id int, input *models.QuestionInput) (*models.Question, error) {
	item := &models.Question{}
	baseQuery := repo.
		Model(item).
		Context(ctx).
		Returning("*").
		Where(gopgutil.BuildConditionEquals(gopgutil.AddAliasToColumnName("id", "question")), id).
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

func (repo *pgRepository) Delete(ctx context.Context, f *models.QuestionFilter) ([]*models.Question, error) {
	items := []*models.Question{}
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

func (repo *pgRepository) Fetch(ctx context.Context, cfg *question.FetchConfig) ([]*models.Question, int, error) {
	var err error
	var items []*models.Question
	total := 0
	query := repo.
		Model(&items).
		Context(ctx).
		Limit(cfg.Limit).
		Offset(cfg.Offset).
		Apply(gopgutil.OrderAppender{
			Relations: map[string]gopgutil.OrderAppenderRelation{
				"qualification": {
					Name: "Qualification",
				},
			},
			Orders: cfg.Sort,
		}.Apply).
		Apply(cfg.Filter.Where)

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

func (repo *pgRepository) GenerateTest(ctx context.Context, cfg *question.GenerateTestConfig) ([]*models.Question, error) {
	subquery := repo.
		Model(&models.Question{}).
		Column("id").
		Where(gopgutil.BuildConditionArray("qualification_id"), pg.Array(cfg.Qualifications)).
		OrderExpr("random()").
		Limit(cfg.Limit)
	items := []*models.Question{}
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
