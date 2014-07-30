package core

type ObjectID string
type NAME_STRING [NAME_LENGTH]byte

const (
	TAG uint16 = 0
	VERSION uint16 = 1
	HEADER_LENGTH uint16 = 8
	NAME_LENGTH uint16 = 255
)

const (
	SYSTEM_CHAN_ID ObjectID = "systemchanid"
	SYSTEM_USER_CHAN_ID ObjectID = "systemuserchanid"
)
