package core

import (
	"fmt"
	"sync"
	zinf "zinx/interface"
	"zinx/utils/log"
)

type ConnManager struct {
	connMap  map[uint32]zinf.ZinfConnection
	connLock sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connMap: make(map[uint32]zinf.ZinfConnection),
	}
}

func (cmp *ConnManager) Add(conn zinf.ZinfConnection) {
	cmp.connLock.Lock()
	defer cmp.connLock.Unlock()

	cmp.connMap[conn.GetConnID()] = conn

	log.Info("conn %d added", conn.GetConnID())
}

func (cmp *ConnManager) Remove(conn zinf.ZinfConnection) {
	cmp.connLock.Lock()
	defer cmp.connLock.Unlock()

	delete(cmp.connMap, conn.GetConnID())

	log.Info("conn %d removed", conn.GetConnID())
}

func (cmp *ConnManager) Get(connId uint32) (zinf.ZinfConnection, error) {
	cmp.connLock.RLock() // no need for WLock
	defer cmp.connLock.RUnlock()

	conn, ok := cmp.connMap[connId]
	if !ok {
		return nil, fmt.Errorf("not a valid connection connId: %d", connId)
	}
	return conn, nil
}

func (cmp *ConnManager) ClearAll() {
	cmp.connLock.Lock()
	defer cmp.connLock.Unlock()

	for id, conn := range cmp.connMap {
		conn.Stop()
		delete(cmp.connMap, id)
	}

	log.Info("remove all connections")
}

func (cmp *ConnManager) GetConnNum() int {
	cmp.connLock.Lock()
	defer cmp.connLock.Unlock()

	return len(cmp.connMap)
}
