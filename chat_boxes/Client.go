package main

import (
	. "Share"
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

var user Client

// 创建表
// 如果没有该表则创建
func createTable(db *sql.DB) error {
	//创建用户表
	//id 为int类型 自动递增（适用于整数类型） 主键（唯一性，用于数据快速检索，一个表只能有一个）
	//name 二进制字符串（可变长度，必须指定长度M）类型 不允许包含NULL值
	query := `CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(50) NOT NULL,
		password VARCHAR(50) NOT NULL,
        online BOOLEAN DEFAULT false NOT NULL 
	)`
	//执行sql
	//.Exec执行查询，但不返回任何行
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	fmt.Println("Table created successfully")
	return nil
}

// 插入数据
func insertData(db *sql.DB, name string, password string) error {
	query := "INSERT INTO users(name, password,online) VALUES (?, ?,?)"

	_, err := db.Exec(query, name, password, false)
	if err != nil {
		return err
	}
	fmt.Println("Data inserted successfully")
	return nil
}

// 在指定列中指定查询
// 查询用户密码
func conditionSelectData(db *sql.DB, name string) (string, error) {
	query := "SELECT password FROM users WHERE name = ?"

	rows, err := db.Query(query, name)
	if err != nil {
		return "err", err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

			fmt.Println("Error closing rows")
		}
	}(rows)
	var password string
	if rows.Next() {

		err = rows.Scan(&password)
		if err != nil {
			return "err", err
		}
		fmt.Println(password)
		return password, nil
	} else {
		return "No user information!", nil
	}
}

// 查询用户在线状态
func online(db *sql.DB, name string) bool {
	query := "SELECT online FROM users WHERE name = ?"

	rows, err := db.Query(query, name)
	if err != nil {
		return false
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println("Error closing rows")
		}
	}(rows)
	var online bool
	if rows.Next() {
		err = rows.Scan(&online)
		if err != nil {
			return false
		}
		return online
		//fmt.Println(online)
	}
	return false
}

// 更改用户在线状态
func onlineSet(db *sql.DB, online bool, name string) {

	query := "UPDATE users SET online = ? WHERE name = ?"

	_, err := db.Exec(query, online, name)
	if err != nil {
		return
	}
	//fmt.Println("Data inserted successfully")

}

// 更新(修改)数据
// 改密码
func updatePaData(db *sql.DB, name string, newPassword string) error {
	query := "UPDATE users SET password = ? WHERE name = ?"

	_, err := db.Exec(query, newPassword, name)
	if err != nil {
		return err
	}
	fmt.Println("Data updated successfully")
	return nil
}

// 更改数据
// 改用户名
func updateNaData(db *sql.DB, oldName string, newName string) error {
	query := "UPDATE users SET name = ? WHERE name = ?"

	_, err := db.Exec(query, newName, oldName)
	if err != nil {
		return err
	}
	fmt.Println("Data updated successfully")
	return nil
}

// 判断账号密码是否有特殊符号以及是否过长
func judge(NaPa string) bool {
	n1 := strings.ContainsAny(NaPa, "/\\,+=->*%$#^&”()|")
	n2 := len(NaPa)
	if n1 {
		fmt.Println("含有特殊字符！请重新输入：")
		return false
	}
	if n2 > 15 {
		fmt.Println("超过15位！请重新输入:")
		return false
	}
	return true
}

// 注册
func register(db *sql.DB) {
	var name string
	var password string

	for {
		fmt.Println("请输入用户名：")
		//创建一个reader对象
		reader := bufio.NewReader(os.Stdin)

		//ReadBytes
		//当读取到‘\n’时，即一行完
		res, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("Read err")
		}
		name = string(res)
		if judge(name) {
			break
		}
	}

	for {
		fmt.Println("请输入密码：")
		reader := bufio.NewReader(os.Stdin)
		res, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("Read err")
		}
		password = string(res)
		if judge(password) {
			break
		}
	}

	//插入数据
	err := insertData(db, name, password)
	//错误处理
	if err != nil {
		fmt.Println("Error inserting data")
	}
	fmt.Println("注册成功！")
}

// 登录
func log(db *sql.DB) {
	var name string
	var passwordTemp string
	fmt.Println("			登录账号")
	fmt.Println("请输入用户名：")
	reader := bufio.NewReader(os.Stdin)

	res1, err := reader.ReadBytes('\n')
	if err != nil {
		fmt.Println("Read err")
	}
	name = string(res1)

mi:
	fmt.Println("请输入密码：")

	res2, err := reader.ReadBytes('\n')
	if err != nil {
		fmt.Println("Read err")
	}
	passwordTemp = string(res2)

	password, err := conditionSelectData(db, name)
	if err != nil {
		fmt.Println("Condition select error")
	}

	if password == "No user information!" {
		fmt.Println("该用户未注册！")
		log(db)
		return
	} else if passwordTemp != password {
		fmt.Println("密码输入有误！")
		fmt.Println("请重新输入：")
		goto mi
	} else if online(db, name) == true {
		fmt.Println("登陆失败！该用户可能已经在线")
		log(db)
		return
	} else {
		//在给自己的窗口打印上线消息
		fmt.Println(time.Now().Format("2006年01月02日 15:04:05"))
		fmt.Printf("欢迎回来！%s\n", name)
	}

	//更改用户在线状态
	onlineSet(db, true, name)

	//登录成功后将用户名传给客户端用于记录在线用户信息

	arr := ClientE(db, name)
	onlineSet(db, false, name)
	if !arr {
		fmt.Println("您已下线！")
	}

}

func help() {
	fmt.Println("退出:	\\q")
	fmt.Println("修改姓名：	\\rN")
	fmt.Println("修改密码：	\\rP")
	fmt.Println("发送文件：	\\file")
	fmt.Println("帮助	\\help")
	fmt.Println("切换聊天模式：输入@ XXX 或者 @all 对应私聊和公聊模式(默认为公聊)")
	fmt.Println("普通消息：直接输入")
	fmt.Println("姓名密码要求：")
	fmt.Println("		不允许 包含“/\\,+=->*%$#^&”()|")
	fmt.Println("		不允许 长度超过15 ")
}

// 修改用户名
func rN(db *sql.DB) string {
	var Newname string

	for {
		fmt.Println("请输入新的用户名：")
		reader := bufio.NewReader(os.Stdin)

		res1, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("Read err")
		}
		Newname = string(res1)
		if judge(Newname) {
			break
		}
	}
	//修改表中的用户名
	err := updateNaData(db, user.Username, Newname)
	if err != nil {
		fmt.Println("Error updating data")
	}
	//修改本地用户名
	user.Username = Newname
	return Newname
	//修改用户名成功
}

// 修改密码
func rP(db *sql.DB) {
	var NewPassword string
	for {
		fmt.Println("请新输入密码：")
		reader := bufio.NewReader(os.Stdin)
		res, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("Read err")
		}
		NewPassword = string(res)
		//NewPassword := strings.Trim(NewPasswordTemp, "\r\n")
		if judge(NewPassword) {
			break
		}
	}
	//修改表中的密码
	err := updatePaData(db, user.Username, NewPassword)
	if err != nil {
		fmt.Println("Error updating data")
	}
}

// 持续接受消息并打印
func accept(conn net.Conn) {
	for {
		//创建字节数组
		buf := [512]byte{}
		//读取回应的信息
		n, err := conn.Read(buf[:])
		if err != nil && err != io.EOF {
			return
		}
		//打印回应的信息
		if n != 0 {
			fmt.Println(string(buf[:n]))
		}
	}
}

// ClientE 客户端
func ClientE(db *sql.DB, name string) bool {
	//建立连接
	conn, err := net.Dial("tcp", "127.0.0.1:20000")
	//Conn :=Message{Username: name,conn: conn}
	//错误处理
	if err != nil {
		fmt.Println("err :", err)
		return true
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Close conn err :", err)
		}
	}(conn) // 关闭连接

	//消息发送的目标
	aim := "all"

	//将用户信息传给服务端
	_, err = conn.Write([]byte(name)) // 发送数据
	//错误处理
	if err != nil {
		return true
	}
	////在线状态写入在线用户列表
	////创建用户信息结构体
	user = Client{Channel: make(chan string), Username: name, Conn: conn}

	go accept(conn)
	for {
		//输入信息
		var information string
		reader := bufio.NewReader(os.Stdin)
		res, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("Read err")
		}

		informationTemp := string(res)
		information = strings.Trim(informationTemp, "\r\n")
		time.Sleep(1 * time.Second)
		switch information {
		case "@":
			//切换聊天对象
			fmt.Println("请输入要切换的对象：")
			res, err := reader.ReadBytes('\n')
			if err != nil {
				fmt.Println("Read err")
			}
			aimTemp := string(res)
			aim = strings.Trim(aimTemp, "\r\n")

			fmt.Println("切换成功！")
		case "\\q":
			//退出
			//通知服务端
			_, err = conn.Write([]byte("\\q"))
			//错误处理
			if err != nil {
				return true
			}
			onlineSet(db, false, name)
			return false
		case "\\rN":
			//修改姓名
			Newname := rN(db)
			//修改在线用户名
			_, err = user.Conn.Write([]byte("Newname")) // 发送数据
			//错误处理
			if err != nil {
				return true
			}
			//发送信息
			_, err = conn.Write([]byte(Newname))
			//错误处理
			if err != nil {
				return true
			}

		case "\\rP":
			//修改密码
			rP(db)
		case "\\file":
			//发送文件
			_, err = conn.Write([]byte(aim + "\\file"))
			//错误处理
			if err != nil {
				return true
			}

			// 接收用户输入，获取完整文件路径
			fmt.Print("请输入完整文件路径：")
			var filePath string
			res, err := reader.ReadBytes('\n')
			if err != nil {
				fmt.Println("Read err")
			}
			filePathTemp := string(res)
			filePath = strings.Trim(filePathTemp, "\r\n")

			fmt.Println(filePath)
			//提取文件信息
			//返回包含文件名的FileInfo
			fileInfo, err := os.Stat(filePath)
			//错误处理
			if err != nil {
				fmt.Println("os.Stat():", err)
				return true
			}

			//发送文件名
			_, err = conn.Write([]byte(fileInfo.Name()))
			// Name()返回string，如：temp.txt，转换为byte切片在网络中传输
			if err != nil {
				fmt.Println("net.Write err:", err)
				return true
			}

			//发送文件
			sendFile(filePath, conn)

		case "\\help":
			help()
		default:
			//正常发送消息
			//发送消息给服务端，并带上发送对象
			//发送聊天对象
			_, err = conn.Write([]byte(aim))
			//错误处理
			if err != nil {
				return true
			}
			//发送信息
			_, err = conn.Write([]byte(information))
			//错误处理
			if err != nil {
				return true
			}
		}
	}
}

// 发送文件内容
func sendFile(filePath string, conn net.Conn) {
	//按文件路径打开文件
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Open err ")
		return
	}
	//defer关闭文件
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Close file err :", err)
			return
		}
	}(file)

	buffer := make([]byte, 1024)
	// 循环处理
	for {
		n, err := file.Read(buffer) // 本地读取文件，按字节
		if err != nil {
			// 注意：读取完成时，会读取到EOF文件结束标志，错误标志置位为“EOF”，需要进行额外判断
			if err == io.EOF {
				//将文件读取结束的消息发送过去
				fmt.Println("发送完成:", err)
			} else {
				fmt.Println("os.Read err:", err)
			}
			return
		}
		_, err = conn.Write(buffer[:n])
		if err != nil {
			fmt.Println("conn.Write err:", err)
			return
		} // 网络写入文件，按字节
	}
}

func main() {
	// DSN (Data Source Name) 是包含数据库连接信息的字符串
	dsn := "root:413814@tcp(localhost:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local"
	//      用户名：密码@协议（主机名:端口）/数据库名称？设置字符集&让驱动解析时间类型字段为time.Time类型&设置时区为本地时区
	// 创建数据库连接
	//创建数据库对象
	//第一个参数：驱动程序名称，第二个参数告诉驱动程序如何访问基础数据存储
	//sql.Open不会建立于数据库的任何链接，也不会验证驱动程序的任何链接参数
	//与基础数据存储区的第一个实际连接将延迟到第一次需要时建立
	db, err := sql.Open("mysql", dsn)
	//错误处理
	if err != nil {
		panic(err.Error())
	}
	//defer 关闭连接
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println("Close database error")
		}
	}(db)

	// 检查连接是否有效
	//.Ping必要时建立连接
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Connected to the database!")
	}

	//创建表
	//没有该表则创建
	err = createTable(db)
	//错误处理
	if err != nil {
		fmt.Println("Error creating table")
	}

	help()
	fmt.Println("1.注册账号")
	fmt.Println("2.登陆账号")
	fmt.Println("3.退出")
	fmt.Println("请选择要进行的操作")
	var a int
	_, err = fmt.Scanln(&a)
	if err != nil {
		fmt.Println("Scan error")
		return
	}

	switch a {
	case 1:
		//注册
		register(db)
		//注册成功后跳转到登录
		fallthrough
	case 2:
		//登录
		log(db)
	case 3:
		return
	}
}
