package repository

import (
	"context"
	"github.com/Kichiyaki/gopgutil/v10"
	"github.com/pkg/errors"
	"strings"

	"github.com/zdam-egzamin-zawodowy/backend/pkg/util/errorutil"

	"github.com/go-pg/pg/v10"

	"github.com/zdam-egzamin-zawodowy/backend/internal/model"
	"github.com/zdam-egzamin-zawodowy/backend/internal/user"
)

type PGRepositoryConfig struct {
	DB *pg.DB
}

type PGRepository struct {
	*pg.DB
}

var _ user.Repository = &PGRepository{}

func NewPGRepository(cfg *PGRepositoryConfig) (*PGRepository, error) {
	if cfg == nil || cfg.DB == nil {
		return nil, errors.New("cfg.DB is required")
	}
	return &PGRepository{
		cfg.DB,
	}, nil
}

func (repo *PGRepository) Store(ctx context.Context, input *model.UserInput) (*model.User, error) {
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

func (repo *PGRepository) UpdateMany(ctx context.Context, f *model.UserFilter, input *model.UserInput) ([]*model.User, error) {
	if _, err := repo.
		Model(&model.User{}).
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

func (repo *PGRepository) Delete(ctx context.Context, f *model.UserFilter) ([]*model.User, error) {
	items := make([]*model.User, 0)
	if _, err := repo.
		Model(&items).
		Context(ctx).
		Returning("*").
		Apply(f.Where).
		Delete(); err != nil && err != pg.ErrNoRows {
		return nil, errorutil.Wrap(err, messageFailedToDeleteModel)
	}
	return items, nil
}

func (repo *PGRepository) Fetch(ctx context.Context, cfg *user.FetchConfig) ([]*model.User, int, error) {
	var err error
	items := make([]*model.User, 0)
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

func handleInsertAndUpdateError(err error) error {
	if strings.Contains(err.Error(), "email") {
		return errorutil.Wrap(err, messageEmailIsAlreadyTaken)
	}
	return errorutil.Wrap(err, messageFailedToSaveModel)
}
