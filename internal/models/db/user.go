package dbmodel

import (
	"database/sql"
	"time"
)

type User struct {
	Id               int
	Email            string
	RoleId           int
	NativeLanguageId int
	CreatedAt        time.Time
	UpdatedAt        sql.NullTime
	DeletedAt        sql.NullTime
}
