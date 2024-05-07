package zinf

type ZInfServer interface {
	Start()
	Stop()
	Serve()
}
