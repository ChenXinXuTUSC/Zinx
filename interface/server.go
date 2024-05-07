package zinf

type ZinfServer interface {
	Start()
	Stop()
	Serve()
}
