package internal

import (
	"context"
	"github.com/Kichiyaki/gopgutil/v10"
	"strings"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var _ pg.BeforeInsertHook = (*User)(nil)

type User struct {
	tableName struct{} `pg:"alias:user"`

	ID          int       `json:"id" pg:",pk" xml:"id" gqlgen:"id"`
	DisplayName string    `json:"displayName" pg:",use_zero,notnull" xml:"displayName" gqlgen:"displayName"`
	Password    string    `json:"-" gqlgen:"-" xml:"password"`
	Email       string    `json:"email" pg:",unique" xml:"email" gqlgen:"email"`
	CreatedAt   time.Time `json:"createdAt" pg:"default:now()" xml:"createdAt" gqlgen:"createdAt"`
	Role        Role      `json:"role" xml:"role" gqlgen:"role"`
	Activated   *bool     `json:"activated" pg:"default:false,use_zero" xml:"activated" gqlgen:"activated"`
}

func (u *User) BeforeInsert(ctx context.Context) (context.Context, error) {
	u.CreatedAt = time.Now()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return ctx, err
	}
	u.Password = string(hashedPassword)

	return ctx, nil
}

func (u *User) CompareHashAndPassword(password string) error {
	if password == u.Password {
		return nil
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return errors.Wrap(err, "CompareHashAndPassword")
	}
	return nil
}

type UserInput struct {
	DisplayName *string `json:"displayName" xml:"displayName" gqlgen:"displayName"`
	Password    *string `json:"password" xml:"password" gqlgen:"password"`
	Email       *string `json:"email" xml:"email" gqlgen:"email"`
	Role        *Role   `json:"role" xml:"role" gqlgen:"role"`
	Activated   *bool   `json:"activated" xml:"activated" gqlgen:"activated"`
}

func (input *UserInput) IsEmpty() bool {
	return input == nil &&
		input.DisplayName == nil &&
		input.Password == nil &&
		input.Email == nil &&
		input.Role == nil &&
		input.Activated == nil
}

func (input *UserInput) Sanitize() *UserInput {
	if input.DisplayName != nil {
		*input.DisplayName = strings.TrimSpace(*input.DisplayName)
	}
	if input.Password != nil {
		*input.Password = strings.TrimSpace(*input.Password)
	}
	if input.Email != nil {
		*input.Email = strings.ToLower(strings.TrimSpace(*input.Email))
	}

	return input
}

func (input *UserInput) ToUser() *User {
	u := &User{
		Activated: input.Activated,
	}
	if input.DisplayName != nil {
		u.DisplayName = *input.DisplayName
	}
	if input.Password != nil {
		u.Password = *input.Password
	}
	if input.Email != nil {
		u.Email = *input.Email
	}
	if input.Role != nil {
		u.Role = *input.Role
	}
	return u
}

func (input *UserInput) ApplyUpdate(q *orm.Query) (*orm.Query, error) {
	if !input.IsEmpty() {
		if input.DisplayName != nil {
			q = q.Set(gopgutil.BuildConditionEquals("display_name"), *input.DisplayName)
		}

		if input.Password != nil {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
			if err != nil {
				return q, err
			}
			q = q.Set(gopgutil.BuildConditionEquals("password"), string(hashedPassword))
		}

		if input.Email != nil {
			q = q.Set(gopgutil.BuildConditionEquals("email"), *input.Email)
		}

		if input.Role != nil {
			q = q.Set(gopgutil.BuildConditionEquals("role"), *input.Role)
		}

		if input.Activated != nil {
			q = q.Set(gopgutil.BuildConditionEquals("activated"), *input.Activated)
		}
	}

	return q, nil
}

type UserFilterOr struct {
	DisplayNameIEQ   string `json:"displayNameIEQ" xml:"displayNameIEQ" gqlgen:"displayNameIEQ"`
	DisplayNameMATCH string `json:"displayNameMATCH" xml:"displayNameMATCH" gqlgen:"displayNameMATCH"`

	EmailIEQ   string `json:"emailIEQ" xml:"emailIEQ" gqlgen:"emailIEQ"`
	EmailMATCH string `json:"emailMATCH" xml:"emailMATCH" gqlgen:"emailMATCH"`
}

func (f *UserFilterOr) WhereWithAlias(q *orm.Query, alias string) *orm.Query {
	if f == nil {
		return q
	}

	q = q.WhereGroup(func(q *orm.Query) (*orm.Query, error) {
		if !isZero(f.DisplayNameMATCH) {
			q = q.WhereOr(
				gopgutil.BuildConditionMatch("?"),
				gopgutil.AddAliasToColumnName("display_name", alias),
				f.DisplayNameMATCH,
			)
		}
		if !isZero(f.DisplayNameIEQ) {
			q = q.WhereOr(
				gopgutil.BuildConditionIEQ("?"),
				gopgutil.AddAliasToColumnName("display_name", alias),
				f.DisplayNameIEQ,
			)
		}

		if !isZero(f.EmailMATCH) {
			q = q.WhereOr(gopgutil.BuildConditionMatch("?"), gopgutil.AddAliasToColumnName("email", alias), f.EmailMATCH)
		}
		if !isZero(f.EmailIEQ) {
			q = q.WhereOr(gopgutil.BuildConditionIEQ("?"), gopgutil.AddAliasToColumnName("email", alias), f.EmailIEQ)
		}

		return q, nil
	})
	return q
}

type UserFilter struct {
	ID    []int `json:"id" xml:"id" gqlgen:"id"`
	IDNEQ []int `json:"idNEQ" xml:"idNEQ" gqlgen:"idNEQ"`

	Activated *bool `json:"activated" xml:"activated" gqlgen:"activated"`

	DisplayName      []string `json:"displayName" xml:"displayName" gqlgen:"displayName"`
	DisplayNameNEQ   []string `json:"displayNameNEQ" xml:"displayNameNEQ" gqlgen:"displayNameNEQ"`
	DisplayNameIEQ   string   `json:"displayNameIEQ" xml:"displayNameIEQ" gqlgen:"displayNameIEQ"`
	DisplayNameMATCH string   `json:"displayNameMATCH" xml:"displayNameMATCH" gqlgen:"displayNameMATCH"`

	Email      []string `json:"email" xml:"email" gqlgen:"email"`
	EmailNEQ   []string `json:"emailNEQ" xml:"emailNEQ" gqlgen:"emailNEQ"`
	EmailIEQ   string   `json:"emailIEQ" xml:"emailIEQ" gqlgen:"emailIEQ"`
	EmailMATCH string   `json:"emailMATCH" xml:"emailMATCH" gqlgen:"emailMATCH"`

	Role    []Role `json:"role" xml:"role" gqlgen:"role"`
	RoleNEQ []Role `json:"roleNEQ" xml:"roleNEQ" gqlgen:"roleNEQ"`

	CreatedAt    time.Time `json:"createdAt" xml:"createdAt" gqlgen:"createdAt"`
	CreatedAtGT  time.Time `json:"createdAtGT" xml:"createdAtGT" gqlgen:"createdAtGT"`
	CreatedAtGTE time.Time `json:"createdAtGTE" xml:"createdAtGTE" gqlgen:"createdAtGTE"`
	CreatedAtLT  time.Time `json:"createdAtLT" xml:"createdAtLT" gqlgen:"createdAtLT"`
	CreatedAtLTE time.Time `json:"createdAtLTE" xml:"createdAtLTE" gqlgen:"createdAtLTE"`

	Or *UserFilterOr
}

func (f *UserFilter) WhereWithAlias(q *orm.Query, alias string) (*orm.Query, error) {
	if f == nil {
		return q, nil
	}

	if !isZero(f.ID) {
		q = q.Where(gopgutil.BuildConditionArray("?"), gopgutil.AddAliasToColumnName("id", alias), pg.Array(f.ID))
	}
	if !isZero(f.IDNEQ) {
		q = q.Where(gopgutil.BuildConditionNotInArray("?"), gopgutil.AddAliasToColumnName("id", alias), pg.Array(f.IDNEQ))
	}

	if !isZero(f.Activated) {
		q = q.Where(gopgutil.BuildConditionEquals("?"), gopgutil.AddAliasToColumnName("activated", alias), f.Activated)
	}

	if !isZero(f.DisplayName) {
		q = q.Where(gopgutil.BuildConditionArray("?"), gopgutil.AddAliasToColumnName("display_name", alias), pg.Array(f.DisplayName))
	}
	if !isZero(f.DisplayNameNEQ) {
		q = q.Where(gopgutil.BuildConditionNotInArray("?"), gopgutil.AddAliasToColumnName("display_name", alias), pg.Array(f.DisplayNameNEQ))
	}
	if !isZero(f.DisplayNameMATCH) {
		q = q.Where(gopgutil.BuildConditionMatch("?"), gopgutil.AddAliasToColumnName("display_name", alias), f.DisplayNameMATCH)
	}
	if !isZero(f.DisplayNameIEQ) {
		q = q.Where(gopgutil.BuildConditionIEQ("?"), gopgutil.AddAliasToColumnName("display_name", alias), f.DisplayNameIEQ)
	}

	if !isZero(f.Email) {
		q = q.Where(gopgutil.BuildConditionArray("?"), gopgutil.AddAliasToColumnName("email", alias), pg.Array(f.Email))
	}
	if !isZero(f.EmailNEQ) {
		q = q.Where(gopgutil.BuildConditionNotInArray("?"), gopgutil.AddAliasToColumnName("email", alias), pg.Array(f.EmailNEQ))
	}
	if !isZero(f.EmailMATCH) {
		q = q.Where(gopgutil.BuildConditionMatch("?"), gopgutil.AddAliasToColumnName("email", alias), f.EmailMATCH)
	}
	if !isZero(f.EmailIEQ) {
		q = q.Where(gopgutil.BuildConditionIEQ("?"), gopgutil.AddAliasToColumnName("email", alias), f.EmailIEQ)
	}

	if !isZero(f.CreatedAt) {
		q = q.Where(gopgutil.BuildConditionEquals("?"), gopgutil.AddAliasToColumnName("created_at", alias), f.CreatedAt)
	}
	if !isZero(f.CreatedAtGT) {
		q = q.Where(gopgutil.BuildConditionGT("?"), gopgutil.AddAliasToColumnName("created_at", alias), f.CreatedAtGT)
	}
	if !isZero(f.CreatedAtGTE) {
		q = q.Where(gopgutil.BuildConditionGTE("?"), gopgutil.AddAliasToColumnName("created_at", alias), f.CreatedAtGTE)
	}
	if !isZero(f.CreatedAtLT) {
		q = q.Where(gopgutil.BuildConditionLT("?"), gopgutil.AddAliasToColumnName("created_at", alias), f.CreatedAtLT)
	}
	if !isZero(f.CreatedAtLTE) {
		q = q.Where(gopgutil.BuildConditionLTE("?"), gopgutil.AddAliasToColumnName("created_at", alias), f.CreatedAtLTE)
	}

	if f.Or != nil {
		q = f.Or.WhereWithAlias(q, alias)
	}

	return q, nil
}

func (f *UserFilter) Where(q *orm.Query) (*orm.Query, error) {
	return f.WhereWithAlias(q, "user")
}
