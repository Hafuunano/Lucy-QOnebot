package mai

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/FloatTech/floatbox/web"
	"image"
	"image/png"
	"io"
	"net/http"
)

type DivingFishB50 struct {
	QQ       int    `json:"qq"`
	Username string `json:"username"`
	B50      bool   `json:"b50"`
}

func QueryMaiBotDataFromQQ(qq int) (playerdata []byte, err error) {
	// packed json and sent.
	jsonStruct := DivingFishB50{QQ: qq, B50: true}
	jsonStructData, err := json.Marshal(jsonStruct)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "https://www.diving-fish.com/api/maimaidxprober/query/player", bytes.NewBuffer(jsonStructData))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 400 {
		return nil, errors.New("400")
	}
	if resp.StatusCode == 403 {
		return nil, errors.New("403")
	}
	playerData, err := io.ReadAll(resp.Body)
	return playerData, err
}
func QueryMaiBotDataFromUserName(username string) (playerdata []byte, err error) {
	// packed json and sent.
	jsonStruct := DivingFishB50{Username: username, B50: true}
	jsonStructData, err := json.Marshal(jsonStruct)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "https://www.diving-fish.com/api/maimaidxprober/query/player", bytes.NewBuffer(jsonStructData))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 400 {
		return nil, errors.New("400")
	}
	if resp.StatusCode == 403 {
		return nil, errors.New("403")
	}
	playerDataByte, err := io.ReadAll(resp.Body)
	return playerDataByte, err
}

// https://www.diving-fish.com/api/chunithmprober/query/player

func QueryChunDataFromQQ(qq int) (playerdata []byte, err error) {
	// packed json and sent.
	jsonStruct := DivingFishB50{QQ: qq, B50: true}
	jsonStructData, err := json.Marshal(jsonStruct)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "https://www.diving-fish.com/api/chunithmprober/query/player", bytes.NewBuffer(jsonStructData))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 400 {
		return nil, errors.New("400")
	}
	if resp.StatusCode == 403 {
		return nil, errors.New("403")
	}
	playerData, err := io.ReadAll(resp.Body)
	return playerData, err
}

func QueryChunDataFromUsername(username string) (playerdata []byte, err error) {
	// packed json and sent.
	jsonStruct := DivingFishB50{Username: username, B50: true}
	jsonStructData, err := json.Marshal(jsonStruct)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "https://www.diving-fish.com/api/chunithmprober/query/player", bytes.NewBuffer(jsonStructData))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 400 {
		return nil, errors.New("400")
	}
	if resp.StatusCode == 403 {
		return nil, errors.New("403")
	}
	playerData, err := io.ReadAll(resp.Body)
	return playerData, err
}

func GetCoverByMusicID(id string) (image image.Image) {
	data, err := web.GetData("https://www.diving-fish.com/covers/" + id + ".png")
	if err != nil {
		return nil
	}
	imageReader := bytes.NewReader(data)
	image, err = png.Decode(imageReader)
	if err != nil {
		return nil
	}
	return image
}
