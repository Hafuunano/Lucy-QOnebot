// Package funwork_saucenao P站ID/saucenao/ascii2d搜图
package funwork

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	"github.com/FloatTech/AnimeAPI/pixiv"
	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/img/pool"
	"github.com/jozsefsallai/gophersauce"
)

var (
	saucenaocli *gophersauce.Client
)

func init() { // 插件主体
	apikeyfile := engine.DataFolder() + "apikey.txt"
	if file.IsExist(apikeyfile) {
		key, err := os.ReadFile(apikeyfile)
		if err != nil {
			panic(err)
		}
		saucenaocli, err = gophersauce.NewClient(&gophersauce.Settings{
			MaxResults: 1,
			APIKey:     binary.BytesToString(key),
		})
		if err != nil {
			panic(err)
		}
	}
	// 根据 PID 搜图
	engine.OnRegex(`^搜图(\d+)$`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			if !limit.Load(ctx.Event.GroupID).Acquire() {
				ctx.SendChain(message.Text("触发保护机制 请过会使用哦."))
				return
			}
			id, _ := strconv.ParseInt(ctx.State["regex_matched"].([]string)[1], 10, 64)
			ctx.SendChain(message.Text("Lucy正在查询..."))
			// 获取P站插图信息
			illust, err := pixiv.Works(id)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			if illust.Pid > 0 {
				name := strconv.FormatInt(illust.Pid, 10)
				var imgs message.Message
				for i := range illust.ImageUrls {
					f := file.BOTPATH + "/" + illust.Path(i)
					n := name + "_p" + strconv.Itoa(i)
					var m *pool.Image
					if file.IsNotExist(f) {
						m, err = pool.GetImage(n)
						if err == nil {
							imgs = append(imgs, message.Image(m.String()))
							continue
						}
						logrus.Debugln("[sausenao]开始下载", n)
						logrus.Debugln("[sausenao]urls:", illust.ImageUrls)
						err1 := illust.DownloadToCache(i)
						if err1 == nil {
							m.SetFile(f)
							_, _ = m.Push(ctxext.SendToSelf(ctx), ctxext.GetMessage(ctx))
						}
						if err1 != nil {
							logrus.Debugln("[sausenao]下载err:", err1)
						}
					}
					imgs = append(imgs, message.Image("file:///"+f))
				}
				txt := message.Text(
					"标题: ", illust.Title, "\n",
					"插画ID: ", illust.Pid, "\n",
					"画师: ", illust.UserName, "\n",
					"Tags: ", illust.Tags, "\n",
					"直链: ", "https://pixivel.moe/detail?id=", illust.Pid,
				)
				if imgs != nil {
					if zero.OnlyGroup(ctx) {
						ctx.SendGroupForwardMessage(ctx.Event.GroupID, message.Message{
							ctxext.FakeSenderForwardNode(ctx, txt),
							ctxext.FakeSenderForwardNode(ctx, imgs...),
						})
					} else {
						// 发送搜索结果
						ctx.Send(append(imgs, message.Text("\n"), txt))
					}
				} else {
					// 图片下载失败，仅发送文字结果
					ctx.SendChain(txt)
				}
			} else {
				ctx.SendChain(message.Text("图片不存在!"))
			}
		})
	// 以图搜图
	engine.OnKeywordGroup([]string{"以图搜图", "搜索图片", "以图识图"}, zero.OnlyGroup, zero.MustProvidePicture).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			if !limit.Load(ctx.Event.GroupID).Acquire() {
				ctx.SendChain(message.Text("触发保护机制 请过会使用哦."))
				return
			}
			// 开始搜索图片
			ctx.SendChain(message.Text("Lucy正在查询..."))
			for _, pic := range ctx.State["image_url"].([]string) {
				if saucenaocli != nil {
					resp, err := saucenaocli.FromURL(pic)
					if err == nil && resp.Count() > 0 {
						result := resp.First()
						s, err := strconv.ParseFloat(result.Header.Similarity, 64)
						if err == nil {
							rr := reflect.ValueOf(&result.Data).Elem()
							b := binary.NewWriterF(func(w *binary.Writer) {
								r := rr.Type()
								for i := 0; i < r.NumField(); i++ {
									if !rr.Field(i).IsZero() {
										w.WriteString("\n")
										w.WriteString(r.Field(i).Name)
										w.WriteString(": ")
										w.WriteString(fmt.Sprint(rr.Field(i).Interface()))
									}
								}
							})
							resp, err := http.Head(result.Header.Thumbnail)
							msg := make(message.Message, 0, 3)
							if s > 80.0 {
								msg = append(msg, message.Text("我有把握是这个!"))
							} else {
								msg = append(msg, message.Text("也许是这个?"))
							}
							if err == nil && resp.StatusCode == http.StatusOK {
								msg = append(msg, message.Image(result.Header.Thumbnail))
							} else {
								msg = append(msg, message.Image(pic))
							}
							msg = append(msg, message.Text("\n图源: ", result.Header.IndexName, binary.BytesToString(b)))
							id := ctx.SendGroupForwardMessage(
								ctx.Event.GroupID,
								message.Message{ctxext.FakeSenderForwardNode(ctx, msg...)},
							).Get("message_id").Int()
							time.Sleep(time.Second * 20)
							ctx.DeleteMessage(message.NewMessageIDFromInteger(id))
							if s > 80.0 {
								continue
							}
							defer resp.Body.Close()
						}
					}
				} else {
					return
				}
			}
		})
}
