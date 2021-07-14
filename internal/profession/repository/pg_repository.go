package repository

import (
	"context"
	"github.com/Kichiyaki/gopgutil/v10"
	"github.com/pkg/errors"
	"strings"

	"github.com/zdam-egzamin-zawodowy/backend/pkg/util/errorutil"

	"github.com/go-pg/pg/v10"

	"github.com/zdam-egzamin-zawodowy/backend/internal/model"
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
		return nil, errors.New("cfg.DB is required")
	}
	return &pgRepository{
		cfg.DB,
	}, nil
}

func (repo *pgRepository) Store(ctx context.Context, input *model.ProfessionInput) (*model.Profession, error) {
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

func (repo *pgRepository) UpdateMany(ctx context.Context, f *model.ProfessionFilter, input *model.ProfessionInput) ([]*model.Profession, error) {
	if _, err := repo.
		Model(&model.Profession{}).
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

func (repo *pgRepository) Delete(ctx context.Context, f *model.ProfessionFilter) ([]*model.Profession, error) {
	items := make([]*model.Profession, 0)
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

func (repo *pgRepository) Fetch(ctx context.Context, cfg *profession.FetchConfig) ([]*model.Profession, int, error) {
	var err error
	items := make([]*model.Profession, 0)
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

func (repo *pgRepository) GetAssociatedQualifications(
	ctx context.Context,
	ids ...int,
) (map[int][]*model.Qualification, error) {
	m := make(map[int][]*model.Qualification)
	for _, id := range ids {
		m[id] = make([]*model.Qualification, 0)
	}
	var qualificationToProfession []*model.QualificationToProfession
	if err := repo.
		Model(&qualificationToProfession).
		Context(ctx).
		Where(gopgutil.BuildConditionArray("profession_id"), pg.Array(ids)).
		Relation("Qualification").
		Order("qualification.formula ASC", "qualification.code ASC").
		Select(); err != nil {
		return nil, errorutil.Wrap(err, messageFailedToFetchAssociatedQualifications)
	}
	for _, record := range qualificationToProfession {
		m[record.ProfessionID] = append(m[record.ProfessionID], record.Qualification)
	}
	return m, nil
}

func handleInsertAndUpdateError(err error) error {
	if strings.Contains(err.Error(), "name") || strings.Contains(err.Error(), "slug") {
		return errorutil.Wrap(err, messageNameIsAlreadyTaken)
	}
	return errorutil.Wrap(err, messageFailedToSaveModel)
}
