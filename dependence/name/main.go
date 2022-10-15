package name // 获取名字

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// 参考了GO-ATRI计划 https://github.com/Kyomotoi/go-ATRI
func StringInArray(aim string, list []string) bool {
	for _, i := range list {
		if i == aim {
			return true
		}
	}
	return false
}

func StoreUserNickname(userID string, nickname string) error {
	var userNicknameData map[string]interface{}
	filePath := "file:///root/Lucy_Project/workon/main/data/zbp/users.json"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			_ = ioutil.WriteFile(filePath, []byte("{}"), 0777)
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
	_ = ioutil.WriteFile(filePath, newData, 0777)
	return nil
}

func LoadUserNickname(userID string) string {
	var d map[string]string
	filePath := "file:///root/Lucy_Project/workon/main/data/zbp/users.json"
	data, err := ioutil.ReadFile(filePath)
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
