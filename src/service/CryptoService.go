package service

import (
	"crypto/rand"
	"math/big"
)

type EncryptKeys struct {
	PublicKey  string
	PrivateKey string
}

// CreateDHKey 创建DH公私钥
func CreateDHKey(prime, generator string) EncryptKeys {
	// 创建私钥
	privateNum, _ := rand.Int(rand.Reader, new(big.Int).Exp(big.NewInt(2), big.NewInt(150), nil))
	privateKey := privateNum.Text(16)

	// 创建公钥
	generatorNum, _ := new(big.Int).SetString(generator, 16)
	primeNum, _ := new(big.Int).SetString(prime, 16)
	publicNum := new(big.Int).Exp(generatorNum, privateNum, primeNum)
	publicKey := publicNum.Text(16)

	return EncryptKeys{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
}

// SharedAESKey 生成共享密钥
func SharedAESKey(publicKey, privateKey, prime string) string {
	publicNum, _ := new(big.Int).SetString(publicKey, 16)
	privateNum, _ := new(big.Int).SetString(privateKey, 16)
	primeNum, _ := new(big.Int).SetString(prime, 16)
	aesKeyNum := new(big.Int).Exp(publicNum, privateNum, primeNum)
	aesKey := aesKeyNum.Text(16)

	// 取后32位
	if len(aesKey) > 32 {
		aesKey = aesKey[len(aesKey)-32:]
	}
	return aesKey
}
