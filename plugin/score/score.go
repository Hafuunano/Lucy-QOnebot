package score

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/img/writer"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/img"
	"github.com/FloatTech/zbputils/img/text"
	"github.com/fogleman/gg"
	_ "github.com/fumiama/sqlite3" // import sql
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	backgroundURL = "https://img.moehu.org/pic.php?id=pc&return=json"
	signinMax     = 1
	// SCOREMAX 分数上限定为600
	SCOREMAX = 600
)

var (
	levelArray = [...]int{0, 1, 2, 5, 10, 20, 35, 55, 75, 100, 120, 180, 260, 360, 480, 600}
	sdb        *scoredb
	engine     = control.Register("score", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "Hi NekoPachi!\n说明书: https://manual-lucy.himoyo.cn",
		PrivateDataFolder: "score",
	})
	pgs        = make(pg, 256)
	RateLimit  = rate.NewManager[int64](time.Second*60, 9)
	CheckLimit = rate.NewManager[int64](time.Minute*1, 4)
	CatchLimit = rate.NewManager[int64](time.Hour*1, 8)
)

type partygame struct {
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Coins string `json:"coins"`
}

type pg = map[string]partygame

// scoredb 分数数据库
type scoredb gorm.DB

// scoretable 分数结构体
type scoretable struct {
	UID   int64 `gorm:"column:uid;primary_key"`
	Score int   `gorm:"column:score;default:0"`
	Coins int   `gorm:"column:coins;default:0"`
}

// TableName ...
func (scoretable) TableName() string {
	return "score"
}

// signintable 签到结构体
type signintable struct {
	UID       int64 `gorm:"column:uid;primary_key"`
	Count     int   `gorm:"column:count;default:0"`
	Coins     int   `gorm:"column:coins;default:0"`
	UpdatedAt time.Time
}

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
				ctx.SendChain(message.Text("ERROR:", err, "\nLoading Pic ERROR"))
				return
			}
			if siUpdateTimeStr != today {
				_ = sdb.InsertOrUpdateSignInCountByUID(uid, 0, 0)
			}
			coinsGet := rand.Intn(100)
			_ = sdb.InsertOrUpdateSignInCountByUID(uid, si.Count+1, si.Coins+coinsGet)
			// 避免图片过大，最大 1280*720
			back = img.Limit(back, 1280, 720)

			canvas := gg.NewContext(back.Bounds().Size().X, int(float64(back.Bounds().Size().Y)*1.7))
			canvas.SetRGB255(137, 207, 240)
			canvas.Clear()
			canvas.DrawImage(back, 0, 0)
			monthWord := now.Format("01/02")
			hourWord, handleMsg := getHourWord(now)
			_, err = file.GetLazyData(text.BoldFontFile, false)
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
			_, err = file.GetLazyData(text.FontFile, false)
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
			_ = sdb.InsertOrUpdateScoreByUID(uid, score, coinsGet)
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
	engine.OnPrefix("获得签到背景", zero.OnlyGroup).SetBlock(true).
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
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("今天你还没有签到哦w"))
				return
			}
			ctx.SendChain(message.Image("file:///" + file.BOTPATH + "/" + picFile))

		})

		// 借鉴了其他bot的功能 编写而成
	engine.OnFullMatch("柠檬片总数", zero.OnlyGroup).Handle(func(ctx *zero.Ctx) {
		uid := ctx.Event.UserID
		si := sdb.GetSignInByUID(uid)
		ctx.SendChain(message.Text("您的柠檬片数量一共是: ", si.Coins))
	})
	engine.OnFullMatch("抽奖", zero.OnlyGroup).Handle(func(ctx *zero.Ctx) {
		if !CatchLimit.Load(ctx.Event.UserID).Acquire() {
			ctx.SendChain(message.Text("太贪心了哦~过会试试吧"))
			return
		}
		var mutex sync.RWMutex // 添加读写锁以保证稳定性
		mutex.Lock()
		uid := ctx.Event.UserID
		data, err := os.ReadFile(engine.DataFolder() + "loads.json")
		if err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			return
		}
		err = json.Unmarshal(data, &pgs)
		if err != nil {
			panic(err)
		}
		si := sdb.GetSignInByUID(uid) // 获取用户目前状况信息
		userCurrentCoins := si.Coins  // loading coins status
		checkEnoughCoins := checkUserCoins(userCurrentCoins)
		if checkEnoughCoins == false {
			ctx.SendChain(message.Reply(uid), message.Text("本次参与的柠檬片不够哦~请多多打卡w"))
			return
		}
		all := rand.Intn(38) // 一共37种可能性
		referpg := pgs[(strconv.Itoa(all))]
		getName := referpg.Name
		getCoinsStr := referpg.Coins
		getCoinsInt, _ := strconv.Atoi(getCoinsStr)
		getDesc := referpg.Desc
		addNewCoins := si.Coins + getCoinsInt - 15
		err = sdb.InsertOrUpdateSignInCountByUID(uid, 0, addNewCoins)
		if err != nil {
			ctx.SendChain(message.Text("ERR:", err))
			return
		}
		msgid := ctx.SendChain(message.At(uid), message.Text(" ~好哦w 来抽奖 看看能抽到什么东西呢w"))
		time.Sleep(time.Second * 3)
		ctx.SendChain(message.Reply(msgid), message.Text("好哦~让咱看看你抽到了什么东西ww\n"),
			message.Text("你抽到的是~ ", getName, "\n", "获得了积分 ", getCoinsInt, "\n", getDesc, "\n目前的柠檬片总数为：", addNewCoins))
		mutex.Unlock()
	})
	// 一次最多骗 50 柠檬片,失败概率较大,失败会被反吞柠檬片

	engine.OnRegex(`^骗(\[CQ:at,qq=(\d+)\]\s?|(\d+))的柠檬片`, zero.OnlyGroup).Handle(func(ctx *zero.Ctx) {
		if !CatchLimit.Load(ctx.Event.UserID).Acquire() {
			ctx.SendChain(message.Text("太贪心了哦~一小时后再来试试吧"))
			return
		}
		TargetInt, _ := strconv.ParseInt(ctx.State["regex_matched"].([]string)[2], 10, 64)
		uid := ctx.Event.UserID
		siEventUser := sdb.GetSignInByUID(uid)        // 获取主用户目前状况信息
		SiTargetUser := sdb.GetSignInByUID(TargetInt) // 获得被抢用户目前情况信息
		switch {
		case siEventUser.Coins < 50:
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("貌似没有足够的柠檬片去准备哦~请多多打卡w"))
			return
		case SiTargetUser.Coins < 50:
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("太坏了~试图的对象貌似没有足够多的柠檬片~"))
			return
		}
		eventUserName := ctx.CardOrNickName(uid)
		eventTargetName := ctx.CardOrNickName(TargetInt)
		modifyCoins := rand.Intn(50)
		if rand.Intn(10)/7 != 0 { // 6成失败概率
			_ = sdb.InsertOrUpdateSignInCountByUID(siEventUser.UID, 0, siEventUser.Coins-modifyCoins)
			time.Sleep(time.Second * 2)
			_ = sdb.InsertOrUpdateSignInCountByUID(SiTargetUser.UID, 0, SiTargetUser.Coins+modifyCoins)
			time.Sleep(time.Second * 2)
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("试着去拿走 ", eventTargetName, " 的柠檬片时,被发现了.\n所以 ", eventUserName, " 失去了 ", modifyCoins, " 个柠檬片\n\n同时 ", eventTargetName, " 得到了 ", modifyCoins, " 个柠檬片"))
			return
		}
		_ = sdb.InsertOrUpdateSignInCountByUID(siEventUser.UID, 0, siEventUser.Coins+modifyCoins)
		time.Sleep(time.Second * 2)
		_ = sdb.InsertOrUpdateSignInCountByUID(SiTargetUser.UID, 0, SiTargetUser.Coins-modifyCoins)
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("试着去拿走 ", eventTargetName, " 的柠檬片时,被成功了.\n所以 ", eventTargetName, " 得到了 ", modifyCoins, " 个柠檬片\n\n同时 ", eventUserName, " 失去了 ", modifyCoins, " 个柠檬片"))
	})
}
