package main

import (
	"fmt"
	"net"
	"strings"
)

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
	if msg == "who" { //指令who
		for _, user := range this.server.onlineMap {
			msg := "[user id = " + user.id + ", user addr = " + user.addr + ",username = " + user.userName + "] exist\n"
			this.sendSingleUsr(msg)
		}
	} else if len(msg) > 7 && msg[:7] == "rename|" { //指令 rename|newName
		newName := msg[7:]
		this.userName = newName
		this.sendSingleUsr("your username has reset:" + newName)
	} else if len(msg) > 3 && msg[:3] == "to|" {
		arr := strings.Split(msg, "|")
		id := arr[1]
		talkMessage := arr[2]
		if id == "" {
			fmt.Println("the id is not validate")
		}
		user := this.server.onlineMap[id]
		if user == nil {
			fmt.Println("the id is not existed in onlineMap")
		}
		message := "userid[" + this.id + "],username=[" + this.userName + "] send message to you:" + talkMessage
		user.sendSingleUsr(message)

	} else {
		this.server.BroadMessage(this, msg)
	}
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

func (this *User) sendSingleUsr(msg string) {
	this.conn.Write([]byte(msg))
}

func (this *User) listenMessage() {
	for {

		msg := <-this.c //从管道中读取数据
		this.conn.Write([]byte(msg + "\n"))
	}
}
