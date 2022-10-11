package score

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/img/writer"
	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/img"
	"github.com/FloatTech/zbputils/img/text"
	"github.com/fogleman/gg"
	_ "github.com/fumiama/sqlite3" // import sql
	"github.com/golang/freetype"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/wcharczuk/go-chart/v2"
	zero "github.com/wdvxdr1123/ZeroBot"
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
)

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
	go func() {
		sdb = initialize(engine.DataFolder() + "score.db")
		log.Println("[score]加载score数据库")
	}()
	engine.OnFullMatch("签到", zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByGroup).
		Handle(func(ctx *zero.Ctx) {
			uid := ctx.Event.UserID
			now := time.Now()
			today := now.Format("20060102")
			si := sdb.GetSignInByUID(uid)
			drawedFile := cachePath + strconv.FormatInt(uid, 10) + today + "signin.png"
			picFile := cachePath + strconv.FormatInt(uid, 10) + today + ".png"
			initPic(picFile)
			siUpdateTimeStr := si.UpdatedAt.Format("20060102")
			if si.Count >= signinMax && siUpdateTimeStr == today {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("酱~ 你今天已经抢到过了哦w"))
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
			canvas.SetRGB(1, 1, 1)
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
			canvas.SetRGB(0, 0, 0)
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
			canvas.DrawString("LEVEL:"+strconv.FormatInt(int64(level), 10), float64(back.Bounds().Size().X)*0.1, float64(back.Bounds().Size().Y)*1.5)
			canvas.DrawRectangle(float64(back.Bounds().Size().X)*0.1, float64(back.Bounds().Size().Y)*1.55, float64(back.Bounds().Size().X)*0.6, float64(back.Bounds().Size().Y)*0.1)
			canvas.SetRGB255(150, 150, 150)
			canvas.Fill()
			var nextLevelScore int
			if level < 10 {
				nextLevelScore = levelArray[level+1]
			} else {
				nextLevelScore = SCOREMAX
			}
			canvas.SetRGB255(0, 0, 0)
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
			ctx.SendChain(message.At(ctx.Event.UserID), message.Text("\n"), message.Text(handleMsg))
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
			os.Remove(picFile)
		})
	engine.OnFullMatch("查看签到排名", zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			today := time.Now().Format("20060102")
			drawedFile := cachePath + today + "scoreRank.png"
			if file.IsExist(drawedFile) {
				ctx.SendChain(message.Image("file:///" + file.BOTPATH + "/" + drawedFile))
				return
			}
			st, err := sdb.GetScoreRankByTopN(10)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			if len(st) == 0 {
				ctx.SendChain(message.Text("ERROR:目前还没有人签到过"))
				return
			}
			_, err = file.GetLazyData(text.FontFile, true)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			b, err := os.ReadFile(text.FontFile)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			font, err := freetype.ParseFont(b)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			bars := make([]chart.Value, len(st))
			for i, v := range st {
				bars[i] = chart.Value{
					Value: float64(v.Score),
					Label: ctx.CardOrNickName(v.UID),
				}
			}
			graph := chart.BarChart{
				Font:  font,
				Title: "签到天数",
				Background: chart.Style{
					Padding: chart.Box{
						Top: 40,
					},
				},
				Height:   500,
				BarWidth: 50,
				Bars:     bars,
			}
			f, err := os.Create(drawedFile)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			err = graph.Render(chart.PNG, f)
			_ = f.Close()
			if err != nil {
				_ = os.Remove(drawedFile)
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			ctx.SendChain(message.Image("file:///" + file.BOTPATH + "/" + drawedFile))
		})

	engine.OnFullMatch("柠檬片总数", zero.OnlyGroup).Handle(func(ctx *zero.Ctx) {
		uid := ctx.Event.UserID
		si := sdb.GetSignInByUID(uid)
		ctx.SendChain(message.Text("您的柠檬片数量一共是: ", si.Coins))
	})
}

func getHourWord(t time.Time) (sign string, reply string) {
	switch {
	case 5 <= t.Hour() && t.Hour() < 12:
		sign = "早上好"
		reply = "今天又是新的一天呢ww"
	case 12 <= t.Hour() && t.Hour() < 14:
		sign = "中午好"
		reply = "吃饭了嘛w 如果没有快去吃饭哦w"
	case 14 <= t.Hour() && t.Hour() < 19:
		sign = "下午好"
		reply = "记得多去补点水哦~ 适当的时候可以出去转转ww"
	case 19 <= t.Hour() && t.Hour() < 24:
		sign = "晚上好"
		reply = "今天过的开心嘛ww"
	case 0 <= t.Hour() && t.Hour() < 5:
		sign = "凌晨好"
		reply = "快去睡觉哦w 不然Lucy会生气的x"
	default:
		sign = ""
	}
	return
}

func getLevel(count int) int {
	for k, v := range levelArray {
		if count == v {
			return k
		} else if count < v {
			return k - 1
		}
	}
	return -1
}

func initPic(picFile string) {
	if file.IsNotExist(picFile) {
		data, err := web.GetData(backgroundURL)
		if err != nil {
			log.Errorln("[score]", err)
		}
		picURL := gjson.Get(string(data), "acgurl").String()
		data, err = web.GetData(picURL)
		if err != nil {
			log.Errorln("[score]", err)
		}
		err = os.WriteFile(picFile, data, 0666)
		if err != nil {
			log.Errorln("[score]", err)
		}
	}
}

// TableName ...
func (signintable) TableName() string {
	return "sign_in"
}

// initialize 初始化ScoreDB数据库
func initialize(dbpath string) *scoredb {
	var err error
	if _, err = os.Stat(dbpath); err != nil || os.IsNotExist(err) {
		// 生成文件
		f, err := os.Create(dbpath)
		if err != nil {
			return nil
		}
		defer f.Close()
	}
	gdb, err := gorm.Open("sqlite3", dbpath)
	if err != nil {
		panic(err)
	}
	gdb.AutoMigrate(&scoretable{}).AutoMigrate(&signintable{})
	return (*scoredb)(gdb)
}

// Close ...
func (sdb *scoredb) Close() error {
	db := (*gorm.DB)(sdb)
	return db.Close()
}

// GetScoreByUID 取得分数
func (sdb *scoredb) GetScoreByUID(uid int64) (s scoretable) {
	db := (*gorm.DB)(sdb)
	db.Debug().Model(&scoretable{}).FirstOrCreate(&s, "uid = ? ", uid)
	return s
}

// InsertOrUpdateScoreByUID 插入或更新打卡累计 + 签到分数(随机化)
func (sdb *scoredb) InsertOrUpdateScoreByUID(uid int64, score int, coins int) (err error) {
	db := (*gorm.DB)(sdb)
	s := scoretable{
		UID:   uid,
		Score: score,
		Coins: coins,
	}
	if err = db.Debug().Model(&scoretable{}).First(&s, "uid = ? ", uid).Error; err != nil {
		// error handling...
		if gorm.IsRecordNotFoundError(err) {
			err = db.Debug().Model(&scoretable{}).Create(&s).Error // newUser not user
		}
	} else {
		err = db.Debug().Model(&scoretable{}).Where("uid = ? ", uid).Update(
			map[string]interface{}{
				"score": score,
			}).Error
	}
	return
}

// GetSignInByUID 取得签到次数
func (sdb *scoredb) GetSignInByUID(uid int64) (si signintable) {
	db := (*gorm.DB)(sdb)
	db.Debug().Model(&signintable{}).FirstOrCreate(&si, "uid = ? ", uid)
	return si
}

// InsertOrUpdateSignInCountByUID 插入或更新签到次数
func (sdb *scoredb) InsertOrUpdateSignInCountByUID(uid int64, count int, coins int) (err error) {
	db := (*gorm.DB)(sdb)
	si := signintable{
		UID:   uid,
		Count: count,
		Coins: coins,
	}
	if err = db.Debug().Model(&signintable{}).First(&si, "uid = ? ", uid).Error; err != nil {
		// error handling...
		if gorm.IsRecordNotFoundError(err) {
			db.Debug().Model(&signintable{}).Create(&si) // newUser not user
		}
	} else {
		err = db.Debug().Model(&signintable{}).Where("uid = ? ", uid).Update(
			map[string]interface{}{
				"count": count,
				"Coins": coins,
			}).Error
	}
	return
}
func (sdb *scoredb) GetScoreRankByTopN(n int) (st []scoretable, err error) {
	db := (*gorm.DB)(sdb)
	err = db.Model(&scoretable{}).Order("score desc").Limit(n).Find(&st).Error
	return
}
