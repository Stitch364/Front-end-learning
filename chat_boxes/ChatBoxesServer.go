package main

import (
	. "Share"
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

// 创建在线成员map
var OnlineMap = make(map[string]Client)

// 创建全局通道
var Message = make(chan string) //用于全局通信

// 客户端监听自己的通道
func listenMessage(cl Client) {
	for {
		//获取用户通道的消息
		msg := <-cl.Channel
		//发送给客户端
		_, err := cl.Conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("Write Error:", err)
			return
		} // 发送数据
	}
}

// 广播消息
// 将全局通道的消息写入msg，再循环发送给所有在线用户
func broadcast() {
	//onlineMap = make(map[string]Client)
	for {
		//循环监听全局通道，获取全局message
		msg := <-Message
		//循环将全局消息广播至所有在线用户
		for _, cl := range OnlineMap {
			//遍历所有在线用户
			//将msg写入用户的channel
			cl.Channel <- msg
		}
	}
}

// 接收用户名（接收信息）
func info(conn net.Conn) string {

	var buf [1024]byte
	reader := bufio.NewReader(conn)
	n, err := reader.Read(buf[:]) // 读取用户信息

	//错误处理
	if err != nil && err != io.EOF {
		fmt.Println("read from client failed, err:", err)
		return ""
	}
	//n为读取到的字节数
	//转换为string类型
	if n != 0 {
		nameTemp := string(buf[:n])
		name := strings.Trim(nameTemp, "\r\n")

		return name
	}
	return ""
}

// 接受文件
func recvFile(fileName string, conn net.Conn, serveConn net.Conn) {
	file, err := os.Create(fileName) // 在当前路径下创建文件
	// 返回File类型
	if err != nil {
		fmt.Println("os.Create err:", err)
		return
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("file close err:", err)
			return
		}
	}(file) // 延迟关闭文件

	buffer := make([]byte, 1024)
	// 循环处理
	for {
		readCount1, err := serveConn.Read(buffer) // 网络读取文件，按字节

		if err != nil {
			// 读取完成时，会读取到EOF文件结束标志，错误标志置位为“EOF”，需要进行额外判断
			if err == io.EOF {
				fmt.Println("接收完成:", err)
			} else {
				fmt.Println("net.Read():", err)
			}
			return
		}
		//if readCount1 == 0 {
		//	return
		//}

		//fmt.Println(string(buffer[:readCount]))
		_, err = file.Write(buffer[:readCount1])
		if err != nil {
			fmt.Println("file.Write err:", err)
			return
		} // 本地写入文件，按字节
		if readCount1 < 1024 {
			return
		}
	}
}

// 循环接收并处理消息
func Info(user Client) bool {
	for {
		aim := info(user.Conn)
		if aim == "\\q" {
			//退出
			return true
		} else if strings.Contains(aim, "\\file") {
			//发送文件
			aim = strings.Trim(aim, "\\file")
			fmt.Println(aim)
			//接收文件名
			buffer := make([]byte, 1024)
			n, err := user.Conn.Read(buffer)
			if err != nil {
				fmt.Println("net.Read():", err)
				return true
			}

			fileName := string(buffer[:n]) // 获取文件名
			//转换为字符串

			//发送，接收到文件名后的响应信息
			_, err = user.Conn.Write([]byte("ok"))
			if err != nil {
				fmt.Println("Write err:", err)
				return true
			}

			//recvFile(fileName, OnlineMap[aim].Conn, user.Conn)
			//Message <- "收到一份文件\n"

			//接收文件内容
			if aim == "all" {
				for name, usern := range OnlineMap {
					if name != user.Username {
						recvFile(fileName, usern.Conn, user.Conn)
						user.Channel <- "收到一份文件\n"
					}
				}
			} else {
				recvFile(fileName, OnlineMap[aim].Conn, user.Conn)
				OnlineMap[aim].Channel <- "收到一份文件\n"
			}
			continue
		} else if aim == "Newname" {
			Newname := info(user.Conn)
			delete(OnlineMap, user.Username)
			user.Username = Newname
			OnlineMap[Newname] = user
		}
		information := info(user.Conn)
		if aim == "all" {
			//广播
			//将信息写入全局通道
			Message <- Make1Msg(user, information)
		} else {
			a, ok := OnlineMap[aim]
			if ok {
				a.Channel <- Make2Msg(user, information)
			} else {
				user.Channel <- "该用户未在线！"
			}
		}
	}
}

// 处理函数
func process(conn net.Conn) {
	//defer关闭连接

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error closing connection")
		}
	}(conn) // 关闭连接

	name := info(conn)
	fmt.Println("登录成功:", name)
	//创建用户信息结构体
	user := Client{Channel: make(chan string), Username: name, Conn: conn}

	//写入在线用户map
	OnlineMap[user.Username] = user

	//持续监听用户自己的通道
	go listenMessage(user)
	//广播用户上线消息给所有在线用户
	//广播上线消息
	Message <- Make0Msg(user)
	//广播当前在线人数
	Message <- Make3Msg(OnlineMap)
	//开始聊天
	q := Info(user)

	//判断聊天是否结束
	if q {
		//下线
		//通知所有用户
		Message <- Make00Msg(user)
		//删除在线用户记录
		delete(OnlineMap, name)
		Message <- Make3Msg(OnlineMap)

		return
	}
}

// 服务端
func main() {
	//1.监听端口
	//127.0.0.1  本地环回地址  本地IP
	//0.0.0.0  广播地址（任意地址）（默认地址）
	listen, err := net.Listen("tcp", "127.0.0.1:20000")
	//错误处理
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}

	//打开广播
	go broadcast()
	fmt.Println("listening on", listen.Addr())
	for {
		//接受请求 阻塞等待
		conn, err := listen.Accept() // 建立连接
		//错误处理
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		//处理请求
		go process(conn) // 启动一个goroutine处理连接
	}
}
