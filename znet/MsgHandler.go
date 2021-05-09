package znet

import (
	"fmt"
	"zinx/zface"
)

type MsgHandler struct {
	// 存放message和router的映射
	MsgRouterPattern map[uint32]zface.IRouter

	/*
		一共WorkerPoolSize个worker
		一个worker对应一个消息队列([]chan zface.IRequest);
		一个消息队列最大能存放1024个任务

	*/
	// 负责worker取任务的消息队列  ===> 因为任务处理主要就是IRequest对象； 和 work数量一致
	TaskQueue []chan zface.IRequest
	// 业务工作worker池的work数量
	WorkerPoolSize uint32 //  工作池中 worker的数量

}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		MsgRouterPattern: make(map[uint32]zface.IRouter),
		TaskQueue:        make([]chan zface.IRequest, 10), // 10个TaskQueue， 和 woker数量一致
		WorkerPoolSize:   uint32(10),                      // 默认workerPool可容纳10个worker
	}
}

// 1. 从map中获取router, 并执行
func (mh *MsgHandler) DoMsgHandler(request zface.IRequest) {
	// 通过msgId获取router
	handler, ok := mh.MsgRouterPattern[request.GetMsgId()]
	if !ok {
		panic("")
	}
	handler.PreHandler(request)
	handler.Handler(request)
	handler.PostHandler(request)
}

// 2. 添加router
func (mh *MsgHandler) AddRouter(msgId uint32, router zface.IRouter) {
	// 1. 判断映射是否存在
	if _, ok := mh.MsgRouterPattern[msgId]; ok {
		panic("Message 和 Router 映射关系已存在！")
	}
	// 2. 添加映射关系
	mh.MsgRouterPattern[msgId] = router
}

// 3. 启动一个worker工作池, 对外暴漏; 一个框架只能有一个worker工作池
func (mh *MsgHandler) StartWorkerPool() {
	/*
		1. 根据workPoolSize，创建worker， 比如5个worker池，开辟5个goroutine
		2. 每个work(也就是goroutine)，创建对应的消息队列(TaskQueue)，并且指定消息队列最大可容纳任务数量(chan缓冲大小值)
		3. 启动worker
	*/
	// 根据workerPoolSize分别开启Worker， 每个Worker用一个go来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ { // 根据工作池的大小，创建worker，并且为每一个worker创建消息队列(最大可缓冲1024个任务)
		/*
			[[1024]chan, [1024]chan, .... [1024]chan]
		*/
		mh.TaskQueue[i] = make(chan zface.IRequest, 1024)
		fmt.Println(mh.TaskQueue)
		mh.StartOneWorker(i, mh.TaskQueue[i])
	}

}

// 4. 由工作池来启动一个worker工作流程
func (mh *MsgHandler) StartOneWorker(id int, taskQueue chan zface.IRequest) {
	/*

	 */
	fmt.Println("当前workerID：", id, "已经启动了.....")
	for true {
		select {
		// 如果有消息，出列的就是一个客户端的request
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

func (mh *MsgHandler) SendReqToTaskQueue(request zface.IRequest) {
	/*
		1. 将消息平均分配给不同的TaskQueue
			根据request的id，或者客户端的id取模运算
		2. 将消息发送给对应的worker的TaskQueue即可
	*/
	workId := uint32(request.GetConnection().GetConnID()) % mh.WorkerPoolSize
	fmt.Println("链接id:", request.GetConnection().GetConnID(),
		"request id:", request.GetMsgId(),
		"work id:", workId)
	mh.TaskQueue[workId] <- request

}
