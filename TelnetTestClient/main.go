// TelnetTestClient project main.go
package main

import (
	"TestProject/utils"
	"Testproject/auth"
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync/atomic"
)

const (
	DEFAULT_ADDRESS = "localhost:8086"
)

var connectionClosed atomic.Value

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Application requires at least 2 parameters: user password")
		return
	}

	conn, err := tls.Dial("tcp", getAddress(), &tls.Config{InsecureSkipVerify: true})
	panicError(err)

	defer conn.Close()

	sendCredentials(conn)

	go readMessages(conn)
	handleUserMessages(conn)

}

func sendCredentials(conn net.Conn) {
	user := &auth.User{Name: os.Args[1], Pass: auth.EncryptPass(os.Args[2]), IsMachine: false}

	data, err := auth.UserToJson(user)
	panicError(err)

	data = append(data, '\n')
	conn.Write(data)
}

func readMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		msg, _, err := reader.ReadLine()
		if err == io.EOF {
			connectionClosed.Store(true)
			return
		} else if err != nil {
			connectionClosed.Store(true)
			panic(err)
		}

		fmt.Println(string(msg))
	}
}

func handleUserMessages(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			return
		} else if err != nil {
			panic(err)
		}

		if connectionClosed.Load() != nil {
			return
		}

		line = append(line, '\n')
		conn.Write(line)
	}
}

func getAddress() string {
	if len(os.Args) > 4 {
		if utils.CheckPort(os.Args[4]) != nil {
			panic(errors.New("Wrong port defined"))
		}
		return os.Args[3] + ":" + os.Args[4]
	}
	return DEFAULT_ADDRESS
}

func panicError(err error) {
	if err != nil {
		panic(err)
	}
}
