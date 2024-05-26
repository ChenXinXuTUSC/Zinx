package service

import (
	"runtime"
	zinf "zinx/interface"
	"zinx/utils/log"
)

func DoConnStart(conn zinf.ZinfConnection) {
	log.Dbug("hook on conn start")
	sendErr := conn.SendBioMsg(2, []byte("invoke start hook"))
	if sendErr != nil {
		log.Erro("conn start hook error: %s", sendErr.Error())
	}
}

func DoConnStop(conn zinf.ZinfConnection) {
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	log.Dbug(string(buf[:n]))
	log.Dbug("hook on conn stop empty")
}
