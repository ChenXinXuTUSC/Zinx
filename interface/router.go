package zinf

type ZinfRouter interface {
	PreProcess(ZinfRequest)
	Handle(ZinfRequest)
	PostProcess(ZinfRequest)
}
