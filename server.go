package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
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
	user := NewUser(connection, this)
	user.online()
	isLive := make(chan bool)
	go func() {
		buf := make([]byte, 4069)

		for {
			n, err := connection.Read(buf)
			if n == 0 { //当读到为0时，表示客户端已经关闭了
				user.offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err :", err)
				return
			}

			msg := string(buf[:n-1])

			user.doMessage(msg)

			isLive <- true //用户活跃

		}
	}()
	//阻塞
	for {
		select {
		case <-isLive:

		case <-time.After(time.Second * 10):
			user.sendSingleUsr("you have terminated because of expiration")
			close(user.c)
			user.conn.Close()
			return

		}
	}

}

// 将广播信息发送到channel当中，另一个协成会去channel中拿数据
func (this *Server) BroadMessage(user *User, msg string) {
	res := "[" + "user id = " + user.id + ",userName = " + user.userName + "]" + ":" + msg
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
		go this.doHandler(connection)
	}

	//close
}
