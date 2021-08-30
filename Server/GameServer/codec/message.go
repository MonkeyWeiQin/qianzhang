package codec

import (
	"battle_rabbit/iface"
	"battle_rabbit/service/log"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

const (
	MsgHeadLen    uint32 = 8
	MaxDataLength        = 10240
)

// TCP网关与客户端的消息定义接收者
// 采用[大端序]为字节存储顺序
// 采用类似LengthFieldBasedFrameDecoder的编码器
// 数据流 : |_①_|_②_|____③______|______....
// ① ==> uint32 数据总长度
// ② ==> uint32 消息ID
// ③ ==> []byte 消息体
type Message struct {
	Id   uint32 //消息的ID
	Data []byte //消息的内容
}

//创建一个Message消息包
func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		Id:   id,
		Data: data,
	}
}

//获取消息数据段长度
func (msg *Message) GetDataLen() uint32 {
	return uint32(len(msg.Data))
}

//获取消息ID
func (msg *Message) GetMsgId() uint32 {
	return msg.Id
}

//获取消息内容
func (msg *Message) GetData() []byte {
	return msg.Data
}

//设计消息ID
func (msg *Message) SetMsgId(msgId uint32) {
	msg.Id = msgId
}

//设计消息内容
func (msg *Message) SetData(data []byte) {
	msg.Data = data
}

//封包并写入数据
func (msg *Message) WriterPack(writer io.Writer) error {
	return WritePackager(msg, writer)
}

//读取数据并拆包
func (msg *Message) ReadPack(reader io.Reader) error {
	return ReadPackager(msg, reader)
}

func WritePackager(msg iface.IMessage, writer io.Writer) (err error) {
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//写dataLen
	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen())
	if err != nil {
		return err
	}
	//写msgID
	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId())
	if err != nil {
		return err
	}

	//写data数据
	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetData())
	if err != nil {
		return err
	}
	n, err := writer.Write(dataBuff.Bytes())
	log.Debug("+++++++ 发送数据长度: %d +++++++", n)
	return
}

func ReadPackager(msg iface.IMessage, reader io.Reader) (err error) {
	// 读取头
	var head = make([]byte, MsgHeadLen)
	_, err = io.ReadFull(reader, head)
	if err != nil {
		return err
	}
	//ioReader
	headBuff := bytes.NewReader(head)
	//dataLen
	var dataLen, mid uint32
	err = binary.Read(headBuff, binary.LittleEndian, &dataLen)
	if err != nil {
		return err
	}

	//msgID
	err = binary.Read(headBuff, binary.LittleEndian, &mid)
	if err != nil {
		return err
	}

	//log.Debug("dataLen ====== ",dataLen)
	//log.Debug("mid ====== ",mid)

	msg.SetMsgId(mid)
	//读取body
	if dataLen > 0 {
		var data = make([]byte, dataLen)
		_, err = io.ReadFull(reader, data)
		if err != nil {
			return err
		}
		if dataLen > MaxDataLength {
			return errors.New("[ readPackager ] message too long!  len > 4096")
		}
		msg.SetData(data)
	}
	//log.Debug("data ====== ",string(msg.GetData()))
	return
}
//
//type MsgPush struct {
//	Code int         `json:"code"`
//	Data interface{} `json:"data"`
//}
//
//func (msg *Message) Err(mid uint32) (err error) {
//	resp := &MsgPush{Code: 500}
//	msg.Id = mid
//	msg.Data, err = jsoniter.Marshal(resp)
//	return err
//}
//
//func (msg *Message) ErrCode(mid uint32, ErrCode int) (err error) {
//	resp := &MsgPush{Code: ErrCode}
//	msg.Id = mid
//	msg.Data, err = jsoniter.Marshal(resp)
//	return err
//}
//
//func (msg *Message) ErrData(mid uint32, ErrCode int, data interface{}) (err error) {
//	resp := &MsgPush{
//		Code: ErrCode,
//		Data: data,
//	}
//
//	msg.Id = mid
//	msg.Data, err = jsoniter.Marshal(resp)
//	return err
//}
//
//func (msg *Message) Success(mid uint32) (err error) {
//	resp := &MsgPush{Code: 200}
//	msg.Id = mid
//	msg.Data, err = jsoniter.Marshal(resp)
//	return err
//}
//
//func (msg *Message) SuccessData(mid uint32, data interface{}) (err error) {
//	resp := &MsgPush{Code: 200, Data: data}
//	msg.Id = mid
//	msg.Data, err = jsoniter.Marshal(resp)
//	return err
//}
//
//
//// 创建推送数据
//func MakePushMsg(mid uint32, data interface{}) (msg *Message, err error) {
//	resp := &MsgPush{Code: define.MsgCode200, Data: data}
//	msg = new(Message)
//	msg.Id = mid
//	msg.Data, err = jsoniter.Marshal(resp)
//	return msg, err
//}
//
//// 创建请求数据
//func MakeReqMsg(mid uint32, data interface{}) (msg *Message, err error) {
//	msg = new(Message)
//	msg.Id = mid
//	msg.Data, err = jsoniter.Marshal(data)
//	return msg, err
//}