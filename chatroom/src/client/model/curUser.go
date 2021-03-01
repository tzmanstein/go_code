package model

import (
	"common/message"
	"net"
)

//在客户端，多场景使用curUser，因此需要将此变量声明为全局
type CurUser struct {
	Conn net.Conn
	message.User //匿名结构体，实现继承
}
