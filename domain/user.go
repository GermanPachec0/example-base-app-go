package domain

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64     `json:"id"`
	Uuid      uuid.UUID `json:"uuid"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreateAt  string    `json:"create_at"`
	UpdateAt  string    `json:"update_at"`
	Type      string    `json:"type"`
}

func (u *User) ValidatePassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw)) == nil
}

func NewUser(e, p, fn, ln string) (*User, error) {
	encp, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	uid := uuid.New()
	return &User{
		Uuid:      uid,
		Email:     e,
		Password:  string(encp),
		FirstName: fn,
		LastName:  ln,
	}, nil
}

func (u *User) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Email, validation.Required, validation.Length(3, 20)),
		validation.Field(&u.Password, validation.Required, validation.Length(4, 100)),
		validation.Field(&u.FirstName, validation.Required, validation.Length(4, 9)),
		validation.Field(&u.LastName, validation.Required, validation.Length(4, 9)),
	)
}

type UserRepository interface {
	GetById(ctx context.Context, uuid uuid.UUID) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	GetAll(ctx context.Context) ([]User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context) (*User, error)
}
