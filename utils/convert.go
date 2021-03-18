package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
)

/**
*将int64类型的数据转换为[]byte类型
 */
func Int2Byte(num int64)([]byte,error) {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	return buff.Bytes(),err
}
func SHA256HashBlock(data []byte)[]byte  {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

//gob 序列化编码
func GobEncode(v interface{})([]byte,error)  {
	buff := new(bytes.Buffer)
	err := gob.NewEncoder(buff).Encode(v)
	return buff.Bytes(), err
}

func GobDecode(data []byte,v *interface{})(interface{},error)  {

	err := gob.NewDecoder(bytes.NewReader(data)).Decode(v)
	return v, err

}