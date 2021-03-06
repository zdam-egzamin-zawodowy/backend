package repository

import (
	"context"
	"fmt"
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
		return nil, fmt.Errorf("qualification/pg_repository: *pg.DB is required")
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
			if strings.Contains(err.Error(), "name") {
				return errorutils.Wrap(err, messageNameIsAlreadyTaken)
			} else if strings.Contains(err.Error(), "code") {
				return errorutils.Wrap(err, messageCodeIsAlreadyTaken)
			}
			return errorutils.Wrap(err, messageFailedToSaveModel)
		}

		for _, professionID := range input.AssociateProfession {
			tx.
				Model(&models.QualificationToProfession{
					QualificationID: item.ID,
					ProfessionID:    professionID,
				}).
				Insert()
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
		if _, err := tx.
			Model(&items).
			Context(ctx).
			Returning("*").
			Apply(input.ApplyUpdate).
			Apply(f.Where).
			Update(); err != nil && err != pg.ErrNoRows {
			if strings.Contains(err.Error(), "name") {
				return errorutils.Wrap(err, messageNameIsAlreadyTaken)
			} else if strings.Contains(err.Error(), "code") {
				return errorutils.Wrap(err, messageCodeIsAlreadyTaken)
			}
			return errorutils.Wrap(err, messageFailedToSaveModel)
		}

		qualificationIDs := make([]int, len(items))
		for index, item := range items {
			qualificationIDs[index] = item.ID
		}

		if len(qualificationIDs) > 0 {
			if len(input.DissociateProfession) > 0 {
				tx.
					Model(&models.QualificationToProfession{}).
					Where(sqlutils.BuildConditionArray("profession_id"), input.DissociateProfession).
					Where(sqlutils.BuildConditionArray("qualification_id"), qualificationIDs).
					Delete()
			}

			if len(input.AssociateProfession) > 0 {
				toInsert := []*models.QualificationToProfession{}
				for _, professionID := range input.AssociateProfession {
					for _, qualificationID := range qualificationIDs {
						toInsert = append(toInsert, &models.QualificationToProfession{
							ProfessionID:    professionID,
							QualificationID: qualificationID,
						})
					}
				}
				tx.Model(&toInsert).Insert()
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
