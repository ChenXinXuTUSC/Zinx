package core

import (
	zinf "zinx/interface"
)

type Request struct {
	conn zinf.ZinfConnection
	data []byte
}

func (rp *Request) GetConnection() zinf.ZinfConnection {
	return rp.conn
}

func (rp *Request) GetData() []byte {
	return rp.data
}
