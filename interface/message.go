package zinf

type ZinfMessage interface {
	GetDataLen() uint32
	SetDataLen(uint32)
	
	GetMsgId() uint32
	SetMsgId(uint32)

	GetData() []byte
	SetData([]byte)
}

type ZinfDataHandler interface {
	GetHeadLen() uint32
	DataPack(ZinfMessage) ([]byte, error)
	DataUnpack([]byte) (ZinfMessage, error)
}
