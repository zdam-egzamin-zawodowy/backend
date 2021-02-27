package models

import (
	"context"
	"time"

	sqlutils "github.com/zdam-egzamin-zawodowy/backend/pkg/utils/sql"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/gosimple/slug"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	tableName struct{} `pg:"alias:user"`

	ID                            int       `json:"id" pg:",pk" xml:"id" gqlgen:"id"`
	Slug                          string    `json:"slug" pg:",unique" xml:"slug" gqlgen:"slug"`
	DisplayName                   string    `json:"displayName" pg:",use_zero" xml:"displayName" gqlgen:"displayName"`
	Password                      string    `json:"-" gqlgen:"-" xml:"password"`
	Email                         string    `json:"email" pg:",unique" xml:"email" gqlgen:"email"`
	CreatedAt                     time.Time `json:"createdAt" pg:"default:now()" xml:"createdAt" gqlgen:"createdAt"`
	Role                          Role      `json:"role" xml:"role" gqlgen:"role"`
	Activated                     *bool     `json:"activated" pg:"default:false,use_zero" xml:"activated" gqlgen:"activated"`
	ActivationToken               string    `json:"-" gqlgen:"-" xml:"activationToken"`
	ActivationTokenGeneratedAt    time.Time `json:"-" gqlgen:"-" pg:"default:now()" xml:"activationTokenGeneratedAt"`
	ResetPasswordToken            string    `json:"-" gqlgen:"-" xml:"resetPasswordToken"`
	ResetPasswordTokenGeneratedAt time.Time `json:"-" gqlgen:"-" pg:"default:now()" xml:"resetPasswordTokenGeneratedAt"`
}

func (u *User) BeforeInsert(ctx context.Context) (context.Context, error) {
	u.CreatedAt = time.Now()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return ctx, err
	}
	u.Password = string(hashedPassword)
	u.Slug = slug.Make(u.DisplayName)

	return ctx, nil
}

func (u *User) BeforeUpdate(ctx context.Context) (context.Context, error) {
	if cost, _ := bcrypt.Cost([]byte(u.Password)); u.Password != "" && cost == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return ctx, err
		}
		u.Password = string(hashedPassword)
	}

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

type UserFilter struct {
	ID    []int `json:"id" xml:"id" gqlgen:"id"`
	IDNEQ []int `json:"idNEQ" xml:"idNEQ" gqlgen:"idNEQ"`

	Slug    []string `json:"slug" xml:"slug" gqlgen:"slug"`
	SlugNEQ []string `json:"slugNEQ" xml:"slugNEQ" gqlgen:"slugNEQ"`

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
}

func (f *UserFilter) WhereWithAlias(q *orm.Query, alias string) (*orm.Query, error) {
	if f == nil {
		return q, nil
	}

	if !isZero(f.ID) {
		q = q.Where(sqlutils.BuildConditionArray(sqlutils.AddAliasToColumnName("id", alias)), pg.Array(f.ID))
	}
	if !isZero(f.IDNEQ) {
		q = q.Where(sqlutils.BuildConditionNotInArray(sqlutils.AddAliasToColumnName("id", alias)), pg.Array(f.IDNEQ))
	}

	if !isZero(f.Slug) {
		q = q.Where(sqlutils.BuildConditionArray(sqlutils.AddAliasToColumnName("slug", alias)), pg.Array(f.Slug))
	}
	if !isZero(f.SlugNEQ) {
		q = q.Where(sqlutils.BuildConditionNotInArray(sqlutils.AddAliasToColumnName("slug", alias)), pg.Array(f.SlugNEQ))
	}

	if !isZero(f.Activated) {
		q = q.Where(sqlutils.BuildConditionEquals(sqlutils.AddAliasToColumnName("activated", alias)), f.Activated)
	}

	if !isZero(f.DisplayName) {
		q = q.Where(sqlutils.BuildConditionArray(sqlutils.AddAliasToColumnName("display_name", alias)), pg.Array(f.DisplayName))
	}
	if !isZero(f.DisplayNameNEQ) {
		q = q.Where(sqlutils.BuildConditionNotInArray(sqlutils.AddAliasToColumnName("display_name", alias)), pg.Array(f.DisplayNameNEQ))
	}
	if !isZero(f.DisplayNameMATCH) {
		q = q.Where(sqlutils.BuildConditionMatch(sqlutils.AddAliasToColumnName("display_name", alias)), f.DisplayNameMATCH)
	}
	if !isZero(f.DisplayNameIEQ) {
		q = q.Where(sqlutils.BuildConditionIEQ(sqlutils.AddAliasToColumnName("display_name", alias)), f.DisplayNameIEQ)
	}

	if !isZero(f.Email) {
		q = q.Where(sqlutils.BuildConditionArray(sqlutils.AddAliasToColumnName("email", alias)), pg.Array(f.Email))
	}
	if !isZero(f.EmailNEQ) {
		q = q.Where(sqlutils.BuildConditionNotInArray(sqlutils.AddAliasToColumnName("email", alias)), pg.Array(f.EmailNEQ))
	}
	if !isZero(f.EmailMATCH) {
		q = q.Where(sqlutils.BuildConditionMatch(sqlutils.AddAliasToColumnName("email", alias)), f.EmailMATCH)
	}
	if !isZero(f.EmailIEQ) {
		q = q.Where(sqlutils.BuildConditionIEQ(sqlutils.AddAliasToColumnName("email", alias)), f.EmailIEQ)
	}

	if !isZero(f.CreatedAt) {
		q = q.Where(sqlutils.BuildConditionEquals(sqlutils.AddAliasToColumnName("created_at", alias)), f.CreatedAt)
	}
	if !isZero(f.CreatedAtGT) {
		q = q.Where(sqlutils.BuildConditionGT(sqlutils.AddAliasToColumnName("created_at", alias)), f.CreatedAtGT)
	}
	if !isZero(f.CreatedAtGTE) {
		q = q.Where(sqlutils.BuildConditionGTE(sqlutils.AddAliasToColumnName("created_at", alias)), f.CreatedAtGTE)
	}
	if !isZero(f.CreatedAtLT) {
		q = q.Where(sqlutils.BuildConditionLT(sqlutils.AddAliasToColumnName("created_at", alias)), f.CreatedAtLT)
	}
	if !isZero(f.CreatedAtLTE) {
		q = q.Where(sqlutils.BuildConditionLTE(sqlutils.AddAliasToColumnName("created_at", alias)), f.CreatedAtLTE)
	}

	return q, nil
}

func (f *UserFilter) Where(q *orm.Query) (*orm.Query, error) {
	return f.WhereWithAlias(q, "user")
}
