package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	ip   string
	port int

	//在线user Map
	onlineMap map[string]*User
	mapLock   sync.RWMutex

	//广播管道
	message chan string
}

// 构造函数
func NewServer(ip string, port int) *Server {
	server := &Server{
		ip:        ip,
		port:      port,
		onlineMap: make(map[string]*User),
		message:   make(chan string),
	}
	return server
}

func (this *Server) doHandler(connection net.Conn) {
	fmt.Println("start to exec handler,this conn:", connection)
	user := NewUser(connection)
	this.mapLock.Lock()
	this.onlineMap[user.id] = user
	this.mapLock.Unlock()

	this.BroadMessage(user, "已经上线") //广播消息

	//阻塞
	select {}
}

// 将广播信息发送到channel当中，另一个协成会去channel中拿数据
func (this *Server) BroadMessage(user *User, msg string) {
	res := "user id = " + user.id + ",userName = " + user.userName + ":" + msg
	this.message <- res
}

func (this *Server) listenMessage() {
	for {
		msg := <-this.message
		this.mapLock.Lock()
		for _, value := range this.onlineMap {
			value.c <- msg
		}
		this.mapLock.Unlock()
	}
}

// 服务器启动
func (this *Server) Start() {
	//listen
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.ip, this.port)) //监听
	if err != nil {
		fmt.Println("net.listen error = ", err)
	}
	defer listen.Close()
	go this.listenMessage()
	for {
		//accept
		connection, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.accept error = ", err)
			continue
		}

		//do handler
		this.doHandler(connection)
	}

	//close
}
