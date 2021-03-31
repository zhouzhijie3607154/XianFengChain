package utils

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"golang.org/x/crypto/ripemd160"
	"math/big"
)

/**
 * 将int类型的数据转换为[]byte类型
 */
func Int2Byte(num int64) ([]byte, error) {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	return buff.Bytes(), err
}

/**
  * 公钥字节切片转换为结构体
 */
func PubBytesToEcdsaPubKey(curve elliptic.Curve,pubKeyBytes[]byte)*ecdsa.PublicKey {
	x, y := elliptic.Unmarshal(elliptic.P256(), pubKeyBytes)
	return &ecdsa.PublicKey{
		elliptic.P256(),
		x, y,
	}
}
func SignBytesToSignature(signBytes []byte) (r,s *big.Int) {
	rBytes := signBytes[:len(signBytes)/2]
	sBytes := signBytes[len(signBytes)/2:]
	r = new(big.Int).SetBytes(rBytes)
	s = new(big.Int).SetBytes(sBytes)
	return r,s
}

/**
 * gob编码序列化
 */
func Encode(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(v)
	return buf.Bytes(), err
}

/**
 * gob反编码
 */
func Decode(data []byte, v interface{}) (interface{}, error) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(v)
	return v, err
}

/**
 * 将json格式的数组转换为对应的字符串类型的切片
 */
func JSONArray2String(array string) ([]string, error) {
	var stringSlice []string
	err := json.Unmarshal([]byte(array), &stringSlice)
	return stringSlice, err
}

/**
 * 将json格式的数组转换为对应的浮点型数据的切片
 */
func JSONArray2Float(array string) ([]float64, error) {
	var floatSlice []float64
	err := json.Unmarshal([]byte(array), &floatSlice)
	return floatSlice, err
}
/**
SHA256 哈希计算
 */
func Hash256(data []byte)([]byte)  {
	sha := sha256.New()
	sha.Write(data)
	return sha.Sum(nil)
}
/**
Ripemd160 的哈希计算
 */
func HashRipemd160(data []byte) []byte {
	hash := ripemd160.New()
	hash.Write(data)
	return hash.Sum(nil)
}
