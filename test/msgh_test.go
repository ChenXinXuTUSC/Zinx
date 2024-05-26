package test

import (
	"sync"
	"testing"
	"time"
	"zinx/core"
	"zinx/service"
)

func TestMsgHandler(t *testing.T) {
	s := core.NewServer()
	s.AddRouter(0, &service.PingRouter{})
	s.AddRouter(1, &service.VersionRouter{})
	go s.Serve()
	time.Sleep(1 * time.Second) // wait for server launch

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := MockClient(0, 0, []byte("try msg type 0")); err != nil {
			t.Fail()
		}

	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := MockClient(1, 0, []byte("try msg type 1")); err != nil {
			t.Fail()
		}
	}()

	wg.Wait()
}
