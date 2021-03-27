package repository

import (
	"context"
	"github.com/pkg/errors"
	"strings"

	sqlutils "github.com/zdam-egzamin-zawodowy/backend/pkg/utils/sql"

	errorutils "github.com/zdam-egzamin-zawodowy/backend/pkg/utils/error"

	"github.com/go-pg/pg/v10"
	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
	"github.com/zdam-egzamin-zawodowy/backend/internal/qualification"
)

type pgRepository struct {
	*pg.DB
}

type PGRepositoryConfig struct {
	DB *pg.DB
}

func NewPGRepository(cfg *PGRepositoryConfig) (qualification.Repository, error) {
	if cfg == nil || cfg.DB == nil {
		return nil, errors.New("qualification/pg_repository: *pg.DB is required")
	}
	return &pgRepository{
		cfg.DB,
	}, nil
}

func (repo *pgRepository) Store(ctx context.Context, input *models.QualificationInput) (*models.Qualification, error) {
	item := input.ToQualification()
	err := repo.RunInTransaction(ctx, func(tx *pg.Tx) error {
		if _, err := tx.
			Model(item).
			Context(ctx).
			Returning("*").
			Insert(); err != nil {
			return handleInsertAndUpdateError(err)
		}

		if len(input.AssociateProfession) > 0 {
			if err := repo.associateQualificationWithProfession(tx, []int{item.ID}, input.AssociateProfession); err != nil {
				return handleInsertAndUpdateError(err)
			}
		}

		return nil
	})
	return item, err
}

func (repo *pgRepository) UpdateMany(
	ctx context.Context,
	f *models.QualificationFilter,
	input *models.QualificationInput,
) ([]*models.Qualification, error) {
	items := []*models.Qualification{}
	err := repo.RunInTransaction(ctx, func(tx *pg.Tx) error {
		if input.HasBasicDataToUpdate() {
			if _, err := tx.
				Model(&models.Qualification{}).
				Context(ctx).
				Apply(input.ApplyUpdate).
				Apply(f.Where).
				Update(); err != nil && err != pg.ErrNoRows {
				return handleInsertAndUpdateError(err)
			}
		}

		if err := tx.
			Model(&items).
			Context(ctx).
			Apply(f.Where).
			Select(); err != nil && err != pg.ErrNoRows {
			return handleInsertAndUpdateError(err)
		}

		qualificationIDs := make([]int, len(items))
		for index, item := range items {
			qualificationIDs[index] = item.ID
		}

		if len(qualificationIDs) > 0 {
			if len(input.DissociateProfession) > 0 {
				_, err := tx.
					Model(&models.QualificationToProfession{}).
					Where(sqlutils.BuildConditionArray("profession_id"), pg.Array(input.DissociateProfession)).
					Where(sqlutils.BuildConditionArray("qualification_id"), pg.Array(qualificationIDs)).
					Delete()
				if err != nil {
					return handleInsertAndUpdateError(err)
				}
			}

			if len(input.AssociateProfession) > 0 {
				if err := repo.associateQualificationWithProfession(tx, qualificationIDs, input.AssociateProfession); err != nil {
					return handleInsertAndUpdateError(err)
				}
			}
		}

		return nil
	})
	return items, err
}

func (repo *pgRepository) Delete(ctx context.Context, f *models.QualificationFilter) ([]*models.Qualification, error) {
	items := []*models.Qualification{}
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

func (repo *pgRepository) Fetch(ctx context.Context, cfg *qualification.FetchConfig) ([]*models.Qualification, int, error) {
	var err error
	items := []*models.Qualification{}
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

func (repo *pgRepository) GetSimilar(ctx context.Context, cfg *qualification.GetSimilarConfig) ([]*models.Qualification, int, error) {
	var err error
	subquery := repo.
		Model(&models.QualificationToProfession{}).
		Context(ctx).
		Where(sqlutils.BuildConditionEquals("qualification_id"), cfg.QualificationID).
		Column("profession_id")
	qualificationIDs := []int{}
	err = repo.
		Model(&models.QualificationToProfession{}).
		Context(ctx).
		Column("qualification_id").
		With("prof", subquery).
		Where(sqlutils.BuildConditionIn("profession_id"), pg.Safe("SELECT profession_id FROM prof")).
		Where(sqlutils.BuildConditionNEQ("qualification_id"), cfg.QualificationID).
		Select(&qualificationIDs)
	if err != nil {
		return nil, 0, errorutils.Wrap(err, messageFailedToFetchModel)
	}

	return repo.Fetch(ctx, &qualification.FetchConfig{
		Sort:   cfg.Sort,
		Limit:  cfg.Limit,
		Offset: cfg.Offset,
		Filter: &models.QualificationFilter{
			ID: qualificationIDs,
		},
		Count: cfg.Count,
	})
}

func (repo *pgRepository) associateQualificationWithProfession(tx *pg.Tx, qualificationIDs, professionIDs []int) error {
	toInsert := []*models.QualificationToProfession{}
	for _, professionID := range professionIDs {
		for _, qualificationID := range qualificationIDs {
			toInsert = append(toInsert, &models.QualificationToProfession{
				ProfessionID:    professionID,
				QualificationID: qualificationID,
			})
		}
	}
	_, err := tx.Model(&toInsert).Insert()
	return err
}

func handleInsertAndUpdateError(err error) error {
	if strings.Contains(err.Error(), "name") {
		return errorutils.Wrap(err, messageNameIsAlreadyTaken)
	} else if strings.Contains(err.Error(), "code") || strings.Contains(err.Error(), "slug") {
		return errorutils.Wrap(err, messageCodeIsAlreadyTaken)
	}
	return errorutils.Wrap(err, messageFailedToSaveModel)
}
