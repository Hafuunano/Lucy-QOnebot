package pgr

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"os"
)

func DecHashToRaw(raw string) string {
	formatData := CBCDecrypt(raw, os.Getenv("hashkey"))
	return formatData
}

func CBCDecrypt(ciphertext string, key string) string {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return ""
	}
	ciphercode, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return ""
	}
	iv := ciphercode[:aes.BlockSize]
	ciphercode = ciphercode[aes.BlockSize:]
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphercode, ciphercode)
	plaintext := string(ciphercode)
	return plaintext[:len(plaintext)-int(plaintext[len(plaintext)-1])]
}
