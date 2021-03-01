package process

import (
	"client/utils"
	"common/message"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

//显示登录成功后的界面...
func ShowMenu() {
	//1.显示在线用户列表
	//2.发送信息
	//3.信息列表
	//4.退出系统

	fmt.Printf("-------恭喜%s登录成功-------\n", curUser.UserName)
	fmt.Println("-------1.显示在线用户列表-------")
	fmt.Println("-------2.发送信息-------")
	fmt.Println("-------3.信息列表-------")
	fmt.Println("-------4.退出系统-------")
	fmt.Println("请选择(1-4):")

	var key int
	var content string
	//经常使用SmsProcess实例，定义在switch外部，防止反复实例化
	smsPro := &SmsProcess{}

	fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		//fmt.Println("显示在线用户列表-")
		outputOnlienUser()
	case 2:
		fmt.Println("发送消息-")
		fmt.Scanf("%s\n", &content)

		smsPro.SendGroupMsg(content)

	case 3:
		fmt.Println("信息列表-")
	case 4:
		fmt.Println("选择退出系统-")
		os.Exit(0)
	default:
		fmt.Println("啥也不是，重选")

	}
}

func serverProcessMsg(conn net.Conn) {

	//创建一个transfer实例，不停的读取服务器发送的消息
	//1.
	tf := &utils.Transfer{
		Conn: conn,
	}

	//只要对方不关闭链接一直处于待机状态
	for {
		//服务器端保持链接
		fmt.Println("客户端正在等待读取服务器发送的消息")
		msg, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.ReadPkg err=", err)
			return
		}
		//如果读取到消息，又是下一步处理逻辑
		//fmt.Printf("msg=%v\n", msg)

		switch msg.Type {
			case message.NotifyUserStatusMsgType:	//有人上线
				//处理
				//1.取出NotifyStatusMsg
				//2. 将取得用户加入到客户端维护的用户列表中（map[int]User形式保存
				var notifyUserStatuMsg message.NotifyUserStatusMsg
				json.Unmarshal([]byte(msg.Data), &notifyUserStatuMsg)
				updateUserStatus(&notifyUserStatuMsg)
			case message.SmsMsgType: //有人群发消息
				outputGroupMsg(&msg)

			default:
				fmt.Println("Unknown msg")

		}
	}

}