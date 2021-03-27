package repository

import (
	"context"
	"github.com/pkg/errors"
	sqlutils "github.com/zdam-egzamin-zawodowy/backend/pkg/utils/sql"
	"strings"

	errorutils "github.com/zdam-egzamin-zawodowy/backend/pkg/utils/error"

	"github.com/go-pg/pg/v10"
	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
	"github.com/zdam-egzamin-zawodowy/backend/internal/profession"
)

type pgRepository struct {
	*pg.DB
}

type PGRepositoryConfig struct {
	DB *pg.DB
}

func NewPGRepository(cfg *PGRepositoryConfig) (profession.Repository, error) {
	if cfg == nil || cfg.DB == nil {
		return nil, errors.New("profession/pg_repository: *pg.DB is required")
	}
	return &pgRepository{
		cfg.DB,
	}, nil
}

func (repo *pgRepository) Store(ctx context.Context, input *models.ProfessionInput) (*models.Profession, error) {
	item := input.ToProfession()
	if _, err := repo.
		Model(item).
		Context(ctx).
		Returning("*").
		Insert(); err != nil {
		return nil, handleInsertAndUpdateError(err)
	}
	return item, nil
}

func (repo *pgRepository) UpdateMany(ctx context.Context, f *models.ProfessionFilter, input *models.ProfessionInput) ([]*models.Profession, error) {
	if _, err := repo.
		Model(&models.Profession{}).
		Context(ctx).
		Apply(input.ApplyUpdate).
		Apply(f.Where).
		Update(); err != nil && err != pg.ErrNoRows {
		return nil, handleInsertAndUpdateError(err)
	}
	items, _, err := repo.Fetch(ctx, &profession.FetchConfig{
		Count:  false,
		Filter: f,
	})
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *pgRepository) Delete(ctx context.Context, f *models.ProfessionFilter) ([]*models.Profession, error) {
	items := []*models.Profession{}
	if _, err := repo.
		Model(&items).
		Context(ctx).
		Returning("*").
		Apply(f.Where).
		Delete(); err != nil && err != pg.ErrNoRows {
		return nil, errorutils.Wrap(err, messageFailedToDeleteModel)
	}
	return items, nil
}

func (repo *pgRepository) Fetch(ctx context.Context, cfg *profession.FetchConfig) ([]*models.Profession, int, error) {
	var err error
	items := []*models.Profession{}
	total := 0
	query := repo.
		Model(&items).
		Context(ctx).
		Limit(cfg.Limit).
		Offset(cfg.Offset).
		Apply(cfg.Filter.Where).
		Order(cfg.Sort...)

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

func (repo *pgRepository) GetAssociatedQualifications(
	ctx context.Context,
	ids ...int,
) (map[int][]*models.Qualification, error) {
	m := make(map[int][]*models.Qualification)
	for _, id := range ids {
		m[id] = []*models.Qualification{}
	}
	qualificationToProfession := []*models.QualificationToProfession{}
	if err := repo.
		Model(&qualificationToProfession).
		Context(ctx).
		Where(sqlutils.BuildConditionArray("profession_id"), pg.Array(ids)).
		Relation("Qualification").
		Select(); err != nil {
		return nil, errorutils.Wrap(err, messageFailedToFetchAssociatedQualifications)
	}
	for _, record := range qualificationToProfession {
		m[record.ProfessionID] = append(m[record.ProfessionID], record.Qualification)
	}
	return m, nil
}

func handleInsertAndUpdateError(err error) error {
	if strings.Contains(err.Error(), "name") || strings.Contains(err.Error(), "slug") {
		return errorutils.Wrap(err, messageNameIsAlreadyTaken)
	}
	return errorutils.Wrap(err, messageFailedToSaveModel)
}
