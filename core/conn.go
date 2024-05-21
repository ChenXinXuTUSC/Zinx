package core

import (
	"errors"
	"io"
	"net"
	zinf "zinx/interface"
	"zinx/utils/log"
)

type Connection struct {
	Conn   *net.TCPConn
	ConnID uint32
	Exit   chan bool // inform the conn has finished

	isClosed   bool
	msgHandler zinf.ZinfMsgHandler
}

func NewConnection(conn *net.TCPConn, connID uint32, msgHandler zinf.ZinfMsgHandler) *Connection {
	cp := &Connection{
		Conn:       conn,
		ConnID:     connID,
		Exit:       make(chan bool, 1),
		isClosed:   false,
		msgHandler: msgHandler,
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

	var dataHandler = NewDataHandler()
	for {
		var headData []byte = make([]byte, dataHandler.GetHeadLen())
		if _, readHeadErr := io.ReadFull(cp.GetTCPConnection(), headData); readHeadErr != nil {
			log.Erro("read head error: %s", readHeadErr.Error())
			cp.Exit <- true
			continue
		}

		msgp, unpackErr := dataHandler.DataUnpack(headData)
		if unpackErr != nil {
			log.Erro("unpack error: %s", unpackErr.Error())
			cp.Exit <- true
			continue
		}

		var data []byte = make([]byte, 0)
		if msgp.GetDataLen() > 0 {
			data = make([]byte, msgp.GetDataLen())
			if _, readDataErr := io.ReadFull(cp.GetTCPConnection(), data); readDataErr != nil {
				log.Erro("read data error: %s", readDataErr.Error())
				cp.Exit <- true
				continue
			}
		}
		msgp.SetData(data)

		req := Request{
			conn: cp,
			msgi: msgp,
		}

		go cp.msgHandler.DoMsgHandler(&req)
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

func (cp *Connection) SendMsg(msgId uint32, data []byte) error {
	if cp.isClosed {
		log.Erro("write to a closed conn")
		return errors.New("this connection has already been closed")
	}

	var dp = NewDataHandler()
	packedData, packErr := dp.DataPack(NewMsg(msgId, data))
	if packErr != nil {
		log.Erro("pack msg data error: %s", packErr.Error())
		return packErr
	}

	if _, txErr := cp.Conn.Write(packedData); txErr != nil {
		log.Erro("erro sending data: %s", txErr.Error())
		cp.Exit <- true
		return txErr
	}

	return nil
}
