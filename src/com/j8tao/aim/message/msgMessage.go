package message

import (
	. "com/j8tao/aim/core"
	. "com/j8tao/aim/user"
	"reflect"
)

func initNetMsgCreator() {
	isSuccess := RegisterMsgFunc(CMD_TALK, createNetMsg)
	LogInfo("Registor message", CMD_TALK)
	if !isSuccess {
		LogError("Registor MsgMessage faild")
	}
}

func createNetMsg(cmdData *Command) NetMsg {
	netMsg := &MsgMessage{}
	netMsg.CreateByBytes(cmdData.Message.([]byte))
	return netMsg
}

type MsgMessage struct {
	SenderID ObjectID
	SenderName NAME_STRING
	Receiver NAME_STRING
	Message []byte
}

func (this *MsgMessage) GetNetBytes() ([]byte, bool) {
	return GenNetBytes(uint16(CMD_TALK), reflect.ValueOf(this))
}

func (this *MsgMessage) CreateByBytes(bytes []byte) (bool, int) {
	return Byte2Struct(reflect.ValueOf(this), bytes)
}

func (this *MsgMessage) Process(p interface{}) {
	puser, ok := p.(*User)
	if !ok {
		return
	}
	if this.SenderID != puser.ID {
		LogError("Message sender is not the processor")
		return
	}
	UserMgr.BroadcastMessage(this)
}
