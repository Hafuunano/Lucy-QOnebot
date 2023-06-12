package pgr

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"strconv"
	"time"
)

func HashResolveSessionByUid(userID int64, md5 string, hash string) string {
	// hash resolver
	// get ThisDay
	userIDToStr := strconv.FormatInt(userID, 10)
	keyFormat := md5 + "+" + time.Now().Weekday().String() + "+" + userIDToStr
	deckey := generateMD5(keyFormat)
	getRawData, err := decryptAES(hash, deckey)
	if err != nil {
		panic(err)
	}
	return getRawData
}

func decryptAES(encryptedData string, key string) (string, error) {
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	iv := make([]byte, aes.BlockSize)
	stream := cipher.NewCBCDecrypter(block, iv)

	decryptedBytes := make([]byte, len(encryptedBytes))
	stream.CryptBlocks(decryptedBytes, encryptedBytes)

	// 去除填充字节
	padding := int(decryptedBytes[len(decryptedBytes)-1])
	decryptedBytes = decryptedBytes[:len(decryptedBytes)-padding]

	return string(decryptedBytes), nil
}

func generateMD5(input string) string {
	hash := md5.New()
	hash.Write([]byte(input))
	hashBytes := hash.Sum(nil)

	// 将哈希字节转换为十六进制表示
	hashString := hex.EncodeToString(hashBytes)

	return hashString
}
