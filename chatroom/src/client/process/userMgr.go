package process

import (
	"client/model"
	"common/message"
	"fmt"
)

//go语言核心编程重点
//服务器，服务器中核心数据通道
// 通讯，协程，协议 后台服务器

//客户端维护的map，登录是进行初始化
var onlineUsersClient map[int]*message.User = make(map[int]*message.User, 10)
var curUser model.CurUser //登录成功后，完成对curUser初始化

//客户端显示当前在线用户
func outputOnlienUser() {
	// 遍历
	fmt.Println("当前在线用户列表：")
	for id, user := range onlineUsersClient {
		fmt.Println("用户id：\t", id)
		fmt.Println("用户=\t", user)
	}
}

//编写一个方法，处理返回的NotifyUserStatusMsg
func updateUserStatus(notifyUserStatusMsg *message.NotifyUserStatusMsg) {

	// 适当优化，当前用户是否存在检查
	user, ok := onlineUsersClient[notifyUserStatusMsg.UserId]
	if !ok {
		user = &message.User{
			UserId: notifyUserStatusMsg.UserId,
			UserStatus: notifyUserStatusMsg.Status,
		}
	}

	user.UserStatus = notifyUserStatusMsg.Status
	onlineUsersClient[notifyUserStatusMsg.UserId] = user

	outputOnlienUser()
}

