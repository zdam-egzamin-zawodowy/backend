package errorutil

import (
	"github.com/pkg/errors"
	"github.com/zdam-egzamin-zawodowy/backend/pkg/mode"
)

func Wrap(details error, message string) error {
	if mode.Get() == mode.ProductionMode {
		return errors.New(message)
	}
	return errors.Wrap(details, message)
}

func Wrapf(details error, message string, args ...interface{}) error {
	if mode.Get() == mode.ProductionMode {
		return errors.Errorf(message, args...)
	}
	return errors.Wrapf(details, message, args...)
}
