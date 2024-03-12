package cryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

func AesEncrypt(plainText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	plainText = pkcs7Padding(plainText, blockSize)

	cipherText := make([]byte, len(plainText))

	mode := cipher.NewCBCEncrypter(block, key[:blockSize])
	mode.CryptBlocks(cipherText, plainText)

	return cipherText, nil
}

func pkcs7Padding(plainText []byte, blockSize int) []byte {
	paddingLength := blockSize - len(plainText)%blockSize
	paddingText := bytes.Repeat([]byte{byte(paddingLength)}, paddingLength)
	return append(plainText, paddingText...)
}

func AesDecrypt(cipherText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()

	mode := cipher.NewCBCDecrypter(block, key[:blockSize])
	mode.CryptBlocks(cipherText, cipherText)

	cipherText = pkcs7Unpadding(cipherText)

	return cipherText, nil
}

func pkcs7Unpadding(cipherText []byte) []byte {
	length := len(cipherText)
	unpaddingLength := int(cipherText[length-1])
	return cipherText[:(length - unpaddingLength)]
}

func Base64Encode(cipherText []byte) string {
	return base64.StdEncoding.EncodeToString(cipherText)
}

func Base64Decode(cipherText string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(cipherText)
}
