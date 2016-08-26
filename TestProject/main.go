// TestProject project main.go
package main

import (
	"TestProject/auth"
	"TestProject/cache"
	"TestProject/utils"
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

const (
	TIMEOUT = 60000000000 //1 minute

	DEFAULT_PORT = "8086"

	NEED_HELP = "Please use \"help\" command to find the available commands."

	HELP_KEYS = "keys - operation to display cached keys. Ex. keys [startIndex] [endIndex]"
	HELP_TTL  = "ttl - operation to update time to live attribute of any cached value. Ex. ttl key ttlInSeconds"

	HELP_GET    = "get - operation to get cached value if it exists. Ex. get key"
	HELP_SET    = "set - operation to set a new cached string value. Ex. set key value [ttl]"
	HELP_UPDATE = "update - operation to exchange an existing cached string value. Ex. update key oldValue newValue [ttlInSeconds]"
	HELP_DELETE = "delete - operation to remove an existing cached value. Ex. delete key [value]"
	HELP_SIZE   = "size - operation to find out the number of cached values. Ex. size"
	HELP_EXIT   = "exit - operation to close the connection with the server. Ex. exit"

	HELP_LSIZE   = "lsize - operation to check the size of a list. Ex. lsize key"
	HELP_LGET    = "lget - operation to get a value from cached list. Ex. lget key index"
	HELP_LAPPEND = "lappend - operation to add a new value into the cached list. Ex. lappend key value [ttlInSeconds]"
	HELP_LDELETE = "ldelete - operation to remove a value from a list by index. Ex. ldelete key index"

	HELP_DSIZE   = "dsize - operation to check the size of a dictionary. Ex. lsize key"
	HELP_DGET    = "dget - operation to get a value from cached dictionary by key. Ex. dget key dictKey"
	HELP_DSET    = "dset - operation to set a key-value pair into a cached dictionary. Ex. dset key dictKey dictValue"
	HELP_DAPPEND = "dappend - operation to add a value to the dictionary. Ex. dappend key dictKey value [ttlInSeconds]"
	HELP_DDELETE = "ddelete - opeartion to remove a value from cached dictionary. Ex. ddelete key dictKey"
)

var existingCaches cache.Cache = cache.NewCache()
var port = DEFAULT_PORT

var stopped atomic.Value

var users map[string]string

func main() {

	var err error

	users, err = auth.ReadUsers()

	if err != nil {
		fmt.Printf("Error [%v] happened", err)
		return
	}

	port = getPort()
	listener := startListenOn(port)

	for {
		conn, err := listener.Accept()

		if stopped.Load() != nil {
			return
		}

		if err != nil {
			fmt.Printf("Error [%v] happened", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {

	defer conn.Close()

	var log utils.Logger
	log = utils.NewTelnetLogger(conn)

	reader := bufio.NewReader(conn)

	credentials, _, err := reader.ReadLine()
	if err != nil {
		fmt.Printf("Error [%v] happened", err)
		return
	}
	user, err := auth.JsonToUser(credentials)
	if err != nil {
		log.Logf("Error [%v] happened", err)
		return
	}

	pass := users[user.Name]
	if pass == "" || pass != user.Pass {
		if user.IsMachine {
			cache.WriteErrorResponse(conn, errors.New("User/password pair is incorrect"))
		} else {
			log.Log("User/password pair is incorrect")
		}
		return
	} else {
		if user.IsMachine {
			cache.WriteResponse(conn, "Ok", nil)
		}
	}

	if user.IsMachine {
		handleMachineConnection(conn, reader, utils.NewConsoleLogger())
	} else {
		handleHumanConnection(conn, reader, log)
	}

}

func printHelp(log utils.Logger) {
	log.Logln(HELP_GET, HELP_SET, HELP_UPDATE, HELP_DELETE, HELP_EXIT, HELP_LGET, HELP_LAPPEND, HELP_LDELETE, HELP_LSIZE, HELP_DGET, HELP_DSET, HELP_DAPPEND, HELP_DDELETE, HELP_KEYS, HELP_TTL)
}

func getCache(id string) cache.Cache {
	existingCache := existingCaches.Get(id)
	if existingCache == nil {
		newCache := cache.NewCache()
		existingCache = existingCaches.PutIfAbsent(id, newCache)
		if existingCache == nil {
			existingCache = newCache
		}
	}
	return existingCache.(cache.Cache)
}

func handleMachineConnection(conn net.Conn, reader *bufio.Reader, log utils.Logger) {

	c, err := readCache(reader)
	if err != nil {
		log.Logf("Error [%v] happened", err)
		cache.WriteErrorResponse(conn, err)
		return
	} else if c == nil {
		return
	}

	cmds := cache.NewJsonCacheCommands(c, conn, log)

	for {

		cmd, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Logf("Error [%v] happened", err)
			log.Log("Connection will be closed")
			return
		}

		command := strings.Trim(string(cmd), " ")
		splitCommand := strings.Split(command, " ")
		command = splitCommand[0]
		params := splitCommand[1:]

		handleCommand(cmds, command, params)
	}
}

func handleHumanConnection(conn net.Conn, reader *bufio.Reader, log utils.Logger) {

	log.Log("You've been connected to In-memory cache. Connection idle timeout is ", int64(float64(TIMEOUT)/float64(time.Second)), "s.")
	log.Log("Please enter first command: \"stop-server\" or \"connect-to\" <cacheId>")

	conn.SetReadDeadline(time.Now().Add(TIMEOUT))

	c, err := readCache(reader)
	if err != nil {
		log.Logf("Error [%v] happened", err)
		return
	} else if c == nil {
		return
	}

	cmds := cache.NewUserFriendlyCommands(c, log)

	log.Log("Connected")

	for {

		conn.SetReadDeadline(time.Now().Add(TIMEOUT))

		cmd, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Logf("Error [%v] happened", err)
			log.Log("Connection will be closed")
			return
		}

		command := strings.Trim(string(cmd), " ")

		splitCommand := strings.Split(command, " ")
		command = splitCommand[0]
		params := splitCommand[1:]

		switch command {
		case "help":
			printHelp(log)
			break
		case "exit":
			log.Log("Connection closed")
			return
		default:
			err = handleCommand(cmds, command, params)
		}

		if err != nil {
			log.Log(NEED_HELP)
		}
	}
}

func handleCommand(cmds cache.CacheCommands, command string, params []string) error {
	var err error
	switch command {
	case "get":
		_, err = cmds.GetValue(params)
		break
	case "set":
		_, err = cmds.SetValue(params)
		break
	case "update":
		_, err = cmds.UpdateValue(params)
		break
	case "delete":
		_, _, err = cmds.RemoveValue(params)
		break
	case "keys":
		_, err = cmds.GetKeys(params)
		break
	case "size":
		cmds.GetSize()
	case "lget":
		_, err = cmds.GetListValue(params)
		break
	case "lappend":
		err = cmds.AppendListValue(params)
		break
	case "ldelete":
		_, err = cmds.DeleteListValue(params)
		break
	case "lsize":
		_, err = cmds.GetListSize(params)
		break
	case "dget":
		_, err = cmds.GetDictValue(params)
		break
	case "dset":
		_, err = cmds.SetDictValue(params)
		break
	case "dappend":
		_, err = cmds.AppendDictValue(params)
		break
	case "ddelete":
		_, err = cmds.DeleteDictValue(params)
		break
	case "dsize":
		_, err = cmds.GetDictSize(params)
		break
	case "ttl":
		_, err = cmds.UpdateTTL(params)
		break
	default:
		err = errors.New("Unknown command")
	}

	return err
}

func stopServer() {
	stopped.Store(true)
}

func getPort() string {
	if len(os.Args) > 1 {
		err := utils.CheckPort(os.Args[1])
		if err != nil {
			panic(errors.New(fmt.Sprint("Invalid port [", os.Args[1], "] is provided.")))
		}
		return os.Args[1]
	}
	return DEFAULT_PORT
}

func startListenOn(port string) net.Listener {

	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		panic(errors.New(fmt.Sprintf("Error [%v] happened", err)))
	}

	config := tls.Config{Certificates: []tls.Certificate{cert}}
	listener, err := tls.Listen("tcp", ":"+port, &config)
	if err != nil {
		panic(fmt.Sprintf("Error [%v] happened", err))
	}

	return listener
}

func readCache(reader *bufio.Reader) (cache.Cache, error) {
	firstCmd, _, err := reader.ReadLine()
	if err != nil {
		return nil, err
	}

	cmd := string(firstCmd)
	if cmd == "stop-server" {
		stopServer()
		return nil, nil
	}

	if !strings.HasPrefix(cmd, "connect-to") {
		return nil, errors.New("Invalid command")
	}

	params := strings.Split(cmd, " ")
	if len(params) != 2 {
		return nil, errors.New("Invalid parameters")
	}

	return getCache(params[1]), nil
}
