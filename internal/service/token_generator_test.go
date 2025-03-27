package service_test

import (
	"testing"

	"gochat/internal/service"

	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	testCases := []struct {
		name      string
		plaintext string
		key       [][]byte
		wantErr   bool
	}{
		{"默认密钥正常加解密", "test123", nil, false},
		{"空字符串加解密", "", nil, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 加密测试
			ciphertext, err := service.EncryptString(tc.plaintext, tc.key...)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, ciphertext)

			// 解密测试
			decrypted, err := service.DecryptString(ciphertext, tc.key...)
			assert.NoError(t, err)
			assert.Equal(t, tc.plaintext, decrypted)
		})
	}

	// 额外错误场景测试
	t.Run("错误密钥解密", func(t *testing.T) {
		ciphertext, _ := service.EncryptString("secret", []byte("goodkey123456789"))

		_, err := service.DecryptString(ciphertext, []byte("wrongkey12345678"))
		assert.Error(t, err)
	})

	t.Run("无效密文格式", func(t *testing.T) {
		_, err := service.DecryptString("invalid_base64!@#")
		assert.Error(t, err)
	})
}
