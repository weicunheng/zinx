package znet

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

	Router zface.IRouter

	// 告知该链接已经退出、停止的channel
	ExitBuffChan chan bool
}

func (c *Connection) StartReader() {

	defer fmt.Printf("id:%s read is exit , remode addr is %s\n", c.ConnId, c.GetRemoteAddr().String())
	defer c.Stop()

	for {
		/*
			buf := make([]byte, 512)
				_, err := c.Conn.Read(buf)
				if err != nil {
					// 读取失败
					panic("读取失败， " + err.Error())
				}

				// 并不是直接调用c.Conn.Write(), 而是调用模板回调方法
				//if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
				//	panic("回写失败" + err.Error())
				//}

				// 当Connection读取完数据之后，我们将Connection链接对象和数据 封装成原子包Request
				// Request作为Router的输入
				r := Request{
					c,
					buf,
				}
				// 根据路由调用方法
				go func(req zface.IRequest) {
					c.Router.PreHandler(req)
					c.Router.Handler(req)
					c.Router.PostHandler(req)
				}(&r)
		*/
		// 使用Message代替
		// 1. 创建拆包对象
		dp := NewDataPack()

		// 2. 读取客户端msg head  8个字节二进制数据
		headByteData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headByteData); err != nil {
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
			if _, err := io.ReadFull(c.GetTCPConnection(), messageBody); err != nil{
				panic("读取链接中Message Body信息 失败！失败原因:" + err.Error())
			}

		}

		msg.SetMsgData(messageBody)
		// 5. 得到当前conn数据Request请求数据


		// 当Connection读取完数据之后，我们将Connection链接对象和数据 封装成原子包Request
		// Request作为Router的输入
		r := Request{
			c,
			msg,
		}
		// 根据路由调用方法
		go func(req zface.IRequest) {
			c.Router.PreHandler(req)
			c.Router.Handler(req)
			c.Router.PostHandler(req)
		}(&r)
		
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
func NewConnection(conn *net.TCPConn, connId uint32, router zface.IRouter) *Connection {
	return &Connection{
		conn,
		connId,
		false,
		router,
		make(chan bool, 1),
	}
}
