package codec

import (
	"battle_rabbit/utils/encrypt"
)

// 默认加解密
type (
	AESCoder struct {
		Key []byte
	}
)

//encode
func (aes *AESCoder)Encode(data []byte)([]byte,error)  {
	return encrypt.AesEncrypt(data,aes.Key)
}

// decode
func (aes *AESCoder)Decode(data []byte)([]byte,error)  {
	return encrypt.AesDecrypt(data,aes.Key)
}
