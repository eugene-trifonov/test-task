package auth

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//Defines a serializable structure that can be sent between client and server,
//or that can be stored.
type User struct {
	Name      string
	Pass      string
	IsMachine bool
}

//Reads the list of users allowed to connect to the In-memory cache
func ReadUsers() (map[string]string, error) {
	file, err := os.Open("users")
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	users := []User{}
	err = json.Unmarshal(data, &users)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	for i := range users {
		m[users[i].Name] = users[i].Pass
	}

	return m, nil
}

//Converts Json string to User structure
func JsonToUser(jsonUserData []byte) (*User, error) {
	var user User

	err := json.Unmarshal(jsonUserData, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

//Converts User structure to Json string
func UserToJson(user *User) ([]byte, error) {
	data, err := json.Marshal(user)

	if err != nil {
		return nil, err
	}

	return data, nil
}
