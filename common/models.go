package common

import (
	"time"
)

type Action struct {
	ID int

	Trigger string
	Image   string
	EnvVars string

	Match string `gorm:"-"`

	AccountID int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Account struct {
	ID int

	Name string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type User struct {
	ID int

	Email string
	Token string

	AccountID int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Listener struct {
	ID int

	Class string

	AccountID int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (u *User) FindOrCreateFromToken(token string) {

}
