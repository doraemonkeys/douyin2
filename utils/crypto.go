package utils

import (
	"errors"
	"strconv"

	"github.com/wumansgy/goEncrypt/aes"
)

type CbcAESCrypt struct {
	// contains filtered or unexported fields
	secretKey []byte
}

// NewAESCrypt 创建AES加密器, HexSecretKey为16进制字符串
// CBC模式，PKCS7填充
func NewAESCrypt(HexSecretKey string) (*CbcAESCrypt, error) {
	secretKey, err := HexStrToBytes(HexSecretKey)
	return &CbcAESCrypt{secretKey: secretKey}, err
}

// 加密为16进制字符串
func (a *CbcAESCrypt) Encrypt(plainText string) (string, error) {
	if plainText == "" {
		return "", errors.New("plainText is empty")
	}
	return aes.AesCbcEncryptHex([]byte(plainText), a.secretKey, nil)
}

func (a *CbcAESCrypt) Decrypt(cipherTextHex string) (string, error) {
	if cipherTextHex == "" {
		return "", errors.New("cipherTextHex is empty")
	}
	plaintext, err := aes.AesCbcDecryptByHex(cipherTextHex, a.secretKey, nil)
	return string(plaintext), err
}

// HexStrToBytes convert hex string to bytes
func HexStrToBytes(hexStr string) ([]byte, error) {
	if len(hexStr)%2 != 0 {
		return nil, errors.New("invalid hex string")
	}
	bytes := make([]byte, len(hexStr)/2)
	for i := 0; i < len(hexStr); i += 2 {
		b, err := strconv.ParseUint(hexStr[i:i+2], 16, 8)
		if err != nil {
			return nil, err
		}
		bytes[i/2] = byte(b)
	}
	return bytes, nil
}
