package postgres

import (
	"context"
	"github.com/Kichiyaki/goutil/envutil"
	"github.com/sirupsen/logrus"

	"github.com/Kichiyaki/go-pg-logrus-query-logger/v10"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-password/password"

	"github.com/zdam-egzamin-zawodowy/backend/internal/model"
)

var log = logrus.WithField("package", "internal/postgres")

type Config struct {
	LogQueries bool
}

func init() {
	orm.RegisterTable((*model.QualificationToProfession)(nil))
}

func Connect(cfg *Config) (*pg.DB, error) {
	db := pg.Connect(prepareOptions())

	if cfg != nil {
		if cfg.LogQueries {
			db.AddQueryHook(querylogger.Logger{
				Log:            log,
				MaxQueryLength: 5000,
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
		User:     envutil.GetenvString("DB_USER"),
		Password: envutil.GetenvString("DB_PASSWORD"),
		Database: envutil.GetenvString("DB_NAME"),
		Addr:     envutil.GetenvString("DB_HOST") + ":" + envutil.GetenvString("DB_PORT"),
		PoolSize: envutil.GetenvInt("DB_POOL_SIZE"),
	}
}

func createSchema(db *pg.DB) error {
	return db.RunInTransaction(context.Background(), func(tx *pg.Tx) error {
		modelsToCreate := []interface{}{
			(*model.User)(nil),
			(*model.Profession)(nil),
			(*model.Qualification)(nil),
			(*model.QualificationToProfession)(nil),
			(*model.Question)(nil),
		}

		for _, model := range modelsToCreate {
			err := tx.Model(model).CreateTable(&orm.CreateTableOptions{
				IfNotExists:   true,
				FKConstraints: true,
			})
			if err != nil {
				return errors.Wrap(err, "couldn't create the table")
			}
		}

		total, err := tx.Model(modelsToCreate[0]).Where("role = ?", model.RoleAdmin).Count()
		if err != nil {
			return errors.Wrap(err, "couldn't count admins")
		}
		if total == 0 {
			activated := true
			pswd, err := password.Generate(15, 6, 0, true, false)
			if err != nil {
				return errors.Wrap(err, "couldn't generate a password for the new admin account")
			}
			email := "admin@admin.com"
			_, err = tx.
				Model(&model.User{
					DisplayName: "admin",
					Email:       email,
					Role:        model.RoleAdmin,
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
