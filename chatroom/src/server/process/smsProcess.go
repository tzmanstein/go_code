package process

import (
	"common/message"
	"encoding/json"
	"fmt"
	"net"
	"server/utils"
)

type SmsProcess struct {
	//.. 暂时不需要字段，
}

//转发消息 lesson340
func (this *SmsProcess) SendGroupMsg(msg *message.Message) {

	//遍历服务器端的onlineUsers map[int]*UserProcess,
	//将消息转发出去。

	//取出msg的内容SmsMsg
	var smsMsg message.SmsMsg
	err := json.Unmarshal([]byte(msg.Data), &smsMsg)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	for id, up := range userMgr.onlineUsers {

		//需要把自己过滤掉，不要发给自己
		if id == smsMsg.UserId {
			continue
		}

		this.SendMsgToSingleOnlineUser(data, up.Conn)

	}
}

func (this *SmsProcess) SendMsgToSingleOnlineUser (data []byte, conn net.Conn) {

	tf := &utils.Transfer{
		Conn: conn,
	}
	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("转发消息失败 err=", err)
	}

}
