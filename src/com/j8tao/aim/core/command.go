package core

type Command struct {
	Cmd uint16
	Message interface {}
	RetChan chan *Command
	OtherInfo interface{}
}

type PackHeader struct {
	Tag uint16 "TAG"
	Version uint16 "VERSION"
	Length uint16
	Cmd uint16
}

const (
	CMD_CONNECTION uint16 = 1001
	CMD_REGISTER uint16 = 1002
	CMD_LOGIN uint16 = 1003
	CMD_TALK uint16 = 1006
	CMD_BUILD uint16 = 1010
)

const (
	CMD_SYSTEM_MAIN_CLOSE uint16 = 10001
	CMD_SYSTEM_USER_OFFLINE uint16 = 10002
	CMD_SYSTEM_USER_LOGIN uint16 = 10005
	CMD_SYSTEM_BROADCAST uint16 = 10006
)

