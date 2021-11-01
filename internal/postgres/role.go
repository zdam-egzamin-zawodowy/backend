package postgres

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

func (role Role) IsValid() bool {
	switch role {
	case RoleAdmin,
		RoleUser:
		return true
	}
	return false
}

func (role Role) String() string {
	return string(role)
}
