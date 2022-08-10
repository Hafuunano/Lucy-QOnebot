// Package partygame 真心话大冒险
package funwork

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type tdata struct {
	Version string `json:"version"`
	Data    []struct {
		Name string   `json:"name"`
		Des  string   `json:"des"`
		Tags []string `json:"tags"`
	} `json:"data"`
}

var (
	action    tdata
	question  tdata
	punishmap sync.Map
)

func init() { // 插件主体
	actionData, err := os.ReadFile(engine.DataFolder() + "action.json")
	if err != nil {
		panic(err)
	}
	questionData, err := os.ReadFile(engine.DataFolder() + "question.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(actionData, &action)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(questionData, &question)
	if err != nil {
		panic(err)
	}
	engine.OnRegex(`^(真心话大冒险)`).Handle(func(ctx *zero.Ctx) {
		key := fmt.Sprintf("%v-%v", ctx.Event.GroupID, ctx.Event.UserID)
		v, ok := punishmap.Load(key)
		if ok {
			ctx.SendChain(message.At(ctx.Event.UserID), message.Text("\n咕噜咕噜~Lucy可不满意,请做以下的方式进行赎罪哦 ʕ·ᴥ·ʔ ", v))
			return
		}
		puid := ctx.Event.UserID
		if len(ctx.Event.Message) > 1 && ctx.Event.Message[1].Type == "at" {
			puid, _ = strconv.ParseInt(ctx.Event.Message[1].Data["qq"], 10, 64)
		}
		ctx.Event.UserID = puid
		_ = getTruthOrDare(ctx)
	})
	engine.OnRegex(`^(饶恕|原谅)`, zero.AdminPermission, zero.OnlyGroup, zero.OnlyToMe).Handle(func(ctx *zero.Ctx) {
		puid := ctx.Event.UserID
		if len(ctx.Event.Message) > 1 && ctx.Event.Message[1].Type == "at" {
			puid, _ = strconv.ParseInt(ctx.Event.Message[1].Data["qq"], 10, 64)
		}
		key := fmt.Sprintf("%v-%v", ctx.Event.GroupID, puid)
		punishmap.Delete(key)
		ctx.SendChain(message.At(puid), message.Text("恭喜你~得到了Lucy的原谅 ʕ·ᴥ·ʔ "))
	})
	engine.OnRegex(`^(惩罚)`, zero.AdminPermission, zero.OnlyGroup, zero.OnlyToMe).Handle(func(ctx *zero.Ctx) {
		puid := ctx.Event.UserID
		if len(ctx.Event.Message) > 1 && ctx.Event.Message[1].Type == "at" {
			puid, _ = strconv.ParseInt(ctx.Event.Message[1].Data["qq"], 10, 64)
		}
		ctx.Event.UserID = puid
		_ = getTruthOrDare(ctx)
	})
	engine.OnRegex(`^(反省|检查罪行)`, zero.OnlyGroup, zero.OnlyToMe).Handle(func(ctx *zero.Ctx) {
		puid := ctx.Event.UserID
		if len(ctx.Event.Message) > 1 && ctx.Event.Message[1].Type == "at" {
			puid, _ = strconv.ParseInt(ctx.Event.Message[1].Data["qq"], 10, 64)
		}
		key := fmt.Sprintf("%v-%v", ctx.Event.GroupID, puid)
		v, ok := punishmap.Load(key)
		if ok {
			ctx.SendChain(message.At(puid), message.Text("哼~你可以通过这样的方式进行赎罪哦~w", v))
		}
	})
}

func getAction() string {
	return action.Data[rand.Intn(len(action.Data))].Name
}

func getQuestion() string {
	return question.Data[rand.Intn(len(question.Data))].Name
}

func getActionOrQuestion() string {
	if time.Now().UnixNano()%2 == 0 {
		return getAction()
	}
	return getQuestion()
}

func getTruthOrDare(ctx *zero.Ctx) bool {
	next := zero.NewFutureEvent("message", 999, false, ctx.CheckSession(), zero.FullMatchRule("真心话", "大冒险")).Next()
	key := fmt.Sprintf("%v-%v", ctx.Event.GroupID, ctx.Event.UserID)
	ctx.SendChain(message.At(ctx.Event.UserID), message.Text("请选择你要接受的挑战呢~ 真心话 or 大冒险"))
	for {
		select {
		case <-time.After(time.Second * 20):
			ctx.SendChain(message.Text("时间太久啦！", zero.BotConfig.NickName[0], "帮你选择"))
			p := getActionOrQuestion()
			punishmap.Store(key, p)
			ctx.SendChain(message.At(ctx.Event.UserID), message.Text("恭喜你获得\"", p, "\"的惩罚"))
			return true
		case c := <-next:
			msg := c.Event.Message.ExtractPlainText()
			var p string
			if msg == "真心话" {
				p = getQuestion()
			} else if msg == "大冒险" {
				p = getAction()
			}
			punishmap.Store(key, p)
			ctx.SendChain(message.At(ctx.Event.UserID), message.Text("恭喜你获得\"", p, "\"的惩罚"))
			return true
		}
	}
}
