package user

import (
	. "com/j8tao/aim/core"
	"fmt"
)

func CreateMessage(cmdData *Command) NetMsg {
	defer func() {
		if err:= recover(); err != nil {
			fmt.Println(err)
		}
	}()
	var netMsg NetMsg
	switch cmdData.Cmd{
	case CMD_BUILD:
		netMsg = &MsgBuild{}
		netMsg.CreateByBytes(cmdData.Message.([]byte))
	}
	return netMsg
}
