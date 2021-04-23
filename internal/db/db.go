package db

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	envutils "github.com/zdam-egzamin-zawodowy/backend/pkg/utils/env"

	gopglogrusquerylogger "github.com/Kichiyaki/go-pg-logrus-query-logger/v10"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-password/password"
	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
)

var log = logrus.WithField("package", "internal/db")

type Config struct {
	LogQueries bool
}

func init() {
	orm.RegisterTable((*models.QualificationToProfession)(nil))
}

func New(cfg *Config) (*pg.DB, error) {
	db := pg.Connect(prepareOptions())

	if cfg != nil {
		if cfg.LogQueries {
			db.AddQueryHook(gopglogrusquerylogger.QueryLogger{
				Entry: log,
			})
		}
	}

	if err := createSchema(db); err != nil {
		return nil, err
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
		modelsToCreate := []interface{}{
			(*models.User)(nil),
			(*models.Profession)(nil),
			(*models.Qualification)(nil),
			(*models.QualificationToProfession)(nil),
			(*models.Question)(nil),
		}

		for _, model := range modelsToCreate {
			err := tx.Model(model).CreateTable(&orm.CreateTableOptions{
				IfNotExists:   true,
				FKConstraints: true,
			})
			if err != nil {
				return errors.Wrap(err, "createSchema")
			}
		}

		total, err := tx.Model(modelsToCreate[0]).Where("role = ?", models.RoleAdmin).Count()
		if err != nil {
			return errors.Wrap(err, "createSchema")
		}
		if total == 0 {
			activated := true
			pswd, err := password.Generate(15, 4, 2, true, false)
			if err != nil {
				return errors.Wrap(err, "createSchema")
			}
			email := "admin@admin.com"
			_, err = tx.
				Model(&models.User{
					DisplayName: "admin",
					Email:       email,
					Role:        models.RoleAdmin,
					Activated:   &activated,
					Password:    pswd,
				}).
				Insert()
			if err != nil {
				return errors.Wrap(err, "createSchema")
			}
			log.
				WithField("email", email).
				WithField("password", pswd).
				Info("Admin account has been created")
		}

		return nil
	})
}
