package bottle // zbp driftbottle 魔改

import (
	"fmt"
	"hash/crc64"
	"strconv"
	"sync"
	"time"

	"github.com/FloatTech/floatbox/binary"
	sql "github.com/FloatTech/sqlite"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type sea struct {
	ID   int64  `db:"id"`   // ID qq_grp_name_msg 的 crc64 hashCheck.
	QQ   int64  `db:"qq"`   // Get current user(Who sends this)
	Name string `db:"Name"` //  his or her name at that time:P
	Msg  string `db:"msg"`  // What he or she sent to Lucy?
	Grp  int64  `db:"grp"`  // which group sends this msg?
	Time int64  `db:"time"` // we need to know the current time,master>
}

var seaSide = &sql.Sqlite{}
var seaLocker sync.RWMutex

// We need a container to receive what we need :(

func init() {
	engine := control.Register("bottle", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "基于driftbottle的魔改版漂流瓶\n",
		PrivateDataFolder: "bottle",
	})
	seaSide.DBPath = engine.DataFolder() + "sea.db"
	err := seaSide.Open(time.Hour * 24)
	if err != nil {
		panic(err)
	}
	_ = CreateChannel(seaSide)
	engine.OnFullMatch("pick", zero.OnlyToMe).Limit(ctxext.LimitByGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		be, err := fetchBottle(seaSide)
		if err != nil {
			ctx.SendChain(message.Text("ERR:", err))
		}
		SenderTime := time.Unix(be.Time, 0).Format("2006-01-01 12:33:36")
		ctx.SendChain(message.Text("~Lucy试着帮你捞出来了这个~\nID: ", be.ID, "\n投递人: ", be.Name, "(", be.QQ, ")", "\n群号: ", be.Grp, "\n时间: ", SenderTime, "\n内容: \n", be.Msg))
	})
	engine.OnFullMatch("throw", zero.OnlyToMe).Limit(ctxext.LimitByGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		var relief extension.CommandModel
		err := ctx.Parse(&relief)
		if err != nil {
		}
		if relief.Args == "" {
			ctx.Send(message.Text("~需要咱帮你丢什么呢w"))
			nextstep := ctx.FutureEvent("message", ctx.CheckSession())
			recv, cancel := nextstep.Repeat()
			for i := range recv {
				msg := i.MessageString()
				{
					if msg != "" {
						relief.Args = msg
						cancel()
						continue
					}
				}
			}
		}
		// check current needs and Prepare to throw bottle.
		err = globalbottle(
			ctx.Event.UserID,
			ctx.Event.GroupID,
			ctx.Event.Time,
			ctx.CardOrNickName(ctx.Event.UserID),
			relief.Args,
		).throw(seaSide)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("已经帮你丢出去了哦~")))
	})
}

func globalbottle(qq, grp, time int64, name, msg string) *sea { // Check as if the User is available?
	id := int64(crc64.Checksum(binary.StringToBytes(fmt.Sprintf("%d_%d_%d_%s_%s", grp, time, qq, name, msg)), crc64.MakeTable(crc64.ISO)))
	return &sea{ID: id, Grp: grp, Time: time, QQ: qq, Name: name, Msg: msg}
}

func (be *sea) throw(db *sql.Sqlite) error {
	seaLocker.Lock()
	defer seaLocker.Unlock()
	return db.Insert("global", be)
}

func (be *sea) destory(db *sql.Sqlite) error {
	seaLocker.Lock()
	defer seaLocker.Unlock()
	return db.Del("global", "WHERE id="+strconv.FormatInt(be.ID, 10))
}

func fetchBottle(db *sql.Sqlite) (*sea, error) {
	seaLocker.RLock()
	defer seaLocker.RLock()
	be := new(sea)
	return be, db.Find("global", "be", "where content LIKE '%ID%'"+" ORDER BY RANDOM() limit 1")
}

func CreateChannel(db *sql.Sqlite) error {
	seaLocker.Lock()
	defer seaLocker.Unlock()
	return db.Create("global", &sea{})
}
