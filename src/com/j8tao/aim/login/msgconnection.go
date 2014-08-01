package login

import (
	. "com/j8tao/aim/core"
	"reflect"
)

type MsgConnection struct {
	AccountID NAME_STRING
}

func createBuildNetMsg(cmdData *Command) NetMsg {
	netMsg := &MsgConnection{}
	netMsg.CreateByBytes(cmdData.Message.([]byte))
	return netMsg
}

func (this *MsgConnection) GetNetBytes() ([]byte, bool) {
	return GenNetBytes(uint16(CMD_CONNECTION), reflect.ValueOf(this))
}

func (this *MsgConnection) CreateByBytes(bytes []byte) (bool, int) {
	return Byte2Struct(reflect.ValueOf(this), bytes)
}

func (this *MsgConnection) Process(p interface{}) {
	tcpClient := p.(*TCPClient)
	if tcpClient == nil {
		return
	}
	tcpClient.processConnection(this)
}
