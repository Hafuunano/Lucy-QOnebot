package arc

import (
	"encoding/json"
	"github.com/FloatTech/floatbox/web"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"strconv"
)

// make struct list.
type ranklist struct {
	Success bool `json:"success"`
	Value   []struct {
		SongId string `json:"song_id"`
		Title  struct {
			Ja string `json:"ja"`
			En string `json:"en"`
		} `json:"title"`
		Artist string `json:"artist"`
		Rank   int    `json:"rank"`
		Status int    `json:"status"`
		Bg     string `json:"bg"`
	} `json:"value"`
}

var rawrank ranklist

func init() {
	engine.OnRegex(`[!|！]arc\srank\s(paid|free)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		// first of all, request data from arc.
		regexPatterns := ctx.State["regex_matched"].([]string)[1]
		rawData, err := web.RequestDataWith(web.NewDefaultClient(), "https://webapi.lowiro.com/webapi/song/rank/"+regexPatterns, "GET", "arcaea.lowiro.com", web.RandUA(), nil)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR:", err))
			return
		}
		// handle raw data and do some format.
		_ = json.Unmarshal(rawData, &rawrank)
		// show only ten. no more than it (
		var originalText string
		for i := 0; i < 10; i++ {
			appendText := rawrank.Value[i].Title.Ja
			num := i + 1
			originalText += strconv.Itoa(num) + " " + appendText + "\n"
		}
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(originalText))
	})
}
