package service

import (
	"net"
	"zinx/utils/log"
)

func EchoBack(conn *net.TCPConn, rxbuf []byte, rxLen int) error {
	log.Dbug("receive bytes: %s", string(rxbuf))
	_, txErr := conn.Write(rxbuf[:rxLen])
	if txErr != nil {
		log.Erro(txErr.Error())
		return txErr
	}

	return nil
}
