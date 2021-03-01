package process

import "fmt"

//因为UserMgr 实例在服务器端有且只有一个
//因为在很多地方，都会使用到，因此，我们
//将其定义为全局变量
var (
	userMgr *UserMgr
)

type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

//完成对userMgr初始化工作
func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

//完成对onlineUsers添加&修改
func (this *UserMgr) AddOnlineUser(up *UserProcess) {
	this.onlineUsers[up.UserId] = up

}

//完成对onlineUsers删除
func (this *UserMgr) DelOnlineUser(userId int) {
	delete(this.onlineUsers, userId)

}

//返回当前在线所有用户
func (this *UserMgr) GetAllOnLineUsers() map[int]*UserProcess {
	return this.onlineUsers
}

//通过UserId返回对应UserProcess
func (this *UserMgr) GetOnlineUserById(userId int) (up *UserProcess, err error){

	//如何从map取出一个值，带检测方式
	up, ok := this.onlineUsers[userId]
	if !ok { //说明，查找用户，当前不在线
		//
		err = fmt.Errorf("用户%d 不存在", userId)
		return
	}
	return
}
