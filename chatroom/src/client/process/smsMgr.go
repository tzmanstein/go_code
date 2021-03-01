package process

import (
	"common/message"
	"encoding/json"
	"fmt"
)

func outputGroupMsg(msg *message.Message) {
	//显示
	//1.反序列化message.Message
	var smsMsg message.SmsMsg
	err := json.Unmarshal([]byte(msg.Data), &smsMsg)

	if err != nil {
		fmt.Println("json.Unmarshal err=", err.Error())
		return
	}

	//显示信息
	info := fmt.Sprintf("用户id：\t %d 对所有人说:\t %s", smsMsg.UserId, smsMsg.Content)
	//content := smsMsg.Content
	fmt.Println(info)
	fmt.Println()

}
