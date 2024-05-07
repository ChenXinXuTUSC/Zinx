package zinf

import "net"

type ZinfConnection interface {
	Start()
	Stop()
	GetTCPConnection() *net.TCPConn
	GetConnID() uint32
	RemoteAddr() net.Addr
}

type Handler func(*net.TCPConn, []byte, int) error
