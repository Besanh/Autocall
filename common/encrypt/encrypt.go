package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"math/rand"
)

var IV string = "AJf3QItKM7+Lkh/BZT2xNg=="

func AESEncrypt(plainText, key string) (string, error) {
	iv, _ := base64.StdEncoding.DecodeString(IV)
	passPhrase := make([]byte, 32)
	copy(passPhrase, []byte(key))
	block, err := aes.NewCipher(passPhrase)
	if err != nil {
		return "", err
	}
	ecb := cipher.NewCBCEncrypter(block, iv)
	content := []byte(plainText)
	content = PKCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	return base64.StdEncoding.EncodeToString(crypted), err
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func AESDecrypt(crypt []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	if len(crypt) == 0 {
		return "", err
	}
	iv, _ := base64.StdEncoding.DecodeString(IV)
	ecb := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(crypt))
	ecb.CryptBlocks(decrypted, crypt)
	return base64.StdEncoding.EncodeToString(PKCS5Trimming(decrypted)), err
}

func PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

var letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789~!@#%^&*./?"

func getRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func GenerateEncrypted(plainText string) (string, error) {
	salt := getRandomString(8)
	encrypted, err := AESEncrypt(plainText, salt)
	if err != nil {
		return "", err
	}
	return salt + "$" + encrypted, nil
}
