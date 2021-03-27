package errorutils

import (
	"github.com/pkg/errors"
	"github.com/zdam-egzamin-zawodowy/backend/pkg/mode"
)

func Wrap(details error, message string) error {
	return errors.Wrap(details, message)
}

func Wrapf(details error, message string, args ...interface{}) error {
	if mode.Get() != mode.ProductionMode {
		return errors.Wrapf(details, message, args...)
	}
	return errors.Errorf(message, args...)
}
