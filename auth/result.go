package auth

type result uint8

const (
	resultSuccess                     result = 0x00
	resultBanned                      result = 0x03
	resultUnknownAccount              result = 0x04
	resultIncorrectPassword           result = 0x05
	resultAlreadyOnline               result = 0x06
	resultNoTime                      result = 0x07
	resultDBBusy                      result = 0x08
	resultVersionInvalid              result = 0x09
	resultVersionUpdate               result = 0x0A
	resultInvalidServer               result = 0x0B
	resultSuspended                   result = 0x0C
	resultFailNoAccess                result = 0x0D
	wowSuccessSurvey                  result = 0x0E
	resultParentControl               result = 0x0F
	resultLockedEnforced              result = 0x10
	resultTrialEnded                  result = 0x11
	resultUseBattlenet                result = 0x12
	resultAntiIndulgence              result = 0x13
	resultSessionExpired              result = 0x14
	resultNoGameAccount               result = 0x15
	resultChargeback                  result = 0x16
	resultInternetGameRoomWithoutBnet result = 0x17
	resultGameAccountLocked           result = 0x18
	resultUnlockableLock              result = 0x19
	resultConversionRequired          result = 0x20
	resultDisconnected                result = 0xFF
)
