package common

import "time"

type Action struct {
	ID int

	Trigger string
	Image   string
	EnvVars string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
