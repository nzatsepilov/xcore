package models

type RealmType uint8

const (
	RealmTypeNormal  RealmType = 0
	RealmTypePVP     RealmType = 1
	RealmTypeNormal2 RealmType = 4
	RealmTypeRP      RealmType = 6
	RealmTypeRPPVP   RealmType = 8
	RealmTypeFFAPVP  RealmType = 16
)

const MaxClientRealmType uint8 = 14
