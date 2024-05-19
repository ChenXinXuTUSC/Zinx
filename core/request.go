package core

import (
	zinf "zinx/interface"
)

type Request struct {
	conn zinf.ZinfConnection
	msgi zinf.ZinfMessage
}

func (rp *Request) GetConnection() zinf.ZinfConnection {
	return rp.conn
}

func (rp *Request) GetData() []byte {
	return rp.msgi.GetData()
}

func (rp *Request) GetMsgId() uint32 {
	return rp.msgi.GetMsgId()
}
