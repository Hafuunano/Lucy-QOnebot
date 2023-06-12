package pgr // Package pgr hosted by Phigros-Library
import (
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// too lazy,so this way is to use thrift host server (Working on HiMoYo Cloud.) (replace: now use PUA API)

// update: use PhigrosUnlimitedAPI + Phigros Library as Maintainer.

var (
	engine = control.Register("pgr", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: true,
		Help:             "Hi NekoPachi!\n",
	})
)

func init() {
	engine.OnRegex(`/pgr bind\s+(\w+)\s+(\w+)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		hash := ctx.State["regex_matched"].([]string)[1]
		md5 := ctx.State["regex_matched"].([]string)[2]
		index, err := decryptAES(hash, md5)
		if err != nil {
			panic(err)
		}
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(index))
	})

}
