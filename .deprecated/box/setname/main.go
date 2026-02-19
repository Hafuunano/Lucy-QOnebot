// Package setname 获取名字
package setname

import (
	"encoding/json"
	"os"

	"github.com/FloatTech/floatbox/file"
	"github.com/bytedance/sonic"
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
	filePath := file.BOTPATH + "/data/zbp/users.json"
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			_ = os.WriteFile(filePath, []byte("{}"), 0777)
		} else {
			return err
		}
	}
	_ = sonic.Unmarshal(data, &userNicknameData)
	userNicknameData[userID] = nickname // setdata.
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
	filePath := file.BOTPATH + "/data/zbp/users.json"
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "你"
	}
	err = sonic.Unmarshal(data, &d)
	if err != nil {
		return "你"
	}
	result := d[userID]
	if result == "" {
		result = "你"
	}
	return result
}
