package zinf

import "net"

type ZinfConnection interface {
	Start()
	Stop()
	GetTCPConnection() *net.TCPConn
	GetConnID() uint32
	GetRemoteAddr() net.Addr

	SendMsg(msgId uint32, data []byte) error
}

type Handler func(*net.TCPConn, []byte, int) error
