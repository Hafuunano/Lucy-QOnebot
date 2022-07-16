package funwork

import (
	"encoding/json"
	"errors"
	"net/url"
	"time"

	"fmt"
	"strconv"

	"github.com/FloatTech/zbputils/web"

	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const apiprefix = "https://status.himoyo.cn/api/run/"

type wtf struct {
	name string
	path string
}

var table = [...]*wtf{
	{"龙世界", "aj70Uq"},
	{"希望你没有看到……~ヽ(*。>Д<)o゜", "MiGAyV"},
	{"你的意义是什么?", "mRIFuS"},
	{"【ABO】性別和信息素", "KXyy9"},
	{"测测cp", "ZoGXQd"},
	{"在JOJO世界 你的替身会是什么 ", "lj0a8o"},
	{"稱號產生器", "titlegen"},
	{"成分报告", "2PCeo1"},
	{"測驗你跟你的朋友是攻/受", "LkQXO3"},
	{"人設生產器", "LBtPu5"},
	{"測驗你在ABO世界的訊息素", "SwmdU"},
	{"測測你和哪位名人相似？", "RHQeXu"},
	{"S/M测试", "Ga47oZ"},
	{"测测你是谁", "aV1AEi"},
	{"取個綽號吧", "LTkyUy"},
	{"什麼都不是", "vyrSCb"},
	{"今天中午吃什麼", "LdS4K6"},
	{"測試你的中二稱號", "LwUmQ6"},
	{"神奇海螺", "Lon1h7"},
	{"ABO測試", "H1Tgd"},
	{"女主角姓名產生器", "MsQBTd"},
	{"您是什么人", "49PwSd"},
	{"如果你成为了干员", "ok5e7n"},
	{"abo人设生成~", "Di8enA"},
	{"小說大綱生產器", "Lnstjz"},
	{"他会喜欢你吗？", "pezX3a"},
	{"抽签！你明年的今天会干什么", "IF31kS"},
	{"如果你是受，會是哪種受呢？", "Dr6zpF"},
	{"cp文梗", "vEO2KD"},
	{"您是什么人？", "TQ5qyl"},
	{"你成為......的機率", "g0uoBL"},
	{"ABO性別與信息素", "KFPju"},
	{"異國名稱產生器(國家、人名、星球...)", "OBpu4"},
	{"對方到底喜不喜歡你", "JSLoZC"},
	{"○○喜歡你嗎？", "S6Uceo"},
	{"测测你的sm属性", "dOtcO5"},
	{"你/妳究竟是攻還是受呢?", "RXALH"},
	{"神秘藏书阁", "tDRyET"},
	{"十年后 你cp的结局是", "VUwnXQ"},
	{"高维宇宙与常数的你", "6Zql97"},
	{"色色的東東", "o2eg74"},
	{"文章標題產生器", "Ky25WO"},
	{"你的成績怎麼樣", "6kZv69"},
	{"智能SM偵測器ヾ(*ΦωΦ)ツ", "9pY6HQ"},
	{"你的使用注意事項", "La4Gir"},
	{"戀愛指數", "Jsgz0"},
	{"他對你...", "CJxHMf"},
	{"你的明日方舟人际关系", "u5z4Mw"},
	{"日本姓氏產生器", "JJ5Ctb"},
	{"當你轉生到了異世界，你將成為...", "FTpwK"},
	{"未來男朋友", "F3dSV"},
	{"ABO與信息素", "KFOGA"},
	{"你必將就這樣一事無成啊アホ", "RWw9oX"},
	{"用習慣舉手的方式測試你的戀愛運!<3", "wv5bzA"},
	{"攻受", "RaKmY"},
	{"你和你喜歡的人的微h寵溺段子XD", "LdQqGz"},
	{"我的藝名", "LBaTx"},
	{"你是什麼神？", "LqZORE"},
	{"你的起源是什麼？", "HXWwC"},
	{"測你喜歡什麼", "Sue5g2"},
	{"看看朋友的秘密", "PgKb8r"},
	{"你在動漫裡的名字", "Lz82V7"},
	{"小說男角名字產生器", "LyGDRr"},
	{"測試短文", "S48yA"},
	{"創造小故事", "Kjy3AS"},
	{"你的另外一個名字", "LuyYQA"},
	{"與你最匹配的攻君屬性 ", "I7pxy"},
	{"英文全名生產器(女)", "HcYbq"},
	{"BL文章生產器", "LBZMO"},
	{"輕小說書名產生器", "NFucA"},
	{"日本名字產生器（女孩子）", "JRiKv"},
	{"中二技能名產生器", "Ky1BA"},
	{"抽籤", "XqxfuH"},
	{"你的蘿莉控程度全國排名", "IIWh9k"},
	{"假设你是日漫男性角色 你的cv是", "O8Nn3"},
	{"欢迎光临异兽宠物商店", "QNrNL"},
	{"真心告白大冒险ww", "LIXOHJ"},
}

func newWtf(index int) *wtf {
	if index >= 0 && index < len(table) {
		return table[index]
	}
	return nil
}

type resultwtf struct {
	Text string `json:"text"`
	// Path string `json:"path"`
	Ok  bool   `json:"ok"`
	Msg string `json:"msg"`
}

func (w *wtf) predict(names ...string) (string, error) {
	name := ""
	for _, n := range names {
		name += "/" + url.QueryEscape(n)
	}
	u := apiprefix + w.path + name
	r, err := web.GetData(u)
	if err != nil {
		return "", err
	}
	re := new(resultwtf)
	err = json.Unmarshal(r, re)
	if err != nil {
		return "", err
	}
	if re.Ok {
		return "> " + w.name + "\n" + re.Text, nil
	}
	return "", errors.New(re.Msg)
}

func init() {
	engine.OnFullMatch("鬼东西列表").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			s := ""
			for i, w := range table {
				s += fmt.Sprintf("%02d. %s\n", i, w.name)
			}
			ctx.SendChain(message.Text(s))
		})
	engine.OnRegex(`^鬼东西(\d*)`, zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			// 调用接口
			i, err := strconv.Atoi(ctx.State["regex_matched"].([]string)[1])
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			w := newWtf(i)
			if w == nil {
				ctx.SendChain(message.Text("没有这项内容！"))
				return
			}
			// 获取名字
			var name string
			var secondname string
			if len(ctx.Event.Message) > 1 && ctx.Event.Message[1].Type == "at" {
				qq, _ := strconv.ParseInt(ctx.Event.Message[1].Data["qq"], 10, 64)
				secondname = ctx.GetGroupMemberInfo(ctx.Event.GroupID, qq, false).Get("nickname").Str
			}
			name = ctx.Event.Sender.NickName
			var text string
			if secondname != "" {
				text, err = w.predict(name, secondname)
			} else {
				text, err = w.predict(name)
			}
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			// TODO: 可注入
			delete := ctx.Send(text)
			time.Sleep(time.Second * 30)
			ctx.DeleteMessage(delete)
		})
}
