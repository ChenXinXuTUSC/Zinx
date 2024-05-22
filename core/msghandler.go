package core

import (
	zinf "zinx/interface"
	config "zinx/utils/conf"
	"zinx/utils/log"
)

type MsgHandler struct {
	MsgId2Router map[uint32]zinf.ZinfRouter

	NumWorker uint32
	TaskQueue []chan zinf.ZinfRequest
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		MsgId2Router: make(map[uint32]zinf.ZinfRouter),
		NumWorker:    config.GlobalConfig.NumWorker,
		TaskQueue:    make([]chan zinf.ZinfRequest, config.GlobalConfig.NumWorker),
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

func (mhp *MsgHandler) StartOneWorker(workerId uint32, taskQ chan zinf.ZinfRequest) {
	log.Info("msg handler worker routine %d launched", workerId)
	for {
		select {
		case request := <-taskQ:
			mhp.DoMsgHandler(request)
		}
	}
}

func (mhp *MsgHandler) StartWorkerPool() {
	for i := 0; i < int(mhp.NumWorker); i++ {
		mhp.TaskQueue[i] = make(chan zinf.ZinfRequest, config.GlobalConfig.NumTaskMx)
		go mhp.StartOneWorker(uint32(i), mhp.TaskQueue[i])
	}
}

func (mhp *MsgHandler) SendMsgToTaskQueue(request zinf.ZinfRequest) {
	// roll polling each worker routine
	workerId := request.GetConnection().GetConnID() % mhp.NumWorker
	log.Info(
		"assign conn %d bussiness to worker %d",
		request.GetConnection().GetConnID(),
		workerId,
	)
	mhp.TaskQueue[workerId] <- request
}
