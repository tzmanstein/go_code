package process

import (
	"common/message"
	"fmt"
	"io"
	"net"
	"server/utils"
)

//先创建一个Processor的结构体
type Processor struct {
	Conn net.Conn
}

//Process2
func (this *Processor) ProcessSwitcher() (err error) {
	//循环读取客户端发送信息
	for {
		//将读取数据包，直接封装成一个函数readPkg(),返回Message, Err
		//创建一个Transfer实例完成读包任务。
		tf := &utils.Transfer{
			Conn: this.Conn,
		}

		msg, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出，服务器正常退出")
				return err
			} else {
				fmt.Println("readPkg err=", err)
				return err
			}
		}
		//fmt.Println("msg=", msg)
		err = this.serverProcessMsg(&msg)

		if err != nil {
			return err
		}
	}
}

//编写一个ServerProcessMsg函数
//功能：根据客户端发送消息种类不同，决定调用那个函数来处理
func (this *Processor) serverProcessMsg(msg *message.Message) (err error) {

	switch msg.Type {
	case message.LoginMsgType:
		//处理登录逻辑
		// 创建一个UserProcess
		up := &UserProcess{
			Conn : this.Conn,
		}
		err = up.serverProcessLogin(msg)

	case message.RegisterMsgType:
		//处理注册
		up := &UserProcess{
			Conn : this.Conn,
		}
		err = up.serverProcessRegister(msg)

	case message.SmsMsgType:
		fmt.Println("群发消息msg=", msg)
		smsProcess := SmsProcess{}
		smsProcess.SendGroupMsg(msg)

	default:
		fmt.Println("消息类型不存在，无法处理...")
	}
	return
}
