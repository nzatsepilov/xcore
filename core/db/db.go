package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"xcore/config"
)

type DB struct {
	*gorm.DB
}

func New(c *config.DBConfig) (*DB, error) {
	s := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v", c.Host, c.Port, c.User, c.Password, c.DBName)
	db, err := gorm.Open("postgres", s)
	if err != nil {
		return nil, err
	}
	return &DB{DB: db}, nil
}
