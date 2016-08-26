package auth

import (
	"crypto/sha512"
	"encoding/base64"
)

//Encrypts user's password before it is sent to server
func EncryptPass(password string) string {
	b := sha512.Sum512([]byte(password))
	return string(base64.StdEncoding.EncodeToString(b[:]))
}
