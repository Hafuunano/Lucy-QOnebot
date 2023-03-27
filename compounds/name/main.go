// Package name 获取名字
package name

import (
	"encoding/json"
	"os"
)

// StringInArray 检查列表是否有关键词 https://github.com/Kyomotoi/go-ATRI
func StringInArray(aim string, list []string) bool {
	for _, i := range list {
		if i == aim {
			return true
		}
	}
	return false
}

// StoreUserNickname Store names in jsons
func StoreUserNickname(userID string, nickname string) error {
	var userNicknameData map[string]interface{}
	filePath := "./data/zbp/users.json"
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			_ = os.WriteFile(filePath, []byte("{}"), 0777)
		} else {
			return err
		}
	}
	err = json.Unmarshal(data, &userNicknameData)
	if err != nil {
		return err
	}
	userNicknameData[userID] = nickname
	newData, err := json.Marshal(userNicknameData)
	if err != nil {
		return err
	}
	_ = os.WriteFile(filePath, newData, 0777)
	return nil
}

// LoadUserNickname Load UserNames to work it well.
func LoadUserNickname(userID string) string {
	var d map[string]string
	// read main files
	filePath := "./data/zbp/users.json"
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "你"
	}
	err = json.Unmarshal(data, &d)
	if err != nil {
		return "你"
	}
	result := d[userID]
	if result == "" {
		result = "你"
	}
	return result
}
