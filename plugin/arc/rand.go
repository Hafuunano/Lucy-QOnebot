package arc

import (
	"encoding/json"
	"github.com/FloatTech/zbputils/ctxext"
	aua "github.com/MoYoez/Go-ArcaeaUnlimitedAPI"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type randsong struct {
	Status  int `json:"status"`
	Content struct {
		Id          string `json:"id"`
		RatingClass int    `json:"ratingClass"`
		Songinfo    struct {
			NameEn         string  `json:"name_en"`
			NameJp         string  `json:"name_jp"`
			Artist         string  `json:"artist"`
			Bpm            string  `json:"bpm"`
			BpmBase        float64 `json:"bpm_base"`
			Set            string  `json:"set"`
			SetFriendly    string  `json:"set_friendly"`
			Time           int     `json:"time"`
			Side           int     `json:"side"`
			WorldUnlock    bool    `json:"world_unlock"`
			RemoteDownload bool    `json:"remote_download"`
			Bg             string  `json:"bg"`
			Date           int     `json:"date"`
			Version        string  `json:"version"`
			Difficulty     int     `json:"difficulty"`
			Rating         int     `json:"rating"`
			Note           int     `json:"note"`
			ChartDesigner  string  `json:"chart_designer"`
			JacketDesigner string  `json:"jacket_designer"`
			JacketOverride bool    `json:"jacket_override"`
			AudioOverride  bool    `json:"audio_override"`
		} `json:"songinfo"`
	} `json:"content"`
}

type arcGPT struct {
	En     []string `json:"en"`
	Ja     []string `json:"ja"`
	ZhHans []string `json:"zh-Hans"`
}

var (
	randgpt     arcGPT
	arcrandsong randsong
	getArcName  string
)

func init() {
	arcGPTJson, err := os.ReadFile(engine.DataFolder() + "arcgpt.json")
	if err != nil {
		panic(err)
	}
	engine.OnFullMatch("!test arc randgpt").SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Sending response to ArcGPT...Hold on("))
		auaRandSong, err := aua.GetSongRandom(os.Getenv("aualink"), os.Getenv("auakey"), strconv.Itoa(0), strconv.Itoa(12))
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Cannot get message from arcgpt(", err))
			return
		}
		auaRandSongBytes := helper.StringToBytes(auaRandSong)
		err = json.Unmarshal(arcGPTJson, &randgpt)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Cannot get message from arcgpt(", err))
			return
		}
		err = json.Unmarshal(auaRandSongBytes, &arcrandsong)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Cannot get message from arcgpt(", err))
			return
		}
		// get length.
		getGPTNums := rand.Intn(len(randgpt.ZhHans))
		// first check the jp name, if not, use eng name.
		getArcName = arcrandsong.Content.Songinfo.NameJp
		getSongArtist := arcrandsong.Content.Songinfo.Artist
		if getArcName == "" {
			getArcName = arcrandsong.Content.Songinfo.NameEn
		}
		// handle texts.
		handledSongGPTtextDone := strings.ReplaceAll(randgpt.ZhHans[getGPTNums], "歌曲名称", getArcName)
		handledSongGPTtext := strings.ReplaceAll(handledSongGPTtextDone, "作曲家", getSongArtist)
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(handledSongGPTtext))
	})

}
