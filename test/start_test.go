package test

import (
	"fmt"
	"net"
	"sync"
	"testing"
	"time"
	core "zinx/core"
	"zinx/service"
	"zinx/utils/toolbox"
)

func TestMockClient(t *testing.T) {
	s := core.NewServer("testserver", 7777)
	var wg sync.WaitGroup
	wg.Add(1)
	go s.Serve()

	time.Sleep(1 * time.Second)

	conn, connErr := net.Dial("tcp", "127.0.0.1:7777")
	if connErr != nil {
		panic(connErr.Error())
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
		toolbox.RandomFill(txbuf, tc.txlen)
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

func TestRouter(t *testing.T) {
	s := core.NewServer("testserver", 7777)
	s.AddRouter(&service.PingRouter{})
	var wg sync.WaitGroup
	wg.Add(1)
	go s.Serve()

	time.Sleep(1 * time.Second)

	conn, connErr := net.Dial("tcp", "127.0.0.1:7777")
	if connErr != nil {
		panic(connErr.Error())
	}

	// table driven test
	cases := []struct {
		name  string
		txlen int
	}{
		{"test1", 1},
		{"test2", 2},
		{"test3", 4},
		{"test4", 8},
	}

	var rxbuf []byte = make([]byte, 512)
	var txbuf []byte = make([]byte, 512)
	for _, tc := range cases {
		toolbox.RandomFill(txbuf, tc.txlen)
		_, txErr := conn.Write(txbuf[:tc.txlen])
		if txErr != nil {
			t.Fail()
		}
		rxLen, rxErr := conn.Read(rxbuf)
		if rxErr != nil {
			t.Fail()
		}
		fmt.Println(string(rxbuf[:rxLen]))
	}

	wg.Done()
	wg.Wait()
}
