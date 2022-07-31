package ws

type MongoDBMsg struct {
	Content   string `bson:"content"`   // 内容
	StartTime string `bson:"startTime"` //创建时间
	EndTime   string `bson:"endTime"`   //过期时间
	Read      uint   `bson:"read"`      //已读状态
}
