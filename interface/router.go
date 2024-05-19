package zinf

type ZinfRouter interface {
	PreProcess(ZinfRequest) error
	Handle(ZinfRequest) error
	PostProcess(ZinfRequest) error
}
