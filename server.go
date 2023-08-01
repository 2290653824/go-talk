package main

import (
	"fmt"
	"net"
)

type Server struct {
	ip   string
	port int
}

// 构造函数
func NewServer(ip string, port int) *Server {
	server := &Server{
		ip:   ip,
		port: port,
	}
	return server
}

func (this *Server) doHandler(connection net.Conn) {
	fmt.Println("start to exec handler,this conn:", connection)
}

// 服务器启动
func (this *Server) Start() {
	//listen
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.ip, this.port)) //监听
	if err != nil {
		fmt.Println("net.listen error = ", err)
	}
	defer listen.Close()
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
