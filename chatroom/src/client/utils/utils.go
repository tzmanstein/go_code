package utils

import (
	"common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

//这里将这些方法关联到结构体中
type Transfer struct {
	//需要哪些字段
	Conn net.Conn
	Buf [8096]byte	//传输时，使用缓存
}

func (this *Transfer) ReadPkg()(msg message.Message, err error) {
	//buf := make([]byte, 8096)
	fmt.Println("读取客户发送的数据...")
	//conn.Read在conn没有关闭的情况下，才会阻塞
	//
	_, err = this.Conn.Read(this.Buf[0:4])
	if err != nil {
		fmt.Println("conn.Read err=", err)
		//err = errors.New("read pkg header error")
		return
	}
	fmt.Println("读到的buf=", this.Buf[:4])

	//根据bug[:4]转换成1个 uint32类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.Buf[0:4])

	//根据 pkgLen读取消息内容
	n, err := this.Conn.Read(this.Buf[:pkgLen]) //从conn中读取pkgLen长度字节放入buf中
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Read fail err=", err)
		//err = errors.New("read pkg body error")
		return
	}

	//把pkgLen，反序列化为 -> Message
	//var msg message.Message
	err = json.Unmarshal(this.Buf[:pkgLen], &msg) //msg必须传递地址,否则为空不能正常取得
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		//err = errors.New("read json.Unmarshal error")
		return
	}

	return
}

func (this *Transfer)WritePkg(data []byte) (err error) {
	//先发送一个长度给对方
	var pkgLen uint32 //几乎约等于4g容量
	pkgLen = uint32(len(data))
	//var bytes [4]byte
	//var bytes [4]byte
	//bufLen := make([]byte, 4)
	binary.BigEndian.PutUint32(this.Buf[0:4], pkgLen)

	//发送长度
	n, err := this.Conn.Write(this.Buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return err
	}

	//发送data本身
	n, err = this.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return err
	}

	return
}
