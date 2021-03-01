package process

import (
	"common/message"
	"encoding/json"
	"fmt"
	"net"
	"server/model"
	"server/utils"
)

type UserProcess struct {
	//字段定义
	Conn net.Conn
	//增加一个字段，表示该Conn是哪个用户
	UserId int

}

//加入通知所有在线用户的处理
//当前上线userID通知其他
func (this *UserMgr) NotifyOthersOnlineUser(userId int) {

	//遍历 onlineUsers, 然后一个接一个返送NotifyUsersStatusMsg
	for id, up := range this.onlineUsers {

		//非通知自身
		if id == userId {
			continue
		}
		//开始通知，每个在线用户单独定义个方法实现
		up.NotifyMeOnline(userId)
	}
}

func (this *UserProcess) NotifyOthersOnlineUser(userId int) {

	//遍历 onlineUsers, 然后一个接一个返送NotifyUsersStatusMsg
	for id, up := range userMgr.onlineUsers {

		//非通知自身
		if id == userId {
			continue
		}
		//开始通知，每个在线用户单独定义个方法实现
		up.NotifyMeOnline(userId)
	}
}



func (this *UserProcess) NotifyMeOnline(userId int) {
	//组装NotifiyUserStatusMsg
	var msg message.Message
	msg.Type = message.NotifyUserStatusMsgType

	var notifyUserStatusMsg message.NotifyUserStatusMsg
	notifyUserStatusMsg.UserId = userId
	notifyUserStatusMsg.Status = message.UserOnline

	//将notifyUserStausMsg序列化

	data ,err := json.Marshal(notifyUserStatusMsg)
	if err != nil  {
		fmt.Println("json.Marshal err= ", err)
		return
	}

	// 将序列化后的notifyUserStatusMsg赋值给msg.data
	msg.Data = string(data)

	//对msg再次序列化，准备发送
	data, err = json.Marshal(msg)
	if err != nil {
		fmt.Println("json.Marshal err= ", err)
		return
	}

	//发送
	tf := &utils.Transfer{
		Conn: this.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyMeOnlie err=", err)
		return
	}

}


func (this *UserProcess) serverProcessRegister(msg *message.Message) (err error) {

	// 1. 既存msg中取出msg.Data,并直接反序列化成registerMsg
	var registerMsg message.RegisterMsg
	err = json.Unmarshal([]byte(msg.Data), &registerMsg)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return err
	}

	//1. 先声明一个rspMsg
	var rspMsg message.Message
	rspMsg.Type = message.RegisterRspMsgType

	//2. 再声明一个LoginResMsg
	var registerRspMsg message.RegisterRspMsg

	// 先去数据库redis完成验证
	//1. 使用model.MyUserDao去redis去验证
	err = model.MyUserDao.Register(&registerMsg.User)

	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerRspMsg.Code = 505
			registerRspMsg.Error = model.ERROR_USER_EXISTS.Error()

		} else {
			registerRspMsg.Code = 506
			registerRspMsg.Error = "Unknwon Error"
		}
	} else {
		registerRspMsg.Code = 200
	}

	//3. 将loginRspMsg序列化
	data, err := json.Marshal(registerRspMsg)
	if err != nil {
		fmt.Println("json.Marshal fail err=", err)
		return
	}

	//4. 将data赋值给resMsg
	rspMsg.Data = string(data)

	//5. 对rspMsg进行序列化，准备发送
	data, err = json.Marshal(rspMsg)
	if err != nil {
		fmt.Println("json.Marshal rspMSg fail err=", err)
		return
	}
	//6. 发送data，将其封装到writePkg函数中
	// 因为十一哦那个分层模式(mvc)，先创建一个Transfer实例，然后读取。
	tf := &utils.Transfer{
		Conn : this.Conn,
	}
	err = tf.WritePkg(data)
	return
}

//编写一个函数serverProcessLogin函数，专门处理登录请求
func (this *UserProcess)serverProcessLogin(msg *message.Message) (err error) {

	//核心代码
	// 1. 既存msg中取出msg.Data,并直接反序列化成LoginMsg
	var loginMsg message.LoginMsg
	err = json.Unmarshal([]byte(msg.Data), &loginMsg)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return err
	}
	//1. 先声明一个rspMsg
	var rspMsg message.Message
	rspMsg.Type = message.LoginResMsgType

	//2. 再声明一个LoginResMsg
	var loginRspMsg message.LoginResMsg

	// 先去数据库redis完成验证
	//1. 使用model.MyUserDao去redis去验证
	user, err := model.MyUserDao.Login(loginMsg.UserId, loginMsg.UserPwd)

	if err != nil {
		//loginRspMsg.Code = 500
		//loginRspMsg.Error = "用户不存在，请注册再使用..."
		//测试成功后，返回具体的错误信息

		if err == model.ERROR_USER_NOTEXISTS {
			loginRspMsg.Code = 500
			loginRspMsg.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginRspMsg.Code = 403
			loginRspMsg.Error = err.Error()
		} else {
			loginRspMsg.Code = 505
			loginRspMsg.Error = "服务器内部错误..."
		}

	} else {
		loginRspMsg. Code = 200
		//用户登录成功后，把登录成功的用户放入到UserMgr中
		//将登录成功的用户userId 赋给this
		this.UserId = loginMsg.UserId
		userMgr.AddOnlineUser(this)

		//userMgr.NotifyOthersOnlineUser(this.UserId)
		this.NotifyOthersOnlineUser(this.UserId)

		//将当前用户的id，放入到loginRspMsg.UsersId
		// 遍历userMgr.onLineUsers
		for id, _ := range userMgr.onlineUsers {
			loginRspMsg.UserIds = append(loginRspMsg.UserIds, id)
		}

		fmt.Println(user, "登录成功!")
	}

	//3. 将loginRspMsg序列化
	data, err := json.Marshal(loginRspMsg)
	if err != nil {
		fmt.Println("json.Marshal fail err=", err)
		return
	}

	//4. 将data赋值给resMsg
	rspMsg.Data = string(data)

	//5. 对rspMsg进行序列化，准备发送
	data, err = json.Marshal(rspMsg)
	if err != nil {
		fmt.Println("json.Marshal rspMSg fail err=", err)
		return
	}
	//6. 发送data，将其封装到writePkg函数中
	// 因为十一哦那个分层模式(mvc)，先创建一个Transfer实例，然后读取。
	tf := &utils.Transfer{
		Conn : this.Conn,
	}
	err = tf.WritePkg(data)
	return
}

