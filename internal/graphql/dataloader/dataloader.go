package dataloader

import (
	"context"
	"time"

	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
	"github.com/zdam-egzamin-zawodowy/backend/internal/qualification"
)

type Config struct {
	QualificationRepo qualification.Repository
}

type DataLoader struct {
	QualificationByID *QualificationLoader
}

func New(cfg Config) *DataLoader {
	return &DataLoader{
		QualificationByID: NewQualificationLoader(QualificationLoaderConfig{
			Wait: 2 * time.Millisecond,
			Fetch: func(ids []int) ([]*models.Qualification, []error) {
				qualificationsNotInOrder, _, err := cfg.QualificationRepo.Fetch(context.Background(), &qualification.FetchConfig{
					Filter: &models.QualificationFilter{
						ID: ids,
					},
					Count: false,
				})
				if err != nil {
					return nil, []error{err}
				}
				qualificationByID := make(map[int]*models.Qualification)
				for _, qualification := range qualificationsNotInOrder {
					qualificationByID[qualification.ID] = qualification
				}
				qualifications := make([]*models.Qualification, len(ids))
				for i, id := range ids {
					qualifications[i] = qualificationByID[id]
				}
				return qualifications, nil
			},
		}),
	}
}
