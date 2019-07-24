package auth

type accountFlags uint32

const (
	accountFlagGM      accountFlags = 0x00000001
	accountFlagTrial   accountFlags = 0x00000008
	accountFlagPropass accountFlags = 0x00800000
)
