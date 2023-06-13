package pgr // Package pgr hosted by Phigros-Library
import (
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"time"
)

// too lazy,so this way is to use thrift host server (Working on HiMoYo Cloud.) (replace: now use PUA API)

// update: use PhigrosUnlimitedAPI + Phigros Library as Maintainer.

var (
	engine = control.Register("pgr", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "Hi NekoPachi!\n",
		PrivateDataFolder: "phi",
	})
)

func init() {
	engine.OnRegex(`^\/pgr\sbind\s(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		hash := ctx.State["regex_matched"].([]string)[1]
		userInfo := GetUserInfoFromDatabase(ctx.Event.UserID)
		if userInfo.Time+(12*60*60) > time.Now().Unix() {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("12小时内仅允许绑定一次哦"))
			return
		}
		indexReply := DecHashToRaw(hash)
		if indexReply == "" {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("请前往 https://pgr.impart.icu 获取绑定码进行绑定 "))
			return
		}
		// get session.
		getQQID, getSessionID := RawJsonParse(indexReply)
		if getQQID != ctx.Event.UserID {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("请求Hash中QQ号不一致，请使用自己的号重新申请"))
			return
		}
		_ = FormatUserDataBase(getQQID, getSessionID, time.Now().Unix()).BindUserDataBase()
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("绑定成功～"))
	})
	engine.OnFullMatch("/pgr b19").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		data := GetRequestInfo(ctx.Event.UserID)
		getFormatLenth := len(data.Content.BestList.Best)
	})
}
