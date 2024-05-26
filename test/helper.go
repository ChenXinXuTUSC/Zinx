package test

import (
	"errors"
	"io"
	"math/rand"
	"net"
	"zinx/core"
	"zinx/utils/log"
)

func RandomFill(buf []byte, n int) {
	for i := 0; i < min(n, len(buf)); i++ {
		if rand.Int31n(2) == 0 {
			buf[i] = byte(65 + rand.Int31n(26))
		} else {
			buf[i] = byte(97 + rand.Int31n(26))
		}
	}
}

func MockClient(clientId, msgId uint32, msgData []byte) error {
	// simulate core.connection on client side
	conn, dialErr := net.Dial("tcp", "127.0.0.1:7777")
	if dialErr != nil {
		log.Erro("dial error: %s", dialErr.Error())
		return dialErr
	}
	defer conn.Close()

	dp := core.NewDataHandler()

	packedData, packErr := dp.DataPack(core.NewMsg(msgId, msgData))
	if packErr != nil {
		log.Erro(packErr.Error())
		return packErr
	}
	if _, txErr := conn.Write(packedData); txErr != nil {
		log.Erro("send error: %s", txErr.Error())
		return txErr
	}

	// read echo
	for {
		var headData []byte = make([]byte, dp.GetHeadLen())
		if _, rxErr := io.ReadFull(conn, headData); rxErr != nil {
			log.Erro("recv error: %s", rxErr.Error())
			return rxErr
		}

		msg, unpackErr := dp.DataUnpack(headData)
		if unpackErr != nil {
			log.Erro("unpack error: %s", unpackErr.Error())
			return unpackErr
		}
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, rxErr := io.ReadFull(conn, data); rxErr != nil {
				log.Erro("recv error: %s", rxErr)
				return rxErr
			}
		}
		msg.SetData(data)

		if data == nil {
			log.Erro("no valid data")
			return errors.New("no valid data")
		}
		if msg.GetMsgId() == 2 {
			log.Info("client#%d recv: type %d, len %d", clientId, msg.GetMsgId(), msg.GetDataLen())
			break
		}
		if msg.GetMsgId() != msgId {
			log.Erro("msg type mismatched")
			return errors.New("msg type mismatched")
		}

		log.Info("client#%d recv: type %d, len %d", clientId, msg.GetMsgId(), msg.GetDataLen())
	}

	return nil
}
