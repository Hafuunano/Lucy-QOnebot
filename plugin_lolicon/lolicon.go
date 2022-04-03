// Package lolicon 基于 https://api.lolicon.app 随机图片
package lolicon

import (
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"github.com/wdvxdr1123/ZeroBot/message"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/FloatTech/AnimeAPI/imgpool"
	control "github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/math"
	"github.com/FloatTech/zbputils/process"

	"github.com/FloatTech/ZeroBot-Plugin/order"
)

const (
	api      = "https://api.lolicon.app/setu/v2?proxy=i.pixiv.re"
	capacity = 10
)

var (
	queue = make(chan string, capacity)
)

func init() {
	limit := rate.NewManager(time.Minute, 10)
	control.Register("lolicon", order.PrioLolicon, &control.Options{
		DisableOnDefault: false,
		Help: "lolicon\n" +
			"- 我要一份色图",
	}).OnFullMatch("我要一份色图").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			if !limit.Load(ctx.Event.GroupID).Acquire() {
				ctx.SendChain(message.Text("好孩子不能一次性看太多的涩图哦~注意身体."))
				return
			}
			go func() {
				for i := 0; i < math.Min(cap(queue)-len(queue), 2); i++ {
					resp, err := http.Get(api)
					if err != nil {
						ctx.SendChain(message.Text("ERROR: ", err))
						continue
					}
					if resp.StatusCode != http.StatusOK {
						ctx.SendChain(message.Text("ERROR: code ", resp.StatusCode))
						continue
					}
					data, _ := ioutil.ReadAll(resp.Body)
					resp.Body.Close()
					json := gjson.ParseBytes(data)
					if e := json.Get("error").Str; e != "" {
						ctx.SendChain(message.Text("ERROR: ", e))
						continue
					}
					url := json.Get("data.0.urls.original").Str
					url = strings.ReplaceAll(url, "i.pixiv.cat", "i.pixiv.re")
					name := url[strings.LastIndex(url, "/")+1 : len(url)-4]
					m, err := imgpool.GetImage(name)
					if err != nil {
						m.SetFile(url)
						_, err = m.Push(ctxext.SendToSelf(ctx), ctxext.GetMessage(ctx))
						process.SleepAbout1sTo2s()
					}
					if err == nil {
						queue <- m.String()
					} else {
						queue <- url
					}
				}
			}()
			select {
			case <-time.After(time.Minute):
				ctx.SendChain(message.Text("ERROR: 等待填充，请稍后再试......"))
			case url := <-queue:
				ctx.SendChain(message.Image(url))
			}
		})
}
