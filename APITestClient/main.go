// TestClient project main.go
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
)

const (
	DEFAULT_ADDRESS = "localhost:8086"
)

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Application requires at least 2 parameters: user password")
		return
	}

	conn, err := connect(getAddress())
	panicError(err)
	defer conn.Close()

	connectToCache(conn, "TestCache")

	cache := cache.NewRemoteCache(conn)

	value, _ := cache.Put("A", "B")
	fmt.Println("cache.Put(\"A\", \"B\")")
	fmt.Println("Replaced: ", value)
	value, _ = cache.Get("A")
	fmt.Println("cache.Get(\"A\")")
	fmt.Println("Got: ", value)
	value, _ = cache.ReplaceValue("A", "B", "C")
	fmt.Println("cache.ReplaceValue(\"A\", \"B\", \"C\")")
	fmt.Println("Replaced: ", value)
}

func connect(address string) (net.Conn, error) {
	conn, err := tls.Dial("tcp", address, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func connectToCache(conn net.Conn, cacheId string) {

	sendCredentials(conn)

	data := []byte("connect-to " + cacheId)
	data = append(data, '\n')
	conn.Write(data)
}

func sendCredentials(conn net.Conn) {
	user := &auth.User{Name: os.Args[1], Pass: auth.EncryptPass(os.Args[2]), IsMachine: true}

	data, err := auth.UserToJson(user)
	panicError(err)

	data = append(data, '\n')
	conn.Write(data)

	data, _, err = bufio.NewReader(conn).ReadLine()
	panicError(err)
	response, err := cache.JsonToResponse(data)
	panicError(err)
	if response.Err != "" {
		panic(response.Err)
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
	if err != nil && err != io.EOF {
		panic(err)
	}
}
