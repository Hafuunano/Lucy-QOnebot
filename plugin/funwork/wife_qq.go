package funwork

import (
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/FloatTech/imgfactory"

	fcext "github.com/FloatTech/floatbox/ctxext"
	"github.com/FloatTech/floatbox/math"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	// 数据库
	sql "github.com/FloatTech/sqlite"
	// 定时器
	"github.com/FloatTech/zbputils/img/text"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"

	"github.com/FloatTech/gg"
)

// MarryLogin get marry login struct list.
type MarryLogin struct {
	db   *sql.Sqlite
	dbmu sync.RWMutex
}

// Userinfo 结婚证信息
type Userinfo struct {
	User       int64  // 用户身份证
	Target     int64  // 对象身份证号
	Username   string // 户主名称
	Targetname string // 对象名称
	Updatetime string // 登记时间

}

// mainList的当前时间
type updateinfo struct {
	GID        int64
	Updatetime string // 登记时间

}

func (sql *MarryLogin) checkUpdate(gid int64) (updatetime string, err error) {
	sql.dbmu.Lock()
	defer sql.dbmu.Unlock()
	err = sql.db.Create("updateinfo", &updateinfo{})
	if err != nil {
		return
	}
	gidstr := strconv.FormatInt(gid, 10)
	dbinfo := updateinfo{}
	err = sql.db.Find("updateinfo", &dbinfo, "where gid is "+gidstr) // 获取表格更新的时间
	if err != nil {
		updatetime = time.Now().Format("2006/01/02")
		err = sql.db.Insert("updateinfo", &updateinfo{GID: gid, Updatetime: updatetime})
	}
	updatetime = dbinfo.Updatetime
	return
}

func (sql *MarryLogin) reset(gid string) error {
	sql.dbmu.Lock()
	defer sql.dbmu.Unlock()
	if gid != "ALL" {
		err := sql.db.Drop(gid)
		if err != nil {
			err = sql.db.Create(gid, &Userinfo{})
			return err
		}
		gidint, _ := strconv.ParseInt(gid, 10, 64)
		updateinfo := updateinfo{
			GID:        gidint,
			Updatetime: time.Now().Format("2006/01/02"),
		}
		err = sql.db.Insert("updateinfo", &updateinfo)
		return err
	}
	grouplist, err := sql.db.ListTables()
	if err != nil {
		return err
	}
	for _, gid := range grouplist {
		err = sql.db.Drop(gid)
		if err != nil {
			continue
		}
		gidint, _ := strconv.ParseInt(gid, 10, 64)
		updateinfo := updateinfo{
			GID:        gidint,
			Updatetime: time.Now().Format("2006/01/02"),
		}
		err = sql.db.Insert("updateinfo", &updateinfo)
	}
	return err
}

func (sql *MarryLogin) LeaveWithWife(gid, wife int64) error {
	sql.dbmu.Lock()
	defer sql.dbmu.Unlock()
	gidstr := strconv.FormatInt(gid, 10)
	wifestr := strconv.FormatInt(wife, 10)
	// 先判断用户是否存在
	err := sql.db.Del(gidstr, "where target = "+wifestr)
	return err
}

func (sql *MarryLogin) LeaveWithHusband(gid, husband int64) error {
	sql.dbmu.Lock()
	defer sql.dbmu.Unlock()
	gidstr := strconv.FormatInt(gid, 10)
	husbandstr := strconv.FormatInt(husband, 10)
	// 先判断用户是否存在
	err := sql.db.Del(gidstr, "where user = "+husbandstr)
	return err
}

func (sql *MarryLogin) ReMarried(gid, uid, target int64, username, targetname string) error {
	sql.dbmu.Lock()
	defer sql.dbmu.Unlock()
	gidstr := strconv.FormatInt(gid, 10)
	uidstr := strconv.FormatInt(uid, 10)
	tagstr := strconv.FormatInt(target, 10)
	var info Userinfo
	err := sql.db.Find(gidstr, &info, "where user = "+uidstr)
	if err != nil {
		err = sql.db.Find(gidstr, &info, "where user = "+tagstr)
	}
	if err != nil {
		return err
	}
	updatetime := time.Now().Format("2006/01/02")
	// 更改夫妻信息
	info.User = uid
	info.Username = username
	info.Target = target
	info.Targetname = targetname
	info.Updatetime = updatetime
	// mainList登记数据
	err = sql.db.Insert(gidstr, &info)
	return err
}

func (sql *MarryLogin) MarriedList(gid int64) (list [][4]string, number int, err error) {
	sql.dbmu.Lock()
	defer sql.dbmu.Unlock()
	gidstr := strconv.FormatInt(gid, 10)
	err = sql.db.Create(gidstr, &Userinfo{})
	if err != nil {
		return
	}
	number, err = sql.db.Count(gidstr)
	if err != nil || number <= 0 {
		return
	}
	var info Userinfo
	list = make([][4]string, 0, number)
	err = sql.db.FindFor(gidstr, &info, "GROUP BY user", func() error {
		if info.Target == 0 {
			return nil
		}
		dbinfo := [4]string{
			info.Username,
			strconv.FormatInt(info.User, 10),
			info.Targetname,
			strconv.FormatInt(info.Target, 10),
		}
		list = append(list, dbinfo)
		return nil
	})
	if len(list) == 0 {
		number = 0
	}
	return
}

func slicename(name string, canvas *gg.Context) (resultname string) {
	usermane := []rune(name) // 将每个字符单独放置
	widthlen := 0
	numberlen := 0
	for i, v := range usermane {
		width, _ := canvas.MeasureString(string(v)) // 获取单个字符的宽度
		widthlen += int(width)
		if widthlen > 350 {
			break // 总宽度不能超过350
		}
		numberlen = i
	}
	if widthlen > 350 {
		resultname = string(usermane[:numberlen-1]) + "......" // 名字切片
	} else {
		resultname = name
	}
	return
}

func (sql *MarryLogin) CheckMarriedList(gid, uid int64) (info Userinfo, status int, err error) {
	sql.dbmu.Lock()
	defer sql.dbmu.Unlock()
	gidstr := strconv.FormatInt(gid, 10)
	uidstr := strconv.FormatInt(uid, 10)
	status = 3
	if err = sql.db.Create(gidstr, &Userinfo{}); err != nil {
		status = 2
		return
	}
	err = sql.db.Find(gidstr, &info, "where user = "+uidstr)
	if err == nil {
		status = 1
		return
	}
	err = sql.db.Find(gidstr, &info, "where target = "+uidstr)
	if err == nil {
		status = 0
		return
	}
	return
}

func (sql *MarryLogin) GetLogined(gid, uid, target int64, username, targetname string) error {
	sql.dbmu.Lock()
	defer sql.dbmu.Unlock()
	gidstr := strconv.FormatInt(gid, 10)
	err := sql.db.Create(gidstr, &Userinfo{})
	if err != nil {
		return err
	}
	updatetime := time.Now().Format("2006/01/02")
	// 填写夫妻信息
	uidinfo := Userinfo{
		User:       uid,
		Username:   username,
		Target:     target,
		Targetname: targetname,
		Updatetime: updatetime,
	}
	// mainList登记数据
	err = sql.db.Insert(gidstr, &uidinfo)
	return err
}

var (
	mainList = &MarryLogin{
		db: &sql.Sqlite{},
	}
	skillCD  = rate.NewManager[string](time.Hour*24, 4)
	sendtext = [...][]string{
		{ // 表白成功
			"是个勇敢的孩子(*/ω＼*) 今天的运气都降临在你的身边~\n\n",
			"貌似很成功的一次尝试呢w~\n\n",
		},
		{ // 表白失败
			"今天的运气有一点背哦~这一次没有成功呢x",
			"_(:з」∠)_下次还有机会 抱抱w",
			"没关系哦，虽然失败了但还有机会呢x",
		},
		{ // ntr成功
			"欸~成功了哦_(:з」∠)_w\n\n",
		},
		{ // 离婚失败
			"为什么要分呢? 这样做可不好哦x 让咱捏捏w",
			"太坏了啦！不许！",
		},
		{ // 离婚成功
			"好呢w 就这样呢(",
			"已经成功了哦w",
		},
	}
)

func init() {
	getdb := fcext.DoOnceOnSuccess(func(ctx *zero.Ctx) bool {
		mainList.db.DBPath = engine.DataFolder() + "groupfun_wife.db"
		err := mainList.db.Open(time.Hour * 24)
		if err != nil {
			ctx.SendChain(message.Text("ERR:", err))
			return false
		}
		return true
	})
	engine.OnFullMatch("娶群友", zero.OnlyGroup, getdb).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			gid := ctx.Event.GroupID
			updatetime, err := mainList.checkUpdate(gid)
			switch {
			case err != nil:
				ctx.SendChain(message.Text("ERR:", err))
				return
			case time.Now().Format("2006/01/02") != updatetime:
				if err := mainList.reset(strconv.FormatInt(gid, 10)); err != nil {
					ctx.SendChain(message.Text("ERR:", err))
					return
				}
			}
			uid := ctx.Event.UserID
			targetinfo, status, err := mainList.CheckMarriedList(gid, uid)
			switch {
			case status == 2:
				ctx.SendChain(message.Text("ERR:", err))
				return
			case status != 3 && targetinfo.Target == 0: // 如果为单身贵族
				ctx.SendChain(message.At(uid), message.Text("\n今天没有老婆哦 ^^ 是不是觉得咱很坏呢 那就下次再来吧(笑"))
				return
			case status == 1: // 娶过别人
				cardTarget := ctx.CardOrNickName(targetinfo.Target)
				ctx.SendChain(
					message.At(uid),
					message.Text("\n今天你已经娶过了，群老婆是"),
					message.Image("https://q4.qlogo.cn/g?b=qq&nk="+strconv.FormatInt(targetinfo.Target, 10)+"&s=640").Add("cache", 0),
					message.Text(
						"\n",
						" [", cardTarget, "]",
						"(", targetinfo.Target, ")哒",
					),
				)
				return
			case status == 0: // 嫁给别人
				cardTarget := ctx.CardOrNickName(targetinfo.User)
				ctx.SendChain(
					message.At(uid),
					message.Text("\n今天你被娶了，群老公是"),
					message.Image("https://q4.qlogo.cn/g?b=qq&nk="+strconv.FormatInt(targetinfo.User, 10)+"&s=640").Add("cache", 0),
					message.Text(
						"\n",
						" [", cardTarget, "]",
						"(", targetinfo.User, ")哒",
					),
				)
				return
			}
			// 无缓存获取群员列表
			temp := ctx.GetThisGroupMemberListNoCache().Array()
			sort.SliceStable(temp, func(i, j int) bool {
				return temp[i].Get("last_sent_time").Int() < temp[j].Get("last_sent_time").Int()
			})
			temp = temp[math.Max(0, len(temp)-30):]
			// 将已经娶过的人剔除
			qqgrouplist := make([]int64, 0, len(temp))
			for k := 0; k < len(temp); k++ {
				usr := temp[k].Get("user_id").Int()
				_, status, _ := mainList.CheckMarriedList(gid, usr)
				if status != 3 {
					continue
				}
				qqgrouplist = append(qqgrouplist, usr)
			}
			// 没有人（只剩自己）的时候
			if len(qqgrouplist) == 1 {
				ctx.SendChain(message.Text("~群里已经没有人是单身了哦 明天再试试叭"))
				return
			}
			// 随机抽娶
			fiancee := qqgrouplist[rand.Intn(len(qqgrouplist))]
			if fiancee == uid { // 如果是自己
				ctx.SendChain(message.At(uid), message.Text("\n^^ 可能运气比较好呢~抽到了自己^^, Lucy再让你抽一下哦"))
				return
			}
			// 去mainList办证
			err = mainList.GetLogined(gid, uid, fiancee, ctx.CardOrNickName(uid), ctx.CardOrNickName(fiancee))
			if err != nil {
				ctx.SendChain(message.Text("ERR:", err))
				return
			}
			// 请大家吃席
			ctx.SendChain(
				message.At(uid),
				message.Text("今天你的群老婆是"),
				message.Image("https://q4.qlogo.cn/g?b=qq&nk="+strconv.FormatInt(fiancee, 10)+"&s=640").Add("cache", 0),
				message.Text(
					"\n",
					"[", ctx.CardOrNickName(fiancee), "]",
					"(", fiancee, ")哒",
				),
			)
		})

	engine.OnRegex(`^(娶|嫁)Lucy`, zero.OnlyGroup, getdb).SetBlock(true).Limit(cdcheck, iscding).Handle(func(ctx *zero.Ctx) {
		// not work with checkdog, so it is useless and need to be rewritten.
		choice := ctx.State["regex_matched"].([]string)[1]
		gid := ctx.Event.GroupID
		uid := ctx.Event.UserID
		var randbook int
		randbook = rand.Intn(2)
		fiancee := ctx.Event.SelfID
		// first of all , judge if it is a dog
		updatetime, err := mainList.checkUpdate(gid)
		switch {
		case err != nil:
			ctx.SendChain(message.Text("check mainList panic():", err))
			return
		case time.Now().Format("2006/01/02") != updatetime:
			if err := mainList.reset(strconv.FormatInt(gid, 10)); err != nil {
				ctx.SendChain(message.Text("check mainList panic():", err))
				return
			}
		}
		// 获取用户信息
		uidtarget, uidstatus, err1 := mainList.CheckMarriedList(gid, uid)
		fianceeinfo, fianceestatus, err2 := mainList.CheckMarriedList(gid, fiancee)
		switch {
		case uidstatus == 2 || fianceestatus == 2:
			ctx.SendChain(message.Text("ERR:", err1, "\n", err2))
			return
		case uidtarget.Target == fiancee: // 如果本就是一块
			ctx.SendChain(message.Text("笨蛋~自己要负责a kora(´･ω･`) "))
			return
		case uidstatus == 1: // 如果如为攻
			ctx.SendChain(message.Text("笨蛋~你已经有老婆了w"))
			return
		case uidstatus == 0: // 如果为受
			ctx.SendChain(message.Text("该是0就是0，当0有什么不好"))
			return
		case fianceestatus != 3 && fianceeinfo.Target == 0:
			ctx.SendChain(message.Text("今天的受是他自己x"))
			return
		case fianceestatus == 1: // 如果如为攻
			ctx.SendChain(message.Text("笨蛋~ta已经有受了w"))
			return
		case fianceestatus == 0: // 如果为受
			ctx.SendChain(message.Text("这是一个纯爱的世界，拒绝NTR"))
			return
		}
		if rand.Intn(5) != 1 {
			ctx.Send(message.Text("笨蛋！不准娶~ ama"))
			return
		}
		ctx.Send("好嘛....就一次哦 哼ama")
		randbook = 1
		if randbook == 0 {
			ctx.SendChain(message.Text(sendtext[1][rand.Intn(len(sendtext[1]))]))
			return
		}
		// 去mainList登记
		var choicetext string
		switch choice {
		case "娶":
			err := mainList.GetLogined(gid, uid, fiancee, ctx.CardOrNickName(uid), ctx.CardOrNickName(fiancee))
			if err != nil {
				ctx.SendChain(message.Text("ERR:", err))
				return
			}
			choicetext = "\n今天你的受是"
		default:
			err := mainList.GetLogined(gid, fiancee, uid, ctx.CardOrNickName(fiancee), ctx.CardOrNickName(uid))
			if err != nil {
				ctx.SendChain(message.Text("ERR:", err))
				return
			}
			choicetext = "\n今天你的攻是"
		}
		// 请大家吃席
		ctx.SendChain(
			message.Text(sendtext[0][rand.Intn(len(sendtext[0]))]),
			message.At(uid),
			message.Text(choicetext),
			message.Image("https://q4.qlogo.cn/g?b=qq&nk="+strconv.FormatInt(fiancee, 10)+"&s=640").Add("cache", 0),
			message.Text(
				"\n",
				"[", ctx.CardOrNickName(fiancee), "]",
				"(", fiancee, ")哒",
			),
		)
	})
	// 单身技能
	engine.OnRegex(`^(娶|嫁)\[CQ:at,qq=(\d+)\]`, zero.OnlyGroup, getdb, checkdog).SetBlock(true).Limit(cdcheck, iscding).
		Handle(func(ctx *zero.Ctx) {
			choice := ctx.State["regex_matched"].([]string)[1]
			fiancee, _ := strconv.ParseInt(ctx.State["regex_matched"].([]string)[2], 10, 64)
			uid := ctx.Event.UserID
			gid := ctx.Event.GroupID
			randbook := rand.Intn(2)
			if uid == fiancee { // 如果是自己
				switch rand.Intn(5) { // 5分之一概率浪费技能
				case 1:
					err := mainList.GetLogined(gid, uid, 0, "", "")
					if err != nil {
						ctx.SendChain(message.Text("ERR:", err))
						return
					}
					ctx.Send(message.ReplyWithMessage(message.Text("貌似Lucy故意添加了 --force 的命令，成功了(笑 ")))
				default:
					ctx.SendChain(message.Text("笨蛋！娶你自己干什么a"))
				}
				return
			}
			if randbook == 0 {
				ctx.SendChain(message.Text(sendtext[1][rand.Intn(len(sendtext[1]))]))
				return
			}
			// 去mainList登记
			var choicetext string
			switch choice {
			case "娶":
				err := mainList.GetLogined(gid, uid, fiancee, ctx.CardOrNickName(uid), ctx.CardOrNickName(fiancee))
				if err != nil {
					ctx.SendChain(message.Text("ERR:", err))
					return
				}
				choicetext = "\n今天你的受是"
			default:
				err := mainList.GetLogined(gid, fiancee, uid, ctx.CardOrNickName(fiancee), ctx.CardOrNickName(uid))
				if err != nil {
					ctx.SendChain(message.Text("ERR:", err))
					return
				}
				choicetext = "\n今天你的攻是"
			}
			// 请大家吃席
			ctx.SendChain(
				message.Text(sendtext[0][rand.Intn(len(sendtext[0]))]),
				message.At(uid),
				message.Text(choicetext),
				message.Image("https://q4.qlogo.cn/g?b=qq&nk="+strconv.FormatInt(fiancee, 10)+"&s=640").Add("cache", 0),
				message.Text(
					"\n",
					"[", ctx.CardOrNickName(fiancee), "]",
					"(", fiancee, ")哒",
				),
			)
		})
	// NTR技能
	engine.OnRegex(`^试着骗(\[CQ:at,qq=(\d+)\]\s?|(\d+))做我的老婆`, zero.OnlyGroup, getdb, checkcp).SetBlock(true).Limit(cdcheck, iscding).
		Handle(func(ctx *zero.Ctx) {
			fid := ctx.State["regex_matched"].([]string)
			fiancee, _ := strconv.ParseInt(fid[2]+fid[3], 10, 64)
			uid := ctx.Event.UserID
			if fiancee == uid {
				ctx.SendChain(message.Text("要骗别人哦~为什么要骗自己呢x"))
				return
			}
			if rand.Intn(10)/4 != 0 { // 十分之三的概率NTR成功
				ctx.SendChain(message.At(uid), message.Text("\n今天没有骗成qwq"))
				return
			}
			gid := ctx.Event.GroupID
			// 判断target是老公还是老婆
			var choicetext string
			_, gender, err := mainList.CheckMarriedList(gid, fiancee)
			switch gender {
			case 3:
				ctx.SendChain(message.Text("他现在还是一人哦~成功可能性更高x 不过看起来今天又浪费了一次机会ww "))
				return
			case 2:
				ctx.SendChain(message.Text("ERR:", err))
				return
			case 1:
				// 和对象结婚登记
				err = mainList.ReMarried(gid, fiancee, uid, ctx.CardOrNickName(fiancee), ctx.CardOrNickName(uid))
				if err != nil {
					ctx.SendChain(message.Text("ERR:", err))
					return
				}
				choicetext = "老公"
			case 0:
				// 和对象结婚登记
				err = mainList.ReMarried(gid, uid, fiancee, ctx.CardOrNickName(uid), ctx.CardOrNickName(fiancee))
				if err != nil {
					ctx.SendChain(message.Text("ERR:", err))
					return
				}
				choicetext = "老婆"
			default:
				ctx.SendChain(message.Text("貌似发生了某种不可预料的错误OxO"))
				return
			}
			// 输出结果
			ctx.SendChain(
				message.Text(sendtext[2][rand.Intn(len(sendtext[2]))]),
				message.At(uid),
				message.Text("今天你的群"+choicetext+"是"),
				message.Image("https://q4.qlogo.cn/g?b=qq&nk="+strconv.FormatInt(fiancee, 10)+"&s=640").Add("cache", 0),
				message.Text(
					"\n",
					"[", ctx.CardOrNickName(fiancee), "]",
					"(", fiancee, ")哒",
				),
			)
		})
	engine.OnFullMatch("群老婆列表", zero.OnlyGroup, getdb).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			list, number, err := mainList.MarriedList(ctx.Event.GroupID)
			if err != nil {
				ctx.SendChain(message.Text("ERR:", err))
				return
			}
			if number <= 0 {
				ctx.SendChain(message.Text("今天还没有人结婚哦"))
				return
			}
			/***********设置图片的大小和底色***********/
			fontSize := 50.0
			if number < 10 {
				number = 10
			}
			canvas := gg.NewContext(1500, int(250+fontSize*float64(number)))
			canvas.SetRGB(1, 1, 1) // 白色
			canvas.Clear()
			/***********设置字体颜色为天蓝色***********/
			canvas.SetRGB255(135, 206, 250)
			/***********设置字体大小,并获取字体高度用来定位***********/
			if err = canvas.LoadFontFace(text.BoldFontFile, fontSize*2); err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			sl, h := canvas.MeasureString("互动列表")
			/***********绘制标题***********/
			canvas.DrawString("互动列表", (1500-sl)/2, 160-h) // 放置在中间位置
			canvas.DrawString("————————————————————", 0, 250-h)
			/***********设置字体大小,并获取字体高度用来定位***********/
			if err = canvas.LoadFontFace(text.BoldFontFile, fontSize); err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			_, h = canvas.MeasureString("坏")
			for i, info := range list {
				canvas.DrawString(slicename(info[0], canvas), 0, float64(260+50*i)-h)
				canvas.DrawString("("+info[1]+")", 350, float64(260+50*i)-h)
				canvas.DrawString("←→", 700, float64(260+50*i)-h)
				canvas.DrawString(slicename(info[2], canvas), 800, float64(260+50*i)-h)
				canvas.DrawString("("+info[3]+")", 1150, float64(260+50*i)-h)
			}
			data, _ := imgfactory.ToBytes(canvas.Image())
			ctx.SendChain(message.ImageBytes(data))
		})

	engine.OnFullMatch("我要离婚", zero.OnlyToMe, zero.OnlyGroup, getdb).SetBlock(true).Limit(cdcheck, iscding2).
		Handle(func(ctx *zero.Ctx) {
			gid := ctx.Event.GroupID
			updatetime, err := mainList.checkUpdate(gid)
			switch {
			case err != nil:
				ctx.SendChain(message.Text("ERR:", err))
				return
			case time.Now().Format("2006/01/02") != updatetime:
				if err := mainList.reset(strconv.FormatInt(gid, 10)); err != nil {
					ctx.SendChain(message.Text("ERR:", err))
					return
				}
				ctx.SendChain(message.Text("很奇怪a~不存在对象~这样发送的意义是什么呢"))
				return
			}

			// 获取用户信息
			uid := ctx.Event.UserID
			info, uidstatus, err := mainList.CheckMarriedList(gid, uid)
			switch uidstatus {
			case 3:
				return
			case 2:
				ctx.SendChain(message.Text("ERR:", err))
				return
			case 1:
				ctx.SendChain(message.Text(sendtext[4][rand.Intn(len(sendtext[4]))]))
				err := mainList.LeaveWithWife(gid, info.Target)
				if err != nil {
					ctx.SendChain(message.Text("ERR:", err))
					return
				}
			case 0:
				ctx.SendChain(message.Text(sendtext[4][rand.Intn(len(sendtext[4]))]))
				err := mainList.LeaveWithHusband(gid, info.User)
				if err != nil {
					ctx.SendChain(message.Text("ERR:", err))
					return
				}
			}
		})
}

// 以群号和昵称为限制
func cdcheck(ctx *zero.Ctx) *rate.Limiter {
	limitID := strconv.FormatInt(ctx.Event.GroupID, 10) + strconv.FormatInt(ctx.Event.UserID, 10)
	return skillCD.Load(limitID)
}

func iscding(ctx *zero.Ctx) {
	ctx.SendChain(message.Text("还在cd哦~请24小时后再来"))
}

// 注入判断 是否为单身
func checkdog(ctx *zero.Ctx) bool {
	fiancee, err := strconv.ParseInt(ctx.State["regex_matched"].([]string)[2], 10, 64)
	if err != nil {
		ctx.SendChain(message.Text("?貌似这个人不存在哦x"))
		return false
	}
	// 判断是否需要重置
	gid := ctx.Event.GroupID
	updatetime, err := mainList.checkUpdate(gid)
	switch {
	case err != nil:
		ctx.SendChain(message.Text("ERR:", err))
		return false
	case time.Now().Format("2006/01/02") != updatetime:
		if err := mainList.reset(strconv.FormatInt(gid, 10)); err != nil {
			ctx.SendChain(message.Text("ERR:", err))
			return false
		}
		return true // 重置后也全是单身
	}
	uid := ctx.Event.UserID
	// 获取用户信息
	uidtarget, uidstatus, err1 := mainList.CheckMarriedList(gid, uid)
	fianceeinfo, fianceestatus, err2 := mainList.CheckMarriedList(gid, fiancee)
	switch {
	case uidstatus == 2 || fianceestatus == 2:
		ctx.SendChain(message.Text("ERR:", err1, "\n", err2))
		return false
	case uidstatus == 3 && fianceestatus == 3: // 必须是两个单身
		return true
	case uidtarget.Target == fiancee: // 如果本就是一块
		ctx.SendChain(message.Text("笨蛋~自己要负责a kora(´･ω･`) "))
		return false
	case uidstatus == 1: // 如果如为攻
		ctx.SendChain(message.Text("笨蛋~你已经有老婆了w"))
		return false
	case uidstatus == 0: // 如果为受
		ctx.SendChain(message.Text("该是0就是0，当0有什么不好"))
		return false
	case fianceestatus != 3 && fianceeinfo.Target == 0:
		ctx.SendChain(message.Text("今天的受是他自己x"))
		return false
	case fianceestatus == 1: // 如果如为攻
		ctx.SendChain(message.Text("笨蛋~ta已经有受了w"))
		return false
	case fianceestatus == 0: // 如果为受
		ctx.SendChain(message.Text("这是一个纯爱的世界，拒绝NTR"))
		return false
	}
	return true
}

// 注入判断 是否满足小三的要求
func checkcp(ctx *zero.Ctx) bool {
	gid := ctx.Event.GroupID
	updatetime, err := mainList.checkUpdate(gid)
	switch {
	case err != nil:
		ctx.SendChain(message.Text("ERR:", err))
		return false
	case time.Now().Format("2006/01/02") != updatetime:
		if err := mainList.reset(strconv.FormatInt(gid, 10)); err != nil {
			ctx.SendChain(message.Text("ERR:", err))
		} else {
			ctx.SendChain(message.Text("哼~果然是个笨蛋 他还是一人呢"))
		}
		return false // 重置后全是单身
	}
	// 检查target
	fid := ctx.State["regex_matched"].([]string)
	fiancee, err := strconv.ParseInt(fid[2]+fid[3], 10, 64)
	if err != nil {
		ctx.SendChain(message.Text("这个人？他存在嘛x"))
		return false
	}
	// 检查用户是否登记过
	uid := ctx.Event.UserID
	userinfo, uidstatus, err := mainList.CheckMarriedList(gid, uid)
	switch {
	case uidstatus == 2:
		ctx.SendChain(message.Text("ERR:", err))
		return false
	case userinfo.Target == fiancee: // 如果本就是一块
		ctx.SendChain(message.Text("笨蛋~这不是你老婆了嘛w"))
		return false
	case uidstatus != 3 && userinfo.Target == 0: // 如果是单身贵族
		ctx.SendChain(message.Text("今天的受是他自己x"))
		return false
	case fiancee == uid: // 自我攻略
		return true
	case uidstatus == 1: // 如果如为攻
		ctx.SendChain(message.Text("哒咩 不给纳小妾！"))
		return false
	case uidstatus == 0: // 如果为受
		ctx.SendChain(message.Text("该是0就是0，当0有什么不好"))
		return false
	}
	fianceeinfo, fianceestatus, err := mainList.CheckMarriedList(gid, fiancee)
	switch {
	case fianceestatus == 2:
		ctx.SendChain(message.Text("ERR:", err))
		return false
	case fianceestatus == 3:
		ctx.SendChain(message.Text("哼~果然是个笨蛋 他还是一人呢"))
		return false
	case fianceeinfo.Target == 0:
		ctx.SendChain(message.Text("今天的受是他自己x"))
		return false
	}
	return true
}
func iscding2(ctx *zero.Ctx) {
	ctx.SendChain(message.Text("哒咩 不可以这样做哦x~cd状态,24小时后再来试试吧x"))
}
