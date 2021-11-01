package postgres

import (
	"context"
	"github.com/Kichiyaki/gopgutil/v10"
	"github.com/zdam-egzamin-zawodowy/backend/internal"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
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

func applyUserInputUpdates(input internal.UserInput) func(q *orm.Query) (*orm.Query, error) {
	return func(q *orm.Query) (*orm.Query, error) {
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

		return q, nil
	}
}

func applyUserFilterOr(f internal.UserFilterOr, alias string) func(q *orm.Query) (*orm.Query, error) {
	return func(q *orm.Query) (*orm.Query, error) {
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

		return q, nil
	}
}

func applyUserFilter(f internal.UserFilter, alias string) func(q *orm.Query) (*orm.Query, error) {
	return func(q *orm.Query) (*orm.Query, error) {
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

		q = q.Apply(applyUserFilterOr(f.Or, alias))

		return q, nil
	}
}
