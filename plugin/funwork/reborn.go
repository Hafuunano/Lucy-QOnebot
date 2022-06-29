// Package funwork 投胎 来自 https://github.com/YukariChiba/tgbot/blob/main/modules/Reborn.py
package funwork

import (
	"fmt"
	"math/rand"

	"encoding/json"
	"github.com/FloatTech/zbputils/file"
	wr "github.com/mroth/weightedrand"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"

	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	areac     *wr.Chooser
	gender, _ = wr.NewChooser(
		wr.Choice{Item: "男孩子", Weight: 33707},
		wr.Choice{Item: "女孩子", Weight: 39292},
		wr.Choice{Item: "雌雄同体", Weight: 1001},
		wr.Choice{Item: "猫猫!", Weight: 10000},
		wr.Choice{Item: "狗狗!", Weight: 10000},
		wr.Choice{Item: "🐉!", Weight: 3000},
		wr.Choice{Item: "龙猫~", Weight: 3000},
	)
)

type ratego []struct {
	Name   string  `json:"name"`
	Weight float64 `json:"weight"`
}

func init() {
	go func() {
		datapath := engine.DataFolder()
		jsonfile := datapath + "ratego.json"
		area := make(ratego, 226)
		err := load(&area, jsonfile)
		if err != nil {
			panic(err)
		}
		choices := make([]wr.Choice, len(area))
		for i, a := range area {
			choices[i].Item = a.Name
			choices[i].Weight = uint(a.Weight * 1e9)
		}
		areac, err = wr.NewChooser(choices...)
		if err != nil {
			panic(err)
		}
		logrus.Printf("[Reborn]读取%d个国家/地区", len(area))
	}()
	engine.OnFullMatchGroup([]string{"reborn", "我要重生", "我要重开"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			if !limit.Load(ctx.Event.GroupID).Acquire() {
				return
			}
			if rand.Int31() > 1<<27 {
				ctx.SendChain(message.At(ctx.Event.UserID), message.Text(fmt.Sprintf("投胎成功！\n您出生在 %s, 是 %s。", randcoun(), randgen())))
			} else {
				ctx.SendChain(message.At(ctx.Event.UserID), message.Text("投胎失败！\n您没能活到出生，祝您下次好运！"))
			}
		})
}

// load 加载rate数据
func load(area *ratego, jsonfile string) error {
	data, err := file.GetLazyData(jsonfile, true)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, area)
}

func randcoun() string {
	return areac.Pick().(string)
}

func randgen() string {
	return gender.Pick().(string)
}
