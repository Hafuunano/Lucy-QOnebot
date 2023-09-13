// Package simia From github.com/FloatTech/Zerobot-Plugin
package simai

import (
	"github.com/FloatTech/ZeroBot-Plugin/compounds/name"
	"github.com/FloatTech/ZeroBot-Plugin/plugin/chat"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"gopkg.in/yaml.v3"
	"math/rand"
	"os"
	"strconv"
	"strings"
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

	engine.OnMessage(zero.OnlyToMe).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		// onload dict.
		msg := ctx.ExtractPlainText()
		var getChartReply []string
		if chat.GetPokeToken(ctx) < 4 {
			getChartReply = data.Proud[msg]
			// if no data
			if getChartReply == nil {
				getChartReply = data.Kawaii[msg]
				if getChartReply == nil {
					// no reply
					return
				}
			}
		} else {
			getChartReply = data.Kawaii[msg]
			// if no data
			if getChartReply == nil {
				getChartReply = data.Proud[msg]
				if getChartReply == nil {
					// no reply
					return
				}
			}
		}
		// Lucy may more pround when poke too much ^^.
		if chat.GetTiredToken(ctx) < 4 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("咱不想说话 好累awww"))
			return
		} else {
			chat.GetCostTiredToken(ctx)
		}
		// show data is existed.
		getReply := getChartReply[rand.Intn(len(getChartReply))]
		getName := name.LoadUserNickname(strconv.FormatInt(ctx.Event.UserID, 10))
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
