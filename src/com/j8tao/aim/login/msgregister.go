package login

import (
	. "com/j8tao/aim/core"
	"reflect"
	"com/j8tao/aim/db"
	"time"
)

type MsgUserRegister struct {
	Name NAME_STRING
	Account NAME_STRING
	UserID ObjectID
}

func (this *MsgUserRegister) GetNetBytes() ([]byte, bool) {
	return GenNetBytes(uint16(CMD_REGISTER), reflect.ValueOf(this))
}

func (this *MsgUserRegister) CreateByBytes(bytes []byte) (bool, int) {
	return Byte2Struct(reflect.ValueOf(this), bytes)
}

func (this *MsgUserRegister) Process(p interface{}) {
	retChan := p.(chan ObjectID)
	rowsNum, err := db.DBMgr.PreQuery("select id from t_bd_user where name = ?", Byte2String(this.Name[:]))
	if err != nil || len(rowsNum) != 0 {
		LogError(err)
		this.UserID = 1
		return
	}

	rowsResult, err1 := db.DBMgr.PreExecute("insert into t_bd_user (id, name, account) values (?,?,?)", 1000000, Byte2String(this.Name[:]), Byte2String(this.Account[:]))
	if err1 != nil {
		LogError(err1)
		this.UserID = 2
		return
	}
	if num, _:= rowsResult.RowsAffected(); num == 0 {
		LogError(err1)
		this.UserID = 2
		return
	}
	this.UserID = 1000000
	select {
	case retChan <- ObjectID(1000000):
	case <- time.After(5*time.Second):
		LogError("MsgUserRegister send error")
	}

}
