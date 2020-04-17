package app

import (
	"fmt"
	"net"
)

type ClientMessage struct {
	Type    string                 `json:"type"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type Request struct {
	Raw  []byte
	Addr *net.UDPAddr
	Body ClientMessage
}

func (request Request) get(key string) interface{} {
	return request.Body.Data[key]
}

func (request Request) getStr(key string) string {
	return fmt.Sprintf("%s", request.Body.Data[key])
}

func (request Request) getDef(key string, def string) interface{} {
	if request.Body.Data[key] != nil {
		return request.Body.Data[key]
	}
	return def
}
