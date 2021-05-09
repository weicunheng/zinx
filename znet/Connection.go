package znet

/*
链接模块，接受客户端的链接，回写客户端数据
*/

import (
	"fmt"
	"io"
	"net"
	"zinx/zface"
)

// 链接对象
type Connection struct {
	// 当前链接的套接字
	Conn *net.TCPConn
	// 当前链接的id， 也可以称为Session ID， 全局唯一
	ConnId uint32
	// 当前链接的状态
	isClosed bool
	// 该链接处理方法api
	//handleAPI zface.HandFunc
	//Router zface.IRouter
	Router zface.IMsgHandler

	// 告知该链接已经退出、停止的channel。 由Reader告知Writer退出
	ExitBuffChan chan bool
	// Reader与Writer通信的chan
	MsgChan chan []byte
}

// 开启写协程
func (conn *Connection) StartWriter() {
	fmt.Println("[**写协程开始执行.....]")
	defer fmt.Printf("id:%s 写协程已经退出 , 客户端地址为:%s\n", conn.ConnId, conn.GetRemoteAddr().String())
	for {
		select {
		case data := <-conn.MsgChan:
			if _, err := conn.Conn.Write(data); err != nil {
				fmt.Println("客户端:", conn.GetRemoteAddr(), "回写数据失败")
				return
			}
		case <-conn.ExitBuffChan: // 表示客户端已经断开，关闭Writer
			return
		}
	}

}

// 开启读协程
func (conn *Connection) StartReader() {
	fmt.Println("[*读协程 开始运行.....]")
	defer fmt.Printf("id:%s 读协程已经退出 , 客户端地址为:%s\n", conn.ConnId, conn.GetRemoteAddr().String())
	defer conn.Stop()

	for {
		// 使用Message代替
		// 1. 创建拆包对象
		dp := NewDataPack()

		// 2. 读取客户端msg head  8个字节二进制数据
		headByteData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn.GetTCPConnection(), headByteData); err != nil {

			panic("读取链接中Message头信息 失败！失败原因:" + err.Error())
		}

		// 3. 拆包得到 msgId和 msgDataLen 放入msg中
		msg, err := dp.UnPack(headByteData)
		if err != nil {
			panic("接卸链接中Message头信息 失败！失败原因:" + err.Error())
		}
		// 4. 根据dataLen再次读取Data， 放入msg.Data中

		var messageBody []byte
		if msg.GetMsgLen() > 0 {
			messageBody = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn.GetTCPConnection(), messageBody); err != nil {
				panic("读取链接中Message Body信息 失败！失败原因:" + err.Error())
			}

		}

		msg.SetMsgData(messageBody)
		// 5. 得到当前conn数据Request请求数据

		// 当Connection读取完数据之后，我们将Connection链接对象和数据 封装成原子包Request
		// Request作为Router的输入
		r := Request{
			conn,
			msg,
		}
		// 处理router和handler
		//go conn.Router.DoMsgHandler(&r) worker协程和请求数量一致，比较耗费资源，优化为req请求任务交给worker pool(工作池)
		
		go conn.Router.SendReqToTaskQueue(&r)
	}
}

// 封包方法
func (conn *Connection) SendMsg(msgId uint32, data []byte) error {

	if conn.isClosed == true {
		panic("链接已经关闭，不能向客户端发送数据！")
	}

	dp := NewDataPack()
	message := NewMessage(msgId, data)

	sendBody, err := dp.Pack(message)
	if err != nil {

		//panic("服务器向客户端发送数据，封包失败！ 失败原因：" + err.Error())
		return err
	}

	//if _, err := conn.Conn.Write(sendBody); err != nil {
	//	//panic("服务器向客户端发送数据，发送失败！ 失败原因：" + err.Error())
	//	return err
	//}
	conn.MsgChan <- sendBody

	return nil
}

func (conn *Connection) Start() {
	fmt.Printf("Conn Start()... id:%d\n", conn.ConnId)
	go conn.StartReader()
	go conn.StartWriter()
}

func (conn *Connection) Stop() {
	if conn.isClosed == true {
		return
	}
	conn.isClosed = true

	conn.ExitBuffChan <- true

	// 关闭管道
	conn.Conn.Close()
	close(conn.ExitBuffChan)
	close(conn.MsgChan)
}

func (conn *Connection) GetTCPConnection() *net.TCPConn {
	return nil
}
func (conn *Connection) GetConnID() uint {
	return 0
}
func (conn *Connection) GetRemoteAddr() net.Addr {
	return conn.Conn.RemoteAddr()
}
func NewConnection(conn *net.TCPConn, connId uint32, router zface.IMsgHandler) *Connection {
	return &Connection{
		conn,
		connId,
		false,
		router,
		make(chan bool, 1),
		make(chan []byte), // 无缓冲chan
	}
}
