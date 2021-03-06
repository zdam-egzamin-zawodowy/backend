package mode

import (
	"os"
)

const (
	EnvKey          = "MODE"
	DevelopmentMode = "development"
	ProductionMode  = "production"
	TestMode        = "test"
)

var mode = DevelopmentMode

func init() {
	Set(os.Getenv(EnvKey))
}

func Set(value string) {
	if value == "" {
		value = DevelopmentMode
	}

	switch value {
	case DevelopmentMode,
		ProductionMode,
		TestMode:
		mode = value
	default:
		panic("unknown mode: " + value)
	}
}

func Get() string {
	return mode
}
