package models

type RealmTimezone uint8

// FIXME: refactor

const (
	RealmTimezoneUnknown      RealmTimezone = 0  // any language
	RealmTimezoneDevelopment  RealmTimezone = 1  // any language
	RealmTimezoneUNITEDSTATES               = 2  // extended-Latin
	RealmTimezoneOCEANIC      RealmTimezone = 3  // extended-Latin
	RealmTimezoneLATINAMERICA               = 4  // extended-Latin
	RealmTimezoneTOURNAMENT5  RealmTimezone = 5  // basic-Latin at create any at login
	RealmTimezoneKOREA        RealmTimezone = 6  // East-Asian
	RealmTimezoneTOURNAMENT7  RealmTimezone = 7  // basic-Latin at create any at login
	RealmTimezoneENGLISH      RealmTimezone = 8  // extended-Latin
	RealmTimezoneGERMAN       RealmTimezone = 9  // extended-Latin
	RealmTimezoneFRENCH       RealmTimezone = 10 // extended-Latin
	RealmTimezoneSPANISH      RealmTimezone = 11 // extended-Latin
	RealmTimezoneRussian      RealmTimezone = 12 // Cyrillic
	RealmTimezoneTOURNAMENT13               = 13 // basic-Latin at create any at login
	RealmTimezoneTAIWAN       RealmTimezone = 14 // East-Asian
	RealmTimezoneTOURNAMENT1                = 15 // basic-Latin at create any at login
	RealmTimezoneCHINA        RealmTimezone = 16 // East-Asian
	RealmTimezoneCN1          RealmTimezone = 17 // basic-Latin at create any at login
	RealmTimezoneCN2          RealmTimezone = 18 // basic-Latin at create any at login
	RealmTimezoneCN3          RealmTimezone = 19 // basic-Latin at create any at login
	RealmTimezoneCN4          RealmTimezone = 20 // basic-Latin at create any at login
	RealmTimezoneCN5          RealmTimezone = 21 // basic-Latin at create any at login
	RealmTimezoneCN6          RealmTimezone = 22 // basic-Latin at create any at login
	RealmTimezoneCN7          RealmTimezone = 23 // basic-Latin at create any at login
	RealmTimezoneCN8          RealmTimezone = 24 // basic-Latin at create any at login
	RealmTimezoneTOURNAMENT25               = 25 // basic-Latin at create any at login
	RealmTimezoneTESTSERVER   RealmTimezone = 26 // any language
	RealmTimezoneTOURNAMENT27               = 27 // basic-Latin at create any at login
	RealmTimezoneQASERVER     RealmTimezone = 28 // any language
	RealmTimezoneCN9          RealmTimezone = 29
)
