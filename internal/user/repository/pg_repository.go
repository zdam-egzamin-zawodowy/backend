package repository

import (
	"context"
	"github.com/pkg/errors"
	"strings"

	errorutils "github.com/zdam-egzamin-zawodowy/backend/pkg/utils/error"

	"github.com/go-pg/pg/v10"
	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
	"github.com/zdam-egzamin-zawodowy/backend/internal/user"
)

type pgRepository struct {
	*pg.DB
}

type PGRepositoryConfig struct {
	DB *pg.DB
}

func NewPGRepository(cfg *PGRepositoryConfig) (user.Repository, error) {
	if cfg == nil || cfg.DB == nil {
		return nil, errors.New("user/pg_repository: *pg.DB is required")
	}
	return &pgRepository{
		cfg.DB,
	}, nil
}

func (repo *pgRepository) Store(ctx context.Context, input *models.UserInput) (*models.User, error) {
	item := input.ToUser()
	if _, err := repo.
		Model(item).
		Context(ctx).
		Returning("*").
		Insert(); err != nil {
		return nil, handleInsertAndUpdateError(err)
	}
	return item, nil
}

func (repo *pgRepository) UpdateMany(ctx context.Context, f *models.UserFilter, input *models.UserInput) ([]*models.User, error) {
	if _, err := repo.
		Model(&models.User{}).
		Context(ctx).
		Apply(input.ApplyUpdate).
		Apply(f.Where).
		Update(); err != nil && err != pg.ErrNoRows {
		return nil, handleInsertAndUpdateError(err)
	}
	items, _, err := repo.Fetch(ctx, &user.FetchConfig{
		Count:  false,
		Filter: f,
	})
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *pgRepository) Delete(ctx context.Context, f *models.UserFilter) ([]*models.User, error) {
	items := []*models.User{}
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

func (repo *pgRepository) Fetch(ctx context.Context, cfg *user.FetchConfig) ([]*models.User, int, error) {
	var err error
	items := []*models.User{}
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

func handleInsertAndUpdateError(err error) error {
	if strings.Contains(err.Error(), "email") {
		return errorutils.Wrap(err, messageEmailIsAlreadyTaken)
	}
	return errorutils.Wrap(err, messageFailedToSaveModel)
}
