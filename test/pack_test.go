package test

import (
	"fmt"
	"io"
	"net"
	"sync"
	"testing"
	"time"
	"zinx/core"
	zinf "zinx/interface"
	"zinx/utils/log"
)


func TestPackUnPack(t *testing.T) {
	listener, listenErr := net.Listen("tcp", "127.0.0.1:7777")
	if listenErr != nil {
		log.Erro(listenErr.Error())
		t.FailNow()
	}

	var wg sync.WaitGroup

	// start server routine
	wg.Add(1)
	go func() {
		defer wg.Done()

		var dataHandler = core.NewDataHandler()
		var msgCnt int = 0
		for i := 0; i < 4; i++{
			conn, acceptErr := listener.Accept()
			if acceptErr != nil {
				log.Erro(acceptErr.Error())
				t.Fail()
				time.Sleep(1 * time.Second)
			}

			wg.Add(1)
			go func(net.Conn) {
				defer wg.Done()

				// loop read
				for {
					// step 1 read out head
					var headData []byte = make([]byte, dataHandler.GetHeadLen())
					_, rxErr := io.ReadFull(conn, headData) // loop read until reach buf len
					if rxErr != nil {
						log.Erro(rxErr.Error())
						return
					}
	
					// step 2 read msg id
					head, unpackErr := dataHandler.DataUnpack(headData)
					if unpackErr != nil {
						log.Erro(unpackErr.Error())
						return
					}
	
					// step 3 continue to read from conn io
					var msgp *core.Message = nil
					if head.GetDataLen() > 0 {
						msgp = head.(*core.Message)
						msgp.Data = make([]byte, head.GetDataLen())
						_, rxErr := io.ReadFull(conn, msgp.Data)
						if rxErr != nil {
							log.Erro(rxErr.Error())
							return
						}
					}
	
					log.Info("[msg#%d, %d]: %s", msgp.Id, msgp.DataLen, string(msgp.Data))
					msgCnt++
				}
			}(conn)
		}
	}()

	cases := []struct {
		Id uint32
		DataLen uint32
		Data []byte
	}{
		{0, 5, []byte("hello")},
		{1, 5, []byte("world")},
		{2, 12, []byte("hello, world")},
	}

	// start client routine
	time.Sleep(1 * time.Second) // wait for server launched
	wg.Add(1)
	go func() {
		defer wg.Done()

		var dataHandler = core.NewDataHandler()

		for _, c := range cases {
			wg.Add(1)
			go func(msg zinf.ZinfMessage) {
				defer wg.Done()

				conn, dialErr := net.Dial("tcp", "127.0.0.1:7777")
				if dialErr != nil {
					t.Fail()
					log.Erro(dialErr.Error())
					return
				}

				data, packErr := dataHandler.DataPack(msg)
				if packErr != nil {
					t.Fail()
					log.Erro(packErr.Error())
					return
				}

				conn.Write(data)
				conn.Close()
			}(&core.Message{
				Id: c.Id,
				DataLen: c.DataLen,
				Data: c.Data,
			})
			
		}
	}()

	// compact test
	var msg1 = &core.Message{
		Id: 3,
		DataLen: 5,
		Data: []byte("abcdefghijklmn"),
	}
	var msg2 = &core.Message{
		Id: 4,
		DataLen: 5,
		Data: []byte("fghijklmnabcde"),
	}
	var dataHandler = core.NewDataHandler()
	data1, packErr1 := dataHandler.DataPack(msg1)
	if packErr1 != nil {
		log.Erro(packErr1.Error())
		t.FailNow()
	}
	data2, packErr2 := dataHandler.DataPack(msg2)
	if packErr2 != nil {
		log.Erro(packErr2.Error())
		t.FailNow()
	}
	fmt.Println(len(data1), len(data2))
	conn, connErr := net.Dial("tcp", "127.0.0.1:7777")
	if connErr != nil {
		log.Erro(connErr.Error())
		t.FailNow()
	}

	conn.Write(append(data1, data2...))
	conn.Close()

	wg.Wait()
}

func TestConnBroke(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		listener, listenErr := net.Listen("tcp", "127.0.0.1:7778")
		if listenErr != nil {
			log.Erro(listenErr.Error())
			t.Fail()
			return
		}

		// block for one connection
		conn, acceptErr := listener.Accept()
		log.Info("accept at %d", time.Now().Unix())
		if acceptErr != nil {
			log.Erro(acceptErr.Error())
			t.Fail()
			return
		}

		conn.Close()
		log.Info("server shut down %d", time.Now().Unix())
	}()


	// client
	time.Sleep(1000 * time.Millisecond)
	log.Info("start client %d", time.Now().Unix())
	conn, dialErr := net.Dial("tcp", "127.0.0.1:7778")
	if dialErr != nil {
		log.Erro(dialErr.Error())
		t.FailNow()
	}
	
	for i := 0; i < 10; i++ {
		txLen, txErr := conn.Write([]byte(fmt.Sprintf("msg #%d", i+1)))
		fmt.Println("msg", i+1, time.Now().Unix(), txLen, txErr)
		time.Sleep(1 * time.Millisecond)
	}
	
	// time.Sleep(2000 * time.Millisecond)
	// log.Info("client write #1 %d", time.Now().Unix())
	// txLen1, txErr1 := conn.Write([]byte("hello, world"))
	// fmt.Println(txLen1, txErr1)

	// time.Sleep(2000 * time.Millisecond)
	// log.Info("client write #2 %d", time.Now().Unix())
	// txLen2, txErr2 := conn.Write([]byte("hello, world"))
	// fmt.Println(txLen2, txErr2)

	wg.Wait()
}
