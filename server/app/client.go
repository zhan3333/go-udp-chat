package app

import (
	"net"
)

type Client struct {
	Name string `json:"name"`
	Addr *net.UDPAddr
}

// 判断两个客户端是否相同
func (client Client) isSame(inputClient Client) bool {
	return client.Addr.String() == inputClient.Addr.String()
}
