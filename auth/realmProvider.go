package auth

import "xcore/config"

type realmProvider struct {
	realms []*config.RealmConfig
}

func NewRealmProvider(c *config.Config) *realmProvider {
	return &realmProvider{
		realms: c.Realms,
	}
}

func (s *realmProvider) GetRealmsCount() int {
	return len(s.realms)
}

func (s *realmProvider) GetRealm(idx int) *config.RealmConfig {
	return s.realms[idx]
}
