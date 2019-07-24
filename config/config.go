package config

import "xcore/core/models"

type Config struct {
	AuthServerAddress  string
	WorldServerAddress string

	DBConfig *DBConfig

	DevAccounts []*DevAccount
	Realms      []*RealmConfig
}

func Current() *Config {
	return &Config{
		AuthServerAddress:  "0.0.0.0:3724",
		WorldServerAddress: "192.168.1.105:8085",

		DBConfig: &DBConfig{
			Host:     "127.0.0.1",
			Port:     "5432",
			User:     "xcore",
			Password: "xcore",
			DBName:   "xcore",
		},
		DevAccounts: []*DevAccount{
			{
				Name:     "dev",
				Password: "123",
			},
		},
		Realms: []*RealmConfig{
			{
				ID:              1,
				Name:            "Test 1",
				Address:         "192.168.1.105:8085",
				IsLocked:        false,
				Type:            models.RealmTypeNormal,
				Flag:            models.RealmFlagNew,
				Timezone:        models.RealmTimezoneDevelopment,
				Population:      models.RealmPopulationLow,
				CharactersCount: 0,
				Version:         "2.4.3.8606",
			},
		},
	}
}
