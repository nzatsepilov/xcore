package config

import "xcore/core/models"

type RealmConfig struct {
	ID              byte
	Name            string
	Address         string
	IsLocked        bool
	Type            models.RealmType
	Flag            models.RealmFlag
	Timezone        models.RealmTimezone
	Population      models.RealmPopulation
	CharactersCount byte
	Version         string
}
