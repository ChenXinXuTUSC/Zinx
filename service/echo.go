package service

import (
	"fmt"
	"strings"
	"zinx/core"
	zinf "zinx/interface"
	"zinx/utils/log"
)

type PingRouter struct {
	core.BaseRouter
}

//	func (rp *PingRouter) PreProcess(request zinf.ZinfRequest) {
//		log.Dbug("Preprocess invoked")
//		_, txErr := request.GetConnection().GetTCPConnection().Write([]byte("before echo..."))
//		if txErr != nil {
//			log.Erro(txErr.Error())
//		}
//	}
//
//	func (rp *PingRouter) Handle(request zinf.ZinfRequest) error {
//		log.Dbug("Handle invoked")
//		_, txErr := request.GetConnection().GetTCPConnection().Write(request.GetData())
//		if txErr != nil {
//			log.Erro(txErr.Error())
//			return txErr
//		}
//		return nil
//	}
func (rp *PingRouter) Handle(request zinf.ZinfRequest) error {
	log.Info("PingRouter invoked")

	// echo back
	txErr := request.GetConnection().SendMsg(
		request.GetMsgId(),
		[]byte(fmt.Sprintf("[server echo] %s", strings.ToUpper(string(request.GetData())))),
	)
	if txErr != nil {
		log.Erro(txErr.Error())
		return txErr
	}

	return nil
}

// func (rp *PingRouter) PostProcess(request zinf.ZinfRequest) error {
// 	log.Dbug("Postprocess invoked")
// 	_, txErr := request.GetConnection().GetTCPConnection().Write([]byte("...echo after"))
// 	if txErr != nil {
// 		log.Erro(txErr.Error())
// 		return txErr
// 	}
// 	return nil
// }

// func EchoBack(conn *net.TCPConn, rxbuf []byte, rxLen int) error {
// 	log.Dbug("receive bytes: %s", string(rxbuf))
// 	_, txErr := conn.Write(rxbuf[:rxLen])
// 	if txErr != nil {
// 		log.Erro(txErr.Error())
// 		return txErr
// 	}

// 	return nil
// }

type ZinxVerRouter struct {
	core.BaseRouter
}
func (zp *ZinxVerRouter) Handle(request zinf.ZinfRequest) error {
	log.Info("ZinxVerRouter invoked")

	// echo back
	txErr := request.GetConnection().SendMsg(
		request.GetMsgId(),
		[]byte(fmt.Sprintf("[server echo] %s", "ZinxVer 0.6")),
	)
	if txErr != nil {
		log.Erro(txErr.Error())
		return txErr
	}

	return nil
}
