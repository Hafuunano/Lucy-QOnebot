// Package simai github.com/FloatTech/Zerobot-Plugin
package simai

import (
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	breaker "github.com/MoYoez/Lucy-QOnebot/box/break"
	"github.com/MoYoez/Lucy-QOnebot/box/setname"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"gopkg.in/yaml.v3"
)

// SimPackData simia Data
type SimPackData struct {
	Proud  map[string][]string `yaml:"傲娇"`
	Kawaii map[string][]string `yaml:"可爱"`
}

func init() {
	engine := control.Register("simai", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		PrivateDataFolder: "simai",
		Help:              "simai - Use simia pre-render dict to make it more clever",
	})
	// onload simia dict
	dictLoaderLocation := engine.DataFolder() + "simai.yml"
	dictLoader, err := os.ReadFile(dictLoaderLocation)
	if err != nil {
		panic(err)
	}
	var data SimPackData
	_ = yaml.Unmarshal(dictLoader, &data)
	onMakeRegex, err := regexp.Compile("[?？!！]")
	if err != nil {
		panic(err)
	}
	engine.OnMessage(zero.OnlyToMe).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		// onload dict.
		msg := ctx.ExtractPlainText()
		var getChartReply []string
		if breaker.GetStringLength(msg) > 50 {
			return
		}
		getList := onMakeRegex.FindAllString(msg, -1)

		if len(getList) >= 1 {
			for _, data := range getList {
				msg = strings.ReplaceAll(msg, data, "")
			}
		} // on ticket saver.
		for dataReply, inner := range data.Proud {
			if msg == dataReply {
				getChartReply = inner
				break
			}
			if strings.Contains(msg, dataReply) && breaker.GetStringLength(dataReply)/breaker.GetStringLength(msg) > 0.45 {
				getChartReply = inner
				break
			}
		}

		if len(getChartReply) == 0 {
			for dataReply, inner := range data.Kawaii {
				if msg == dataReply {
					getChartReply = inner
					break
				}
				if strings.Contains(msg, dataReply) && breaker.GetStringLength(dataReply)/breaker.GetStringLength(msg) > 0.45 {
					getChartReply = inner
					break
				}
			}
		}
		// if no data
		if len(getChartReply) == 0 {
			// no reply
			return
		}

		// show data is existed.
		getReply := getChartReply[rand.Intn(len(getChartReply))]
		getName := setname.LoadUserNickname(strconv.FormatInt(ctx.Event.UserID, 10))
		if getName == "你" {
			getName = ctx.CardOrNickName(ctx.Event.UserID)
		}
		getLucyName := []string{"Lucy", "Lucy酱"}[rand.Intn(2)]
		getReply = strings.ReplaceAll(getReply, "{segment}", " ")
		getReply = strings.ReplaceAll(getReply, "{name}", getName)
		getReply = strings.ReplaceAll(getReply, "{me}", getLucyName)
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(getReply))
	})
}
