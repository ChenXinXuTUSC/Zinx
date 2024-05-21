package test

import (
	"errors"
	"io"
	"net"
	"sync"
	"testing"
	"time"
	"zinx/core"
	"zinx/service"
	"zinx/utils/log"
)

func TestMsgHandler(t *testing.T) {
	s := core.NewServer()
	s.AddRouter(0, &service.PingRouter{})
	s.AddRouter(1, &service.ZinxVerRouter{})
	go s.Serve()
	time.Sleep(1 * time.Second) // wait for server launch

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := MockClient(0, []byte("try msg type 0")); err != nil {
			t.Fail()
		}

	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := MockClient(1, []byte("try msg type 1")); err != nil {
			t.Fail()
		}
	}()

	wg.Wait()
}

func MockClient(msgId uint32, msgData []byte) error {
	// simulate core.connection on client side
	conn, dialErr := net.Dial("tcp", "127.0.0.1:7777")
	if dialErr != nil {
		log.Erro("dial error: %s", dialErr.Error())
		return dialErr
	}
	defer conn.Close()

	dp := core.NewDataHandler()

	packedData, packErr := dp.DataPack(core.NewMsg(msgId, msgData))
	if packErr != nil {
		log.Erro(packErr.Error())
		return packErr
	}
	if _, txErr := conn.Write(packedData); txErr != nil {
		log.Erro("send error: %s", txErr.Error())
		return txErr
	}

	// read echo
	var headData []byte = make([]byte, dp.GetHeadLen())
	if _, rxErr := io.ReadFull(conn, headData); rxErr != nil {
		log.Erro("recv error: %s", rxErr.Error())
		return rxErr
	}

	msg, unpackErr := dp.DataUnpack(headData)
	if unpackErr != nil {
		log.Erro("unpack error: %s", unpackErr.Error())
		return unpackErr
	}
	var data []byte
	if msg.GetDataLen() > 0 {
		data = make([]byte, msg.GetDataLen())
		if _, rxErr := io.ReadFull(conn, data); rxErr != nil {
			log.Erro("recv error: %s", rxErr)
			return rxErr
		}
	}
	msg.SetData(data)

	if data == nil {
		log.Erro("no valid data")
		return errors.New("no valid data")
	}
	if msg.GetMsgId() != msgId {
		log.Erro("msg type mismatched")
		return errors.New("msg type mismatched")
	}

	log.Info("receive server echo: type %d, len %d, data %s", msg.GetMsgId(), msg.GetDataLen(), string(msg.GetData()))

	return nil
}
