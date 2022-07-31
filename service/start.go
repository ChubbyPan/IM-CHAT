package service

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"main.go/pkg/errcode"
)

func (manager *ClientManager) Start() {
	for {
		fmt.Println("----------监听管道通信---------")
		select {
		case conn := <-manager.Register:
			fmt.Printf("有新连接：%v", conn.ID)
			fmt.Println("Manager.Register:", conn)
			Manager.Clients[conn.ID] = conn //把连接放到用户管理上
			replyMsg := ReplyMsg{
				Code:    errcode.WebsocketSuccess,
				Content: errcode.MsgFlags[errcode.WebsocketSuccess],
			}
			msg, _ := json.Marshal(replyMsg)
			_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
		}
	}

}
