package errcode

var MsgFlags = map[int]string{
	SUCCESS:                 "ok",
	UpdatePasswordSuccess:   "修改密码成功",
	NotExistInentifier:      "该第三方账号未绑定",
	ERROR:                   "fail",
	InvalidParams:           "请求参数错误",
	ErrorDatabase:           "数据库操作出错,请重试",
	WebsocketSuccessMessage: "解析content内容信息",
	WebsocketSuccess:        "发送信息成功/请求历史纪录操作成功",
	WebsocketEnd:            "请求历史纪录，但没有更多记录了",
	WebsocketOnlineReply:    "对方在线应答成功",
	WebsocketOfflineReply:   "对方离线回答成功",
	WebsocketLimit:          "请求受到限制",
}
