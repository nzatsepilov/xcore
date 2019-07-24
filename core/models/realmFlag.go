package models

type RealmFlag uint8

const (
	RealmFlagVersionMismatch RealmFlag = 1 << iota
	RealmFlagOffline
	RealmFlagSpecifyBuild
	RealmFlagUnknown1
	RealmFlagUnknown2
	RealmFlagRecommended
	RealmFlagNew
	RealmFlagFull
)

func (r RealmFlag) Has(f RealmFlag) bool {
	return r&f == f
}

func (r *RealmFlag) Append(f RealmFlag) {
	*r |= f
}

func (r *RealmFlag) Remove(f RealmFlag) {
	*r &= ^f
}

func (r *RealmFlag) Toggle(f RealmFlag) {
	*r ^= f
}
