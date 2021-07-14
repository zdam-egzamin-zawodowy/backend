package dataloader

import (
	"context"
	"time"

	"github.com/zdam-egzamin-zawodowy/backend/internal/profession"

	"github.com/zdam-egzamin-zawodowy/backend/internal/model"
	"github.com/zdam-egzamin-zawodowy/backend/internal/qualification"
)

const (
	wait = 2 * time.Millisecond
)

type Config struct {
	ProfessionRepo    profession.Repository
	QualificationRepo qualification.Repository
}

type DataLoader struct {
	QualificationByID            *QualificationLoader
	QualificationsByProfessionID *QualificationSliceByProfessionIDLoader
}

func New(cfg Config) *DataLoader {
	return &DataLoader{
		QualificationByID: NewQualificationLoader(QualificationLoaderConfig{
			Wait: wait,
			Fetch: func(ids []int) ([]*model.Qualification, []error) {
				qualificationsNotInOrder, _, err := cfg.QualificationRepo.Fetch(context.Background(), &qualification.FetchConfig{
					Filter: &model.QualificationFilter{
						ID: ids,
					},
					Count: false,
				})
				if err != nil {
					return nil, []error{err}
				}
				qualificationByID := make(map[int]*model.Qualification)
				for _, qualification := range qualificationsNotInOrder {
					qualificationByID[qualification.ID] = qualification
				}
				qualifications := make([]*model.Qualification, len(ids))
				for i, id := range ids {
					qualifications[i] = qualificationByID[id]
				}
				return qualifications, nil
			},
		}),
		QualificationsByProfessionID: NewQualificationSliceByProfessionIDLoader(QualificationSliceByProfessionIDLoaderConfig{
			Wait: wait,
			Fetch: func(ids []int) ([][]*model.Qualification, []error) {
				m, err := cfg.ProfessionRepo.GetAssociatedQualifications(context.Background(), ids...)
				if err != nil {
					return nil, []error{err}
				}

				qualifications := make([][]*model.Qualification, len(ids))

				for i, id := range ids {
					qualifications[i] = m[id]
				}

				return qualifications, nil
			},
		}),
	}
}
