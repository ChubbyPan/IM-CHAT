package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	logging "github.com/sirupsen/logrus"
	"main.go/cache"
	"main.go/conf"
	"main.go/pkg/errcode"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

//过期时间处理
const month = 60 * 60 * 24 * 30

//发送消息结构体
type SendMsg struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}

//接收消息结构体
type ReplyMsg struct {
	From    string `json:"from"`
	Code    int    `json:"code"`
	Content string `json:"content"`
}

//会话用户信息
type Client struct {
	ID     string
	SendID string
	Socket *websocket.Conn
	Send   chan []byte
}

//广播类，包括广播内容和源用户
type Broadcast struct {
	Client  *Client //源用户
	Massage []byte
	Type    int
}

//用户管理
type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan *Broadcast
	Reply      chan *Client
	Register   chan *Client
	Unregister chan *Client
}

//信息转json
type Massage struct {
	Sender    string `json:"sender.omitempty"`
	Recipient string `json:"recipient.omitempty"`
	Content   string `json:"content.omitempty"`
}

var Manager = ClientManager{
	Clients:    make(map[string]*Client), //参与连接的用户
	Broadcast:  make(chan *Broadcast),
	Reply:      make(chan *Client),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}

//连接标识
func CreateID(uid, toUid string) string {
	return uid + "->" + toUid
}

func WsHandler(c *gin.Context) {
	uid := c.Query("uid")
	toUid := c.Query("toUid")
	fmt.Println(uid, toUid)
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}).Upgrade(c.Writer, c.Request, nil) //升级成ws协议
	if err != nil {
		fmt.Println("ws upgrade failed!")
		logging.Info(err)
		http.NotFound(c.Writer, c.Request)
		return
	}
	//	创建一个用户实例
	client := &Client{
		ID:     CreateID(uid, toUid),
		SendID: CreateID(toUid, uid),
		Socket: conn,
		Send:   make(chan []byte),
	}
	// 用户注册到用户管理上
	Manager.Register <- client
	// fmt.Println(Manager.Register)
	go client.Read()
	go client.Write()

}

//read逻辑：
//1.正常读消息（需要缓存源用户ID，保证连接3个月内不过期） 2.如果是单向通信需要保证发送3条以内的消息来进行防骚扰（关闭连接） 3.消息全部读完需要关闭连接
func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c //用户离线操作
		_ = c.Socket.Close()
	}()
	for {
		c.Socket.PongHandler()
		sendMsg := new(SendMsg)
		// c.Socket.ReadMessage(&sendMsg) // string类型的消息
		err := c.Socket.ReadJSON(&sendMsg) // 序列化消息
		if err != nil {
			logging.Info(err)
			Manager.Unregister <- c //用户离线操作
			_ = c.Socket.Close()
			break
		}
		if sendMsg.Type == 1 { // 发送消息
			//Redis `GET key` command. It returns redis.Nil error when key does not exist.
			r1, _ := cache.RedisClient.Get(c.ID).Result()     // 1->2 caching 1?
			r2, _ := cache.RedisClient.Get(c.SendID).Result() // 2->1 caching 2?
			//防骚扰功能
			if r1 > "3" && r2 == "" { //1->2发了3条，但是2没有回复，停止1发送
				replyMsg := ReplyMsg{
					Code:    errcode.WebsocketLimit,
					Content: errcode.MsgFlags[errcode.WebsocketLimit],
				}
				msg, _ := json.Marshal(replyMsg) // 序列化
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				continue
			} else {
				//如果单向联系在三条之内或者双方已有联系，缓存客户端ID
				cache.RedisClient.Incr(c.ID)
				_, _ = cache.RedisClient.Expire(c.ID, time.Hour*24*30*3).Result() // 三个月过期
				//防止过快分手，建立连接之后三个月过期
			}
			Manager.Broadcast <- &Broadcast{
				Client:  c,
				Massage: []byte(sendMsg.Content),
			}
		} else if sendMsg.Type == 2 { // 获取历史消息
			timeT, err := strconv.Atoi(sendMsg.Content) //string to int
			if err != nil {
				timeT = 999999
			}
			results, _ := FindMore(conf.MongoDBName, c.SendID, c.ID, int64(timeT), 10) // 获取十条历史消息
			if len(results) > 10 {
				results = results[:10]
			} else if len(results) == 0 {
				replyMsg := ReplyMsg{
					Code:    errcode.WebsocketEnd,
					Content: "无消息",
				}
				msg, _ := json.Marshal(replyMsg) // 序列化
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				continue
			}
			for _, result := range results {
				replyMsg := ReplyMsg{
					From:    result.From,
					Content: result.Msg,
				}
				msg, _ := json.Marshal(replyMsg) // 序列化
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			}

		}
	}

}
func (c *Client) Write() {
	defer func() {
		_ = c.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{}) // 出错发送空消息
				return
			}
			replyMsg := ReplyMsg{
				Code:    errcode.WebsocketSuccessMessage,
				Content: fmt.Sprintf("%s", string(message)), //c.send是字节类型
			}
			msg, _ := json.Marshal(replyMsg)
			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)

		}
	}
}
