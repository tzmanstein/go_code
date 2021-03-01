package message

//消息类型定义
const (
	LoginMsgType = "LoginMsg"
	LoginResMsgType = "LoginResMsg"
	RegisterMsgType = "RegisterMsg"
	RegisterRspMsgType = "RegisterRspMsg"
	NotifyUserStatusMsgType = "NotifyUserStatusMsg"
	SmsMsgType = "SmsMsg"
)

//定义用户状态的常量
const (
	UserOnline = iota
	UserOffline
	UserBusyStatus

)

//定义2个消息
type Message struct {
	Type string `json:"type"`	//消息类型
	Data string `json:"data"`	//消息内容
}

type LoginMsg struct {
	UserId int `json:"userId"`		//用户ID
	UserPwd string `json:"userPwd"`	//用户密码
	UserName string `json:"userName"`	//用户名
}

type LoginResMsg struct {
	Code int `json:"code"`		//返回状态码 500表示用户未注册， 200表示登录成功
	UserIds []int 				//增加再丢单，保存用户id的切片
	Error string `json:"error"`	//返回错误信息
}

type RegisterMsg struct {
	//..

	User User `json:"user"`
}

type RegisterRspMsg struct {
	Code int `json:"code"`	//返回状态码 400， 表示该用户已经存在， 200表示注册成功
	Error string `json:"error"`	//返回错误信息
}

//为了配合服务器端推送用户状态变化消息
type NotifyUserStatusMsg struct {
	UserId int `json:"userId"`	//用户id
	Status int `json:"status"`	//用户的状态
}

//增加一个SmsMsg 发送用消息
type SmsMsg struct {
	Content string `json:"content"` //内容
	User //匿名结构体，继承
}

// SmsRspMsg