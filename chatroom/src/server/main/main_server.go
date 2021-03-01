package main

import (
	"fmt"
	"net"
	"server/model"
	"server/process"
	"time"
)

func init() {
	//服务器启动时，初始化redis的连接池
	initPool("localhost:6379", 16, 0, 300*time.Second)
	initUserDao()
}

//声明一个函数，完成对UserDao的初始化任务
func initUserDao() {
	//这里的pool 本身就是一个全局变量
	//需要注意初始化的顺序问题
	//initPool, 在 initUserDao
	model.MyUserDao = model.NewUserDao(pool)

}

func processMain(conn net.Conn) {

	defer conn.Close()

	//调用总控，创建
	processor := &process.Processor{
		Conn : conn,
	}
	err := processor.ProcessSwitcher()
	if err != nil {
		fmt.Println("客户端和服务器端的协程错误 err=", err)
		return
	}
}

func main() {


	fmt.Println("服务器在8889端口监听")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Println("net.Listen err=", err)
		return
	}
	defer listen.Close()
	//监听成功，等待Client链接服务器
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=", err)
		}
		//链接成功，启动一个routine和客户保持通讯。
		go processMain(conn)

	}
}
