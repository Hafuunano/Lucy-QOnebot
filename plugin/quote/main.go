package quote

import (
	"fmt"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
	"hash/crc64"
	"strconv"
	"time"
)

var engine = control.Register("quote", &ctrl.Options[*zero.Ctx]{
	DisableOnDefault:  false,
	Help:              "Make A Quote! 记录群友发言\n说明书: https://lucy.impart.icu",
	PrivateDataFolder: "quote",
})

func init() {
	engine.OnRegex(`^(\[CQ:reply,id=(.*)\])\s/quote$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getPatternUserMessageID := ctx.State["regex_matched"].([]string)[2]
		rsp := ctx.CallAction("get_msg", zero.Params{
			"message_id": getPatternUserMessageID,
		}).Data.String()
		SenderUserID := gjson.Get(rsp, "sender.user_id").Int()
		MessageRaw := gjson.Get(rsp, "message").String()
		formatMessageRaw := message.UnescapeCQCodeText(MessageRaw)
		InsertUserCorruptionTarget(Corruption{
			TrackID:   TrackIDGenerator(ctx.Event.GroupID),
			QQ:        SenderUserID,
			Msg:       formatMessageRaw,
			HandlerQQ: ctx.Event.UserID,
			Time:      time.Now().Unix(),
		}, ctx.Event.GroupID)
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("成功~"))
	})
	engine.OnFullMatch("/quote").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		// make a quote
		getRandGeneratedData := IndexOfDataCorruption(ctx.Event.GroupID)
		if getRandGeneratedData.TrackID == 0 {
			ctx.Send(message.Text("找不到Quote，或许是本群没有x"))
			return
		}
		format := fmt.Sprintf("TrackerID: %d\nTrackerUser: %s(%d)\nTime: %s\n%s(%d): \n%s", getRandGeneratedData.TrackID, ctx.CardOrNickName(getRandGeneratedData.HandlerQQ), getRandGeneratedData.HandlerQQ, time.Unix(getRandGeneratedData.Time, 0).Format("2006-01-02 15:04:05"), ctx.CardOrNickName(getRandGeneratedData.QQ), getRandGeneratedData.QQ, getRandGeneratedData.Msg)
		ctx.Send(format)
	})
}

func TrackIDGenerator(groupID int64) int64 {
	return int64(crc64.Checksum(helper.StringToBytes(strconv.FormatInt(groupID+time.Now().Unix(), 10)), crc64.MakeTable(crc64.ISO)))
}
