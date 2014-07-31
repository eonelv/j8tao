package main

import (
	"net"
	"os"
	. "com/j8tao/aim/core"
	"com/j8tao/aim/cfg"
	. "com/j8tao/aim/db"
	. "com/j8tao/aim/user"
	_ "com/j8tao/aim/cfg"
	"runtime"
	"fmt"
)

func main() {
	Start()
}

func Start() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	cfgOK, cfgErr := cfg.LoadCfg()
	LogInfo("")
	LogInfo("---------------------------------")
	if !cfgOK {
		LogInfo("Load server config error.", cfgErr)
		os.Exit(100)
	}
	LogInfo("Load server config success.")
	if !CreateDBMgr(cfg.GetServerHome() + "/" + cfg.GetDBName()) {
		LogError("Connect dataBase error")
		os.Exit(101)
	}
	CreateChanMgr()
	CreateUserMgr()
	sysChan := make(chan *Command)
	RegisterChan(SYSTEM_CHAN_ID, sysChan)
	go processTCP()
	LogInfo("Server bootup success.")
	for {
		select {
		case msg := <-sysChan:
			LogInfo("main recv msg:", msg.Cmd)
			if msg.Cmd == CMD_SYSTEM_MAIN_CLOSE {
				return
			}
		}
	}
}

func checkError(err error){
	if err != nil {
		LogError(err)
		os.Exit(0)
	}
}

func processTCP() {
	service := fmt.Sprintf(":%d",  cfg.GetServerPort())
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		conn,err := listener.AcceptTCP()
		if err != nil {
			continue
		}
		go processConnect(conn)
	}
}

func processConnect(conn *net.TCPConn){
	client := &TCPClient{}
	objID := conn.RemoteAddr().String()
	client.ID = ObjectID(objID)
	client.Conn = conn
	client.Sender = CreateTCPSender(conn)
	go ProcessRecv(client)
}






