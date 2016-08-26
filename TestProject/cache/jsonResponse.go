package cache

import (
	"encoding/json"
	"net"
)

//Serializable response that is used to communicate between client and server
//for Cache API
type JsonResponse struct {
	Value interface{}
	Err   string
}

//Writes a JsonResponse serializable structure with passed error
func WriteErrorResponse(conn net.Conn, err error) error {
	return WriteResponse(conn, nil, err)
}

//Writes a JsonResponse serializable structure with both passed value and error
func WriteResponse(conn net.Conn, value interface{}, err error) error {
	var e string
	if err != nil {
		e = err.Error()
	}

	bytes, marshalErr := json.Marshal(&JsonResponse{Value: value, Err: e})
	if marshalErr != nil {
		err := handleMarshalError(conn, marshalErr)
		if err != nil {
			return err
		}
	}
	writeErr := writeBytes(conn, append(bytes, '\n'))
	if writeErr != nil {
		return writeErr
	}
	return nil
}

func handleMarshalError(conn net.Conn, err error) error {
	bytes, marshalErr := json.Marshal(&JsonResponse{nil, err.Error()})
	if marshalErr != nil {
		return marshalErr
	}
	return writeBytes(conn, bytes)
}

func writeBytes(conn net.Conn, bytes []byte) error {
	_, err := conn.Write(bytes)
	return err
}

//Converts Json string to JsonResponse structure
func JsonToResponse(data []byte) (*JsonResponse, error) {
	response := new(JsonResponse)
	err := json.Unmarshal(data, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
