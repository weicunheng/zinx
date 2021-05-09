package znet

import "zinx/zface"

type MsgHandler struct {
	// 存放message和router的映射
	MsgRouterPattern map[uint32] zface.IRouter
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		MsgRouterPattern: make(map[uint32] zface.IRouter),
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
	if _, ok := mh.MsgRouterPattern[msgId]; ok{
		panic("Message 和 Router 映射关系已存在！")
	}
	// 2. 添加映射关系
	mh.MsgRouterPattern[msgId] = router
}
