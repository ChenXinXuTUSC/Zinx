package test

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"testing"
	"time"
	service "zinx/service"
)

func TestMockClient(t *testing.T) {

	s := service.NewServer("testserver", 7777)
	var wg sync.WaitGroup
	wg.Add(1)
	go s.Serve()

	time.Sleep(1 * time.Second)

	conn, connErr := net.Dial("tcp", "127.0.0.1:7777")
	if connErr != nil {
		panic(connErr.Error())
	}

	var randomGen = func(buf []byte, n int) {
		for i := 0; i < min(n, len(buf)); i++ {
			if rand.Int31n(2) == 0 {
				buf[i] = byte(65 + rand.Int31n(26))
			} else {
				buf[i] = byte(97 + rand.Int31n(26))
			}
		}
	}

	// table driven test
	cases := []struct {
		name  string
		txlen int
	}{
		{"test1", 1},
		{"test2", 16},
		{"test3", 64},
		{"test4", 128},
	}

	var rxbuf []byte = make([]byte, 512)
	var txbuf []byte = make([]byte, 512)
	for _, tc := range cases {
		randomGen(txbuf, tc.txlen)
		txLen, txErr := conn.Write(txbuf[:tc.txlen])
		if txErr != nil {
			t.Fail()
		}
		rxLen, rxErr := conn.Read(rxbuf)
		if rxErr != nil {
			t.Fail()
		}
		if txLen != rxLen {
			t.Fail()
		}
		if string(txbuf[:tc.txlen]) != string(rxbuf[:rxLen]) {
			t.Fail()
		}
		fmt.Printf("tx %s, rx %s", string(txbuf[:tc.txlen]), string(rxbuf[:rxLen]))
	}

	wg.Done()
	wg.Wait()
}
