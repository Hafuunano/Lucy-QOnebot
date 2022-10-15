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
	pgs       = make(pg, 256)
	RateLimit = rate.NewManager[int64](time.Second*60, 9)
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
	engine.OnFullMatch("签到", zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByGroup).
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
	engine.OnFullMatch("柠檬片总数", zero.OnlyGroup).Handle(func(ctx *zero.Ctx) {
		uid := ctx.Event.UserID
		si := sdb.GetSignInByUID(uid)
		ctx.SendChain(message.Text("您的柠檬片数量一共是: ", si.Coins))
	})
	engine.OnFullMatch("抽奖").Handle(func(ctx *zero.Ctx) {
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
		if checkEnoughCoins == true {
			ctx.SendChain(message.Reply(uid), message.Text("酱~抽奖哦~~w"))
			time.Sleep(time.Second * 3)
		} else {
			ctx.SendChain(message.Reply(uid), message.Text("本次参与的柠檬片不够哦~请多多打卡w"))
			return
		}
		err = sdb.InsertOrUpdateSignInCountByUID(uid, 0, si.Coins-15)
		if err != nil {
			ctx.SendChain(message.Text("ERR: ", err))
			return
		}
		all := rand.Intn(38) // 一共37种可能性
		referpg := pgs[(strconv.Itoa(all))]
		getName := referpg.Name
		getCoinsStr := referpg.Coins
		getCoinsInt, _ := strconv.Atoi(getCoinsStr)
		getDesc := referpg.Desc
		err = sdb.InsertOrUpdateSignInCountByUID(uid, 0, si.Coins+getCoinsInt)
		if err != nil {
			ctx.SendChain(message.Text("ERR:", err))
			return
		}
		ctx.SendChain(message.Reply(uid), message.Text("好哦~让咱看看你抽到了什么东西ww\n"),
			message.Text("你抽到的是~", getName, "\n", "获得了积分", getCoinsInt, "\n", getDesc, "\n目前的柠檬片总数为：", si.Coins+getCoinsInt))
		mutex.Unlock()
	})
}
