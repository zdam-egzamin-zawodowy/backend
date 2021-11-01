package dataloader

import (
	"context"
	"github.com/zdam-egzamin-zawodowy/backend/internal"
	"time"

	"github.com/zdam-egzamin-zawodowy/backend/internal/profession"

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
			Fetch: func(ids []int) ([]*internal.Qualification, []error) {
				qualificationsNotInOrder, _, err := cfg.QualificationRepo.Fetch(context.Background(), &qualification.FetchConfig{
					Filter: &internal.QualificationFilter{
						ID: ids,
					},
					Count: false,
				})
				if err != nil {
					return nil, []error{err}
				}
				qualificationByID := make(map[int]*internal.Qualification)
				for _, qualification := range qualificationsNotInOrder {
					qualificationByID[qualification.ID] = qualification
				}
				qualifications := make([]*internal.Qualification, len(ids))
				for i, id := range ids {
					qualifications[i] = qualificationByID[id]
				}
				return qualifications, nil
			},
		}),
		QualificationsByProfessionID: NewQualificationSliceByProfessionIDLoader(QualificationSliceByProfessionIDLoaderConfig{
			Wait: wait,
			Fetch: func(ids []int) ([][]*internal.Qualification, []error) {
				m, err := cfg.ProfessionRepo.GetAssociatedQualifications(context.Background(), ids...)
				if err != nil {
					return nil, []error{err}
				}

				qualifications := make([][]*internal.Qualification, len(ids))

				for i, id := range ids {
					qualifications[i] = m[id]
				}

				return qualifications, nil
			},
		}),
	}
}
