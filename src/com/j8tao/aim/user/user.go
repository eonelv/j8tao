package user

import (
	"fmt"
	. "com/j8tao/aim/core"
)

const (
	USER_STATUS_INIT int16 = 0
	USER_STATUS_ONLINE int16 = 100
	USER_STATUS_OFFLINE int16 = 200
)

type User struct {
	ID ObjectID
	recvChan chan *Command
	innerChan chan *Command
	netChan chan *Command
	Sender *TCPSender
	Status int16
}

func CreateUser(id ObjectID) *User{
	user := &User{}
	user.ID = id
	user.recvChan = make(chan *Command)
	user.innerChan = make(chan *Command)
	RegisterChan(id, user.innerChan)

	user.Status = USER_STATUS_INIT
	go startUserRecv(user)

	return user
}

func startUserRecv(user *User) {
	for {
		select {
		case msg:=<-user.recvChan:
			if msg == nil && user.Status == USER_STATUS_OFFLINE{
				return
			}
			user.processClientMessage(msg)
		case msg:= <-user.innerChan:
			if msg == nil && user.Status == USER_STATUS_OFFLINE{
				return
			}
			user.processInnerMessage(msg)
		}
	}
}

func (user *User) processClientMessage(msg *Command) {
	if msg == nil {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			LogError("User processClientMsg failed:", err, " cmd:", msg.Cmd)
		}
	}()
	netMsg := CreateNetMsg(msg)
	netMsg.Process(user)
}

func (user *User) processInnerMessage(msg *Command) {
	 if msg == nil {
		 return
	 }
	switch msg.Cmd{
	case CMD_SYSTEM_USER_LOGIN:
		user.userLogin(msg)
	case CMD_SYSTEM_USER_OFFLINE:
		user.userOffline(msg)
	}
}

func (user *User) userLogin(msg *Command) {
	user.Status = USER_STATUS_ONLINE
	user.netChan = msg.RetChan
	user.Sender = msg.OtherInfo.(*TCPSender)

	msg.RetChan = user.recvChan
	user.netChan <- msg
}

func (user *User) userOffline(msg *Command) {
	defer func() {
		if err := recover(); err != nil {
			LogError("User processClientMsg failed:", err, " cmd:", msg.Cmd)
		}
	}()
	if msg.RetChan != user.netChan {
		return
	}
	UnRegisterChan(user.ID)
	close(user.recvChan)
	close(user.innerChan)
	user.Sender.Close()
	user.Status = USER_STATUS_OFFLINE
	fmt.Println("User offline", user.ID)
}
