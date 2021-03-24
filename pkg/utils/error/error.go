package errorutils

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/zdam-egzamin-zawodowy/backend/pkg/mode"
)

func Wrap(details error, message string) error {
	return Wrapf(details, message)
}

func Wrapf(details error, message string, args ...interface{}) error {
	if mode.Get() != mode.ProductionMode {
		return errors.Wrapf(details, message, args...)
	}
	return fmt.Errorf(message, args...)
}
