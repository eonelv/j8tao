package message

import (
	. "com/j8tao/aim/core"
)

func init() {
	LogInfo("Registor message", CMD_TALK)
	initNetMsgCreator()
}
