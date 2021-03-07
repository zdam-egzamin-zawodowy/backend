package repository

import (
	"context"
	"fmt"
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
		return nil, fmt.Errorf("profession/pg_repository: *pg.DB is required")
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
		if strings.Contains(err.Error(), "name") || strings.Contains(err.Error(), "slug") {
			return nil, errorutils.Wrap(err, messageNameIsAlreadyTaken)
		}
		return nil, errorutils.Wrap(err, messageFailedToSaveModel)
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
		if strings.Contains(err.Error(), "name") || strings.Contains(err.Error(), "slug") {
			return nil, errorutils.Wrap(err, messageNameIsAlreadyTaken)
		}
		return nil, errorutils.Wrap(err, messageFailedToSaveModel)
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
