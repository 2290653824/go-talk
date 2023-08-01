package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
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
	flag       int
}

func NewClient(ip string, port int) *Client {
	client := &Client{
		ServerIp:   ip,
		ServerPort: port,
		flag:       999,
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
	go client.dealResponse()
	client.run()
}

func (this *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更改用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		this.flag = flag
		return true
	} else {
		fmt.Println("client input flag invalid")
		return false
	}
}

func (this *Client) privateChat() {
	this.menuPerson()
	var targetPerson string
	var targetMessage string
	fmt.Println(">>>>>>请选择你想要发送的对象的id,exit退出")
	fmt.Scanln(&targetPerson)

	for targetPerson != "exit" {
		fmt.Println(">>>>>>>请输入你要发送的消息，exit退出")
		fmt.Scanln(&targetMessage)
		for targetMessage != "exit" {
			meg := "to|" + targetPerson + "|" + targetMessage + "\n"
			_, err := this.ServerConn.Write([]byte(meg))
			if err != nil {
				fmt.Println("private send message error=", err)
			}

			targetMessage = ""
			fmt.Println(">>>>>>>请输入你要发送的消息，exit退出")
			fmt.Scanln(&targetMessage)
		}

		targetPerson = ""
		this.menuPerson()
		fmt.Println(">>>>>>请选择你想要发送的对象的id,exit退出")
		fmt.Scanln(&targetPerson)

	}
}

func (this *Client) menuPerson() {
	msg := "who\n"
	_, err := this.ServerConn.Write([]byte(msg))
	if err != nil {
		fmt.Println("menuPerson error=", err)
	}
}

func (this *Client) run() {
	for this.flag != 0 {
		for this.menu() != true {

		}
		switch this.flag {
		case 1:
			this.publicChat()
		case 2:
			this.privateChat()
		case 3:
			this.updateName()
		}
	}
}

func (this *Client) publicChat() {
	var publicMessage string
	fmt.Println(">>>>>>>>>>请输入聊天内容，exit直接退出")
	fmt.Scanln(&publicMessage)

	for publicMessage != "exit" {

		if len(publicMessage) != 0 {
			_, err := this.ServerConn.Write([]byte(publicMessage))
			if err != nil {
				fmt.Println("public message send error=", err)
				break
			}
		}

		publicMessage = ""
		fmt.Println(">>>>>>>>>>请输入聊天内容，exit直接退出")
		fmt.Scanln(&publicMessage)
	}
}

func (this *Client) updateName() bool {
	fmt.Println(">>>>>>>请输入你的名字")
	fmt.Scanln(&this.ServerName)

	msg := "rename|" + this.ServerName + "\n"
	_, err := this.ServerConn.Write([]byte(msg))
	if err != nil {
		fmt.Println("发送更新名字数据失败，error=", err)
		return false
	}
	return true
}

func (this *Client) dealResponse() {
	io.Copy(os.Stdout, this.ServerConn) //永久阻塞
}
