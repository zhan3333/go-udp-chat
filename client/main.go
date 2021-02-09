package main

import (
	"encoding/json"
	"fmt"
	"log"
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

func (client Client) getChannels() {
	client.send(Request{
		Type:    "get_channels",
		Message: "",
		Data:    nil,
	})
}

//客户端主要做以下的事情
//连接到服务器
//选择room
//输入聊天内容
func main() {
	for client.Name == "" {
		fmt.Print("Input Name: ")
		_, _ = fmt.Scanln(&client.Name)
	}
	ip := net.ParseIP("127.0.0.1")
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: ip, Port: 9981}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		log.Panicf("run client failed: %+v", err)
	}
	defer conn.Close()

	client.conn = conn
	client.getChannels()
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
	Type     string                 `json:"type"`
	Message  string                 `json:"message"`
	Channel  string                 `json:"channel"`
	SendUser string                 `json:"send_user"`
	Data     map[string]interface{} `json:"data"`
}

func (client Client) handleRead() {
	bites := make([]byte, 1024)
	for {
		l, _ := client.conn.Read(bites)
		if l > 0 {
			var response Response
			_ = json.Unmarshal(bites[:l], &response)
			switch response.Type {
			case "message":
				fmt.Printf("\n[%s %s] %s \n", response.Channel, response.SendUser, response.Message)
				break
			case "channels":
				fmt.Printf("\n[系统有以下渠道, 输入渠道名称或者自定义名称] \n %v \n", response.Data["channels"])
				channel := ""
				for channel == "" {
					fmt.Scanln(&channel)
				}
				client.Channel = channel
				client.connect() // 连接到服务器
				fmt.Printf("[连接到 %s 通道成功] \n", client.Channel)
			}
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}
