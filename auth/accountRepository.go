package auth

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/jinzhu/gorm"
	"strings"
	"xcore/config"
	"xcore/core/db"
	"xcore/core/models"
)

type AccountRepository interface {
	CreateAccount(name string, password string) error
	GetAccountWithName(name string) (*models.Account, error)
	SaveAccount(a *models.Account) error
}

type accountRepository struct {
	db *db.DB
}

func NewAccountRepository(c *config.Config, db *db.DB) (AccountRepository, error) {
	r := &accountRepository{
		db: db,
	}
	if err := r.init(c); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *accountRepository) init(c *config.Config) error {
	if err := r.migrate(); err != nil {
		return err
	}
	if err := r.createDevAccounts(c.DevAccounts); err != nil {
		return err
	}
	return nil
}

func (r *accountRepository) migrate() error {
	return r.db.AutoMigrate(&models.Account{}).Error
}

func (r *accountRepository) createDevAccounts(accounts []*config.DevAccount) error {
	for _, a := range accounts {
		nameUpper := strings.ToUpper(a.Name)
		exists, err := r.HasAccountWithName(nameUpper)
		if err != nil {
			return nil
		}

		if exists {
			continue
		}

		if err := r.CreateAccount(a.Name, a.Password); err != nil {
			return err
		}
	}
	return nil
}

func (r *accountRepository) CreateAccount(name string, password string) error {
	nameUpper := strings.ToUpper(name)
	passUpper := strings.ToUpper(password)

	passHash := sha1.New()
	passHash.Write([]byte(nameUpper))
	passHash.Write([]byte(":"))
	passHash.Write([]byte(passUpper))
	hashHex := hex.EncodeToString(passHash.Sum(nil))
	acc := models.Account{Name: nameUpper, PasswordHash: hashHex}
	return r.db.Save(&acc).Error
}

func (r *accountRepository) HasAccountWithName(name string) (bool, error) {
	var count int
	r.db.Where("name = ?", name).Find(&models.Account{}).Count(&count)
	if err := r.db.Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *accountRepository) GetAccountWithName(name string) (*models.Account, error) {
	var acc models.Account
	err := r.db.Where("name = ?", strings.ToUpper(name)).First(&acc).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &acc, err
}

func (r *accountRepository) SaveAccount(a *models.Account) error {
	return r.db.Save(*a).Error
}
