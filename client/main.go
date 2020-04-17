package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

var client Client

type Client struct {
	Name    string `json:"name"`
	Channel string `json:"channel"`
	conn    *net.UDPConn
}

type Request struct {
	Type    string                 `json:"type"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

// 发送消息到服务器
func (client Client) send(request Request) {
	bytes, _ := json.Marshal(request)
	_, _ = client.conn.Write(bytes)
}

// 连接到服务器
func (client Client) connect() {
	client.send(Request{
		Type:    "connect",
		Message: "",
		Data: map[string]interface{}{
			"channel": client.Channel,
			"name":    client.Name,
		},
	})
}

func main() {
	client.Channel = "default"
	for client.Name == "" {
		fmt.Print("Input Name: ")
		_, _ = fmt.Scanln(&client.Name)
	}
	ip := net.ParseIP("127.0.0.1")
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: ip, Port: 9981}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		fmt.Println(err)
	}
	client.conn = conn
	defer conn.Close()
	client.connect() // 连接到服务器
	// 读取数据
	go client.handleRead()
	// 发送数据
	client.handleWrite()
}

func (client Client) handleWrite() {
	for {
		var command string
		fmt.Print("Input Command: ")
		_, _ = fmt.Scanln(&command)
		if command != "" {
			client.send(Request{
				Type:    "message",
				Message: command,
			})
		}
	}
}

type Response struct {
	Type     string `json:"type"`
	Message  string `json:"message"`
	Channel  string `json:"channel"`
	SendUser string `json:"send_user"`
}

func (client Client) handleRead() {
	bites := make([]byte, 1024)
	for {
		l, _ := client.conn.Read(bites)
		if l > 0 {
			var response Response
			_ = json.Unmarshal(bites[:l], &response)
			fmt.Printf("\n[%s %s] %s \n", response.Channel, response.SendUser, response.Message)
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}
