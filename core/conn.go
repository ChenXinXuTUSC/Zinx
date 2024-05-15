package core

import (
	"net"
	zinf "zinx/interface"
	"zinx/utils/log"
)

type Connection struct {
	Conn   *net.TCPConn
	ConnID uint32
	Exit   chan bool // inform the conn has finished

	isClosed bool
	router   zinf.ZinfRouter // handler method
}

func NewConnection(conn *net.TCPConn, connID uint32, router zinf.ZinfRouter) *Connection {
	cp := &Connection{
		Conn:     conn,
		ConnID:   connID,
		Exit:     make(chan bool, 1),
		isClosed: false,
		router:   router,
	}
	return cp
}

func (cp *Connection) Start() {
	// launch the data listenning loop
	go cp.ReadLoop()

	for {
		select {
		case <-cp.Exit:
			return
		}
	}
}

func (cp *Connection) Stop() {
	if cp.isClosed {
		return // already closed
	}
	cp.isClosed = true

	// close the tcp connection
	cp.Conn.Close()

	// inform that this conn is terminated
	cp.Exit <- true

	// close channel
	close(cp.Exit)
}

func (cp *Connection) ReadLoop() {
	// receive data from conn and transfer
	// to user's handler callback
	log.Info("conn %d start listenning rx data", cp.ConnID)
	defer log.Info("conn %d finish listenning rx data", cp.ConnID)
	defer cp.Stop()

	var rxbuf []byte = make([]byte, 512)
	for {
		// read at most 512 bytes from conn
		rxLen, rxErr := cp.Conn.Read(rxbuf)
		if rxErr != nil {
			log.Erro("conn %d read error", rxErr.Error())
			cp.Exit <- true
			return
		}

		req := Request{
			conn: cp,
			data: rxbuf[:rxLen],
		}
		go func(request zinf.ZinfRequest) {
			cp.router.PreProcess(request) // will invoke BaseRouter's empty method
			cp.router.Handle(request)
			cp.router.PostProcess(request) // will invoke BaseRouter's empty method
		}(&req)

	}
}

func (cp *Connection) GetTCPConnection() *net.TCPConn {
	return cp.Conn
}
func (cp *Connection) GetConnID() uint32 {
	return cp.ConnID
}
func (cp *Connection) GetRemoteAddr() net.Addr {
	return cp.Conn.RemoteAddr()
}
