package service

import (
	"context"
	"fmt"
	"time"

	"sort"

	logging "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"

	"main.go/conf"
	"main.go/model/ws"
)

type SendSortMsg struct {
	Content  string `json:"content"`
	Read     uint   `json:"read"`
	CreateAt int64  `json:"create_at"`
}

//插入到mongodb数据库当中
func InsertMsg(database string, id string, content string, read uint, expire int64) error {
	// 没有该id的集合则创建该集合,有该id的集合直接使用
	collection := conf.MongoDBClient.Database(database).Collection(id)
	fmt.Printf("this message belongs to %s", id)
	comment := ws.MongoDBMsg{
		Content:   content,
		StartTime: time.Now().Unix(),
		EndTime:   time.Now().Unix() + expire,
		Read:      read,
	}
	_, err := collection.InsertOne(context.TODO(), comment) //必须要传context，不知道传啥就穿context todo
	if err != nil {
		logging.Info(err)
		fmt.Println("insert message failed!")
	}
	return err
}

//获取历史信息
func FindMore(database string, sendID string, id string, time int64, pageSize int) (results []ws.Result, err error) {
	var resultsMe []ws.MongoDBMsg  // id
	var resultsYou []ws.MongoDBMsg // sendID
	sendIdCollection := conf.MongoDBClient.Database(database).Collection(sendID)
	idCollection := conf.MongoDBClient.Database(database).Collection(id)
	sendIdTimeCursor, _ := sendIdCollection.Find(context.TODO(), bson.D{{}})
	idTimeCursor, _ := idCollection.Find(context.TODO(), bson.M{})

	//测试集合中的数据是否正常输出
	for idTimeCursor.Next(context.TODO()) {
		r := new(ws.MongoDBMsg)
		err := idTimeCursor.Decode(&r)
		if err != nil {
			fmt.Printf("解析失败 %s", err)
		}
		fmt.Println("msg:", r)
	}

	err = sendIdTimeCursor.All(context.TODO(), &resultsYou) // sendId 对面发过来的
	if err != nil {
		fmt.Printf("解析失败 %s", err)
	}
	err = idTimeCursor.All(context.TODO(), &resultsMe) // Id 发给对面的
	if err != nil {
		fmt.Printf("解析失败 %s", err)
	}
	fmt.Println("test FindMore", resultsMe, resultsYou)
	results, _ = AppendAndSort(resultsMe, resultsYou)
	// fmt.Println("test AppendAndSort", resultsMe, resultsYou)
	return
}

func AppendAndSort(resultMe, resultYou []ws.MongoDBMsg) (results []ws.Result, err error) {
	for _, r := range resultMe {
		sendSort := SendSortMsg{ //返回结果需要的信息格式
			Content:  r.Content,
			Read:     r.Read,
			CreateAt: r.StartTime,
		}
		result := ws.Result{ // 返回结果需要的所有信息内容,包括传送者
			StartTime: r.StartTime,
			Msg:       fmt.Sprintf("%s", sendSort),
			From:      "me",
		}
		results = append(results, result)
	}
	for _, r := range resultYou {
		sendSort := SendSortMsg{ //返回结果需要的信息格式
			Content:  r.Content,
			Read:     r.Read,
			CreateAt: r.StartTime,
		}
		result := ws.Result{ // 返回结果需要的所有信息内容,包括传送者
			StartTime: r.StartTime,
			Msg:       fmt.Sprintf("%s", sendSort),
			From:      "you",
		}
		results = append(results, result)
	}
	// fmt.Println(results)
	sort.Slice(results, func(i, j int) bool { return results[i].StartTime < results[j].StartTime })
	return results, nil
}
