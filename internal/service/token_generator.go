package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

const (
	defaultKey = "0123456789abcdef" // 16字节默认密钥
)

// EncryptString AES-GCM 加密 (推荐使用)
func EncryptString(plaintext string, key ...[]byte) (string, error) {
	usedKey := getKey(key...)

	block, err := aes.NewCipher(usedKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// DecryptString AES-GCM 解密
func DecryptString(ciphertext string, key ...[]byte) (string, error) {
	usedKey := getKey(key...)

	block, err := aes.NewCipher(usedKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 修复类型声明，保持为字节切片操作
	data, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	if len(data) < gcm.NonceSize() {
		return "", errors.New("invalid ciphertext")
	}

	nonce, ciphertextBytes := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)

	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// 统一密钥处理逻辑
func getKey(key ...[]byte) []byte {
	if len(key) > 0 && len(key[0]) >= 16 { // 允许16/24/32字节密钥
		return key[0]
	}
	return []byte(defaultKey)
}
