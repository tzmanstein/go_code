package main

import (
	"client/process"
	"fmt"
	"os"
)

//定义两个全局变量，一个用户Id，一个用户密码
var userId int
var userPwd string
var userName string

func main() {

	//接收用户选项
	var key int
	//判断是否继续显示
	var loop = true

	for loop {
		fmt.Println("--------------------欢迎登录多人聊天系统--------------------")
		fmt.Println("\t\t\t 1.登录聊天室")
		fmt.Println("\t\t\t 2.注册用户")
		fmt.Println("\t\t\t 3.退出系统")
		fmt.Println("\t\t\t 请选择（1-3)：")

		fmt.Scanf("%d\n", &key)
		switch key {
		case 1 :
			fmt.Println("1.登录聊天室")
			fmt.Println("请输入用户的ID")
			fmt.Scanf("%d\n", &userId)
			//fmt.Scanln(&userId)
			fmt.Println("请输入用户的密码")
			fmt.Scanf("%s\n", &userPwd)
			//loop = false
			up := &process.UserProcess{}
			up.Login(userId, userPwd)

		case 2 :
			fmt.Println("2.注册用户")
			fmt.Println("请输入用户id：")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户密码：")
			fmt.Scanf("%s\n", &userPwd)
			fmt.Println("请输入用户名字(nickName):")
			fmt.Scanf("%s\n", &userName)

			//2. 调用UserProcess，完成注册请求
			up := &process.UserProcess{}
			up.Register(userId, userPwd, userName)

			//loop = false
		case 3 :
			fmt.Println("3.退出系统")
			os.Exit(0)
			//loop = false
		default:
			fmt.Println("啥也不是,重来！")

		}
	}

	////根据用户输入，显示新的提示信息
	//if key == 1 {
	//	//用户登录
	//	fmt.Println("请输入用户的ID")
	//	fmt.Scanf("%d\n", &userId)
	//	//fmt.Scanln(&userId)
	//	fmt.Println("请输入用户的密码")
	//	fmt.Scanf("%s\n", &userPwd)
	//
	//	// 想把登录函数，写到另外一个文件. login.go
	//
	//	login(userId, userPwd) //需要重新调用
	//
	//	//使用新的程序结构，创建
	//
	//	//err := login(userId, userPwd)
	//	//if err != nil {
	//	//	fmt.Println("登录失败")
	//	//} else {
	//	//	fmt.Println("登录成功")
	//	//}
	//
	//} else if key == 2 {
	//	fmt.Println("进行用户注册的逻辑")
	//} else if key == 3 {
	//
	//}
}
