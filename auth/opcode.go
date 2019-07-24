package auth

type opcode uint8

const (
	logonChallengeOpcode     opcode = 0x00
	logonProofOpcode         opcode = 0x01
	reconnectChallengeOpcode opcode = 0x02
	reconnectProofOpcode     opcode = 0x03
	realmlistOpcode          opcode = 0x10
	xferInitiateOpcode       opcode = 0x30
	xferDataOpcode           opcode = 0x31
	xferAcceptOpcode         opcode = 0x32
	xferResumeOpcode         opcode = 0x33
	xferCancelOpcode         opcode = 0x34
)

func (o opcode) String() string {
	switch o {
	case logonChallengeOpcode:
		return "logonChallengeOpcode"
	case logonProofOpcode:
		return "logonProofOpcode"
	case reconnectChallengeOpcode:
		return "reconnectChallengeOpcode"
	case reconnectProofOpcode:
		return "reconnectProofOpcode"
	case realmlistOpcode:
		return "realmlistOpcode"
	case xferInitiateOpcode:
		return "xferInitiateOpcode"
	case xferDataOpcode:
		return "xferDataOpcode"
	case xferAcceptOpcode:
		return "xferAcceptOpcode"
	case xferResumeOpcode:
		return "xferResumeOpcode"
	case xferCancelOpcode:
		return "xferCancelOpcode"
	}
	return "unknown"
}
