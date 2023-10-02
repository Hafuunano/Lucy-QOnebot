// Package funwork music 网易云 点歌 Powered By LemonKoi Vercel Services.
package funwork

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/FloatTech/floatbox/web"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var limitForMusic = rate.NewManager[int64](time.Minute*3, 8)

func init() {
	engine.OnRegex(`^点歌\s?(.{1,25})$`).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			if !limitForMusic.Load(ctx.Event.GroupID).Acquire() {
				ctx.SendChain(message.Text("太快了哦，麻烦慢一点~"))
				return
			}
			requestURL := "https://nem.lemonkoi.one/search?limit=1&type=1&keywords=" + url.QueryEscape(ctx.State["regex_matched"].([]string)[1])
			data, err := web.GetData(requestURL)
			var webStatusCode *http.Response
			if err != nil {
				fmt.Print("ERR:", webStatusCode.StatusCode)
			}
			ctx.SendChain(message.Music("163", gjson.ParseBytes(data).Get("result.songs.0.id").Int()))
		})
}
