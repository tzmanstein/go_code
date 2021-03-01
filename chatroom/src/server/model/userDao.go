package model

import (
	"common/message"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

//我们在服务器启动后，就初始化一个userDao实例
//把次成为全局变量，在需要和redis操作好似，就直接使用
var (
	MyUserDao *UserDao
)


//定义一个UserDao结构体
//完成对User 结构体的各种操作
type UserDao struct {
	pool *redis.Pool

}

//使用工厂模式，创建一个UserDao实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao){
	userDao = &UserDao{
		pool: pool,
	}
	return
}

//思考一下在UserDao，应该提供哪些方法
//1. 根据用户id返回一个User实例和err
func (this *UserDao) getUserById(conn redis.Conn, id int) (user *User, err error){

	//通过给定id 去 redis查询这个用户
	res, err := redis.String(conn.Do("hget", "users", id))
	if err != nil {
		if err == redis.ErrNil {	//表示在 users hash中，没有找到对应id
			err = ERROR_USER_NOTEXISTS
		}
		return
	}
	user = &User{}

	//这里我们需要把res，反序列化成User实例
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	fmt.Println("userid=", user.UserId)
	return
}

//完成登录的校验Login
//1.Login完成对用户的验证
//2.如果用户的id和pwd都正确，则返回一个user实例
//3. 如果用户的id或pwd有错误，则返回对应的错误信息
func (this *UserDao) Login(userId int, userPwd string)(user *User, err error) {

	//先从UserDao 的连接池中取出一个连接
	conn := this.pool.Get()
	defer conn.Close()

	user, err = this.getUserById(conn, userId)
	if err != nil {

		return
	}

	//取得用户信息，进行校验
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}
	return
}


func (this *UserDao) Register(user *message.User) (err error) {
	//先从UserDao的链接池中取出一个连接
	conn := this.pool.Get()
	defer conn.Close()

	_, err = this.getUserById(conn, user.UserId)
	if err == nil {
		//用户存在
		err = ERROR_USER_EXISTS
		return
	}

	//此时，说明id在redis中不存在，可以完成注册
	data, err := json.Marshal(user)
	if err != nil {
		return
	}

	_, err = conn.Do("hset", "users", user.UserId, string(data))
	if err != nil {
		fmt.Println("保存注册用户错误 err=", err)
		return
	}

	return

}
