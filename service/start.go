package service

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	logging "github.com/sirupsen/logrus"
	"main.go/conf"
	"main.go/pkg/errcode"
)

func (manager *ClientManager) Start() {
	for {
		fmt.Println("----------监听管道通信---------")
		select {
		case conn := <-manager.Register:
			fmt.Printf("有新连接：%s \n", conn.ID)
			fmt.Println("Manager.Register:", conn)
			Manager.Clients[conn.ID] = conn //把连接放到用户管理上
			replyMsg := ReplyMsg{
				Code:    errcode.WebsocketSuccess,
				Content: errcode.MsgFlags[errcode.WebsocketSuccess],
			}
			msg, _ := json.Marshal(replyMsg)
			_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
		case conn := <-manager.Unregister:
			fmt.Printf("连接失败： %s \n", conn.ID)
			if _, ok := Manager.Clients[conn.ID]; ok {
				replyMsg := &ReplyMsg{
					Code:    errcode.WebsocketEnd,
					Content: errcode.MsgFlags[errcode.WebsocketEnd],
				}
				msg, _ := json.Marshal(replyMsg)
				_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
				//连接失败， 关闭通道和用户管理对象
				close(conn.Send)
				delete(Manager.Clients, conn.ID)
			}
		case broadcast := <-Manager.Broadcast:
			message := broadcast.Massage
			sendID := broadcast.Client.SendID
			flag := false // 默认消息接收对象不在线
			// 遍历用户管理，查看消息接收对象是否在线
			for id, conn := range manager.Clients {
				if id != sendID {
					continue
				}
				//如果是在线状态
				select {
				case conn.Send <- message:
					flag = true
				default:
					close(conn.Send)
					delete(Manager.Clients, conn.ID)
				}
			}
			id := broadcast.Client.ID //发送消息方ID
			//如果对方在线，返回对方在线应答信息，将历史信息插入到MongoDB中
			if flag {
				replyMsg := &ReplyMsg{
					Code:    errcode.WebsocketOnlineReply,
					Content: errcode.MsgFlags[errcode.WebsocketOnlineReply],
				}
				msg, _ := json.Marshal(replyMsg)
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				err := InsertMsg(conf.MongoDBName, id, string(message), 1, int64(3*month)) // 默认1表示已读状态
				if err != nil {
					fmt.Println("Insert msg error!", err)
					logging.Info(err)
				}
			} else {
				replyMsg := &ReplyMsg{
					Code:    errcode.WebsocketOfflineReply,
					Content: errcode.MsgFlags[errcode.WebsocketOfflineReply],
				}
				msg, _ := json.Marshal(replyMsg)
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				err := InsertMsg(conf.MongoDBName, id, string(message), 0, int64(3*month))
				if err != nil {
					fmt.Println("Insert msg error!", err)
					logging.Info(err)
				}
			}
		}
	}

}
