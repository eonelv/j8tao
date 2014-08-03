package login

import (
	. "com/j8tao/aim/core"
	"reflect"
	"com/j8tao/aim/db"
	"time"
	. "com/j8tao/aim/idmgr"
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

	id := SysIDGenerator.GetNextID(uint64(1001000))

	if id == 0 {
		LogError("generate ID faild")
		this.UserID = 3
		return
	}
	rowsResult, err1 := db.DBMgr.PreExecute("insert into t_bd_user (id, name, account) values (?,?,?)", id, Byte2String(this.Name[:]), Byte2String(this.Account[:]))
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
	this.UserID = id
	select {
	case retChan <- id:
	case <- time.After(20*time.Second):
		LogError("MsgUserRegister send error")
	}

}
