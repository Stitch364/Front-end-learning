package Share

import (
	"net"
	"strconv"
)

type Client struct {
	channel  chan string //用户 通道用于通信
	Username string
	//Addr     string
	conn net.Conn
}

// 创建在线成员map
var onlineMap = make(map[string]Client)

// 创建全局通道
var message = make(chan string) //用于全局通信

// 创建传输用户信息的通道
var userChan = make(chan Client)

// Make0Msg 改消息格式0
// 上线通知
func Make0Msg(cl Client) string {
	a := cl.Username + "已上线！"
	return a
}

// Make1Msg 改消息格式1
// 公聊
func Make1Msg(cl Client, msg string) string {
	a := cl.Username + ":" + msg
	return a
}

// Make2Msg 改消息格式2
// 私聊
func Make2Msg(cl Client, msg string) string {
	a := "[" + "悄悄话" + "]" + cl.Username + ":" + msg
	return a
}

// Make3Msg 改消息格式3
// 在线人数
func Make3Msg() string {
	a := "当前在线人数：" + strconv.Itoa(len(onlineMap))
	return a
}
