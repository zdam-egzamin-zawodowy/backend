package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zdam-egzamin-zawodowy/backend/pkg/filestorage"
	errorutils "github.com/zdam-egzamin-zawodowy/backend/pkg/utils/error"
	sqlutils "github.com/zdam-egzamin-zawodowy/backend/pkg/utils/sql"

	"github.com/go-pg/pg/v10"
	"github.com/zdam-egzamin-zawodowy/backend/internal/db"
	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
	"github.com/zdam-egzamin-zawodowy/backend/internal/question"
)

type pgRepository struct {
	*pg.DB
	*repository
}

type PGRepositoryConfig struct {
	DB          *pg.DB
	FileStorage filestorage.FileStorage
}

func NewPGRepository(cfg *PGRepositoryConfig) (question.Repository, error) {
	if cfg == nil || cfg.DB == nil || cfg.FileStorage == nil {
		return nil, fmt.Errorf("question/pg_repository: *pg.DB and filestorage.FileStorage are required")
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
	if _, err := repo.
		Model(item).
		Context(ctx).
		Returning("*").
		Insert(); err != nil {
		if strings.Contains(err.Error(), "questions_from_content_correct_answer_qualification_id_key") {
			return nil, errorutils.Wrap(err, messageSimilarRecordExists)
		}
		return nil, errorutils.Wrap(err, messageFailedToSaveModel)
	}

	repo.saveImages(item, input)

	return item, nil
}

func (repo *pgRepository) UpdateOneByID(ctx context.Context, id int, input *models.QuestionInput) (*models.Question, error) {
	item := &models.Question{}
	baseQuery := repo.
		Model(item).
		Context(ctx).
		Returning("*").
		Where(sqlutils.BuildConditionEquals(sqlutils.AddAliasToColumnName("id", "question")), id).
		Set("updated_at = ?", time.Now())

	if _, err := baseQuery.
		Clone().
		Apply(input.ApplyUpdate).
		Update(); err != nil && err != pg.ErrNoRows {
		if strings.Contains(err.Error(), "questions_from_content_correct_answer_qualification_id_key") {
			return nil, errorutils.Wrap(err, messageSimilarRecordExists)
		}
		return nil, errorutils.Wrap(err, messageFailedToSaveModel)
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
		return nil, errorutils.Wrap(err, messageFailedToSaveModel)
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
		return nil, errorutils.Wrap(err, messageFailedToDeleteModel)
	}

	go repo.getAllImagesAndDelete(items)

	return items, nil
}

func (repo *pgRepository) Fetch(ctx context.Context, cfg *question.FetchConfig) ([]*models.Question, int, error) {
	var err error
	items := []*models.Question{}
	total := 0
	query := repo.
		Model(&items).
		Context(ctx).
		Limit(cfg.Limit).
		Offset(cfg.Offset).
		Apply(db.Sort{
			Relationships: map[string]string{
				"qualification": "qualification",
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
		return nil, 0, errorutils.Wrap(err, messageFailedToFetchModel)
	}
	return items, total, nil
}

func (repo *pgRepository) GenerateTest(ctx context.Context, cfg *question.GenerateTestConfig) ([]*models.Question, error) {
	subquery := repo.
		Model(&models.Question{}).
		Column("id").
		Where(sqlutils.BuildConditionArray("qualification_id"), pg.Array(cfg.Qualifications)).
		OrderExpr("random()").
		Limit(cfg.Limit)
	items := []*models.Question{}
	if err := repo.
		Model(&items).
		Context(ctx).
		Where(sqlutils.BuildConditionIn("id"), subquery).
		Select(); err != nil && err != pg.ErrNoRows {
		return nil, errorutils.Wrap(err, messageFailedToFetchModel)
	}
	return items, nil
}
