package utils

import (
	"strconv"
)

type Reader interface {
	Read(buffer []byte) (int, error)
}

//Returns minimal integer value.
func Min(a, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

//Checks if the provided string is a port.
//Returns an error if it's not.
func CheckPort(port string) error {
	_, err := strconv.Atoi(port)
	return err
}
