package main

import (
	"flag"
	"fmt"
	"net"
)

var serverIp string
var serverPort int

// 包初始化函数
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "server的ip地址") //可以给定默认值
	flag.IntVar(&serverPort, "port", 8888, "server的port")
}

type Client struct {
	ServerIp   string
	ServerPort int
	ServerName string
	ServerConn net.Conn
}

func NewClient(ip string, port int) *Client {
	client := &Client{
		ServerIp:   ip,
		ServerPort: port,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println("client connection error")
		return nil
	}
	client.ServerConn = conn
	return client
}

func main() {
	flag.Parse() //./client -h可以查看帮助

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("连接服务器失败")
		return
	}

	fmt.Println("连接服务器成功")
	select {}
}
