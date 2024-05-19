package core

import (
	"bytes"
	"encoding/binary"
	"errors"
	zinf "zinx/interface"
	config "zinx/utils/conf"
	"zinx/utils/log"
)

type Message struct {
	Id      uint32
	DataLen uint32
	Data    []byte
}

func NewMsg(id uint32, data []byte) *Message {
	return &Message{
		Id: id,
		DataLen: uint32(len(data)),
		Data: data,
	}
}
func NewMsgWithLen(id, dataLen uint32, data []byte) *Message {
	return &Message{
		Id: id,
		DataLen: dataLen,
		Data: data,
	}
}

func (mp *Message) GetMsgId() uint32 {
	return mp.Id
}
func (mp *Message) SetMsgId(id uint32) {
	mp.Id = id
}

func (mp *Message) GetDataLen() uint32 {
	return mp.DataLen
}
func (mp *Message) SetDataLen(dataLen uint32) {
	mp.DataLen = dataLen
}

func (mp *Message) GetData() []byte {
	return mp.Data
}
func (mp *Message) SetData(data []byte) {
	mp.Data = data
}


// message is sent in TCP stream, no message boundary
// thus each message need to be specified with length
// for a broken TCP conn, the peer will send  an  RST
// to shut down the conn, which results in EOF  error
// during block read
type DataHandler struct {}
func NewDataHandler() *DataHandler {
	return &DataHandler{}
}
func (dp *DataHandler) GetHeadLen() uint {
	return 8
}
func (dp *DataHandler) DataPack(msg zinf.ZinfMessage) ([]byte, error) {
	var buf bytes.Buffer

	// store data len
	if wrErr := binary.Write(&buf, binary.LittleEndian, msg.GetDataLen()); wrErr != nil {
		log.Erro(wrErr.Error())
		return nil, wrErr
	}
	// store msg id
	if wrErr := binary.Write(&buf, binary.LittleEndian, msg.GetMsgId()); wrErr != nil {
		log.Erro(wrErr.Error())
		return nil, wrErr
	}
	// store data
	if wrErr := binary.Write(&buf, binary.LittleEndian, msg.GetData()[:msg.GetDataLen()]); wrErr != nil {
		log.Erro(wrErr.Error())
		return nil, wrErr
	}

	return buf.Bytes(), nil
}
func (dp *DataHandler) DataUnpack(data []byte) (zinf.ZinfMessage, error) {
	var buf = bytes.NewReader(data)

	var msgp = new(Message)

	// read data len
	if rdErr := binary.Read(buf, binary.LittleEndian, &msgp.DataLen); rdErr != nil {
		log.Erro(rdErr.Error())
		return nil, rdErr
	}
	// read msg id
	if rdErr := binary.Read(buf, binary.LittleEndian, &msgp.Id); rdErr != nil {
		log.Erro(rdErr.Error())
		return nil, rdErr
	}
	// read data
	if config.GlobalConfig.MaxPacketSize > 0 && msgp.DataLen > config.GlobalConfig.MaxPacketSize {
		log.Warn("message packet size too large (%d > %d)", msgp.DataLen, config.GlobalConfig.MaxPacketSize)
		return nil, errors.New("message packet size too large")
	}

	return msgp, nil
}
