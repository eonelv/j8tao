package cfg

import (
	"os"
	"bufio"
	"strings"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
)

type ServerConfig struct {
	ServerPort int
	ServerHome string
	DBName string
	BuildCmd string
	BuildCmdSP string
	isDebug bool
}

var srvCfg ServerConfig
var serverCfg map[string]string
func LoadCfg() (bool, error) {
	userFile := "server.cfg"

	file,err := os.OpenFile(userFile, os.O_RDONLY, os.ModeAppend)

	if err != nil {
		return false, err
	}
	defer file.Close()

	serverCfg = make(map[string]string)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// expression match
		lineArray := strings.Split(line, "=")
		if len(lineArray) >= 2 {
			serverCfg[lineArray[0]] = lineArray[1]
		}
	}
	srvCfg = ServerConfig{}
	srvCfg.ServerHome = serverCfg["SERVER_HOME"]
	srvCfg.ServerPort, _ = strconv.Atoi(serverCfg["SERVER_PORT"])
	srvCfg.BuildCmd = serverCfg["BUILD_CMD"]
	srvCfg.BuildCmdSP = serverCfg["BUILD_CMD_SP"]
	srvCfg.DBName = serverCfg["DB_NAME"]
	srvCfg.isDebug, _ = strconv.ParseBool(serverCfg["IS_DEBUG"])

	return true, nil
}

func  GetServerPort() int {
	return srvCfg.ServerPort
}

func GetServerHome() string {
	return srvCfg.ServerHome
}

func GetCmd() string {
	return srvCfg.BuildCmd
}

func GetCmdSP() string {
	return srvCfg.BuildCmdSP
}

func GetDBName() string {
	return srvCfg.DBName
}

func IsDebug() bool {
	return srvCfg.isDebug
}


