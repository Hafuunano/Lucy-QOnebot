package name // Package name 获取名字

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
			panic(err)
			return err
		}
	}
	err = json.Unmarshal(data, &userNicknameData)
	if err != nil {
		panic(err)
		return err
	}
	userNicknameData[userID] = nickname
	newData, err := json.Marshal(userNicknameData)
	if err != nil {
		panic(err)
		return err
	}
	_ = os.WriteFile(filePath, newData, 0777)
	return nil
}

// LoadUserNickname Load UserNames(if had.)
func LoadUserNickname(userID string) string {
	var d map[string]string
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
