package process

import (
	"client/utils"
	"common/message"
	"encoding/json"
	"fmt"
)

type SmsProcess struct {

}


//发送群聊消息
func (this *SmsProcess) SendGroupMsg(content string) (err error) {

	//1. 创建一个Msg
	var msg message.Message
	msg.Type = message.SmsMsgType

	//2. 创建一个SmsMsg实例
	var smsMsg message.SmsMsg
	smsMsg.Content = content //内容
	smsMsg.UserId = curUser.UserId
	smsMsg.UserStatus = curUser.UserStatus

	//3.序列化 smsMsg
	data ,err := json.Marshal(smsMsg)
	if err != nil {
		fmt.Println("SendGroupMsg smsMsg json.Marshal fail =", err.Error())
		return
	}

	msg.Data = string(data)
	//4. 再次序列化
	data, err = json.Marshal(msg)
	if err != nil {
		fmt.Println("SendGroupMsg msg json.Marshal fail =", err.Error())
		return
	}

	//5.发送给服务器
	tf := &utils.Transfer{
		Conn : curUser.Conn,
	}

	//6. 发送
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("SendGroupMsg err=", err)
		return
	}

	return
}