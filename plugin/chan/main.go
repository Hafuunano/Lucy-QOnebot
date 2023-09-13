package _chan

import (
	"encoding/json"
	"github.com/FloatTech/ZeroBot-Plugin/compounds/name"
	"github.com/FloatTech/floatbox/process"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"github.com/wdvxdr1123/ZeroBot/message"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type kimo = map[string]*[]string

var (
	limit  = rate.NewManager[int64](time.Minute*3, 28) // 回复限制
	engine = control.Register("chan", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: true,
		Help:             "chan -- Easy list talk.",
	})
)

func init() {
	go func() {
		data, err := os.ReadFile(engine.DataFolder() + "kimoi.json")
		if err != nil {
			panic(err)
		}
		kimomap := make(kimo, 256)
		err = json.Unmarshal(data, &kimomap)
		if err != nil {
			panic(err)
		}
		chatList := make([]string, 0, 256)
		for k := range kimomap {
			chatList = append(chatList, k)
		}
		engine.OnFullMatchGroup(chatList, zero.OnlyToMe).SetBlock(true).Handle(
			func(ctx *zero.Ctx) {
				switch {
				case limit.Load(ctx.Event.UserID).AcquireN(3):
					key := ctx.MessageString()
					val := *kimomap[key]
					text := val[rand.Intn(len(val))]
					userID := strconv.FormatInt(ctx.Event.UserID, 10)
					userNickName := name.LoadUserNickname(userID)
					result := strings.ReplaceAll(text, "你", userNickName)
					process.SleepAbout1sTo2s()
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(result)) // 来自于 https://github.com/Kyomotoi/AnimeThesaurus 的回复 经过二次修改
				case limit.Load(ctx.Event.UserID).Acquire():
					process.SleepAbout1sTo2s()
					ctx.Send(message.Text("咱不想说话~好累qwq"))
					return
				default:
				}
			})
	}()
}

func GetTiredToken(ctx *zero.Ctx) float64 {
	return limit.Load(ctx.Event.UserID).Tokens()
}

func GetCostTiredToken(ctx *zero.Ctx) bool {
	return limit.Load(ctx.Event.UserID).AcquireN(3)
}
