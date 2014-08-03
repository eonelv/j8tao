package main

import (
	_ "com/j8tao/aim/core"
	_ "com/j8tao/aim/build"
	"net"
	"os"
	. "com/j8tao/aim/core"
	"com/j8tao/aim/cfg"
	. "com/j8tao/aim/db"
	_ "com/j8tao/aim/cfg"
	. "com/j8tao/aim/user"
	. "com/j8tao/aim/login"
	. "com/j8tao/aim/idmgr"
	"runtime"
	"fmt"
	"reflect"
	"regexp"
)

func main() {
	Start()
}

func Start() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	cfgOK, cfgErr := cfg.LoadCfg()
	LogInfo("")
	LogInfo("------------------start server-----------------------")
	if !cfgOK {
		LogInfo("Load server config error.", cfgErr)
		os.Exit(100)
	}

	if !CreateDBMgr(cfg.GetServerHome() + "/" + cfg.GetDBName()) {
		LogError("Connect dataBase error")
		os.Exit(101)
	}
	InitGenerator()
	CreateChanMgr()
	if ok, err := CreateUserMgr(); !ok{
		LogError("Create user manager error.", err)
		return
	}

	sysChan := make(chan *Command)
	RegisterChan(SYSTEM_CHAN_ID, sysChan)
	go processTCP()

	for {
		select {
		case msg := <-sysChan:
			LogInfo("system command :", msg.Cmd)
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
	defer func() {
		if err := recover(); err != nil {
			LogError(err)    //这里的err其实就是panic传入的内容
		}
	}()
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
	defer func() {
		if err := recover(); err != nil {
			LogError(err)    //这里的err其实就是panic传入的内容
		}
	}()
	client := &TCPClient{}
	objID := conn.RemoteAddr().String()
	re := regexp.MustCompile(`(\d+\.\d+\.\d+\.\d+)|(\d+)`)
	ips := re.FindStringSubmatch(objID)
	CopyArray(reflect.ValueOf(&client.AccountID), []byte(ips[0]))
	client.Conn = conn
	client.Sender = CreateTCPSender(conn)
	go ProcessRecv(client)
}






