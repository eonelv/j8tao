package idmgr

import (
	"sync"
	. "com/j8tao/aim/core"
	. "com/j8tao/aim/db"
	"time"
)

type IDGenerator struct {
	IDChanSet map[uint64]chan ObjectID
	mutex sync.RWMutex
}

var SysIDGenerator IDGenerator
func InitGenerator() {
	SysIDGenerator = IDGenerator{}
	SysIDGenerator.IDChanSet = make(map[uint64]chan ObjectID)
	SysIDGenerator.mutex = sync.RWMutex{}
	SysIDGenerator.create()
}

func generateID(idMax ObjectID, ch chan ObjectID) {
	for {
		ch <- idMax
		LogInfo("写入ID", idMax)
		idMax++
	}

}

func (this *IDGenerator) create() bool{
	rows, err := DBMgr.PreQuery("select (id / ?) as serverplat, max(id) as maxid from t_bd_user group by (id / ?)", ID_SEPARATOR, ID_SEPARATOR)
	if err != nil || len(rows) == 0 {
		LogInfo("Generator create faild")
		return false
	}

	for _, row := range rows {
		serverPlat := row.GetUint64("serverplat")
		idMax := row.GetObjectID("maxid")
		this.IDChanSet[serverPlat] = make(chan ObjectID, 1)
		go generateID(idMax + 1, this.IDChanSet[serverPlat])
		LogInfo("创建Channel", serverPlat)
	}
	return true
}

func (this *IDGenerator) GetNextID(serverPlat uint64) ObjectID {
	var ch chan ObjectID
	var ok bool
	this.mutex.RLock()
	ch, ok = this.IDChanSet[serverPlat]
	this.mutex.RUnlock()
	if !ok {
		this.mutex.Lock()
		ch, ok = this.IDChanSet[serverPlat]
		if !ok {
			idMax := ObjectID(serverPlat * ID_SEPARATOR + 1)
			ch = make(chan ObjectID, 1)
			this.IDChanSet[serverPlat] = ch
			go generateID(idMax, ch)
		}
		this.mutex.Unlock()
	}
	var id ObjectID
	select {
	case id = <-ch:
	case <- time.After(20*time.Second):
		LogError("get Next ID time out")
		return ObjectID(0)
	}

	return id
}
