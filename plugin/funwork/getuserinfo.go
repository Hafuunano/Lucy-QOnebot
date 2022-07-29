package funwork

import (
	"strconv"
	"time"

	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// debug
var limitinfo = rate.NewManager[int64](time.Minute*5, 1)

func init() {
	engine.OnRegex(`^查找信息(\d+)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		if !limitinfo.Load(ctx.Event.UserID).Acquire() {
			return
		}
		reachinfo, _ := strconv.ParseInt(ctx.State["regex_matched"].([]string)[1], 10, 64)
		getUserInfo := ctx.GetGroupMemberInfo(ctx.Event.GroupID, reachinfo, false)
		tempUserInfo := getUserInfo.String()
		userName := ctx.CardOrNickName(reachinfo)
		userSexInfo := gjson.Get(tempUserInfo, "sex")
		userJoinTimeUnix := gjson.Get(tempUserInfo, "join_time").Int()
		userJoinTime := time.Unix(userJoinTimeUnix, 0).Format("2006-01-02 03:04:05 PM")
		lastSendTimeUnix := gjson.Get(tempUserInfo, "last_sent_time").Int()
		userLastSendTIme := time.Unix(lastSendTimeUnix, 0).Format("2006-01-02 03:04:05 PM")
		userUnfriendly := gjson.Get(tempUserInfo, "unfriendly")
		userHonorTitle := gjson.Get(tempUserInfo, "title")
		// userArea := gjson.Get(tempUserInfo, "area")
		//	userAge := ctx.Event.Sender.Age
		// userLevel := gjson.Get(tempUserInfo, "level")
		ctx.SendChain(message.Text("你查询的人为: ", userName, "\n性别:", userSexInfo, "\n最后一次发送信息时间 :", userLastSendTIme, "\n加入时间:", userJoinTime, "\n是否有不友好记录:", userUnfriendly, "\n头衔:", userHonorTitle))
	})
}
