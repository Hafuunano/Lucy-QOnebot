package nsfw

import (
	"math/rand"
	"time"

	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	api = "https://api.lolicon.app/setu/v2"
)

var (
	engine = control.Register("nsfw", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: true,
		Help:             "Hi NekoPachi!\n",
	})

	limit = rate.NewManager[int64](time.Minute*3, 8)
)

func init() {
	engine.OnFullMatch("涩涩", zero.OnlyToMe).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		if !limit.Load(ctx.Event.UserID).Acquire() {
			return
		}
		if rand.Intn(2) == 1 {
			data, err := web.GetData(api)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			picURL := gjson.Get(string(data), "data.0.urls.original").String()
			messageID := ctx.SendChain(message.Text(picURL))
			time.Sleep(time.Second * 20)
			ctx.DeleteMessage(messageID)
		} else {
			ctx.Send(message.Text([]string{"看什么看！咱没有涩图 哼!", "只有笨蛋才看涩图", "好孩子是不会看涩图的", "敲~笨蛋 不许色色", "为什么要看涩图呢?"}[rand.Intn(5)]))
		}
	})
}
