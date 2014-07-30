package core

import (
	"reflect"
)

type CreateNetMsgFunc func(cmdData *Command) NetMsg

type NetMsg interface {
	GetNetBytes() ([]byte, bool)
	CreateByBytes(bytes []byte) (bool, int)
	Process(p interface{})
}

func GenNetBytes(cmd uint16, values reflect.Value) ([]byte,bool) {
	datas,ok := Struct2Bytes(values)
	if !ok {
		return nil,false
	}
	length := uint16(HEADER_LENGTH + uint16(len(datas)))
	header := &PackHeader{TAG, VERSION, length, cmd}
	headerDatas,okh := Struct2Bytes(reflect.ValueOf(header))
	if !okh {
		return nil, false
	}
	return append(headerDatas, datas...), true
}
