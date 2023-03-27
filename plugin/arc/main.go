// Package arc for arc render b30
package arc

import (
	"encoding/json"
	"os"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	aua "github.com/MoYoez/Go-ArcaeaUnlimitedAPI"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

type arcaea struct {
	Status  int `json:"status"`
	Content struct {
		Best30Avg   float64 `json:"best30_avg"`
		Recent10Avg float64 `json:"recent10_avg"`
		AccountInfo struct {
			Code                   string `json:"code"`
			Name                   string `json:"name"`
			UserID                 int    `json:"user_id"`
			IsMutual               bool   `json:"is_mutual"`
			IsCharUncappedOverride bool   `json:"is_char_uncapped_override"`
			IsCharUncapped         bool   `json:"is_char_uncapped"`
			IsSkillSealed          bool   `json:"is_skill_sealed"`
			Rating                 int    `json:"rating"`
			JoinDate               int64  `json:"join_date"`
			Character              int    `json:"character"`
		} `json:"account_info"`
		Best30List []struct {
			Score             int     `json:"score"`
			Health            int     `json:"health"`
			Rating            float64 `json:"rating"`
			SongID            string  `json:"song_id"`
			Modifier          int     `json:"modifier"`
			Difficulty        int     `json:"difficulty"`
			ClearType         int     `json:"clear_type"`
			BestClearType     int     `json:"best_clear_type"`
			TimePlayed        int64   `json:"time_played"`
			NearCount         int     `json:"near_count"`
			MissCount         int     `json:"miss_count"`
			PerfectCount      int     `json:"perfect_count"`
			ShinyPerfectCount int     `json:"shiny_perfect_count"`
		} `json:"best30_list"`
		Best30Overflow []struct {
			Score             int     `json:"score"`
			Health            int     `json:"health"`
			Rating            float64 `json:"rating"`
			SongID            string  `json:"song_id"`
			Modifier          int     `json:"modifier"`
			Difficulty        int     `json:"difficulty"`
			ClearType         int     `json:"clear_type"`
			BestClearType     int     `json:"best_clear_type"`
			TimePlayed        int64   `json:"time_played"`
			NearCount         int     `json:"near_count"`
			MissCount         int     `json:"miss_count"`
			PerfectCount      int     `json:"perfect_count"`
			ShinyPerfectCount int     `json:"shiny_perfect_count"`
		} `json:"best30_overflow"`
		Best30Songinfo []struct {
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
		} `json:"best30_songinfo"`
		Best30OverflowSonginfo []struct {
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
		} `json:"best30_overflow_songinfo"`
	} `json:"content"`
}

var (
	engine = control.Register("arcaea", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "Hi NekoPachi!\n说明书: https://manual-lucy.himoyo.cn",
		PrivateDataFolder: "arcaea",
	})
)

func init() {
	// no database support yet due to i am lazy.
	engine.OnRegex(`^!test arc\s*(\d+)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		id := ctx.State["regex_matched"].([]string)[1]
		// get token
		// get player info
		playerdata, err := aua.Best30(os.Getenv("aualink"), os.Getenv("auakey"), id)
		if err != nil {
			ctx.SendChain(message.Text("获取玩家信息失败"))
			return
		}
		playerdataByte := helper.StringToBytes(playerdata)
		var r arcaea
		_ = json.Unmarshal(playerdataByte, &r)
		// get first song
		ctx.SendChain(message.Text("test\n"), message.Text(r.Content.Best30List[0].SongID))
	})
}
