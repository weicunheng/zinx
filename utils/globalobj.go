package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/zface"
)

/*
存储 有关zinx框架相关的全局参数
此模块供其他模块使用

如果用户提供了zinx。json，也会使用用户的配置
*/

type GlobalObj struct {
	TCPServer zface.IServer // zinx全局Server对象
	Host      string
	TcpPort   string
	Name      string
	Version   string

	MaxPackageSize uint32
	MaxConn        int
}

// 定义一个全局变量， 目的就是使其他模块都能访问里面的参数
var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("./config/zinx.json")
	if err != nil{
		panic("config/zinx.json读取失败！失败原因:" + err.Error())
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil{
		panic("json反序列化失败！ 失败原因:" + err.Error())
	}
}

// 提供一个 init 方法， 目的是初始化GlobalObject对象，和加载服务器应用配置文件conf/zinx.json
func init() {
	// 初始化GlobalObject
	GlobalObject := &GlobalObj{
		Name:           "ZinxServerApp",
		Version:        "V0.1",
		TcpPort:        "8081",
		Host:           "127.0.0.1",
		MaxConn:        12000,
		MaxPackageSize: 4096,
	}
	GlobalObject.Reload()
}
