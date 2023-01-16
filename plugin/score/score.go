package score // Package score

import (
	"github.com/FloatTech/ZeroBot-Plugin/plugin/funwork"
	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/web"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/fogleman/gg"
	"github.com/tidwall/gjson"
	"github.com/wdvxdr1123/ZeroBot/message"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	_ "github.com/FloatTech/sqlite" // import sql
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
)

const (
	backgroundURL = "https://iw233.cn/api.php?sort=iw233&type=json"
	signinMax     = 1
)

var (
	levelArray = [...]int{0, 1, 2, 5, 10, 20, 35, 55, 75, 100, 120, 180, 260, 360, 480, 600}
	sdb        *scoredb
	engine     = control.Register("score", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "Hi NekoPachi!\n说明书: https://manual-lucy.himoyo.cn",
		PrivateDataFolder: "score",
	})
)

func init() {
	cachePath := engine.DataFolder() + "scorecache/"
	engine.OnFullMatchGroup([]string{"签到", "打卡"}, zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByGroup).
		Handle(func(ctx *zero.Ctx) {
			if !RateLimit.Load(ctx.Event.GroupID).Acquire() {
				return
			}
			uid := ctx.Event.UserID
			now := time.Now()
			today := now.Format("20060102")
			si := sdb.GetSignInByUID(uid)
			drawedFile := cachePath + strconv.FormatInt(uid, 10) + today + "signin.png"
			siUpdateTimeStr := si.UpdatedAt.Format("20060102")
			if si.Count >= signinMax && siUpdateTimeStr == today {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("酱~ 你今天已经签到过了哦w"))
				if file.IsExist(drawedFile) {
					ctx.SendChain(message.Image("file:///" + file.BOTPATH + "/" + drawedFile))
				}
				return
			}
			coinsGet := 200 + rand.Intn(200)
			_ = sdb.InsertUserCoins(uid, si.Coins+coinsGet)
			_ = sdb.InsertOrUpdateSignInCountByUID(uid, si.Count+1) // 柠檬片获取
			score := sdb.GetScoreByUID(uid).Score
			score += 1 //  每日+1
			_ = sdb.InsertOrUpdateScoreByUID(uid, score)
			CurrentCountTable := sdb.GetCurrentCount(today)
			handledTodayNum := CurrentCountTable.Counttime + 1
			_ = sdb.UpdateUserTime(handledTodayNum, today) // 总体计算 隔日清零
			if now.Hour() > 6 && now.Hour() < 19 {
				// package for test draw.
				getTimeReplyMsg := getHourWord(time.Now()) // get time and msg
				currentTime := time.Now().Format("2006-01-02 15:04:05")
				// time day.
				dayTimeImg, _ := gg.LoadImage(engine.DataFolder() + "BetaScoreDay.png")
				dayGround := gg.NewContext(1920, 1080)
				dayGround.DrawImage(dayTimeImg, 0, 0)
				_ = dayGround.LoadFontFace(engine.DataFolder()+"dyh.ttf", 60)
				dayGround.SetRGB(0, 0, 0)

				// draw something with cautions Only (
				dayGround.DrawString(currentTime, 1270, 950)            // draw time
				dayGround.DrawString(getTimeReplyMsg, 50, 930)          // draw text.
				dayGround.DrawString(ctx.CardOrNickName(uid), 310, 110) // draw name :p why I should do this???
				_ = dayGround.LoadFontFace(engine.DataFolder()+"dyh.ttf", 60)
				dayGround.DrawStringWrapped(strconv.Itoa(handledTodayNum), 350, 255, 1, 1, 0, 1.3, gg.AlignCenter)   // draw first part
				dayGround.DrawStringWrapped(strconv.Itoa(si.Count+1), 1000, 255, 1, 1, 0, 1.3, gg.AlignCenter)       // draw second part
				dayGround.DrawStringWrapped(strconv.Itoa(coinsGet), 220, 370, 1, 1, 0, 1.3, gg.AlignCenter)          // draw third part
				dayGround.DrawStringWrapped(strconv.Itoa(si.Coins+coinsGet), 720, 370, 1, 1, 0, 1.3, gg.AlignCenter) // draw forth part
				// level array with rectangle work.
				rankNum := getLevel(score)
				RankGoal := rankNum + 1
				achieveNextGoal := levelArray[RankGoal]
				achievedGoal := levelArray[rankNum]
				currentNextGoalMeasure := achieveNextGoal - score  // measure rest of the num. like 20 - currentLink(TestRank 15)
				measureGoalsLens := achieveNextGoal - achievedGoal // like 20 - 10
				currentResult := float64(currentNextGoalMeasure) / float64(measureGoalsLens)
				// draw this part
				dayGround.SetRGB255(180, 255, 254)        // aqua color
				dayGround.DrawRectangle(70, 570, 600, 50) // draw rectangle part1
				dayGround.Fill()
				dayGround.SetRGB255(130, 255, 254)
				dayGround.DrawRectangle(70, 570, 600*currentResult, 50) // draw rectangle part2
				dayGround.Fill()
				dayGround.SetRGB255(0, 0, 0)
				dayGround.DrawString("Lv. "+strconv.Itoa(rankNum)+" 签到天数 + 1", 80, 490)
				_ = dayGround.LoadFontFace(engine.DataFolder()+"dyh.ttf", 40)
				dayGround.DrawString(strconv.Itoa(currentNextGoalMeasure)+"/"+strconv.Itoa(measureGoalsLens), 710, 610)
				_ = dayGround.SavePNG(drawedFile)
				ctx.SendChain(message.Image("file:///" + file.BOTPATH + "/" + drawedFile))
				time.Sleep(time.Second * 5)
				data, err := web.RequestDataWith(web.NewDefaultClient(), backgroundURL, "GET", funwork.Referer, web.RandUA(), nil)
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				picURLRaw := gjson.Get(string(data), "pic.0").String()
				replaceRegexp := regexp.MustCompile(`https://[0-9a-zA-Z]+.sinaimg.cn/`)
				picURL := replaceRegexp.ReplaceAllString(picURLRaw, "https://simg.himoyo.cn/")
				deleteThisOne := ctx.SendChain(message.Image(picURL))
				time.Sleep(time.Second * 40)
				ctx.DeleteMessage(deleteThisOne)
			} else {
				// nightVision
				// package for test draw.
				getTimeReplyMsg := getHourWord(time.Now()) // get time and msg
				currentTime := time.Now().Format("2006-01-02 15:04:05")
				nightTimeImg, _ := gg.LoadImage(engine.DataFolder() + "BetaScoreNight.png")
				nightGround := gg.NewContext(1886, 1080)
				nightGround.DrawImage(nightTimeImg, 0, 0)
				_ = nightGround.LoadFontFace(engine.DataFolder()+"dyh.ttf", 60)
				nightGround.SetRGB255(255, 255, 255)
				// draw something with cautions Only (
				nightGround.DrawString(currentTime, 1360, 910)            // draw time
				nightGround.DrawString(getTimeReplyMsg, 60, 930)          // draw text.
				nightGround.DrawString(ctx.CardOrNickName(uid), 350, 140) // draw name :p why I should do this???
				_ = nightGround.LoadFontFace(engine.DataFolder()+"dyh.ttf", 60)
				nightGround.DrawStringWrapped(strconv.Itoa(handledTodayNum), 345, 275, 1, 1, 0, 1.3, gg.AlignCenter)   // draw first part
				nightGround.DrawStringWrapped(strconv.Itoa(si.Count+1), 990, 275, 1, 1, 0, 1.3, gg.AlignCenter)        // draw second part
				nightGround.DrawStringWrapped(strconv.Itoa(coinsGet), 225, 360, 1, 1, 0, 1.3, gg.AlignCenter)          // draw third part
				nightGround.DrawStringWrapped(strconv.Itoa(si.Coins+coinsGet), 720, 360, 1, 1, 0, 1.3, gg.AlignCenter) // draw forth part
				// level array with rectangle work.
				rankNum := getLevel(score)
				RankGoal := rankNum + 1
				achieveNextGoal := levelArray[RankGoal]
				achievedGoal := levelArray[rankNum]
				currentNextGoalMeasure := achieveNextGoal - score  // measure rest of the num. like 20 - currentLink(TestRank 15)
				measureGoalsLens := achieveNextGoal - achievedGoal // like 20 - 10
				currentResult := float64(currentNextGoalMeasure) / float64(measureGoalsLens)
				// draw this part
				nightGround.SetRGB255(49, 86, 157)          // aqua color
				nightGround.DrawRectangle(70, 570, 600, 50) // draw rectangle part1
				nightGround.Fill()
				nightGround.SetRGB255(255, 255, 255)
				nightGround.DrawRectangle(70, 570, 600*currentResult, 50) // draw rectangle part2
				nightGround.Fill()
				nightGround.SetRGB255(255, 255, 255)
				nightGround.DrawString("Lv. "+strconv.Itoa(rankNum)+" 签到天数 + 1", 80, 490)
				_ = nightGround.LoadFontFace(engine.DataFolder()+"dyh.ttf", 40)
				nightGround.DrawString(strconv.Itoa(currentNextGoalMeasure)+"/"+strconv.Itoa(measureGoalsLens), 710, 610)
				_ = nightGround.SavePNG(drawedFile)
				ctx.SendChain(message.Image("file:///" + file.BOTPATH + "/" + drawedFile))
				time.Sleep(time.Second * 5)
				data, err := web.RequestDataWith(web.NewDefaultClient(), backgroundURL, "GET", funwork.Referer, web.RandUA(), nil)
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				picURLRaw := gjson.Get(string(data), "pic.0").String()
				replaceRegexp := regexp.MustCompile(`https://[0-9a-zA-Z]+.sinaimg.cn/`)
				picURL := replaceRegexp.ReplaceAllString(picURLRaw, "https://simg.himoyo.cn/")
				deleteThisOne := ctx.SendChain(message.Image(picURL))
				time.Sleep(time.Second * 40)
				ctx.DeleteMessage(deleteThisOne)
			}
		})
}

/*
		cachePath := engine.DataFolder() + "scorecache/"
		engine.OnFullMatchGroup([]string{"签到", "打卡"}, zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByGroup).
			Handle(func(ctx *zero.Ctx) {
				if !RateLimit.Load(ctx.Event.GroupID).Acquire() {
					return
				}
				uid := ctx.Event.UserID
				now := time.Now()
				today := now.Format("20060102")
				si := sdb.GetSignInByUID(uid)
				drawedFile := cachePath + strconv.FormatInt(uid, 10) + today + "signin.png"
				picFile := cachePath + strconv.FormatInt(uid, 10) + today + ".png"
				initPic(picFile)
				siUpdateTimeStr := si.UpdatedAt.Format("20060102")
				if si.Count >= signinMax && siUpdateTimeStr == today {
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("酱~ 你今天已经签到过了哦w"))
					if file.IsExist(drawedFile) {
						ctx.SendChain(message.Image("file:///" + file.BOTPATH + "/" + drawedFile))
					}
					return
				}
				back, err := gg.LoadImage(picFile)
				if err != nil {
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERROR:", err, "\nloading picture err , cannot draw."))
					return
				}
				if siUpdateTimeStr != today {
					_ = sdb.InsertOrUpdateSignInCountByUID(uid, 0)
				}
				coinsGet := 200 + rand.Intn(200)
				_ = sdb.InsertUserCoins(uid, si.Coins+coinsGet)
				_ = sdb.InsertOrUpdateSignInCountByUID(uid, si.Count+1)
				// 避免图片过大，最大 1280*720
				back = img.Limit(back, 1280, 720)

				canvas := gg.NewContext(back.Bounds().Size().X, int(float64(back.Bounds().Size().Y)*1.7))
				canvas.SetRGB255(137, 207, 240)
				canvas.Clear()
				canvas.DrawImage(back, 0, 0)
				monthWord := now.Format("01/02")
				hourWord, handleMsg := getHourWord(now)
				_, err = file.GetLazyData(text.BoldFontFile, control.Md5File, false)
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				if err = canvas.LoadFontFace(text.BoldFontFile, float64(back.Bounds().Size().X)*0.1); err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				canvas.SetRGB(255, 191, 205)
				canvas.DrawString(hourWord, float64(back.Bounds().Size().X)*0.1, float64(back.Bounds().Size().Y)*1.2)
				canvas.DrawString(monthWord, float64(back.Bounds().Size().X)*0.6, float64(back.Bounds().Size().Y)*1.2)
				nickName := ctx.CardOrNickName(uid)
				_, err = file.GetLazyData(text.FontFile, control.Md5File, false)
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				if err = canvas.LoadFontFace(text.FontFile, float64(back.Bounds().Size().X)*0.04); err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				add := 1
				canvas.DrawString(nickName+fmt.Sprintf(" 签到天数+%d", add), float64(back.Bounds().Size().X)*0.1, float64(back.Bounds().Size().Y)*1.3)
				score := sdb.GetScoreByUID(uid).Score
				score += add
				_ = sdb.InsertOrUpdateScoreByUID(uid, score)
				level := getLevel(score)
				canvas.DrawString("当前签到天数:"+strconv.FormatInt(int64(score), 10)+"  |  柠檬片 + "+strconv.FormatInt(int64(coinsGet), 10)+" 片", float64(back.Bounds().Size().X)*0.1, float64(back.Bounds().Size().Y)*1.4)
				canvas.DrawString("LEVEL:"+strconv.FormatInt(int64(level), 10)+" | "+handleMsg, float64(back.Bounds().Size().X)*0.1, float64(back.Bounds().Size().Y)*1.5)
				canvas.DrawRectangle(float64(back.Bounds().Size().X)*0.1, float64(back.Bounds().Size().Y)*1.55, float64(back.Bounds().Size().X)*0.6, float64(back.Bounds().Size().Y)*0.1)
				canvas.SetRGB255(150, 150, 150)
				canvas.Fill()
				var nextLevelScore int
				if level < 10 {
					nextLevelScore = levelArray[level+1]
				} else {
					nextLevelScore = SCOREMAX
				}
				canvas.SetRGB255(255, 191, 205)
				canvas.DrawRectangle(float64(back.Bounds().Size().X)*0.1, float64(back.Bounds().Size().Y)*1.55, float64(back.Bounds().Size().X)*0.6*float64(score)/float64(nextLevelScore), float64(back.Bounds().Size().Y)*0.1)
				canvas.SetRGB255(102, 102, 102)
				canvas.Fill()
				canvas.DrawString(fmt.Sprintf("%d/%d", score, nextLevelScore), float64(back.Bounds().Size().X)*0.75, float64(back.Bounds().Size().Y)*1.62)
				f, err := os.Create(drawedFile)
				if err != nil {
					log.Errorln("[score]", err)
					data, cl := writer.ToBytes(canvas.Image())
					ctx.SendChain(message.ImageBytes(data))
					cl()
					return
				}
				_, err = writer.WriteTo(canvas.Image(), f)
				_ = f.Close()
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				ctx.SendChain(message.Image("file:///" + file.BOTPATH + "/" + drawedFile))
				time.Sleep(time.Second * 5)
			})

	 engine.OnPrefixGroup([]string{"获得打卡背景", "获得签到背景"}, zero.OnlyGroup).SetBlock(true).
			Handle(func(ctx *zero.Ctx) {
				param := ctx.State["args"].(string)
				var uidStr string
				if len(ctx.Event.Message) > 1 && ctx.Event.Message[1].Type == "at" {
					uidStr = ctx.Event.Message[1].Data["qq"]
				} else if param == "" {
					uidStr = strconv.FormatInt(ctx.Event.UserID, 10)
				}
				picFile := cachePath + uidStr + time.Now().Format("20060102") + ".png"
				if file.IsNotExist(picFile) {
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("今天你还没有打卡哦w"))
					return
				}
				ctx.SendChain(message.Image("file:///" + file.BOTPATH + "/" + picFile))
			})
*/
