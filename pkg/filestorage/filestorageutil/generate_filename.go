package filestorageutil

import (
	"strings"

	"github.com/google/uuid"
)

func GenerateFilename(ext string) string {
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	return uuid.New().String() + ext
}
