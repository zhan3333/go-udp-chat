package app

import (
	"encoding/json"
	"fmt"
	"net"
)

type Server struct {
	Listener   *net.UDPConn
	Host       string
	Port       int
	ChannelObj ChannelObj
	ClientMap  map[string]*Client
}

func (server *Server) CreateChannel(name string, creator string) {
	server.ChannelObj.add(&Channel{
		Name:      name,
		Creator:   creator,
		Clients:   []*Client{},
		ClientMap: map[string]*Client{},
	})
}

func (server *Server) Start() error {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(server.Host), Port: server.Port})
	server.Listener = listener
	return err
}

func (server Server) Receive() {
	data := make([]byte, 1024)
	for {
		l, remoteAddr, err := server.Listener.ReadFromUDP(data)
		if err != nil {
			_ = fmt.Errorf("[无效的UDP数据包] %s", err.Error())
		}
		var message ClientMessage
		if err := json.Unmarshal(data[:l], &message); err != nil {
			// 解码错误
			fmt.Printf("[解码失败 %s] %s \n", data, err.Error())
		}
		request := Request{
			Raw:  data,
			Addr: remoteAddr,
			Body: message,
		}
		server.handleMessage(request)
	}
}

func (server Server) getClient(addr string) *Client {
	return server.ClientMap[addr]
}

func (server *Server) addClient(client *Client) {
	server.ClientMap[client.Addr.String()] = client
}

func (server Server) send(response Response, client Client) {
	bytes, _ := json.Marshal(response)
	l, err := server.Listener.WriteToUDP(bytes, client.Addr)
	if err != nil {
		_ = fmt.Errorf("[发送消息 (%s) 失败到 %s]: %s", string(bytes), client.Addr, err.Error())
		return
	}
	fmt.Printf("[Send success %d bites] \n", l)
}

func (server Server) handleMessage(request Request) {
	fmt.Printf("[接收到消息 %s] %s \n", request.Addr, string(request.Raw))

	switch request.Body.Type {
	case "message":
		client := server.getClient(request.Addr.String())
		if client == nil {
			fmt.Printf("client: %v \n", client)
		}
		server.broadcast(request.Body.Message, client)
		break
	case "connect":
		name := request.getStr("name")
		channel := request.getStr("channel")
		if name == "" || channel == "" {
			fmt.Printf("Name: %s, Channel: %s 必须传入", name, channel)
			break
		}
		client := Client{
			Name: name,
			Addr: request.Addr,
		}
		server.addClient(&client)
		server.ChannelObj.getOrAdd(channel).addClient(&client)
		break
	case "get_channels":
		client := Client{
			Name: "",
			Addr: request.Addr,
		}
		names := server.GetChannelsNames()
		server.send(Response{
			Type:    "channels",
			Message: "",
			Data: map[string]interface{}{
				"channels": names,
			},
		}, client)
		break
	default:
		break
	}
}

type ChannelObj struct {
	Channels map[string]*Channel
}

func (channelObj ChannelObj) get(name string) *Channel {
	return channelObj.Channels[name]
}

func (channelObj *ChannelObj) getOrAdd(name string) *Channel {
	channel := channelObj.Channels[name]
	if channel == nil {
		newChannel := Channel{
			Name:      name,
			Creator:   "system",
			Clients:   []*Client{},
			ClientMap: map[string]*Client{},
		}
		channelObj.add(&newChannel)
		channel = &newChannel
	}
	return channel
}

func (channelObj ChannelObj) del(name string) {
	delete(channelObj.Channels, name)
}

func (channelObj *ChannelObj) add(channel *Channel) {
	channelObj.Channels[channel.Name] = channel
}

// 找到所在的通道
func (server *Server) findChannel(client *Client) *Channel {
	for _, channel := range server.ChannelObj.Channels {
		if channel.ClientMap[client.Addr.String()] != nil {
			client.Name = channel.ClientMap[client.Addr.String()].Name
			return channel
		}
	}
	// 未找到通道, 放入到默认通道中
	defaultChannel := server.ChannelObj.getOrAdd("default")
	defaultChannel.addClient(client)
	return defaultChannel
}

func (server Server) broadcast(message string, exceptClient *Client) {
	channel := server.findChannel(exceptClient)
	fmt.Printf("[广播消息 %s => %s] \n", channel.Name, message)
	fmt.Printf("[Clients] %v \n", channel.getClientsNames())
	response := Response{
		Type:     "message",
		Message:  message,
		Channel:  channel.Name,
		SendUser: exceptClient.Name,
	}
	for _, client := range channel.Clients {
		if !client.isSame(*exceptClient) {
			fmt.Printf("[Send to %s %s]: %s \n", client.Name, client.Addr.String(), message)
			server.send(response, *client)
		}
	}
}

func (server Server) GetChannels() map[string]*Channel {
	return server.ChannelObj.Channels
}

// 获取通道名称列表
func (server Server) GetChannelsNames() []string {
	var names []string
	for name, _ := range server.ChannelObj.Channels {
		names = append(names, name)
	}
	return names
}
