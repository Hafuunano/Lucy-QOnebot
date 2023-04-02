// Package arc for arc render b30
package arc

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/jpeg"
	"os"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	aua "github.com/MoYoez/Go-ArcaeaUnlimitedAPI"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

var (
	r      arcaea
	engine = control.Register("arcaea", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "Hi NekoPachi!\n说明书: https://lucy.impart.icu",
		PrivateDataFolder: "arcaea",
	})
)

func init() {
	// arc b30 is still in test(
	engine.OnRegex(`^!test arc\s*(\d+)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		id := ctx.State["regex_matched"].([]string)[1]
		// get player info
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Ok, trying to get "+id+" data."))
		playerdata, err := aua.Best30(os.Getenv("aualink"), os.Getenv("auakey"), id)
		if err != nil {
			ctx.SendChain(message.Text("cannot get user data."))
			return
		}
		playerdataByte := helper.StringToBytes(playerdata)
		_ = json.Unmarshal(playerdataByte, &r)
		checkStatus := r.Status
		if checkStatus != 0 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Status code is not valid."))
			return
		}
		mainBG, _ := os.ReadFile(arcaeaRes + "/resource/b30/B30.png")
		mainBGDecoded, _, _ := image.Decode(bytes.NewReader(mainBG))
		basicBG := DrawMainUserB30(mainBGDecoded, r)
		tureResult := FinishedFullB30(basicBG, r)
		var buf bytes.Buffer
		err = jpeg.Encode(&buf, tureResult, nil)
		if err != nil {
			panic(err)
		}
		base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image("base64://"+base64Str))
	})

	engine.OnRegex(`!test\sarc\sbind\s(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {

	})
}
