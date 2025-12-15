package dbmodel

import "time"

type User struct {
	Id               int
	Email            string
	RoleId           int
	NativeLanguageId int
	CreatedAt        time.Time
	UpdatedAt        *time.Time
	DeletedAt        *time.Time
}
