package models

import (
	"database/sql"
	"github.com/jinzhu/gorm"
)

type Account struct {
	gorm.Model
	Name         string `gorm:"size:16; unique;"`
	PasswordHash string
	PasswordKey  sql.NullString
	SessionKey   sql.NullString
}
