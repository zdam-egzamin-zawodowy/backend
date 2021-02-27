package db

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/sirupsen/logrus"
)

type DebugHook struct {
	Entry *logrus.Entry
}

var _ pg.QueryHook = (*DebugHook)(nil)

func (logger DebugHook) BeforeQuery(ctx context.Context, evt *pg.QueryEvent) (context.Context, error) {
	q, err := evt.FormattedQuery()
	if err != nil {
		return nil, err
	}

	if evt.Err != nil {
		logger.Entry.Errorf("%s executing a query:\n%s\n", evt.Err, q)
	} else {
		logger.Entry.Info(string(q))
	}

	return ctx, nil
}

func (DebugHook) AfterQuery(context.Context, *pg.QueryEvent) error {
	return nil
}
