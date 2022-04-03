// Package shadiao 来源于 https://shadiao.app/# 的接口
package shadiao

import (
	"time"

	control "github.com/FloatTech/zbputils/control"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"

	"github.com/FloatTech/ZeroBot-Plugin/order"
)

const (
	chpURL         = "https://chp.shadiao.app/api.php"
	pyqURL         = "https://pyq.shadiao.app/api.php"
	yduanziURL     = "http://www.yduanzi.com/duanzi/getduanzi"
	todayURL       = "http://www.ipip5.com/today/api.php?type=txt"
	yiyanURL       = "https://v1.hitokoto.cn/?encode=text"
	ua             = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36"
	uas            = "Mozilla/5.0 (Linux; Android 10; MIX 2S) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.98 Mobile Safari/537.36"
	chpReferer     = "https://chp.shadiao.app/"
	pyqReferer     = "https://pyq.shadiao.app/"
	yduanziReferer = "http://www.yduanzi.com/?utm_source=shadiao.app"
	todayReferer   = "http://www.ipip5.com/"
	yiyanReferer   = "http://baidu.com/"
)

var (
	engine = control.Register("shadiao", order.PrioShaDiao, &control.Options{
		DisableOnDefault: false,
		Help: " 沙雕app + OvOoa API \n" +
			"- 哄我\n- - 讲个段子 +....addons",
	})
	limit = rate.NewManager(time.Minute, 60)
)
