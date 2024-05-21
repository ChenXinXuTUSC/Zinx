package core

import (
	zinf "zinx/interface"
	"zinx/utils/log"
)

type MsgHandler struct {
	MsgId2Router map[uint32]zinf.ZinfRouter
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		MsgId2Router: make(map[uint32]zinf.ZinfRouter),
	}
}

func (mhp *MsgHandler) DoMsgHandler(request zinf.ZinfRequest) {
	handler, hok := mhp.MsgId2Router[request.GetMsgId()]
	if !hok {
		log.Erro("no valid handler for msg type %d", request.GetMsgId())
		return
	}

	handler.PreProcess(request)
	handler.Handle(request)
	handler.PostProcess(request)
}

func (mhp *MsgHandler) AddRouter(msgId uint32, router zinf.ZinfRouter) {
	if _, ok := mhp.MsgId2Router[msgId]; ok {
		log.Warn("handler for msg type %d already exist", msgId)
		return
	}
	mhp.MsgId2Router[msgId] = router
	log.Info("router for msg type %d added successfully", msgId)
}
