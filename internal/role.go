package internal

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"strconv"
	"strings"
)

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

func (role *Role) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return errors.New("enums must be strings")
	}

	*role = Role(strings.ToLower(str))
	if !role.IsValid() {
		return errors.Errorf("%s is not a valid Role", str)
	}
	return nil
}

func (role Role) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(role.String()))
}
