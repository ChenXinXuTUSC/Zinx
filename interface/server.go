package zinf

type ZinfServer interface {
	Start()
	Stop()
	Serve()

	AddRouter(msgId uint32, router ZinfRouter)
}
