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
	iv := ciphercode[:aes.BlockSize]        // 密文的前 16 个字节为 iv
	ciphercode = ciphercode[aes.BlockSize:] // 正式密文
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphercode, ciphercode)
	plaintext := string(ciphercode) // ↓ 减去 padding
	return plaintext[:len(plaintext)-int(plaintext[len(plaintext)-1])]
}
