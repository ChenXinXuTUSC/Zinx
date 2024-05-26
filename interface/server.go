package zinf

type ZinfServer interface {
	Start()
	Stop()
	Serve()

	AddRouter(msgId uint32, router ZinfRouter)
	GetConnMgr() ZinfConnManager

	// set hook for connection
	SetHookOnConnStart(func (ZinfConnection))
	SetHookOnConnStop(func (ZinfConnection))
	CallHookOnConnStart(ZinfConnection)
	CallHookOnConnStop(ZinfConnection)
}
