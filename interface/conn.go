package zinf

import "net"

type ZinfConnection interface {
	Start()
	Stop()
	GetTCPConnection() *net.TCPConn
	GetConnID() uint32
	GetRemoteAddr() net.Addr

	SendBioMsg(msgId uint32, data []byte) error
	SendNioMsg(msgId uint32, data []byte) error
}

type Handler func(*net.TCPConn, []byte, int) error
