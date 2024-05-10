package mai

import (
	"bytes"
	"encoding/json"
	"fmt"
	Stringbreaker "github.com/MoYoez/Lucy-QOnebot/box/break"
	"image"
	rand2 "math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/web"
	"github.com/FloatTech/gg"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/img/text"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	engine = control.Register("maidx", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "Hi NekoPachi!\n",
		PrivateDataFolder: "maidx",
	})
)

type MaiSongData []struct {
	Id     string    `json:"id"`
	Title  string    `json:"title"`
	Type   string    `json:"type"`
	Ds     []float64 `json:"ds"`
	Level  []string  `json:"level"`
	Cids   []int     `json:"cids"`
	Charts []struct {
		Notes   []int  `json:"notes"`
		Charter string `json:"charter"`
	} `json:"charts"`
	BasicInfo struct {
		Title       string `json:"title"`
		Artist      string `json:"artist"`
		Genre       string `json:"genre"`
		Bpm         int    `json:"bpm"`
		ReleaseDate string `json:"release_date"`
		From        string `json:"from"`
		IsNew       bool   `json:"is_new"`
	} `json:"basic_info"`
}

func init() {
	engine.OnRegex(`^[！! /]chun$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		uid := ctx.Event.UserID
		dataPlayer, err := QueryChunDataFromQQ(int(uid))
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
			return
		}
		txt := HandleChunDataByUsingText(dataPlayer)
		base64Font, err := text.RenderToBase64(txt, text.BoldFontFile, 1920, 45)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
			return
		}
		ctx.SendChain(message.Image("base64://" + binary.BytesToString(base64Font)))
	})
	engine.OnRegex(`^[! ！/](mai|b50)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		uid := ctx.Event.UserID
		if GetUserSwitcherInfoFromDatabase(uid) {
			// use lxns checker service.
			getUserData := RequestBasicDataFromLxns(uid)
			if getUserData.Code != 200 {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("aw 出现了一点小错误~：\n - 请检查你是否有上传过数据并且绑定了QQ号\n - 请检查你的设置是否允许了第三方查看"))
				return
			}
			getGameUserData := RequestB50DataByFriendCode(getUserData.Data.FriendCode)
			if getGameUserData.Code != 200 {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("aw 出现了一点小错误~：\n - 请检查你是否有上传过数据并且绑定了QQ号\n - 请检查你的设置是否允许了第三方查看"))
				return
			}
			getImager, _ := ReFullPageRender(getGameUserData, getUserData, ctx)
			_ = gg.NewContextForImage(getImager).SavePNG(engine.DataFolder() + "save/" + "LXNS_" + strconv.Itoa(int(ctx.Event.UserID)) + ".png")
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Render User B50 : "+getUserData.Data.Name), message.Image(Saved+"LXNS_"+strconv.Itoa(int(ctx.Event.UserID))+".png"))
		} else {
			dataPlayer, err := QueryMaiBotDataFromQQ(int(uid), true)
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
				return
			}
			var data player
			_ = json.Unmarshal(dataPlayer, &data)
			renderImg, plateStat := FullPageRender(data, ctx)
			tipPlate := ""
			getRand := rand2.Intn(10)
			if getRand == 8 {
				if !plateStat {
					tipPlate = "tips: 可以使用 ！mai plate xxx 来绑定称号~\n"
				}
			}
			_ = gg.NewContextForImage(renderImg).SavePNG(engine.DataFolder() + "save/" + strconv.Itoa(int(ctx.Event.UserID)) + ".png")
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Render User B50 : "+data.Username+"\n"+tipPlate), message.Image(Saved+strconv.Itoa(int(ctx.Event.UserID))+".png"))

		}
	})
	engine.OnRegex(`^[! ！/](b40)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		uid := ctx.Event.UserID
		dataPlayer, err := QueryMaiBotDataFromQQ(int(uid), false)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
			return
		}
		var data player
		_ = json.Unmarshal(dataPlayer, &data)
		renderImg, plateStat := FullPageRender(data, ctx)
		tipPlate := ""
		getRand := rand2.Intn(10)
		if getRand == 8 {
			if !plateStat {
				tipPlate = "tips: 可以使用 ！mai plate xxx 来绑定称号~\n"
			}
		}
		_ = gg.NewContextForImage(renderImg).SavePNG(engine.DataFolder() + "save/b40_" + strconv.Itoa(int(ctx.Event.UserID)) + ".png")
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Render User B40 : "+data.Username+"\n"+tipPlate), message.Image(Saved+"b40_"+strconv.Itoa(int(ctx.Event.UserID))+".png"))
	})
	engine.OnRegex(`^[! ！/](mai|b50)\sswitch$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getBool := GetUserSwitcherInfoFromDatabase(ctx.Event.UserID)
		err := FormatUserSwitcher(ctx.Event.UserID, !getBool).ChangeUserSwitchInfoFromDataBase()
		if err != nil {
			panic(err)
		}
		var getEventText string
		// due to it changed, so reverse.
		if !getBool {
			getEventText = "Lxns查分"
		} else {
			getEventText = "Diving Fish查分"
		}
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经修改为"+getEventText))
	})
	engine.OnRegex(`^[! ！/](mai|b50)\splate\s(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getPlateInfo := ctx.State["regex_matched"].([]string)[2]
		_ = FormatUserDataBase(ctx.Event.UserID, getPlateInfo, GetUserDefaultinfoFromDatabase(ctx.Event.UserID)).BindUserDataBase()
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经将称号绑定上去了哦w"))
	})
	engine.OnRegex(`^[! ！/](mai|b50)\supload`, PictureHandler).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getPic := ctx.State["image_url"].([]string)[0]
		imageData, err := web.GetData(getPic)
		if err != nil {
			return
		}
		getRaw, _, err := image.Decode(bytes.NewReader(imageData))
		if err != nil {
			panic(err)
		}
		// pic Handler
		getRenderPlatePicRaw := gg.NewContext(1260, 210)
		getRenderPlatePicRaw.DrawRoundedRectangle(0, 0, 1260, 210, 10)
		getRenderPlatePicRaw.Clip()
		getHeight := getRaw.Bounds().Dy()
		getLength := getRaw.Bounds().Dx()
		var getHeightHandler, getLengthHandler int
		switch {
		case getHeight < 210 && getLength < 1260:
			getRaw = Resize(getRaw, 1260, 210)
			getHeightHandler = 0
			getLengthHandler = 0
		case getHeight < 210:
			getRaw = Resize(getRaw, getLength, 210)
			getHeightHandler = 0
			getLengthHandler = (getRaw.Bounds().Dx() - 1260) / 3 * -1
		case getLength < 1260:
			getRaw = Resize(getRaw, 1260, getHeight)
			getHeightHandler = (getRaw.Bounds().Dy() - 210) / 3 * -1
			getLengthHandler = 0
		default:
			getLengthHandler = (getRaw.Bounds().Dx() - 1260) / 3 * -1
			getHeightHandler = (getRaw.Bounds().Dy() - 210) / 3 * -1
		}
		getRenderPlatePicRaw.DrawImage(getRaw, getLengthHandler, getHeightHandler)
		getRenderPlatePicRaw.Fill()
		_ = getRenderPlatePicRaw.SavePNG(userPlate + strconv.Itoa(int(ctx.Event.UserID)) + ".png")
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经存入了哦w"))
	})
	engine.OnRegex(`^[! ！/](mai|b50)\sremove`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		_ = os.Remove(userPlate + strconv.Itoa(int(ctx.Event.UserID)) + ".png")
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经删掉了哦w"))
	})
	engine.OnRegex(`^[! ！/](mai|b50)\sdefault\splate\s(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getDefaultInfo := ctx.State["regex_matched"].([]string)[2]
		if getDefaultInfo == "" {
			_ = FormatUserDataBase(ctx.Event.UserID, GetUserInfoFromDatabase(ctx.Event.UserID), getDefaultInfo).BindUserDataBase()
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经恢复了正常~"))
			return
		}
		_, err := GetDefaultPlate(getDefaultInfo)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("设定的预设不正确"))
			return
		}
		_ = FormatUserDataBase(ctx.Event.UserID, GetUserInfoFromDatabase(ctx.Event.UserID), getDefaultInfo).BindUserDataBase()
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经设定好了哦~"))

	})
	engine.OnRegex(`^[! ！/](mai|b50)\sbind\s(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getDefaultInfo := ctx.State["regex_matched"].([]string)[2]
		indexReply := DecHashToRaw(getDefaultInfo)
		// get session.
		if indexReply == "" {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("请前往 https://mai.lemonkoi.one 获取绑定码进行绑定"))
			return
		}
		getQID, getSessionID := RawJsonParse(indexReply)
		if getQID != ctx.Event.UserID {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("请求Hash中QQ号不一致，请使用自己的号重新申请"))
			return
		}
		// check id
		getID := GetWahlapUserID(getSessionID)
		if getID == -1 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ID 无效或者是过期 ，请使用新的ID或者再次尝试"))
			return
		}
		// login.
		err := FormatUserIDDatabase(ctx.Event.UserID, strconv.Itoa(int(getID))).BindUserIDDataBase()
		if err != nil {
			panic(err)
		}
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("绑定成功~"))
	})
	engine.OnRegex(`^[! ！/](mai|b50)\sunbind$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("取消绑定成功~"))
		RemoveUserIdFromDatabase(ctx.Event.UserID)
	})
	engine.OnRegex(`^[! ！/](mai|b50)\sunlock$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getMaiID := GetUserIDFromDatabase(ctx.Event.UserID)
		if getMaiID.Userid == "" {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("请前往 https://mai.lemonkoi.one 获取绑定码进行绑定"))
			return
		}
		getCodeStat, err := web.RequestDataWithHeaders(http.DefaultClient, "https://maihook.lemonkoi.one/api/idunlocker?userid="+getMaiID.Userid, "GET", func(request *http.Request) error {
			request.Header.Set("valid", os.Getenv("validauth"))
			return nil
		}, nil)
		if err != nil {
			return
		}
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(string(getCodeStat)))
	})
	engine.OnRegex(`^[! ！/](mai|b50)\stokenbind\s(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getDefaultInfo := ctx.State["regex_matched"].([]string)[2]
		FormatUserToken(strconv.FormatInt(ctx.Event.UserID, 10), getDefaultInfo).BindUserToken()
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("绑定成功~"))
	})
	engine.OnRegex(`^[! ！/](mai|b50)\supdate$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getID := ctx.Event.UserID
		getMaiID := GetUserIDFromDatabase(getID)
		if getMaiID.Userid == "" {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("请前往 https://mai.lemonkoi.one 获取绑定码进行绑定"))
			return
		}
		getTokenId := GetUserToken(strconv.FormatInt(getID, 10))
		if getTokenId == "" {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("请使用 /mai tokenbind <tokenid> 绑定水鱼查分器，其中 TokenID 从 https://www.diving-fish.com/maimaidx/prober 用户设置中拿到, 格式：' /mai tokenbind xxx '，不需要 < > "))
			return
		}
		checkertoken := CheckTheTicketIsValid(getTokenId)
		if checkertoken == "" {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("无法校验 Token 可用性, 可能 Lucy 连接到 水鱼的服务已宕机"))
			return
		}
		if checkertoken == "non-exist" {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("此 Token 不合法 ，请重新绑定"))
			return
		}
		getData, err := web.RequestDataWithHeaders(http.DefaultClient, "https://maihook.lemonkoi.one/api/getMusicList?userid="+getMaiID.Userid+"&index=0", "GET", func(request *http.Request) error {
			request.Header.Set("valid", os.Getenv("validauth"))
			return nil
		}, nil)
		if err != nil {
			panic(err)
		}
		// update by path.
		var unmashellData UserMusicListStruct
		json.Unmarshal(getData, &unmashellData)
		resp := UpdateHandler(unmashellData, getTokenId)
		if unmashellData.NextIndex != 0 {
			for i := unmashellData.NextIndex; i > 0; {
				var unmashellDataS UserMusicListStruct
				iStr := strconv.Itoa(i)
				getDataS, err := web.RequestDataWithHeaders(http.DefaultClient, "https://maihook.lemonkoi.one/api/getMusicList?userid="+getMaiID.Userid+"&index="+iStr, "GET", func(request *http.Request) error {
					request.Header.Set("valid", os.Getenv("validauth"))
					return nil
				}, nil)
				if err != nil {
					panic(err)
				}
				json.Unmarshal(getDataS, &unmashellDataS)
				UpdateHandler(unmashellDataS, getTokenId)
				i = unmashellDataS.NextIndex
			}
		}
		if resp == 200 || resp == 500 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("上传成功了哦~ QWQ"))
		} else {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("貌似出现了一点点问题, 要不重新试试(?"))
		}

	})
	engine.OnRegex(`^[! ！/](mai|b50)\sregion$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getID := ctx.Event.UserID
		getMaiID := GetUserIDFromDatabase(getID)
		if getMaiID.Userid == "" {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("请前往 https://mai.lemonkoi.one 获取绑定码进行绑定"))
			return
		}
		getReplyMsg, _ := web.RequestDataWithHeaders(http.DefaultClient, "https://maihook.lemonkoi.one/api/getRegion?userid="+getMaiID.Userid, "GET", func(request *http.Request) error {
			request.Header.Set("valid", os.Getenv("validauth"))
			return nil
		}, nil)
		if !strings.Contains(string(getReplyMsg), "{") {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("返回了错误.png, ERROR:"+string(getReplyMsg)))
			return
		}
		var MixedMagic GetUserRegionStruct
		json.Unmarshal(getReplyMsg, &MixedMagic)
		var returnText string
		for _, onlistLoader := range MixedMagic.UserRegionList {
			returnText = returnText + MixedRegionWriter(onlistLoader.RegionId-1, onlistLoader.PlayCount, onlistLoader.Created) + "\n\n"
		}
		if returnText == "" {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("目前 Lucy 没有查到您的游玩记录哦~"))
			return
		}
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("目前查询到您的游玩记录如下: \n\n"+returnText))
	})
	engine.OnRegex(`^[! ！/](mai|b50)\sautosync`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		// AutoUpdate : ReMaimai Project
		// Check UserStatus
		getQuery := QueryUserTokenDataByQQ(ctx.Event.UserID).UserID
		if getQuery != 0 {
			RemoveUserTokenData(ctx.Event.UserID)
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经帮你关闭了自动同步了~"))
			return
		}
		getUserID := GetUserIDFromDatabase(ctx.Event.UserID)
		if getUserID.Userid == "" {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("请前往 https://mai.lemonkoi.one 获取绑定码进行绑定"))
			return
		}
		getToken := GetUserToken(strconv.FormatInt(ctx.Event.UserID, 10))
		if getToken == "" {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("请使用 /mai tokenbind <tokenid> 绑定水鱼查分器，其中 TokenID 从 https://www.diving-fish.com/maimaidx/prober 用户设置中拿到"))
			return
		}
		queryIntUserID, _ := strconv.ParseInt(getUserID.Userid, 10, 64)
		// Write
		InsertQueryParamsToStruct(UserAutoSync{QQ: ctx.Event.UserID, UserID: queryIntUserID})
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("启用自动同步成功~ 在 Lucy 参与加速项目下可自动同步当前游玩歌曲~ "))
	})
	engine.OnRegex(`^(网咋样|[! ！/](mai|b50)\sstatus$)`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		//	getWebStatus := ReturnWebStatus()
		getZlibError := ReturnZlibError()
		// 20s one request.
		var getLucyRespHandler int
		if getZlibError.Full.Field3 < 180 {
			getLucyRespHandler = getZlibError.Full.Field3
		} else {
			getLucyRespHandler = getZlibError.Full.Field3 - 180
		}
		getLucyRespHandlerStr := strconv.Itoa(getLucyRespHandler)

		getZlibWord := "Zlib 压缩跳过率: \n" + "10mins (" + ConvertZlib(getZlibError.ZlibError.Field1, getZlibError.Full.Field1) + " Loss)\n" + "30mins (" + ConvertZlib(getZlibError.ZlibError.Field2, getZlibError.Full.Field2) + " Loss)\n" + "60mins (" + ConvertZlib(getZlibError.ZlibError.Field3, getZlibError.Full.Field3) + " Loss)\n"
		PingStatus := "以下为 Ping Status, 会返回当前服务器响应时间，如果遇到无反馈可能是炸了( :"
		//getRealStatus := "\n以下数据来源于mai机台的数据反馈\n"
		// CICI SERVICE IS NOT WORK YET. // ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("* Zlib 压缩跳过率可以很好的反馈当前 MaiNet (Wahlap Service) 当前负载的情况，根据样本 + Lucy处理情况 来判断 \n* 错误率收集则来源于 机台游玩数据，反应各地区真实mai游玩错误情况 \n* 在 1小时 内，Lucy 共处理了 "+getLucyRespHandlerStr+"次 请求💫，其中详细数据如下:\n\n"+getZlibWord+getRealStatus+"\n"+ConvertRealPlayWords(playerStatus)))
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("* Zlib 压缩跳过率可以很好的反馈当前 MaiNet (Wahlap Service) 当前负载的情况，根据样本 + Lucy处理情况 来判断 \n* 错误率收集则来源于 机台游玩数据，反应各地区真实mai游玩错误情况 \n* 在 1小时 内，Lucy 共处理了 "+getLucyRespHandlerStr+"次 请求💫，其中详细数据如下:\n\n"+getZlibWord+PingStatus+"\n"+Reply(AsyncPinger())))
	})
	engine.OnRegex(`^[! ！/](mai|b50)\squery\s(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getDefaultInfo := ctx.State["regex_matched"].([]string)[2]
		// CASE: if User Trigger This command, check other settings.
		// getQuery:
		// level_index | song_type
		getLength, getSplitInfo := Stringbreaker.SplitCommandTo(getDefaultInfo, 2)
		userSettingInterface := map[string]string{}
		var settedSongAlias string
		if getLength > 1 { // prefix judge.
			settedSongAlias = getSplitInfo[1]
			for i, returnLevelValue := range []string{"绿", "黄", "红", "紫", "白"} {
				if strings.Contains(getSplitInfo[0], returnLevelValue) {
					userSettingInterface["level_index"] = strconv.Itoa(i)
					break
				}
			}
			switch {
			case strings.Contains(getSplitInfo[0], "dx"):
				userSettingInterface["song_type"] = "dx"
			case strings.Contains(getSplitInfo[0], "标"):
				userSettingInterface["song_type"] = "standard"
			}
		} else {
			// no other infos. || default setting ==> dx Master | std Master | dx expert | std expert (as the highest score)
			settedSongAlias = getSplitInfo[0]
		}
		getUserID := ctx.Event.UserID
		getBool := GetUserSwitcherInfoFromDatabase(getUserID)
		// checker id | song
		var isIDChecker bool
		var songIDList []int
		// first read the config.
		getLevelIndex := userSettingInterface["level_index"]
		getSongType := userSettingInterface["song_type"]
		var getReferIndexIsOn bool
		var accStat bool
		if getLevelIndex != "" { // use custom diff
			getReferIndexIsOn = true
		}
		switch {
		case strings.HasPrefix(settedSongAlias, "id"):
			// useID checker.
			isIDChecker = true
			getParse, err := strconv.ParseInt(strings.Replace(settedSongAlias, "id", "", 1), 10, 64)
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ID 查找参数非法"))
				return
			}
			songIDList = []int{int(getParse)}
		default:
			isIDChecker = false
			queryStatus, songIDLists, accStats, returnListHere := QueryReferSong(settedSongAlias, getBool)
			songIDList = songIDLists
			accStat = accStats
			if !queryStatus {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("未找到对应歌曲，可能是数据库未收录（"))
				return
			}
			if accStat {
				// Handler Which Song user played.
				var FullList []int
				for _, list := range returnListHere {
					for _, listInsider := range list {
						FullList = append(FullList, listInsider)
					}
				}

				FullList = removeIntDuplicates(FullList)
				// make both song Handler, check this song is from sd | DX pattern.
				var sampleListShown []int
				// sometimes list maybe contain dx | SD, but they are same song.
				if len(FullList) == 2 {
					for _, listSample := range FullList {
						sampleListShown = append(sampleListShown, simpleNumHandler(listSample, false)) //  convert to DX pattern.
					}
					sampleListShown = removeIntDuplicates(sampleListShown)
				}

				if len(sampleListShown) == 1 {
					songIDList = []int{simpleNumHandler(songIDList[0], false)}
				} else {
					// varies handler machine,means it has songs.
					// query them.
					songIDList = FullList
				}

			}
		}

		// get SongID, render.

		if getBool { // lxns service.
			getFriendID := RequestBasicDataFromLxns(getUserID)
			if getFriendID.Data.FriendCode == 0 {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("没有绑定哦～ 请查看你是否在 maimai.lxns.net 上绑定了qq并且允许通过qq查看w "))
				return
			}
			if !getReferIndexIsOn { // no refer then return the last one.
				var getReport LxnsMaimaiRequestUserReferBestSong
				switch {
				case getSongType == "standard":
					for _, songIdInt := range songIDList {
						getReport = RequestReferSong(getFriendID.Data.FriendCode, int64(songIdInt), true)
						if getReport.Code == 200 && len(getReport.Data) != 0 {
							break
						}
					}
					if len(getReport.Data) == 0 {
						ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("没有发现 SD 谱面～ 如不确定可以忽略请求参数, Lucy会自动识别"))
						return
					}
				case getSongType == "dx":
					for _, songIdInt := range songIDList {
						getReport = RequestReferSong(getFriendID.Data.FriendCode, int64(songIdInt), false)
						if getReport.Code == 200 && len(getReport.Data) != 0 {
							break
						}
					}
					if len(getReport.Data) == 0 {
						ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("没有发现 DX 谱面～ 如不确定可以忽略请求参数, Lucy会自动识别"))
						return
					}

				default:
					for _, songIdInt := range songIDList {
						getReport = RequestReferSong(getFriendID.Data.FriendCode, int64(songIdInt), false)
						fmt.Print(getReport.Code)
						if getReport.Code == 200 && len(getReport.Data) != 0 {
							break
						}

					}

					if getReport.Code != 200 {
						for _, songIdInt := range songIDList {
							getReport = RequestReferSong(getFriendID.Data.FriendCode, int64(songIdInt), true)
							if getReport.Code == 200 && len(getReport.Data) != 0 {
								break
							}
						}
					}
				}
				getReturnTypeLength := len(getReport.Data)
				if getReturnTypeLength == 0 {
					if !isIDChecker {
						ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Lucy 似乎没有查询到你的游玩数据呢（"))
					} else {
						ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Lucy 查找了对应ID 但是没有发现数据～"))
					}
					return
				}
				// DataGet, convert To MaiPlayData Render.
				maiRenderPieces := LxnsMaimaiRequestDataPiece{
					Id:           getReport.Data[len(getReport.Data)-1].Id,
					SongName:     getReport.Data[len(getReport.Data)-1].SongName,
					Level:        getReport.Data[len(getReport.Data)-1].Level,
					LevelIndex:   getReport.Data[len(getReport.Data)-1].LevelIndex,
					Achievements: getReport.Data[len(getReport.Data)-1].Achievements,
					Fc:           getReport.Data[len(getReport.Data)-1].Fc,
					Fs:           getReport.Data[len(getReport.Data)-1].Fs,
					DxScore:      getReport.Data[len(getReport.Data)-1].DxScore,
					DxRating:     getReport.Data[len(getReport.Data)-1].DxRating,
					Rate:         getReport.Data[len(getReport.Data)-1].Rate,
					Type:         getReport.Data[len(getReport.Data)-1].Type,
					UploadTime:   getReport.Data[len(getReport.Data)-1].UploadTime,
				}
				getFinalPic := ReCardRenderBase(maiRenderPieces, 0, true)
				_ = gg.NewContextForImage(getFinalPic).SavePNG(engine.DataFolder() + "save/" + "LXNS_PIC_" + strconv.Itoa(songIDList[0]) + "_" + strconv.Itoa(int(getUserID)) + ".png")
				if accStat {
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Lucy 查询到多个别名，此处默认为您返回了 "+getReport.Data[0].SongName+"谱面"), message.Image(Saved+"LXNS_PIC_"+strconv.Itoa(songIDList[0])+"_"+strconv.Itoa(int(getUserID))+".png"))
				} else {
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image(Saved+"LXNS_PIC_"+strconv.Itoa(songIDList[0])+"_"+strconv.Itoa(int(getUserID))+".png"))
				}

			} else {
				var getReport LxnsMaimaiRequestUserReferBestSongIndex
				getLevelIndexToint, _ := strconv.ParseInt(getLevelIndex, 10, 64)
				switch {
				case getSongType == "standard":
					for _, p := range songIDList {
						getReport = RequestReferSongIndex(getFriendID.Data.FriendCode, int64(p), getLevelIndexToint, true)
						if getReport.Code == 200 && getReport.Data.SongName != "" {
							break
						}
					}
					if getReport.Code == 404 {
						ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("没有发现 SD 谱面～ 如不确定可以忽略请求参数, Lucy会自动识别"))
						return
					}
				case getSongType == "dx":
					for _, p := range songIDList {
						getReport = RequestReferSongIndex(getFriendID.Data.FriendCode, int64(p), getLevelIndexToint, false)
						if getReport.Code == 200 && getReport.Data.SongName != "" {
							break
						}
					}
					if getReport.Code == 404 {
						ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("没有发现 DX 谱面～ 如不确定可以忽略请求参数, Lucy会自动识别"))
						return
					}
				default:
					for _, p := range songIDList {
						getReport = RequestReferSongIndex(getFriendID.Data.FriendCode, int64(p), getLevelIndexToint, false)
						if getReport.Code == 200 && getReport.Data.SongName != "" {
							break
						}
					}
					if getReport.Code != 200 {
						for _, p := range songIDList {
							getReport = RequestReferSongIndex(getFriendID.Data.FriendCode, int64(p), getLevelIndexToint, true)
							if getReport.Code == 200 && getReport.Data.SongName != "" {
								break
							}
						}
					}
				}
				if getReport.Data.SongName == "" { // nil pointer.
					if !isIDChecker {
						ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Lucy 似乎没有查询到你的游玩数据呢（"))
					} else {
						ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Lucy 查找了对应ID 但是没有发现数据～"))
					}
				}
				maiRenderPieces := LxnsMaimaiRequestDataPiece{
					Id:           getReport.Data.Id,
					SongName:     getReport.Data.SongName,
					Level:        getReport.Data.Level,
					LevelIndex:   getReport.Data.LevelIndex,
					Achievements: getReport.Data.Achievements,
					Fc:           getReport.Data.Fc,
					Fs:           getReport.Data.Fs,
					DxScore:      getReport.Data.DxScore,
					DxRating:     getReport.Data.DxRating,
					Rate:         getReport.Data.Rate,
					Type:         getReport.Data.Type,
					UploadTime:   getReport.Data.UploadTime,
				}
				getFinalPic := ReCardRenderBase(maiRenderPieces, 0, true)
				_ = gg.NewContextForImage(getFinalPic).SavePNG(engine.DataFolder() + "save/" + "LXNS_PIC_" + strconv.Itoa(songIDList[0]) + "_" + strconv.Itoa(int(getUserID)) + ".png")
				if accStat {
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Lucy 查询到多个别名，此处默认为您返回了 "+getReport.Data.SongName+"谱面"), message.Image(Saved+"LXNS_PIC_"+strconv.Itoa(songIDList[0])+"_"+strconv.Itoa(int(getUserID))+".png"))

				} else {
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image(Saved+"LXNS_PIC_"+strconv.Itoa(songIDList[0])+"_"+strconv.Itoa(int(getUserID))+".png"))

				}
			}
		} else {
			// DivingFish Checker
			toint := strconv.Itoa(int(ctx.Event.UserID))
			fullDevData := QueryDevDataFromDivingFish(toint)
			// default setting ==> dx Master | std Master | dx expert | std expert (as the highest score)
			var ReferSongTypeList []int
			switch {
			case getSongType == "standard":
				// roll songIDList first.
				for _, songID := range songIDList {
					if !isIDChecker {
						songID = simpleNumHandler(songID, false)
					}
					for numPosition, index := range fullDevData.Records {
						if index.SongId == songID && index.Type == "SD" {
							ReferSongTypeList = append(ReferSongTypeList, numPosition)
						}
					}
					if len(ReferSongTypeList) != 0 {
						break
					}

				}
				if len(ReferSongTypeList) == 0 {
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("没有发现游玩过的 SD 谱面～ 如不确定可以忽略请求参数, Lucy会自动识别"))
					return
				}
			case getSongType == "dx":
				for _, songID := range songIDList {
					if !isIDChecker {
						songID = simpleNumHandler(songID, true)
					}
					for numPosition, index := range fullDevData.Records {
						if index.SongId == songID && index.Type == "DX" {
							ReferSongTypeList = append(ReferSongTypeList, numPosition)
						}
					}
					if len(ReferSongTypeList) != 0 {
						break
					}
				}
				if len(ReferSongTypeList) == 0 {
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("没有发现游玩过的 DX 谱面～ 如不确定可以忽略请求参数, Lucy会自动识别"))
					return
				}
			default: // no settings.
				for _, songID := range songIDList {
					if !isIDChecker {
						songID = simpleNumHandler(songID, true)
					}
					for numPosition, index := range fullDevData.Records {
						if index.SongId == songID && index.Type == "DX" {
							ReferSongTypeList = append(ReferSongTypeList, numPosition)
						}
					}
					if len(ReferSongTypeList) != 0 {
						break
					}
				}
				if len(ReferSongTypeList) == 0 {
					for _, songID := range songIDList {
						if !isIDChecker {
							songID = simpleNumHandler(songID, false)
						}
						for numPosition, index := range fullDevData.Records {
							if index.SongId == songID && index.Type == "SD" {
								ReferSongTypeList = append(ReferSongTypeList, numPosition)
							}
						}
						if len(ReferSongTypeList) != 0 {
							break
						}
					}
				}
				if len(ReferSongTypeList) == 0 {
					if !isIDChecker {
						ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Lucy 似乎没有查询到你的游玩数据呢（"))
					} else {
						ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Lucy 查找了对应ID 但是没有发现数据～"))
					}
					return
				}
			}

			if !getReferIndexIsOn {
				// index a map =>  level_index = "record_diff"
				levelIndexMap := map[int]string{}
				for _, dataPack := range ReferSongTypeList {
					levelIndexMap[fullDevData.Records[dataPack].LevelIndex] = strconv.Itoa(dataPack)
				}
				var trulyReturnedData string
				for i := 4; i >= 0; i-- { // divingfish is reverse.
					if levelIndexMap[i] != "" {
						trulyReturnedData = levelIndexMap[i]
						break
					}
				}
				getNum, _ := strconv.Atoi(trulyReturnedData)
				// getNum ==> 0
				returnPackage := fullDevData.Records[getNum]
				_ = gg.NewContextForImage(RenderCard(returnPackage, 0, true)).SavePNG(engine.DataFolder() + "save/" + strconv.Itoa(songIDList[0]) + "_" + strconv.Itoa(int(getUserID)) + ".png")
				if accStat {
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Lucy 查询到多个别名，此处默认为您返回了 "+returnPackage.Title+" 谱面"), message.Image(Saved+strconv.Itoa(songIDList[0])+"_"+strconv.Itoa(int(getUserID))+".png"))
				} else {
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image(Saved+strconv.Itoa(songIDList[0])+"_"+strconv.Itoa(int(getUserID))+".png"))

				}
			} else {
				levelIndexMap := map[int]string{}
				for _, dataPack := range ReferSongTypeList {
					levelIndexMap[fullDevData.Records[dataPack].LevelIndex] = strconv.Itoa(dataPack)
				}
				getDiff, _ := strconv.Atoi(userSettingInterface["level_index"])
				if levelIndexMap[getDiff] != "" {
					getNum, _ := strconv.Atoi(levelIndexMap[getDiff])
					returnPackage := fullDevData.Records[getNum]
					_ = gg.NewContextForImage(RenderCard(returnPackage, 0, true)).SavePNG(engine.DataFolder() + "save/" + strconv.Itoa(songIDList[0]) + "_" + strconv.Itoa(int(getUserID)) + ".png")
					if accStat {
						ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Lucy 查询到多个别名，此处默认为您返回了 "+returnPackage.Title+" 谱面"), message.Image(Saved+strconv.Itoa(songIDList[0])+"_"+strconv.Itoa(int(getUserID))+".png"))

					} else {
						ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image(Saved+strconv.Itoa(songIDList[0])+"_"+strconv.Itoa(int(getUserID))+".png"))

					}
				} else {
					if !isIDChecker {
						ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Lucy 貌似你没有玩过这个难度的曲子哦～"))
					} else {
						ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Lucy 查找了对应ID 在这个难度下没有发现数据～"))
					}
				}
			}
		}

	})
	engine.OnRegex(`^[! ！/](mai|b50)\saliasupdate$`, zero.SuperUserPermission).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		UpdateAliasPackage()
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("成功～"))
	})
	engine.OnRegex(`^[! ！/](mai|b50)\srefer\s(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		matched := ctx.State["regex_matched"].([]string)[2]
		dataPlayer, err := QueryMaiBotDataFromUserName(matched)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
			return
		}
		var data player
		_ = json.Unmarshal(dataPlayer, &data)
		renderImg, plateStat := FullPageRender(data, ctx)
		tipPlate := ""
		if !plateStat {
			tipPlate = "tips: 可以使用 ！mai plate xxx 来绑定称号~"
		}
		_ = gg.NewContextForImage(renderImg).SavePNG(engine.DataFolder() + "save/" + strconv.Itoa(int(ctx.Event.UserID)) + ".png")
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Render User B50 : "+data.Username+"\n"+tipPlate+"\n"), message.Image(Saved+strconv.Itoa(int(ctx.Event.UserID))+".png"))
	})
}
