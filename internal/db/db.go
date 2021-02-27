package db

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	envutils "github.com/zdam-egzamin-zawodowy/backend/pkg/utils/env"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/pkg/errors"
	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
)

const (
	extensions = `
		CREATE EXTENSION IF NOT EXISTS tsm_system_rows;
	`
)

func init() {
	orm.RegisterTable((*models.QualificationToProfession)(nil))
}

type Config struct {
	DebugHook bool
}

func New(cfg *Config) (*pg.DB, error) {
	db := pg.Connect(prepareOptions())
	if err := createSchema(db); err != nil {
		return nil, err
	}

	if cfg != nil {
		if cfg.DebugHook {
			db.AddQueryHook(DebugHook{
				Entry: logrus.WithField("package", "internal/db"),
			})
		}
	}

	return db, nil
}

func prepareOptions() *pg.Options {
	return &pg.Options{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
		Addr:     os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
		PoolSize: envutils.GetenvInt("DB_POOL_SIZE"),
	}
}

func createSchema(db *pg.DB) error {
	return db.RunInTransaction(context.Background(), func(tx *pg.Tx) error {
		if _, err := tx.Exec(extensions); err != nil {
			return errors.Wrap(err, "createSchema")
		}

		models := []interface{}{
			(*models.User)(nil),
			(*models.Profession)(nil),
			(*models.Qualification)(nil),
			(*models.QualificationToProfession)(nil),
			(*models.Question)(nil),
		}

		for _, model := range models {
			err := tx.Model(model).CreateTable(&orm.CreateTableOptions{
				IfNotExists: true,
			})
			if err != nil {
				return errors.Wrap(err, "createSchema")
			}
		}

		return nil
	})
}
