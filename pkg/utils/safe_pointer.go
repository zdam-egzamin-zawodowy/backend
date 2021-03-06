package utils

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
