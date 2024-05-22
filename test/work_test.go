package test

import (
	"sync"
	"testing"
	"time"
	"zinx/core"
	"zinx/service"
)

func TestWorkerPool(t *testing.T) {
	s := core.NewServer()
	s.AddRouter(0, &service.PingRouter{})
	s.AddRouter(1, &service.ZinxVerRouter{})
	go s.Serve()
	time.Sleep(1 * time.Second) // wait for server launch

	var wg sync.WaitGroup

	for i := 0; i < int(1e3); i++ {
		wg.Add(1)
		go func(clid int) {
			defer wg.Done()
			var buf []byte = make([]byte, 1024)
			RandomFill(buf, len(buf))
			if err := MockClient(uint32(clid), uint32(clid % 2), buf); err != nil {
				t.Fail()
			}
		}(i)
	}

	wg.Wait()
}
