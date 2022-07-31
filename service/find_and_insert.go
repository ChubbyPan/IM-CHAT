package service

import (
	"context"
	"fmt"
	"time"

	logging "github.com/sirupsen/logrus"
	"main.go/conf"
	"main.go/model/ws"
)

//插入到mongodb数据库当中
func InsertMsg(database string, id string, content string, read uint, expire int64) error {
	// 没有该id的集合则创建该集合,有该id的集合直接使用
	collection := conf.MongoDBClient.Database(database).Collection(id)
	fmt.Printf("this message belongs to %s", id)
	comment := ws.MongoDBMsg{
		Content:   content,
		StartTime: string(time.Now().Unix()),
		EndTime:   string(time.Now().Unix() + expire),
		Read:      read,
	}
	_, err := collection.InsertOne(context.TODO(), comment) //必须要传context，不知道传啥就穿context todo
	if err != nil {
		logging.Info(err)
		fmt.Println("insert message failed!")
	}
	return err
}
