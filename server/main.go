package main

import (
	"fmt"
	"udp-chat/server/app"
)

func main() {
	var server = app.Server{
		Listener: nil,
		Host:     "127.0.0.1",
		Port:     9981,
		ChannelObj: app.ChannelObj{
			Channels: map[string]*app.Channel{},
		},
		ClientMap: map[string]*app.Client{},
	}
	err := server.Start()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Local: <%s> \n", server.Listener.LocalAddr().String())
	server.CreateChannel("default", "admin")
	server.Receive()
}
