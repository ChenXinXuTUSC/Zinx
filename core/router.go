package core

import "zinx/interface"

// abstract class stub
type BaseRouter struct{}

func (rp *BaseRouter) PreProcess(req zinf.ZinfRequest) {}
func (rp *BaseRouter) Handle(req zinf.ZinfRequest) {}
func (rp *BaseRouter) PostProcess(req zinf.ZinfRequest) {}


