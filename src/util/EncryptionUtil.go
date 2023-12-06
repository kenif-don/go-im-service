package util

import (
	utils "IM-Service/src/configs/err"
	"crypto/aes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"math/big"
	"strconv"
	"strings"
	"time"
)

var dealKeys = []int{
	0x07, 0xB6, 0x79, 0x56, 0x7A, 0x5C, 0x4A, 0xBE, 0x1D, 0xF1, 0xB2, 0x10, 0x3C, 0x5E, 0xDC, 0xA6,
	0x56, 0xE7, 0x88, 0x25, 0x87, 0x95, 0xD5, 0x85, 0x76, 0x7D, 0xEA, 0x66, 0xF5, 0x0A, 0xC3, 0xA8,
	0x55, 0x28, 0x67, 0x14, 0x06, 0xE7, 0xCB, 0x68, 0xAC, 0x2E, 0x00, 0x36, 0x57, 0x2F, 0xD2, 0xE2,
	0x54, 0xE9, 0xC6, 0xA3, 0x03, 0xC6, 0x07, 0x33, 0xBD, 0xF1, 0x6D, 0x46, 0x62, 0xFD, 0x82, 0xCF,
	0xA3, 0x50, 0x15, 0xB2, 0x53, 0xA4, 0x9C, 0x93, 0x98, 0x55, 0x8E, 0xF8, 0xC1, 0x0C, 0x15, 0x71,
	0x42, 0x6A, 0xA4, 0xF1, 0x5D, 0x72, 0xB1, 0xC4, 0xF6, 0xF0, 0x56, 0xAE, 0xCA, 0x77, 0x44, 0x45,
	0x21, 0x1B, 0x93, 0x40, 0x49, 0x89, 0x52, 0x76, 0x2C, 0x64, 0xB8, 0x3B, 0xF9, 0x8D, 0x51, 0xA5,
	0x80, 0x2C, 0x92, 0x39, 0xF7, 0xAD, 0xAF, 0x59, 0x1F, 0x06, 0xDE, 0x5A, 0x1D, 0x91, 0x1C, 0xDB,
	0x6F, 0xAD, 0xC1, 0xE8, 0xE5, 0xD4, 0xB4, 0x7C, 0x3E, 0x61, 0x73, 0x2D, 0xCE, 0xCD, 0x01, 0xDF,
	0x5E, 0xCE, 0x60, 0xB7, 0x83, 0xD1, 0x39, 0xA9, 0xF3, 0x35, 0x05, 0xBA, 0x88, 0x78, 0x97, 0xFC,
	0x3D, 0x2F, 0xF9, 0x36, 0x2A, 0x38, 0xB0, 0x25, 0x16, 0xA7, 0x08, 0x8C, 0xF6, 0x21, 0xC8, 0x22,
	0xBC, 0x90, 0x48, 0x35, 0x9A, 0x0D, 0x1A, 0xD9, 0xFA, 0xCC, 0x70, 0xAA, 0x42, 0x3F, 0xB6, 0xE1,
	0xBB, 0x41, 0x17, 0x74, 0xC2, 0x48, 0x7E, 0x80, 0xD6, 0x09, 0xC5, 0x24, 0x60, 0x30, 0x0E, 0xE3,
	0xFA, 0x92, 0x66, 0x43, 0xE1, 0x8A, 0x4D, 0xD7, 0x1B, 0x6B, 0x23, 0x65, 0xA0, 0x12, 0x9D, 0x9B,
	0xE0, 0x93, 0xE5, 0xD2, 0xE3, 0xF4, 0xDC, 0x41, 0xA4, 0x3A, 0x10, 0x2B, 0x96, 0xED, 0x1B, 0x1E,
	0xA9, 0xB4, 0x34, 0x11, 0x94, 0xA6, 0x75, 0x34, 0xD8, 0x89, 0xFC, 0x4F, 0x3B, 0x22, 0xB1, 0xA7,
}

func GetSignByTime(timestamp int64) string {
	timeStr := strconv.FormatInt(timestamp, 10)
	dealKey := strconv.Itoa(dealKeys[int(timestamp%int64(len(dealKeys)))])
	sign := MD5(timeStr + "_" + dealKey)
	return strings.ToUpper(sign)
}
func GetSign() (int64, string) {
	timestamp := time.Now().UnixMilli()
	sign := GetSignByTime(timestamp)
	return timestamp, sign
}
func MD5(v string) string {
	d := []byte(v)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}
func MD5Bytes(data []byte) string {
	m := md5.New()
	m.Write(data)
	return hex.EncodeToString(m.Sum(nil))
}

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
func EncryptAes(src, key string) (string, *utils.Error) {
	origData := []byte(src)
	data, err := EncryptAes2(origData, key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}
func EncryptAes2(data []byte, key string) ([]byte, *utils.Error) {
	cipher, _ := aes.NewCipher([]byte(key))
	length := (len(data) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, data)
	pad := byte(len(plain) - len(data))
	for i := len(data); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted := make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(data); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		if bs > be {
			return nil, utils.ERR_ENCRYPT_FAIL
		}
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}
	return encrypted, nil
}
func DecryptAes(data, key string) (string, *utils.Error) {
	encrypted, _ := base64.StdEncoding.DecodeString(data)
	res, err := DecryptAes2(encrypted, key)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
func DecryptAes2(data []byte, key string) ([]byte, *utils.Error) {
	cipher, _ := aes.NewCipher([]byte(key))
	decrypted := make([]byte, len(data))
	for bs, be := 0, cipher.BlockSize(); bs < len(data); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		if bs > be {
			return nil, utils.ERR_ENCRYPT_FAIL
		}
		if be > len(data) {
			return nil, utils.ERR_ENCRYPT_FAIL
		}
		cipher.Decrypt(decrypted[bs:be], data[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}
	if trim < 0 || trim > len(decrypted) {
		return nil, utils.ERR_ENCRYPT_FAIL
	}
	//再这里引起过解密失败,所以注释掉
	//if !utf8.Valid(decrypted[:trim]) {
	//	return nil, utils.ERR_DECRYPT_FAIL
	//}
	return decrypted[:trim], nil
}
