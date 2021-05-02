package znet

import (
	"fmt"
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
	handleAPI zface.HandFunc
	// 告知该链接已经退出、停止的channel
	ExitBuffChan chan bool
}

func (c *Connection) StartReader() {

	defer fmt.Printf("id:%s read is exit , remode addr is %s\n", c.ConnId, c.GetRemoteAddr().String())
	defer c.Stop()

	for {
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			// 读取失败
			panic("读取失败， " + err.Error())
		}

		// 并不是直接调用c.Conn.Write(), 而是调用模板回调方法
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			panic("回写失败" + err.Error())
		}
	}
}

func (conn *Connection) Start() {
	fmt.Printf("Conn Start()... id:%d\n", conn.ConnId)
	go conn.StartReader()
}
func (conn *Connection) Stop() {
	if conn.isClosed == true {
		return
	}
	conn.isClosed = true

	// 关闭管道
	conn.Conn.Close()
	close(conn.ExitBuffChan)
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
func NewConnection(conn *net.TCPConn, connId uint32, callbackFunc zface.HandFunc) *Connection {
	return &Connection{
		conn,
		connId,
		false,
		callbackFunc,
		make(chan bool, 1),
	}
}
