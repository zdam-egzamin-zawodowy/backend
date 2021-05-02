package safepointer

func SafeBoolPointer(v *bool, def bool) bool {
	if v == nil {
		return def
	}
	return *v
}

func SafeIntPointer(s *int, def int) int {
	if s == nil {
		return def
	}
	return *s
}

func SafeStringPointer(s *string, def string) string {
	if s == nil {
		return def
	}
	return *s
}
