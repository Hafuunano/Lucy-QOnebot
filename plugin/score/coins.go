// Package score 简单的积分系统
package score

import (
	"encoding/json"
	"math"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/FloatTech/floatbox/web"
	"github.com/FloatTech/zbputils/control"
	"github.com/tidwall/gjson"

	fcext "github.com/FloatTech/floatbox/ctxext"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type partygame struct {
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Coins string `json:"coins"`
}

var (
	pgs          = make(pg, 256)
	checkLimit   = rate.NewManager[int64](time.Minute*1, 5) // time setup
	catchLimit   = rate.NewManager[int64](time.Hour*1, 7)   // time setup
	processLimit = rate.NewManager[int64](time.Hour*1, 5)   // time setup
	payLimit     = rate.NewManager[int64](time.Hour*1, 10)  // time setup
)

type pg = map[string]partygame

func init() {
	loadFiles := fcext.DoOnceOnSuccess(func(ctx *zero.Ctx) bool {
		data, err := os.ReadFile(engine.DataFolder() + "loads.json")
		if err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			return false
		}
		err = json.Unmarshal(data, &pgs)
		if err != nil {
			panic(err)
		}
		return true
	})
	// 借鉴了其他bot的功能 编写而成
	engine.OnFullMatch("柠檬片总数", zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		uid := ctx.Event.UserID
		si := sdb.GetSignInByUID(uid)
		ctx.SendChain(message.Text("您的柠檬片数量一共是: ", si.Coins))
	})
	engine.OnRegex(`^查询(\[CQ:at,qq=(\d+)\]\s?|(\d+))的柠檬片`, zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		TargetInt, _ := strconv.ParseInt(ctx.State["regex_matched"].([]string)[2], 10, 64)
		siTargetUser := sdb.GetSignInByUID(TargetInt)
		getTargetName := ctx.CardOrNickName(TargetInt)
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("这位 ( ", getTargetName, ") 的柠檬片为", siTargetUser.Coins, "个"))
	})
	engine.OnFullMatch("抽奖", loadFiles, zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		if !checkLimit.Load(ctx.Event.UserID).Acquire() {
			ctx.SendChain(message.Text("太贪心了哦~过会试试吧"))
			return
		}
		var mutex sync.RWMutex // 添加读写锁以保证稳定性
		mutex.Lock()
		uid := ctx.Event.UserID
		si := sdb.GetSignInByUID(uid) // 获取用户目前状况信息
		userCurrentCoins := si.Coins  // loading coins status
		checkEnoughCoins := checkUserCoins(userCurrentCoins)
		if !checkEnoughCoins {
			ctx.SendChain(message.Reply(uid), message.Text("本次参与的柠檬片不够哦~请多多打卡w"))
			return
		}
		all := rand.Intn(39) // 一共39种可能性
		referpg := pgs[(strconv.Itoa(all))]
		getName := referpg.Name
		getCoinsStr := referpg.Coins
		getCoinsInt, _ := strconv.Atoi(getCoinsStr)
		getDesc := referpg.Desc
		addNewCoins := si.Coins + getCoinsInt - 200
		_ = sdb.InsertUserCoins(uid, addNewCoins)
		msgid := ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(" 嗯哼~来玩抽奖了哦w 看看能抽到什么呢w"))
		time.Sleep(time.Second * 3)
		ctx.SendChain(message.Reply(msgid), message.Text("呼呼~让咱看看你抽到了什么东西ww\n"),
			message.Text("你抽到的是~ ", getName, "\n", "获得了积分 ", getCoinsInt, "\n", getDesc, "\n目前的柠檬片总数为：", addNewCoins))
		mutex.Unlock()
	})
	// 一次最多骗 500 柠檬片,失败概率较大,失败会被反吞柠檬片
	engine.OnRegex(`^抢(\[CQ:at,qq=(\d+)\]\s?|(\d+))的柠檬片`, zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		if !catchLimit.Load(ctx.Event.UserID).Acquire() {
			ctx.SendChain(message.Text("太贪心了哦~一小时后再来试试吧"))
			return
		}
		TargetInt, _ := strconv.ParseInt(ctx.State["regex_matched"].([]string)[2], 10, 64)
		if TargetInt == ctx.Event.UserID {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("哈? 干嘛骗自己的?坏蛋哦"))
			return
		}
		uid := ctx.Event.UserID
		siEventUser := sdb.GetSignInByUID(uid)        // 获取主用户目前状况信息
		siTargetUser := sdb.GetSignInByUID(TargetInt) // 获得被抢用户目前情况信息
		switch {
		case siEventUser.Coins < 500:
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("貌似没有足够的柠檬片去准备哦~请多多打卡w"))
			return
		case siTargetUser.Coins < 500:
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("太坏了~试图的对象貌似没有足够多的柠檬片~"))
			return
		}
		eventUserName := ctx.CardOrNickName(uid)
		eventTargetName := ctx.CardOrNickName(TargetInt)
		modifyCoins := rand.Intn(500)

		if rand.Intn(10)/8 != 0 { // 7成失败概率
			_ = sdb.InsertUserCoins(siEventUser.UID, siEventUser.Coins-modifyCoins)
			_ = sdb.InsertUserCoins(siTargetUser.UID, siTargetUser.Coins+modifyCoins)
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("试着去拿走 ", eventTargetName, " 的柠檬片时,被发现了.\n所以 ", eventUserName, " 失去了 ", modifyCoins, " 个柠檬片\n\n同时 ", eventTargetName, " 得到了 ", modifyCoins, " 个柠檬片"))
			return
		}
		_ = sdb.InsertUserCoins(siEventUser.UID, siEventUser.Coins+modifyCoins)
		_ = sdb.InsertUserCoins(siTargetUser.UID, siTargetUser.Coins-modifyCoins)
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("试着去拿走 ", eventTargetName, " 的柠檬片时,成功了.\n所以 ", eventUserName, " 得到了 ", modifyCoins, " 个柠檬片\n\n同时 ", eventTargetName, " 失去了 ", modifyCoins, " 个柠檬片"))
	})
	engine.OnRegex(`^骗\s?\[CQ:at,qq=(\d+)\]\s(\d+)个柠檬片$`, zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		if !catchLimit.Load(ctx.Event.UserID).Acquire() {
			ctx.SendChain(message.Text("太贪心了哦~一小时后再来试试吧"))
			return
		}
		TargetInt, _ := strconv.ParseInt(ctx.State["regex_matched"].([]string)[1], 10, 64)
		modifyCoins, _ := strconv.Atoi(ctx.State["regex_matched"].([]string)[2])
		if TargetInt == ctx.Event.UserID {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("哈? 干嘛骗自己的?坏蛋哦"))
			return
		}
		switch {
		case modifyCoins <= 50:
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("貌似你是想倒贴别人来着嘛?可以试试多骗一点哦，既然都骗了那就多点吧x"))
			return
		case modifyCoins > 2000:
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("不要太贪心了啦！太坏了 "))
			return
		}
		uid := ctx.Event.UserID
		siEventUser := sdb.GetSignInByUID(uid)        // 获取主用户目前状况信息
		siTargetUser := sdb.GetSignInByUID(TargetInt) // 获得被抢用户目前情况信息
		switch {
		case siTargetUser.Coins < modifyCoins:
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("太坏了~试图的对象貌似没有足够多的柠檬片~"))
			return
		case siEventUser.Coins < modifyCoins:
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("貌似你需要有那么多数量的柠檬片哦w"))
			return
		}
		eventUserName := ctx.CardOrNickName(uid)
		eventTargetName := ctx.CardOrNickName(TargetInt)
		// get random numbers.
		getTargetChanceToDealRaw := math.Round(float64(modifyCoins / 20)) // the total is 0-100，however I don't allow to get 0 chance. lmao. max is 100 if modify is 2000
		getTargetChanceToDealPossibilityKey := rand.Intn(101)
		if getTargetChanceToDealPossibilityKey < int(getTargetChanceToDealRaw) { // failed
			doubledModifyNum := modifyCoins * 2
			if doubledModifyNum > siEventUser.Coins {
				doubledModifyNum = siEventUser.Coins
				_ = sdb.InsertUserCoins(siEventUser.UID, siEventUser.Coins-doubledModifyNum)
				_ = sdb.InsertUserCoins(siTargetUser.UID, siTargetUser.Coins+doubledModifyNum)
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("试着去骗走 ", eventTargetName, " 的柠檬片时,被 ", eventTargetName, " 发现了.\n本该失去 ", modifyCoins*2, "\n但因为 ", eventUserName, " 的柠檬片过少，所以 ", eventUserName, " 失去了 ", doubledModifyNum, " 个柠檬片\n\n同时 ", eventTargetName, " 得到了 ", doubledModifyNum, " 个柠檬片"))
				return
			}
			_ = sdb.InsertUserCoins(siEventUser.UID, siEventUser.Coins-doubledModifyNum)
			_ = sdb.InsertUserCoins(siTargetUser.UID, siTargetUser.Coins+doubledModifyNum)
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("试着去骗走 ", eventTargetName, " 的柠檬片时,被 ", eventTargetName, " 发现了.\n所以 ", eventUserName, " 失去了 ", doubledModifyNum, " 个柠檬片\n\n同时 ", eventTargetName, " 得到了 ", doubledModifyNum, " 个柠檬片"))
			return
		}
		_ = sdb.InsertUserCoins(siEventUser.UID, siEventUser.Coins+modifyCoins)
		_ = sdb.InsertUserCoins(siTargetUser.UID, siTargetUser.Coins-modifyCoins)
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("试着去拿走 ", eventTargetName, " 的柠檬片时,成功了.\n所以 ", eventUserName, " 得到了 ", modifyCoins, " 个柠檬片\n\n同时 ", eventTargetName, " 失去了 ", modifyCoins, " 个柠檬片"))
	})
	engine.OnRegex(`^给\s?\[CQ:at,qq=(\d+)\]\s转(\d+)个柠檬片$`, zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		if !processLimit.Load(ctx.Event.UserID).Acquire() {
			ctx.SendChain(message.Text("请等一会再转账哦w"))
			return
		}
		TargetInt, _ := strconv.ParseInt(ctx.State["regex_matched"].([]string)[1], 10, 64)
		modifyCoins, _ := strconv.Atoi(ctx.State["regex_matched"].([]string)[2])
		if TargetInt == ctx.Event.UserID {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("不可以给自己转账哦w（敲）"))
			return
		}
		uid := ctx.Event.UserID
		siEventUser := sdb.GetSignInByUID(uid)        // 获取主用户目前状况信息
		siTargetUser := sdb.GetSignInByUID(TargetInt) // 获得被转账用户目前情况信息
		if modifyCoins > siEventUser.Coins {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("貌似你的柠檬片数量不够哦~"))
			return
		}
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("转账成功了哦~\n", ctx.CardOrNickName(siEventUser.UID), " 变化为 ", siEventUser.Coins, " - ", modifyCoins, "= ", siEventUser.Coins-modifyCoins, "\n", ctx.CardOrNickName(siTargetUser.UID), " 变化为: ", siTargetUser.Coins, " + ", modifyCoins, "= ", siTargetUser.Coins+modifyCoins))
		_ = sdb.InsertUserCoins(siEventUser.UID, siEventUser.Coins-modifyCoins)
		_ = sdb.InsertUserCoins(siTargetUser.UID, siTargetUser.Coins+modifyCoins)
	})
	engine.OnRegex(`^HandleCoins\s?\[CQ:at,qq=(\d+)\]\s(\d+)$`, zero.SuperUserPermission).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		TargetInt, _ := strconv.ParseInt(ctx.State["regex_matched"].([]string)[1], 10, 64)
		modifyCoins, _ := strconv.Atoi(ctx.State["regex_matched"].([]string)[2])
		siTargetUser := sdb.GetSignInByUID(TargetInt) // get user info
		unModifyCoins := siTargetUser.Coins
		_ = sdb.InsertUserCoins(TargetInt, unModifyCoins+modifyCoins)
		ctx.Send("Handle Coins Successfully.\n")
	})
	engine.OnRegex(`^丢弃(\d+)个柠檬片$`).SetBlock(true).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		modifyCoins, _ := strconv.Atoi(ctx.State["regex_matched"].([]string)[1])
		handleUser := sdb.GetSignInByUID(ctx.Event.UserID)
		currentUserCoins := handleUser.Coins
		if currentUserCoins-modifyCoins < 0 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("貌似你的柠檬片不够处理呢("))
			return
		}
		hadModifyCoins := currentUserCoins - modifyCoins
		_ = sdb.InsertUserCoins(handleUser.UID, hadModifyCoins)
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经帮你扔掉了哦"))
	})
	engine.OnFullMatch("兑换涩图", zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		if !payLimit.Load(ctx.Event.UserID).Acquire() {
			ctx.SendChain(message.Text("坏欸！为什么一个群有这么多人看涩图啊（晕"))
			return
		}
		modified, ok := control.Lookup("nsfw")
		if ok {
			status := modified.IsEnabledIn(ctx.Event.GroupID)
			if status {
				var mutex sync.RWMutex // 添加读写锁以保证稳定性
				mutex.Lock()
				uid := ctx.Event.UserID
				si := sdb.GetSignInByUID(uid) // 获取用户目前状况信息
				userCurrentCoins := si.Coins  // loading coins status
				if userCurrentCoins < 400 {
					ctx.SendChain(message.Reply(uid), message.Text("本次参与的柠檬片不够哦~请多多打卡w，一次兑换最少需要400"))
					return
				}
				img, err := web.GetData("https://api.lolicon.app/setu/v2?r18=1&num=1")
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				picURL := gjson.Get(string(img), "data.0.urls.original").String()
				messageID := ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(picURL)).ID()
				if messageID != 0 { // 保证成功后才扣除
					_ = sdb.InsertUserCoins(si.UID, userCurrentCoins-400)
				}
			} else {
				ctx.SendChain(message.Text("本群并没有开启nsfw哦，不允许使用此功能哦x"))
				return
			}
		}
	})
}
