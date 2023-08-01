package main

import "net"

type User struct {
	id       string
	userName string
	addr     string
	c        chan string //管道中传输字符串
	conn     net.Conn
}

func NewUser(userName string, conn net.Conn) *User {
	addr := conn.RemoteAddr().String()
	user := &User{
		addr,
		userName,
		addr,
		make(chan string), //创建一个管道，没有缓冲区
		conn,
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
