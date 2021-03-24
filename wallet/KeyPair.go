package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

/**
地址所对应的密钥对(私钥+公钥),封装在一个自定义的结构体中
*/
type KeyPair struct {
	Priv *ecdsa.PrivateKey
	Pub  []byte
}
//该函数用于生成并返回一对密钥对
func NewKeyPair() (*KeyPair, error) {
	curve := elliptic.P256()
	pri, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}
	pub := elliptic.Marshal(curve, pri.X, pri.Y)

	return &KeyPair{
		Priv: pri,
		Pub:  pub,
	}, nil
}
