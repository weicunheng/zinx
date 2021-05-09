package znet

import (
	"errors"
	"fmt"
	"net"
	"zinx/utils"
	"zinx/zface"
)

type Server struct {
	// 服务底层协议
	IPVersion string
	// IP地址
	IPAdr string
	// 端口
	Port string
	// 服务名称
	ServerName string

	// 路由功能
	//Router zface.IRouter
	Router zface.IMsgHandler
}

func CallBackClient(c *net.TCPConn, buf []byte, cnt int) error {
	fmt.Println("【服务器回写回调函数】")
	if _, err := c.Write(buf[:cnt]); err != nil {
		//panic("回写失败" + err.Error())
		return errors.New("回调方法回写失败...")
	}
	return nil
}

func (s *Server) Stop() {

}

func (s *Server) accept(listener *net.TCPListener) {
	// 如果有客户端连接， 通过listener可以获取 连接对象
	for {
		connection, err := listener.AcceptTCP()
		if err != nil {
			panic("Accept 失败: " + err.Error())
		}
		// 证明已经连接, 做一些业务逻辑
		var cid uint32
		cid = 0
		dealConn := NewConnection(connection, cid, s.Router)
		cid++

		// 启动链接服务程序
		go dealConn.Start()
		//go func() {
		//	for {
		//		buf := make([]byte, 512)
		//		cnt, err := connection.Read(buf)
		//		if err != nil {
		//			// 读取失败
		//			panic("读取失败， " + err.Error())
		//		}
		//		if _, err := connection.Write(buf[:cnt]); err != nil {
		//			panic("回写失败" + err.Error())
		//		}
		//	}
		//
		//}()
	}
}

func (s *Server) Start() {
	// go默认会创建socket，所以我们只需要获取连接对象即可
	tcpAddr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%s", s.IPAdr, s.Port))
	if err != nil {
		panic("服务启动失败！失败原因:" + err.Error())
	}
	fmt.Printf("Listening %s:%s ....\n", s.IPAdr, s.Port)
	// 监听服务器连接
	listener, err := net.ListenTCP(s.IPVersion, tcpAddr)
	if err != nil {
		panic("监听TCP失败: " + err.Error())
	}
	fmt.Printf("Listen success on %s:%s with tcp4\n", s.IPAdr, s.Port)

	// 会阻塞等待客户端连接
	go s.accept(listener)

}

func (s *Server) Server() {
	// 启动服务
	go s.Start()

	// ... 做一些启动服务后的一些操作
	select {}
}

func (s *Server) AddRouter(msgId uint32, router zface.IRouter) {
	// 添加路由
	//s.Router = r
	s.Router.AddRouter(msgId, router)
}

func NewServer(name string) zface.IServer {

	//utils.GlobalObject.Reload() // 加载配置文件
	settings := utils.GlobalObject
	settings.Reload()
	s := &Server{
		"tcp",
		settings.Host,
		settings.TcpPort,
		name,
		NewMsgHandler(),  // 初始化map
	}
	return s
}
