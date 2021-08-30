package protocol

import (
	"battle_rabbit/codec"
	"battle_rabbit/define"
	jsoniter "github.com/json-iterator/go"
)
var (
	respError ,_ = jsoniter.Marshal(MsgPush{Code: 500})
	respSuccess ,_ = jsoniter.Marshal(MsgPush{Code: 200})
)
type MsgPush struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func Err(mid uint32) *codec.Message {
	return &codec.Message{
		Id:   mid,
		Data: respError,
	}
}

func ErrCode(mid uint32, ErrCode int) *codec.Message {
	data, _ := jsoniter.Marshal(MsgPush{Code: ErrCode})
	return &codec.Message{
		Id:   mid,
		Data: data,
	}
}

func ErrData(mid uint32, ErrCode int, body interface{}) *codec.Message {
	data, _ := jsoniter.Marshal(MsgPush{
		Code: ErrCode,
		Data: body,
	})
	return &codec.Message{
		Id:   mid,
		Data: data,
	}
}

func  Success(mid uint32) *codec.Message {
	return &codec.Message{
		Id:   mid,
		Data: respSuccess,
	}
}

func  SuccessData(mid uint32, body interface{}) *codec.Message {
	data, _ := jsoniter.Marshal(MsgPush{
		Code: 200,
		Data: body,
	})
	return &codec.Message{
		Id:   mid,
		Data: data,
	}
}


//// 创建推送数据
func MakePushMsg(mid uint32, data interface{}) (msg *codec.Message) {
	resp := &MsgPush{Code: define.MsgCode200, Data: data}
	msg = new(codec.Message)
	msg.Id = mid
	msg.Data, _ = jsoniter.Marshal(resp)
	return msg
}
//
//// 创建请求数据
func MakeReqMsg(mid uint32, data interface{}) (msg *codec.Message, err error) {
	msg = new(codec.Message)
	msg.Id = mid
	msg.Data, err = jsoniter.Marshal(data)
	return msg, err
}
