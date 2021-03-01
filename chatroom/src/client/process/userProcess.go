package process

import (
	"client/utils"
	"common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

type UserProcess struct {
	//
}

func (this *UserProcess) Register(userId int,
	userPwd string, userName string) (err error) {
	//1.链接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")	//理论上链接目标地址需要从配置文件中读出
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return err
	}

	defer conn.Close()

	//2.准备通过conn发送消息给服务器
	var msg message.Message
	msg.Type = message.RegisterMsgType

	//3. 创建一个LoginMsg结构体
	var registerMsg message.RegisterMsg
	registerMsg.User.UserId = userId
	registerMsg.User.UserPwd = userPwd
	registerMsg.User.UserName = userName

	// 4.将registerMsg序列化
	data, err := json.Marshal(registerMsg)
	if err != nil {
		fmt.Println("json.Mashal err=", err)
		return err
	}

	// 5. 把data赋值给msg.Data字段
	msg.Data = string(data)

	// 6. 将发送消息进行序列化
	data, err = json.Marshal(msg)
	if err != nil {
		fmt.Println("json.Mashal send err=", err)
		return err
	}

	//创建一个Transfer实例
	tf := &utils.Transfer{	//QA:使用地址符赋值的意思
		Conn : conn,
	}

	//发送data给服务器
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("注册发送信息错误 err=", err)
	}

	msg, err = tf.ReadPkg()

	//将msg的data部分反序列化， LoginResMsg
	if err != nil {
		fmt.Println("readPkg(conn) err=", err)
		return
	}
	var registerRspMsg message.RegisterRspMsg
	err = json.Unmarshal([]byte(msg.Data), &registerRspMsg)
	if registerRspMsg.Code == 200 {
		fmt.Println("注册成功，可重新登录")
		//os.Exit(0)
	} else {
		fmt.Println(registerRspMsg.Error)
		//os.Exit(0)
	}
	return
}

//登录函数，完成登录
func (this *UserProcess)Login(userId int, userPwd string) (err error) {

	////开始定协议
	//fmt.Printf("userId=%d, userPwd=%s\n", userId, userPwd)
	//
	//return nil

	//1.链接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")	//理论上链接目标地址需要从配置文件中读出
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return err
	}

	defer conn.Close()

	//2.准备通过conn发送消息给服务器
	var msg message.Message
	msg.Type = message.LoginMsgType

	//3. 创建一个LoginMsg结构体
	var loginMsg message.LoginMsg
	loginMsg.UserId = userId
	loginMsg.UserPwd = userPwd
	//loginMsg.UserName = "Z.Y,Zhao"

	// 4.将loginMsg序列化
	data, err := json.Marshal(loginMsg)
	if err != nil {
		fmt.Println("json.Mashal err=", err)
		return err
	}
	// 5. 把data赋值给msg.Data字段
	msg.Data = string(data)

	// 6. 将发送消息进行序列化
	sendata, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("json.Mashal send err=", err)
		return err
	}

	// 7.data就是需要发送的消息,发送长度和消息
	// 7.1 先把data长度发送给服务器
	// 先获取到data长度->转换成一个表示长度的byte slice
	//pkgLen := binary.ByteOrder()
	var pkgLen uint32 //几乎约等于4g容量
	pkgLen = uint32(len(sendata))
	//var bytes [4]byte
	//var bytes [4]byte
	bufLen := make([]byte, 4)
	binary.BigEndian.PutUint32(bufLen[0:4], pkgLen)

	//发送长度
	n, err := conn.Write(bufLen)
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return err
	}

	fmt.Printf("客户端，放松消息的长度=%d, 内容=%s\n", len(sendata), string(sendata))

	//发送消息本身
	_, err = conn.Write(sendata)
	if err != nil {
		fmt.Println("conn.Write(senddata) fail", err)
		return err
	}

	//创建一个Transfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}
	msg, err = tf.ReadPkg()

	//将msg的data部分反序列化， LoginResMsg
	if err != nil {
		fmt.Println("readPkg(conn) err=", err)
		return
	}
	var loginRspMsg message.LoginResMsg
	err = json.Unmarshal([]byte(msg.Data), &loginRspMsg)
	if loginRspMsg.Code == 200 {
		//初始化curUser
		curUser.Conn = conn
		curUser.UserId = userId
		curUser.UserStatus = message.UserOnline

		//可以显示当前在线用户列表，遍历loginRspMsg.UsersId
		fmt.Println("当前在线用户类标如下：")
		for _, v := range loginRspMsg.UserIds {
			// 不显示登录用户自身信息
			if v == userId {
				continue
			}
			fmt.Println("用户id:\t", v)
			//完成 客户端 onlineUsersClient初始化
			user := &message.User{
				UserId: v,
				UserStatus: message.UserOnline,
			}
			onlineUsersClient[v] = user


		}
		fmt.Print("\n\n")

		//这里还需要在客户端启动一个协程
		//该协程保持和服务器端的通讯，如果服务器有数据推送
		//则接收并显示在客户端的纵终端
		go serverProcessMsg(conn)

		//1.显示登录成功的菜单[循环]..
		for {
			ShowMenu()
		}
	} else {
		fmt.Println(loginRspMsg.Error)
	}

	return nil

}
