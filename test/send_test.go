package test

import (
	"io"
	"net"
	"sync"
	"testing"
	"time"
	"zinx/core"
	"zinx/service"
	"zinx/utils/log"
)

func TestSendMsg(t *testing.T) {
	s := core.NewServer()
	s.AddRouter(0, &service.PingRouter{})
	go s.Serve()
	time.Sleep(1 * time.Second) // wait for server launch

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		conn, connErr := net.Dial("tcp", "127.0.0.1:7777")
		if connErr != nil {
			log.Erro(connErr.Error())
			t.Fail()
			return
		}

		var dp = core.NewDataHandler()
		for i := 0; i < 3; i++ {
			// send message
			packedMsgData, packErr := dp.DataPack(core.NewMsg(0, []byte("Zinx v0.5 client test")))
			if packErr != nil {
				log.Erro(packErr.Error())
				t.Fail()
				return
			}
			if _, txErr := conn.Write(packedMsgData); txErr != nil {
				log.Erro(txErr.Error())
				t.Fail()
				continue
			}

			// read echo from server
			var headData []byte = make([]byte, dp.GetHeadLen())
			if _, rxErr := io.ReadFull(conn, headData); rxErr != nil {
				log.Erro(rxErr.Error())
				t.Fail()
				continue
			}

			msg, unpackErr := dp.DataUnpack(headData)
			if unpackErr != nil {
				log.Erro(unpackErr.Error())
				t.Fail()
				continue
			}
			var data []byte = make([]byte, 0)
			if msg.GetDataLen() > 0 {
				data = make([]byte, msg.GetDataLen())
				if _, rxErr := io.ReadFull(conn, data); rxErr != nil {
					log.Erro(rxErr.Error())
					t.Fail()
					continue
				}
			}
			msg.SetData(data)

			log.Info("receive server echo: %d %s", msg.GetMsgId(), string(msg.GetData()))
		}
	}()

	wg.Wait()
}
