package arc

import (
	"encoding/json"
	"github.com/FloatTech/zbputils/ctxext"
	aua "github.com/MoYoez/Arcaea_auaAPI"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"github.com/wdvxdr1123/ZeroBot/message"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
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
		} `json:"song_info"`
	} `json:"content"`
}

type arcGPT struct {
	En     []string `json:"en"`
	Ja     []string `json:"ja"`
	ZhHans []string `json:"zh-Hans"`
}

var (
	randGPT        arcGPT
	arcRandSong    randsong
	getArcSongName string
	randLimit      = rate.NewManager[int64](time.Minute*5, 45) // 限制并发
)

func init() {
	arcGPTJson, err := os.ReadFile(engine.DataFolder() + "arcgpt.json")
	if err != nil {
		panic(err)
	}
	engine.OnRegex(`[！!]arc\srandGPT`).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		switch {
		case randLimit.Load(ctx.Event.GroupID).AcquireN(6):
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Sending response to ArcGPT...Hold on("))
			auaRandSongBytes, err := aua.GetSongRandom(os.Getenv("aualink"), os.Getenv("auakey"), strconv.Itoa(0), strconv.Itoa(12))
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Cannot get message from arcgpt(", err))
				return
			}
			err = json.Unmarshal(arcGPTJson, &randGPT)
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Cannot get message from arcgpt(", err))
				return
			}
			err = json.Unmarshal(auaRandSongBytes, &arcRandSong)
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Cannot get message from arcgpt(", err))
				return
			}
			// get length.
			getGPTNums := rand.Intn(len(randGPT.ZhHans))
			// first check the jp name, if not, use eng name.
			getArcSongName = arcRandSong.Content.Songinfo.NameJp
			getSongArtist := arcRandSong.Content.Songinfo.Artist
			if getArcSongName == "" {
				getArcSongName = arcRandSong.Content.Songinfo.NameEn
			}
			// handle texts.
			handledSongGPTtextDone := strings.ReplaceAll(randGPT.ZhHans[getGPTNums], "歌曲名称", getArcSongName)
			handledSongGPTtext := strings.ReplaceAll(handledSongGPTtextDone, "作曲家", getSongArtist)
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(handledSongGPTtext))
		case randLimit.Load(ctx.Event.GroupID).Acquire():
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("暂时无法处理更多请求。\n"))
		default:
		}
	})
}
