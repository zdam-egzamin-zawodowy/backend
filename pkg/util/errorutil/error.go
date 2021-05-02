package errorutil

import (
	"github.com/Kichiyaki/appmode"
	"github.com/pkg/errors"
)

func Wrap(details error, message string) error {
	if appmode.Equals(appmode.ProductionMode) {
		return errors.New(message)
	}
	return errors.Wrap(details, message)
}

func Wrapf(details error, message string, args ...interface{}) error {
	if appmode.Equals(appmode.ProductionMode) {
		return errors.Errorf(message, args...)
	}
	return errors.Wrapf(details, message, args...)
}
