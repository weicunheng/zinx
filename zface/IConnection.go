package zface

import "net"

/*
本模块的作用是，从Server中把  Write 和 Read 抽取出来
并且结合具体的业务，调用不同的callback回调方法处理  当前链接数据
connection, err := listener.AcceptTCP()
		if err != nil {
			panic("Accept 失败: " + err.Error())
		}
		// 证明已经连接, 做一些业务逻辑
		go func() {
			for {
				buf := make([]byte, 512)
				cnt, err := connection.Read(buf)
				if err != nil {
					// 读取失败
					panic("读取失败， " + err.Error())
				}
				if _, err := connection.Write(buf[:cnt]); err != nil {
					panic("回写失败" + err.Error())
				}
			}

		}()

*/
type IConnection interface {
	Start() // 启动链接的方法
	Stop() // 关闭链接的方法
	GetTCPConnection() *net.TCPConn// 从当前链接获取原始的Socket TCP Connection
	GetConnID() uint // 获取当前链接id
	GetRemoteAddr() net.Addr // 获取远程链接地址信息  IP 、Port
}

/*
HandFunc  回调方法， 所有conn链接在处理业务的函数接口
*net.TCPAddr: 当前链接
[]byte 数据
int: 数据的长度

如果我们想要指定一个conn的处理业务，只要定义一个HandFunc类型的函数即可，然后和该链接进行绑定
*/
type HandFunc func(c *net.TCPConn, buf []byte, cnt int) error