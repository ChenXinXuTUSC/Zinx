package core

import (
	"errors"
	"io"
	"net"
	zinf "zinx/interface"
	config "zinx/utils/conf"
	"zinx/utils/log"
)

type Connection struct {
	Conn   *net.TCPConn
	ConnID uint32
	Exit   chan bool // inform the conn has finished

	isClosed   bool
	msgHandler zinf.ZinfMsgHandler
	msgChan    chan []byte
}

func NewConnection(conn *net.TCPConn, connID uint32, msgHandler zinf.ZinfMsgHandler) *Connection {
	cp := &Connection{
		Conn:       conn,
		ConnID:     connID,
		Exit:       make(chan bool, 1),
		isClosed:   false,
		msgHandler: msgHandler,
		msgChan:    make(chan []byte),
	}
	return cp
}

func (cp *Connection) Start() {
	// launch the data listenning loop
	go cp.StartReader()
	go cp.StartWriter()

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

func (cp *Connection) StartReader() {
	// receive data from conn and transfer
	// to user's handler callback
	log.Info("conn %d reader routine start", cp.ConnID)
	defer cp.Stop()

	var dataHandler = NewDataHandler()
	for {
		var headData []byte = make([]byte, dataHandler.GetHeadLen())
		if _, readHeadErr := io.ReadFull(cp.GetTCPConnection(), headData); readHeadErr != nil {
			if readHeadErr.Error() != "EOF" {
				log.Erro("read head error: %s", readHeadErr.Error())
				cp.Exit <- true
			}
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

		// turn over to worker pool or start one routine
		if config.GlobalConfig.NumWorker > 0 {
			cp.msgHandler.SendMsgToTaskQueue(&req)
		} else {
			go cp.msgHandler.DoMsgHandler(&req) // temporary go routine
		}
	}
}

func (cp *Connection) StartWriter() {
	log.Info("conn %d writer routine start", cp.ConnID)
	defer cp.Stop()

	for {
		select {
		case data := <- cp.msgChan:
			if _, txErr := cp.Conn.Write(data); txErr != nil {
				log.Erro(txErr.Error())
				return
			}
		case <- cp.Exit:
			return
		}
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
	// will invoked by other handler to send data
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

	// turn over the write business to writer goroutine
	cp.msgChan <- packedData

	return nil
}
