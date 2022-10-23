package score

import (
	"encoding/json"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
	"unsafe"

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
	RateLimit    = rate.NewManager[int64](time.Second*60, 9) // time setup
	CheckLimit   = rate.NewManager[int64](time.Minute*1, 4)  // time setup
	CatchLimit   = rate.NewManager[int64](time.Hour*1, 8)    // time setup
	processLimit = rate.NewManager[int64](time.Hour*1, 3)    // time setup
)

type pg = map[string]partygame

func init() {
	// 借鉴了其他bot的功能 编写而成
	engine.OnFullMatch("柠檬片总数", zero.OnlyGroup).Handle(func(ctx *zero.Ctx) {
		uid := ctx.Event.UserID
		si := sdb.GetSignInByUID(uid)
		ctx.SendChain(message.Text("您的柠檬片数量一共是: ", si.Coins))
	})
	engine.OnFullMatch("抽奖", zero.OnlyGroup).Handle(func(ctx *zero.Ctx) {
		if !CheckLimit.Load(ctx.Event.UserID).Acquire() {
			ctx.SendChain(message.Text("太贪心了哦~过会试试吧"))
			return
		}
		var mutex sync.RWMutex // 添加读写锁以保证稳定性
		mutex.Lock()
		uid := ctx.Event.UserID
		data, err := os.ReadFile(engine.DataFolder() + "loads.json")
		if err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			return
		}
		err = json.Unmarshal(data, &pgs)
		if err != nil {
			panic(err)
		}
		si := sdb.GetSignInByUID(uid) // 获取用户目前状况信息
		userCurrentCoins := si.Coins  // loading coins status
		checkEnoughCoins := checkUserCoins(userCurrentCoins)
		if !checkEnoughCoins {
			ctx.SendChain(message.Reply(uid), message.Text("本次参与的柠檬片不够哦~请多多打卡w"))
			return
		}
		all := rand.Intn(39) // 一共37种可能性
		referpg := pgs[(strconv.Itoa(all))]
		getName := referpg.Name
		getCoinsStr := referpg.Coins
		getCoinsInt, _ := strconv.Atoi(getCoinsStr)
		getDesc := referpg.Desc
		addNewCoins := si.Coins + getCoinsInt - 15
		err = sdb.InsertOrUpdateSignInCountByUID(uid, 0, addNewCoins)
		if err != nil {
			ctx.SendChain(message.Text("ERR:", err))
			return
		}
		msgid := ctx.SendChain(message.At(uid), message.Text(" ~好哦w 来抽奖 看看能抽到什么东西呢w"))
		time.Sleep(time.Second * 3)
		ctx.SendChain(message.Reply(msgid), message.Text("好哦~让咱看看你抽到了什么东西ww\n"),
			message.Text("你抽到的是~ ", getName, "\n", "获得了积分 ", getCoinsInt, "\n", getDesc, "\n目前的柠檬片总数为：", addNewCoins))
		mutex.Unlock()
	})
	// 一次最多骗 50 柠檬片,失败概率较大,失败会被反吞柠檬片
	engine.OnRegex(`^抢(\[CQ:at,qq=(\d+)\]\s?|(\d+))的柠檬片`, zero.OnlyGroup).Handle(func(ctx *zero.Ctx) {
		if !CatchLimit.Load(ctx.Event.UserID).Acquire() {
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
		case siEventUser.Coins < 50:
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("貌似没有足够的柠檬片去准备哦~请多多打卡w"))
			return
		case siTargetUser.Coins < 50:
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("太坏了~试图的对象貌似没有足够多的柠檬片~"))
			return
		}
		eventUserName := ctx.CardOrNickName(uid)
		eventTargetName := ctx.CardOrNickName(TargetInt)
		modifyCoins := rand.Intn(50)
		if rand.Intn(10)/7 != 0 { // 6成失败概率
			_ = sdb.InsertOrUpdateSignInCountByUID(siEventUser.UID, 0, siEventUser.Coins-modifyCoins)
			time.Sleep(time.Second * 2)
			_ = sdb.InsertOrUpdateSignInCountByUID(siTargetUser.UID, 0, siTargetUser.Coins+modifyCoins)
			time.Sleep(time.Second * 2)
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("试着去拿走 ", eventTargetName, " 的柠檬片时,被发现了.\n所以 ", eventUserName, " 失去了 ", modifyCoins, " 个柠檬片\n\n同时 ", eventTargetName, " 得到了 ", modifyCoins, " 个柠檬片"))
			return
		}
		_ = sdb.InsertOrUpdateSignInCountByUID(siEventUser.UID, 0, siEventUser.Coins+modifyCoins)
		time.Sleep(time.Second * 2)
		_ = sdb.InsertOrUpdateSignInCountByUID(siTargetUser.UID, 0, siTargetUser.Coins-modifyCoins)
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("试着去拿走 ", eventTargetName, " 的柠檬片时,成功了.\n所以 ", eventUserName, " 得到了 ", modifyCoins, " 个柠檬片\n\n同时 ", eventTargetName, " 失去了 ", modifyCoins, " 个柠檬片"))
	})

	engine.OnRegex(`^骗\s?\[CQ:at,qq=(\d+)\]\s(\d+)个柠檬片$`, zero.OnlyGroup).Handle(func(ctx *zero.Ctx) {
		if !CatchLimit.Load(ctx.Event.UserID).Acquire() {
			ctx.SendChain(message.Text("太贪心了哦~一小时后再来试试吧"))
			return
		}
		TargetInt, _ := strconv.ParseInt(ctx.State["regex_matched"].([]string)[1], 10, 64)
		modifyCoinsInt64, _ := strconv.ParseInt(ctx.State["regex_matched"].([]string)[2], 10, 64)
		modifyCoins := *(*int)(unsafe.Pointer(&modifyCoinsInt64))
		if TargetInt == ctx.Event.UserID {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("哈? 干嘛骗自己的?坏蛋哦"))
			return
		}
		switch {
		case modifyCoins <= 0:
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("貌似你是想倒贴别人来着?"))
			return
		case modifyCoins > 2000:
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("不要太贪心了啦！太坏了吧"))
			return
		}
		uid := ctx.Event.UserID
		siEventUser := sdb.GetSignInByUID(uid)        // 获取主用户目前状况信息
		siTargetUser := sdb.GetSignInByUID(TargetInt) // 获得被抢用户目前情况信息
		if siTargetUser.Coins < modifyCoins {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("太坏了~试图的对象貌似没有足够多的柠檬片~"))
			return
		}
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
		if rand.Intn(2100)/modifyCoins != 0 { // failed
			doubledModifyNum := modifyCoins * 2
			if doubledModifyNum > siEventUser.Coins {
				doubledModifyNum = siEventUser.Coins
				_ = sdb.InsertOrUpdateSignInCountByUID(siEventUser.UID, 0, siEventUser.Coins-doubledModifyNum)
				time.Sleep(time.Second * 2)
				_ = sdb.InsertOrUpdateSignInCountByUID(siTargetUser.UID, 0, siTargetUser.Coins+doubledModifyNum)
				time.Sleep(time.Second * 2)
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("试着去骗走 ", eventTargetName, " 的柠檬片时,被 ", eventTargetName, " 发现了.\n本该失去 ", modifyCoins*2, "\n但因为 ", eventUserName, " 的柠檬片过少，所以 ", eventUserName, " 失去了 ", doubledModifyNum, " 个柠檬片\n\n同时 ", eventTargetName, " 得到了 ", doubledModifyNum, " 个柠檬片"))
				return
			}
			_ = sdb.InsertOrUpdateSignInCountByUID(siEventUser.UID, 0, siEventUser.Coins-doubledModifyNum)
			time.Sleep(time.Second * 2)
			_ = sdb.InsertOrUpdateSignInCountByUID(siTargetUser.UID, 0, siTargetUser.Coins+doubledModifyNum)
			time.Sleep(time.Second * 2)
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("试着去骗走 ", eventTargetName, " 的柠檬片时,被 ", eventTargetName, " 发现了.\n所以 ", eventUserName, " 失去了 ", doubledModifyNum, " 个柠檬片\n\n同时 ", eventTargetName, " 得到了 ", doubledModifyNum, " 个柠檬片"))
			return
		}
		_ = sdb.InsertOrUpdateSignInCountByUID(siEventUser.UID, 0, siEventUser.Coins+modifyCoins)
		time.Sleep(time.Second * 2)
		_ = sdb.InsertOrUpdateSignInCountByUID(siTargetUser.UID, 0, siTargetUser.Coins-modifyCoins)
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("试着去拿走 ", eventTargetName, " 的柠檬片时,成功了.\n所以 ", eventUserName, " 得到了 ", modifyCoins, " 个柠檬片\n\n同时 ", eventTargetName, " 失去了 ", modifyCoins, " 个柠檬片"))
	})
	engine.OnRegex(`^给\s?\[CQ:at,qq=(\d+)\]\s转(\d+)个柠檬片$`, zero.OnlyGroup).Handle(func(ctx *zero.Ctx) {
		if !processLimit.Load(ctx.Event.UserID).Acquire() {
			ctx.SendChain(message.Text("请等一会再转账哦w"))
			return
		}
		TargetInt, _ := strconv.ParseInt(ctx.State["regex_matched"].([]string)[1], 10, 64)
		modifyCoinsInt64, _ := strconv.ParseInt(ctx.State["regex_matched"].([]string)[2], 10, 64)
		modifyCoins := *(*int)(unsafe.Pointer(&modifyCoinsInt64))
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
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("转账成功了哦~\n", ctx.CardOrNickName(siEventUser.UID), " 变化为 ", siEventUser.Coins, " - ", modifyCoins, "\n", ctx.CardOrNickName(siTargetUser.UID), " 变化为: ", siTargetUser.Coins, " + ", modifyCoins))
		_ = sdb.InsertOrUpdateSignInCountByUID(siEventUser.UID, 0, siEventUser.Coins-modifyCoins)
		time.Sleep(time.Second * 2)
		_ = sdb.InsertOrUpdateSignInCountByUID(siTargetUser.UID, 0, siTargetUser.Coins+modifyCoins)
	})
}
