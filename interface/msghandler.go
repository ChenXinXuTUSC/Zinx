package zinf

type ZinfMsgHandler interface {
	DoMsgHandler(request ZinfRequest)
	AddRouter(msgId uint32, router ZinfRouter)

	// launch worker poool
	StartWorkerPool()
	SendMsgToTaskQueue(request ZinfRequest)
}
