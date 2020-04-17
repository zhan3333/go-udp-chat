# Golang UDP 聊天室

本项目通过 UDP 通信协议, 建立客户端与服务端之间的通信. 

## 运行流程

- 运行服务端

```shell script
cd server
go run main.go
```

- 运行客户端

```shell script
cd client
go run main.go
```

## 功能

- [x] client 在 channel 中群发消息
- [x] channel 中其它 client 接收消息
- [] client 自定义 channel
- [] 心跳检测 client 是否在线
- [] client 切换当前 channel
- [] client 重名检查