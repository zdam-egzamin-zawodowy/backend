package repository

import (
	"context"
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
	db *pg.DB
}

func NewPGRepository(cfg PGRepositoryConfig) profession.Repository {
	return &pgRepository{
		cfg.db,
	}
}

func (repo *pgRepository) Store(ctx context.Context, input *models.ProfessionInput) (*models.Profession, error) {
	item := input.ToProfession()
	if _, err := repo.
		Model(item).
		Context(ctx).
		Returning("*").
		Insert(); err != nil {
		if strings.Contains(err.Error(), "name") {
			return nil, errorutils.Wrap(err, nameIsAlreadyTaken)
		}
		return nil, errorutils.Wrap(err, failedToSaveModel)
	}
	return item, nil
}

func (repo *pgRepository) Update(ctx context.Context, f *models.ProfessionFilter, input *models.ProfessionInput) ([]*models.Profession, error) {
	items := []*models.Profession{}
	if _, err := repo.
		Model(&items).
		Context(ctx).
		Returning("*").
		Apply(input.ApplyUpdate).
		Apply(f.Where).
		Update(); err != nil {
		if strings.Contains(err.Error(), "name") {
			return nil, errorutils.Wrap(err, nameIsAlreadyTaken)
		}
		return nil, errorutils.Wrap(err, failedToSaveModel)
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
		Delete(); err != nil {
		return nil, errorutils.Wrap(err, failedToDeleteModel)
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
		Apply(cfg.Filter.Where)

	if cfg.Count {
		total, err = query.SelectAndCount()
	} else {
		err = query.Select()
	}
	if err != nil && err != pg.ErrNoRows {
		return nil, 0, errorutils.Wrap(err, failedToFetchModel)
	}
	return items, total, nil
}
