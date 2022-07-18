package arknights

import (
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/img/text"
	"github.com/fogleman/gg"
	zero "github.com/wdvxdr1123/ZeroBot"
)

func init() {
	Fonts, _ = gg.LoadFontFace(text.FontFile, 18)
	engine := control.Register("arknight", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		//	PublicDataFolder: "ArkNights",
		Help: "查公招\n" + "方舟今日资源",
	})
	engine.OnRegex("^查公招$").Handle(recruit)
	engine.OnKeyword("方舟资源").Handle(daily)
}
