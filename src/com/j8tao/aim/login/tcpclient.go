package login

import (
	"net"
	. "com/j8tao/aim/core"
)

import (
	"time"
	"io"
	"reflect"
	"com/j8tao/aim/db"
)

const (
	STR_AUTH_REQ    = "<policy-file-request/>"
	STR_AUTH_RETURN = "<?xml version=\"1.0\"?><cross-domain-policy><site-control permitted-cross-domain-policies='all'/><allow-access-from domain=\"*\" to-ports=\"*\"/></cross-domain-policy>"
)

type TCPClient struct {
	ID ObjectID
	AccountID NAME_STRING
	Conn *net.TCPConn
	Sender *TCPSender
	dataChan     chan *Command // 自身接收用的channel
	userChan chan *Command // 跟自己相关的User的接收channel
	isLogin bool
	IsConnection bool
	UserEncrypt *Encrypt
}

func ProcessRecv(client *TCPClient) {
	defer func() {
		if err := recover(); err != nil {
			LogError(err)    //这里的err其实就是panic传入的内容
		}
	}()
	conn := client.Conn
	defer conn.CloseWrite()
	defer func() {
		if client.dataChan != nil {
			close(client.dataChan)
		}
	}()
	defer client.close()
	client.UserEncrypt = &Encrypt{}
	client.UserEncrypt.InitEncrypt(164, 29, 30, 60, 241, 79, 251, 107)
	for {
		headerBytes := make([]byte, HEADER_LENGTH)
		_, err := io.ReadFull(conn, headerBytes)
		if err != nil {
			LogError("Read Data Error", client.ID, err.Error())
			break
		}

		if headerBytes[0] == STR_AUTH_REQ[0] && !client.isLogin {
			tempbuf := make([]byte, len(STR_AUTH_REQ)-int(HEADER_LENGTH))
			_, err = io.ReadFull(conn, tempbuf)
			if err != nil {
				LogError("HandleUserConnect read rest auth req err", err)
				return
			}
			headerBytes = append(headerBytes, tempbuf...)
			authReq := string(headerBytes)
			if authReq == STR_AUTH_REQ {
				conn.Write(append([]byte(STR_AUTH_RETURN), 0))
			} else {
				LogError("recv wrong auth req:", authReq)
			}
			return
		}
		client.UserEncrypt.Encrypt(headerBytes, 0, len(headerBytes), true)

		header := &PackHeader{}
		Byte2Struct(reflect.ValueOf(header), headerBytes)

		bodyBytes := make([]byte, header.Length-HEADER_LENGTH)
		_, err = io.ReadFull(conn, bodyBytes)
		if err != nil {
			LogError("read data error ", err.Error())
			break
		}

		client.UserEncrypt.Encrypt(bodyBytes, 0, len(bodyBytes), false)
		client.UserEncrypt.Reset()

		client.processClientMessage(header, bodyBytes)
	}
}

func (client *TCPClient)processClientMessage(header *PackHeader, datas []byte) {
	if !client.isLogin {
		client.processLogin(header, datas)
	} else {
		client.routMsgToUser(header, datas)
	}
}

func (client *TCPClient) processLogin(header *PackHeader, datas []byte) {
	defer func() {
		if err := recover(); err != nil {
			LogError(err)
		}
	}()

	if header.Cmd != CMD_CONNECTION && !client.IsConnection {
		LogError("Wrong command", header.Cmd, " should be ", CMD_CONNECTION)
		return
	}
	if !client.IsConnection {
		go client.Sender.Start()
	}

	if header.Cmd == CMD_CONNECTION {
		msgConnection := &MsgConnection{}
		msgConnection.CreateByBytes(datas)
		msgConnection.Process(client)
		return
	}

	var userID ObjectID
	var targetChan chan *Command
	if header.Cmd == CMD_LOGIN {

		msgLogin := &MsgUserLogin{}
		msgLogin.CreateByBytes(datas)
		userID = msgLogin.ID
		targetChan = GetChanByID(userID)
		LogInfo("登录用户ID：", userID)

	} else if header.Cmd == CMD_REGISTER {
		msgUserRegister := &MsgUserRegister{}
		msgUserRegister.CreateByBytes(datas)

		chanRet := make(chan ObjectID)

		go msgUserRegister.Process(chanRet)

		select {
		case id := <- chanRet:
			userID = id
		case <-time.After(20 * time.Second):
			LogError("register put user channel failed:", header.Cmd)
			client.Sender.Send(msgUserRegister)
			return
		}
		client.Sender.Send(msgUserRegister)
		targetChan = GetChanByID(SYSTEM_USER_CHAN_ID)
	}

	client.ID = userID
	client.isLogin = true

	client.dataChan = make(chan *Command)
	msgSend := &Command{CMD_SYSTEM_USER_LOGIN, client.ID, client.dataChan, nil}
	msgSend.OtherInfo = client.Sender

	select {
	case targetChan <- msgSend:
	case <-time.After(5 * time.Second):
		LogError("loginUserToGame put user channel failed:", CMD_SYSTEM_USER_LOGIN)
	}
	client.waitLoginReturn()
}

func (client *TCPClient) processConnection(msgConn *MsgConnection) {
	client.IsConnection = true
	sql := "select * from t_bd_user where account = ?"
	rows, err := db.DBMgr.PreQuery(sql, Byte2String(msgConn.AccountID[:]))
	var msgUserLogin *MsgUserLogin = &MsgUserLogin{}
	if err != nil || len(rows) == 0 {
		client.Sender.Send(msgUserLogin)
		LogInfo("there is no user", err)
		return
	}
	msgUserLogin.ID = rows[0].GetObjectID("id")
	client.Sender.Send(msgUserLogin)
}

func (client *TCPClient) waitLoginReturn() bool {
	msg := <-client.dataChan
	if msg.RetChan == nil {
		return false
	}
	client.userChan = msg.RetChan
	return true
}

// 将消息路由到玩家处理
func (client *TCPClient) routMsgToUser(header *PackHeader, data []byte) bool {
	msg := &Command{header.Cmd, data, nil, nil}

	select {
	case client.userChan <- msg:
	case <-time.After(5 * time.Second):
		LogError("routMsgToUser put user channel failed:", client.ID)
		return false
	}

	return true
}

func (client *TCPClient) close() {
	if !client.isLogin {
		return
	}

	userInnerChan := GetChanByID(client.ID)
	closeMsg := &Command{CMD_SYSTEM_USER_OFFLINE, nil, client.dataChan, nil}
	client.isLogin = false

	select {
	case userInnerChan <- closeMsg:
	case <-time.After(5 * time.Second):
		LogError("sendOffline put user channel failed:", client.ID)
		return
	}
	return
}
