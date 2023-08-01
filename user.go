package main

import "net"

type User struct {
	id       string
	userName string
	addr     string
	c        chan string //管道中传输字符串
	conn     net.Conn
	server   *Server
}

// 用户上线的功能
func (this *User) online() {
	this.server.mapLock.Lock()
	this.server.onlineMap[this.id] = this
	this.server.mapLock.Unlock()

	this.server.BroadMessage(this, "sign in") //广播消息
}

// 用户下线
func (this *User) offline() {

	this.server.mapLock.Lock()
	delete(this.server.onlineMap, this.id)
	this.server.mapLock.Unlock()
	this.server.BroadMessage(this, "sign out")
}

// 用户处理消息
func (this *User) doMessage(msg string) {
	this.server.BroadMessage(this, msg)
}

func NewUser(conn net.Conn, server *Server) *User {
	addr := conn.RemoteAddr().String()
	user := &User{
		addr,
		addr,
		addr,
		make(chan string), //创建一个管道，没有缓冲区
		conn,
		server,
	}
	go user.listenMessage()
	return user
}

func (this *User) listenMessage() {
	for {
		msg := <-this.c //从管道中读取数据
		this.conn.Write([]byte(msg + "\n"))
	}
}
