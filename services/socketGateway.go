package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type jspkg struct {
	Cmd  string      `json:"cmd"`
	Data interface{} `json:"data"`
}

var (
	Socketconnections map[string]net.Conn //全局socket连接
)

func InitSocket() { //异步进程启动socket服务器
	go initSocketServer()
}

func initSocketServer() {
	Socketconnections = make(map[string]net.Conn)

	netListen, err := net.Listen("tcp", ":6666")
	checkErr(err)
	defer netListen.Close()

	log.Println("Socket server initial. Listen on :6666\n")
	hearbeatDuration := 120 //心跳间隔

	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}

		log.Println(conn.RemoteAddr().String(), " connected")
		conn.SetReadDeadline(time.Now().Add(time.Duration(hearbeatDuration) * time.Second))
		go handleConnection(conn, hearbeatDuration)
	}
}

func handleConnection(conn net.Conn, timout int) {
	buffer := make([]byte, 2048)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println(conn.RemoteAddr().String(), " connection err:", err)
			conn.Close()
			return
		}

		data := (buffer[:n])      //获取从socket客户端得到的数据
		hbChan := make(chan byte) //HeartBeat Channel
		dataErrCount := 0         //解析包错误次数

		pkg := jspkg{}
		err = json.Unmarshal(data, &pkg) //将json数据decode
		if err != nil {                  //非json正确格式
			dataErrCount++
			if dataErrCount >= 3 { //解析包错误次数超过2次，将连接视为非法链接，将连接关闭
				conn.Close()
			}
		}

		fmt.Println(pkg.Cmd)
		switch {
		case pkg.Cmd == "login": //用户登录，并且将id与socket连接绑定
			var ctl SocketFunctions
			ctl.Login(conn, pkg.Data)

		case pkg.Cmd == "sendMessage": //发送数据到指定ID
			var ctl SocketFunctions
			ctl.SendMessageToClient(pkg.Data)

		default:
			fmt.Println("default")
		}

		go heartBeat(conn, hbChan, timout)
		go gravelChannel(data, hbChan)
	}
}

func heartBeat(conn net.Conn, hbChan chan byte, timeout int) {
	select {
	case _ = <-hbChan:
		//			fmt.Println(conn.RemoteAddr().String(), "receive data string:", string(fk))
		conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))

	case <-time.After(5 * time.Second):
		fmt.Println("client close")
		conn.Close()
	}
}

func gravelChannel(n []byte, hbChan chan byte) {
	for _, v := range n {
		hbChan <- v
	}
	close(hbChan)
}

func checkErr(err error) {
	if err != nil {
		log.Println(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
