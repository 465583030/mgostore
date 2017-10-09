package mgostore

import (
	"testing"

	"github.com/gsingharoy/mgostore/lib"
	"github.com/stretchr/testify/assert"
)

func Test_encryptFields(t *testing.T) {

	m := &mockModel{
		PlainTextField:  "plain text",
		NumField:        42,
		EncryptedField2: "encrypted text"}
	encryptFields(m)

	assert.Equal(t,
		"plain text",
		m.PlainTextField,
		"Expected plain text field to not be encrypted")
	assert.Equal(t,
		42,
		m.NumField,
		"Expected non string field to not change event though they have encrypt tag")
	assert.Equal(t,
		"encrypted text",
		m.EncryptedField2,
		"Expected invalid type of encrypt to be not encrypted")
	assert.Equal(t,
		0,
		len(m.EncryptedField1),
		"Expect empty field with valid encrypt tag not to be encrypted")

	m.EncryptedField1 = "now encrypt!"
	encryptFields(m)
	key := []byte(testEncryptionSecret)
	decryptedText, _ := lib.AesDecrypt(key, m.EncryptedField1)
	assert.Equal(t, "now encrypt!", decryptedText, "Expected decryption of encrypted text to match")
}

func Test_decryptFields(t *testing.T) {
	m := &mockModel{
		PlainTextField:  "plain text",
		NumField:        42,
		EncryptedField2: "encrypted text"}
	decryptFields(m)
	assert.Equal(t,
		"plain text",
		m.PlainTextField,
		"Expected plain text field to not be decrypted")
	assert.Equal(t,
		42,
		m.NumField,
		"Expected non string field to not change event though they have encrypt tag")
	assert.Equal(t,
		"encrypted text",
		m.EncryptedField2,
		"Expected invalid type of encrypt to be not decrypted")
	assert.Equal(t,
		0,
		len(m.EncryptedField1),
		"Expect empty field with valid encrypt tag to not be decrypted")

	key := []byte(testEncryptionSecret)
	encryptedText, _ := lib.AesEncrypt(key, "encrypt this!")
	m.EncryptedField1 = encryptedText
	decryptFields(m)

	assert.Equal(t, "encrypt this!", m.EncryptedField1, "Expected decryption of encrypted field to match")
}
