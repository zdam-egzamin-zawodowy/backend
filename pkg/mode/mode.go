package mode

import (
	"github.com/zdam-egzamin-zawodowy/backend/pkg/util/envutil"
)

const (
	EnvKey          = "MODE"
	DevelopmentMode = "development"
	ProductionMode  = "production"
	TestMode        = "test"
)

var mode = DevelopmentMode

func init() {
	Set(envutil.GetenvString(EnvKey))
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
