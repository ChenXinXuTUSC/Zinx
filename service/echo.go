package service

import (
	"zinx/core"
	zinf "zinx/interface"
	"zinx/utils/log"
)

type PingRouter struct {
	core.BaseRouter
}

func (rp *PingRouter) PreProcess(request zinf.ZinfRequest) {
	log.Dbug("Preprocess invoked")
	_, txErr := request.GetConnection().GetTCPConnection().Write([]byte("before echo..."))
	if txErr != nil {
		log.Erro(txErr.Error())
	}
}
func (rp *PingRouter) Handle(request zinf.ZinfRequest) {
	log.Dbug("Handle invoked")
	_, txErr := request.GetConnection().GetTCPConnection().Write(request.GetData())
	if txErr != nil {
		log.Erro(txErr.Error())
	}
}
func (rp *PingRouter) PostProcess(request zinf.ZinfRequest) {
	log.Dbug("Postprocess invoked")
	_, txErr := request.GetConnection().GetTCPConnection().Write([]byte("...echo after"))
	if txErr != nil {
		log.Erro(txErr.Error())
	}
}

// func EchoBack(conn *net.TCPConn, rxbuf []byte, rxLen int) error {
// 	log.Dbug("receive bytes: %s", string(rxbuf))
// 	_, txErr := conn.Write(rxbuf[:rxLen])
// 	if txErr != nil {
// 		log.Erro(txErr.Error())
// 		return txErr
// 	}

// 	return nil
// }
