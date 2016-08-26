package utils

import (
	"fmt"
	"net"
)

//Represents a logger interface.
//From time to time the logging shoud happen into a different places: console, file, socket connection etc.
//For this purpose the interface was introduced.
type Logger interface {
	//Formatted log, the same as fmt.Printf().
	Logf(str string, params ...interface{})

	//Log in one line, the same as fmt.Println().
	Log(any ...interface{})

	//Log each value into a separate line.
	Logln(any ...interface{})
}

type ConsoleLogger struct{}

func NewConsoleLogger() Logger {
	return new(ConsoleLogger)
}

func (logger ConsoleLogger) Logf(str string, params ...interface{}) {
	fmt.Printf(str, params)
	fmt.Println()
}

func (logger ConsoleLogger) Log(any ...interface{}) {
	fmt.Println(any)
}

func (logger ConsoleLogger) Logln(any ...interface{}) {
	for i := 0; i < len(any); i++ {
		fmt.Println(any[i])
	}
}

type TelnetLogger struct {
	conn net.Conn
}

func NewTelnetLogger(conn net.Conn) *TelnetLogger {
	logger := new(TelnetLogger)
	logger.conn = conn
	return logger
}

func (logger TelnetLogger) Logf(str string, params ...interface{}) {
	logger.conn.Write([]byte(fmt.Sprintf(str, params...)))
	logger.conn.Write([]byte("\r\n"))
}

func (logger TelnetLogger) Log(any ...interface{}) {
	for i := 0; i < len(any); i++ {
		logger.conn.Write([]byte(fmt.Sprint(any[i])))
	}
	logger.conn.Write([]byte("\r\n"))
}

func (logger TelnetLogger) Logln(any ...interface{}) {
	for i := 0; i < len(any); i++ {
		logger.conn.Write([]byte(fmt.Sprint(any[i])))
		logger.conn.Write([]byte("\r\n"))
	}
}

func ConvertStrings(strings []string) []interface{} {
	results := make([]interface{}, len(strings))
	for i, k := range strings {
		results[i] = k
	}
	return results
}
