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
	log.Info("echo handler invoked")

	// echo back
	for i := 0; i < 2; i++ {
		txErr := request.GetConnection().SendBioMsg(
			request.GetMsgId(),
			[]byte(fmt.Sprintf("[server echo %d/2] %s", i+1, strings.ToUpper(string(request.GetData())))),
		)
		if txErr != nil {
			log.Erro(txErr.Error())
			return txErr
		}
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

type VersionRouter struct {
	core.BaseRouter
}
func (vp *VersionRouter) Handle(request zinf.ZinfRequest) error {
	log.Info("version handler invoked")

	// echo back
	txErr := request.GetConnection().SendBioMsg(
		request.GetMsgId(),
		[]byte(fmt.Sprintf("[server echo] %s", "ZinxVer 0.6")),
	)
	if txErr != nil {
		log.Erro(txErr.Error())
		return txErr
	}

	return nil
}

type HookRouter struct {
	core.BaseRouter
}
func (hp *HookRouter) Handle(request zinf.ZinfRequest) error {
	log.Info("hook handler invoked")

	// echo back
	txErr := request.GetConnection().SendBioMsg(
		request.GetMsgId(),
		[]byte(fmt.Sprintf("[server hook] %s", string(request.GetData()))),
	)
	if txErr != nil {
		log.Erro(txErr.Error())
		return txErr
	}

	return nil
}
