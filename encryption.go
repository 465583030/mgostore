package mgostore

import (
	"reflect"

	"github.com/gsingharoy/mgostore/lib"
)

func encryptFields(m Model) error {
	s := reflect.ValueOf(m).Elem()
	t := reflect.TypeOf(m).Elem()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		vField := s.Field(i)
		if f.Type == reflect.TypeOf("") && len(vField.String()) > 0 {
			// encrypted with aes algorithm
			if f.Tag.Get("encrypt") == "aes" {
				cryptoConfig := m.DBConfig().CryptoConfig
				if cryptoConfig == nil {
					return ErrMissingCryptoSecret
				}
				encryptedvalue, err := lib.AesEncrypt(cryptoConfig.AESSecret, vField.String())
				if err != nil {
					return err
				}
				vField.SetString(encryptedvalue)
			}
		}
	}
	return nil
}

func decryptFields(m Model) error {
	s := reflect.ValueOf(m).Elem()
	t := reflect.TypeOf(m).Elem()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		vField := s.Field(i)
		if f.Type == reflect.TypeOf("") && len(vField.String()) > 0 {
			// encrypted with aes algorithm
			if f.Tag.Get("encrypt") == "aes" {
				cryptoConfig := m.DBConfig().CryptoConfig
				if cryptoConfig == nil {
					return ErrMissingCryptoSecret
				}
				decryptedValue, err := lib.AesDecrypt(cryptoConfig.AESSecret, vField.String())
				if err != nil {
					return err
				}
				vField.SetString(decryptedValue)
			}
		}
	}
	return nil
}
