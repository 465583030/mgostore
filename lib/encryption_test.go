package lib

import (
	"crypto/aes"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testAesKey = "0AD491CF3BC7461F20C9F30DC7CB4A80"
	testAesIv  = "4966CF23AAEF695B28115ED3F2FA9F62"
)

func TestAesEncrypt(t *testing.T) {
	key := []byte("INVALID AES KEY")
	t.Log("When an invalid AES key is sent")
	encryptedText, err := AesEncrypt(key, "Sample text")
	assert.Equal(t, encryptedText, "", "Encrypted text is not empty")
	assert.Equal(t, err, aes.KeySizeError(15), "Expected invalid key size error")

	t.Log("When the AES key is a valid one")
	key = []byte(testAesKey)
	encryptedText, err = AesEncrypt(key, "Sample text")
	assert.Equal(t, err, nil, "Expected error to be nil")
	assert.NotEqual(t, encryptedText, "", "Text is encrypted successfully")
}

func TestAesDecrypt(t *testing.T) {
	key := []byte("INVALID AES KEY")
	t.Log("When an invalid AES key is sent")
	decryptedText, err := AesDecrypt(key, "Sample text")
	assert.Equal(t, err, aes.KeySizeError(15), "Expected invalid key size error")
	assert.Equal(t, decryptedText, "", "Decrypted text is empty")

	t.Log("When the AES key is a valid one")
	key = []byte(testAesKey)
	encryptedText, _ := AesEncrypt(key, "Sample text")
	decryptedText, err = AesDecrypt(key, encryptedText)
	assert.Equal(t, err, nil, "Expected error to be nil")
	assert.Equal(t, "Sample text", decryptedText, "Decryption does not match original plain text")
}
