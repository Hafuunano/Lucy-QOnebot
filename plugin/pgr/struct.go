package pgr

/*
var userinfoB19 PhigrosStruct

	type PhigrosStruct struct {
		Status  bool `json:"status"`
		Content struct {
			RankingScore      float64 `json:"RankingScore"`
			ChallengeModeRank int     `json:"ChallengeModeRank"`
			PlayerID          string  `json:"PlayerID"`
			BestList          struct {
				Phi  bool `json:"phi"`
				Best []struct {
					Songid   string  `json:"songid"`
					Level    string  `json:"level"`
					Songname string  `json:"songname"`
					Rating   float64 `json:"rating"`
					Rks      float64 `json:"rks"`
					Score    int     `json:"score"`
					Acc      float64 `json:"acc"`
					Isfc     bool    `json:"isfc"`
				} `json:"best"`
			} `json:"best_list"`
			BestSonginfo []struct {
				Songid      string `json:"songid"`
				Level       string `json:"level"`
				ChartDetail struct {
					Rating     float64 `json:"rating"`
					Charter    string  `json:"charter"`
					NumOfNotes int     `json:"numOfNotes"`
				} `json:"chartDetail"`
				ChapterCode string `json:"chapterCode"`
				UnlockType  int    `json:"unlockType"`
				SongsName   string `json:"songsName"`
				Composer    string `json:"composer"`
				Illustrator string `json:"illustrator"`
			} `json:"best_songinfo"`
		} `json:"content"`
	}

	func GetRequestInfo(userID int64) PhigrosStruct {
		getPhigrosLink := os.Getenv("pualink")
		getPhigrosKey := os.Getenv("puakey")
		userData := GetUserInfoFromDatabase(userID)
		if userData != nil {
			return PhigrosStruct{}
		}
		getData, err := DrawRequestPhigros(getPhigrosLink+"/user/best19?SessionToken="+userData.PhiSession+"withsonginfo=true", getPhigrosKey, "POST")
		if err != nil {
			panic(err)
		}
		_ = json.Unmarshal(getData, &userinfoB19)
		return userinfoB19

}

// DrawRequestPhigros 发送请求结构体

	func DrawRequestPhigros(workurl string, token string, method string) (reply []byte, err error) {
		replyByte, err := http.NewRequest(method, workurl, nil)
		replyByte.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
		replyByte.Header.Set("Authorization", "Bearer "+token)
		if err != nil {
			return nil, err
		}
		resp, err := http.DefaultClient.Do(replyByte)
		if err != nil {
			panic(err)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(resp.Body)
		replyBack, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		return replyBack, err
	}
*/
type saveURLData struct {
	SaveUrl       string  `json:"saveUrl"`
	SaveVersion   int     `json:"存档版本"`
	HardCoreScore int     `json:"课题分"`
	RKS           float64 `json:"RKS"`
	GameVersion   int     `json:"游戏版本"`
	Avatar        string  `json:"头像"`
	EZ            []int   `json:"EZ"`
	HD            []int   `json:"HD"`
	IN            []int   `json:"IN"`
	AT            []int   `json:"AT"`
}

type b19info []struct {
	SongId string  `json:"songId"`
	Level  string  `json:"level"`
	Score  int     `json:"score"`
	Acc    float64 `json:"acc"`
	Fc     bool    `json:"fc"`
	Rating float64 `json:"定数"`
	Rks    float64 `json:"单曲rks"`
}
