package lovekitten

import (
	"math/rand"
	"sync"
	"time"

	sql "github.com/FloatTech/sqlite"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	engine = control.Register("loveKitten", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help:             "Core of LoveKitten, 应当写在./下的  懒死了",
	})
	userLovePackage = &sql.Sqlite{}
	userLocker      sync.RWMutex
)

func init() {
	userLovePackage.DBPath = engine.DataFolder() + "userLove.db"
	err := userLovePackage.Open(time.Hour * 24)
	if err != nil {
		panic(err)
	}
	_ = createFirstTable(userLovePackage) // 第一次没有就生成这东西
	engine.On("notice/notify/poke", zero.OnlyToMe).SetBlock(false).Handle(func(ctx *zero.Ctx) {

	})
}

func createFirstTable(db *sql.Sqlite) error {
	userLocker.Lock() // 没有就干一个算了w
	defer userLocker.Unlock()
	return db.Create("globalUser", &userLovePackage)
}

func addCurrentNumberToReferred(db *sql.Sqlite) (Num int, err error) {
	userLocker.Lock()
	defer userLocker.Unlock()
	RanNum := rand.Intn(10)
	db.Insert("globalUser", &userLovePackage)
	return ctx.Send(message.Text("貌似因为某种原因~咱对你的好感度变化了", ModifyNum))
}

/*
可能需要被改动的地方 :

Chat插件组 分为两个词库组 一个是 HighKitten 还有一个是 LowKitten

需要互动的地方有 : 戳一戳（随机 + 每日限制次数）功能使用频繁程度（某个人过于频繁调用特殊功能的话 则 down）

好感度的动态应该以每日10以内为准 默认为200 在 100 or lower 为低好感度 300 or higher为高好感度

WDNMD我写的什么几把玩意

*/
