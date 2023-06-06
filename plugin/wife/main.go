package wife

import (
	"fmt"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"math/rand"
	"strconv"
)

var engine = control.Register("wife", &ctrl.Options[*zero.Ctx]{
	DisableOnDefault:  false,
	Help:              "Hi NekoPachi!\n说明书: https://lucy.impart.icu",
	PrivateDataFolder: "wife",
})

/*
StatusID:

Type 1: Normal Mode, nothing happened.

Type 2: Cannot be the Target, Target became initative, so reverse.

(However the target and the initative should be in their position, DO NOT CHANGE. )

Type 3: Something is wrong, you are Target == initative Person. (Drop The Person Before.)

Type 4: Removed.
(When User get others person. || IF REMARRIED, CHANGE IT TO TYPE1.) || (Be check more Time to reduce to err.)

Type 5: NTR Mode
(Tips: NTR means changed their pairkey & TargetID || UserID, need to do some changes. ) ||
(Attempt to do once more every person.)

Type 6: No wife Mod?
Fake - Invisibile person here.
(Lucy Hides this and shows it in the next Time if a person uses NTR,
shows nothing, and Lucy will make it for joke. LMAO)
*/

func init() {
	dict := make(map[string][]string) // this dict is used to reply
	// dict path.
	dict["block"] = []string{"嗯哼？貌似没有找到哦w", "再试试哦w，或许有帮助w", "运气不太好哦，想一下办法呢x"}
	dict["success"] = []string{"Lucky To You~", "恭喜哦ww~ ", "这边来恭喜一下哦w～"}
	dict["failed"] = []string{"今天的运气有一点背哦~这一次没有成功呢x", "_(:з」∠)_下次还有机会 抱抱w", "没关系哦，虽然失败了但还有机会呢x"}
	dict["ntr"] = []string{"嗯哼～这位还是成功了呢x", "aaa 好怪 不过还是让你通过了 ^^ "}
	dict["lost_failed"] = []string{"为什么要分呢? 让咱捏捏w", "太坏了啦！不许！"}
	dict["lost_success"] = []string{"好呢w 就这样呢(", "已经成功了哦w"}
	dict["hide_mode"] = []string{"哼哼～ 哼唧", "喵喵喵？！"}

	engine.OnFullMatch("娶群友").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		/*
			Work:
			- Check the User Status, if the user is 1 or 0 || 10 ,then pause and do this handler.
			- Choose a person, do something acciednt. (if the person had, pause and give one more chance.)
			- Check the banned or Disabled Status (To Target,if had,then stoppped it and give no chance. Others has checked itself too. )
			- add this key.
			- add more feature.
		*/
		uid := ctx.Event.UserID
		gid := ctx.Event.GroupID
		// fast check
		if !CheckTheUserStatusAndDoRepeat(ctx) {
			return
		}
		// ok , go next.
		ChooseAPerson := GetUserListAndChooseOne(ctx)
		if !CheckTheTargetUserStatusAndDoRepeat(ctx, ChooseAPerson) {
			return
		}
		// check the target status.
		getStatusIfBannned := CheckTheUserIsInBlackListOrGroupList(ChooseAPerson, uid, gid)
		/*
			disabled_Target
			blacklist_Target
		*/
		if getStatusIfBannned {
			// blocked.
			GlobalCDModelCost(ctx)
			getReply := dict["block"][rand.Intn(len(dict["block"]))]
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(getReply))
			return
		}
		// reverse check
		reverseCheckTheUserIsDisabled := CheckDisabledListIsExistedInThisGroup(marryList, uid, gid)
		if reverseCheckTheUserIsDisabled {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("你已经禁用了被随机，所以不可以参与娶群友哦w"))
			return
		}
		// go next. do something colorful, pls cost something.
		// go next.
		getExistedToken := GlobalCDModelCostLeastReply(ctx)
		if getExistedToken == 0 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("今天的机会已经使用完了哦～12小时后再来试试吧"))
			return
		}
		// one chance to get himself | herself
		if ChooseAPerson == ctx.Event.UserID {
			// status code 3
			GlobalCDModelCost(ctx)
			// drop target pls.
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(ReplyMentMode("嗯哼哼～抽到了自己，然而 Lucy 还是将双方写成一个人哦w （笑w ", ctx.Event.UserID, 1, ctx)))
			generatePairKey := GenerateMD5(ctx.Event.UserID, ctx.Event.UserID, ctx.Event.GroupID)
			err := InsertUserGlobalMarryList(marryList, ctx.Event.GroupID, ctx.Event.UserID, ctx.Event.UserID, 3, generatePairKey)
			if err != nil {
				panic(err)
			}
		}
		// get failed possibility.
		returnNumber := GetSomeRanDomChoiceProps(ctx)
		switch {
		case returnNumber == 1:
			GlobalCDModelCost(ctx)
			getSuccessMsg := dict["success"][rand.Intn(len(dict["success"]))]
			// normal mode. nothing happened.
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(ReplyMentMode(getSuccessMsg, ChooseAPerson, 1, ctx)))
			generatePairKey := GenerateMD5(ctx.Event.UserID, ChooseAPerson, ctx.Event.GroupID)
			err := InsertUserGlobalMarryList(marryList, ctx.Event.GroupID, ctx.Event.UserID, ChooseAPerson, 1, generatePairKey)
			if err != nil {
				fmt.Print(err)
				return
			}
		case returnNumber == 2:
			GlobalCDModelCost(ctx)
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(ReplyMentMode("貌似很奇怪哦～因为某种奇怪的原因～1变成了0,0变成了1", ChooseAPerson, 0, ctx)))
			generatePairKey := GenerateMD5(ChooseAPerson, ctx.Event.UserID, ctx.Event.GroupID)
			err := InsertUserGlobalMarryList(marryList, ctx.Event.GroupID, ChooseAPerson, ctx.Event.UserID, 2, generatePairKey)
			if err != nil {
				panic(err)
			}
		// reverse Target Mode
		case returnNumber == 3:
			GlobalCDModelCost(ctx)
			// drop target pls.
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(ReplyMentMode("嗯哼哼～发生了一些错误～本来应当抽到别人的变成了自己～所以", ctx.Event.UserID, 1, ctx)))
			generatePairKey := GenerateMD5(ctx.Event.UserID, ctx.Event.UserID, ctx.Event.GroupID)
			err := InsertUserGlobalMarryList(marryList, ctx.Event.GroupID, ctx.Event.UserID, ctx.Event.UserID, 3, generatePairKey)
			if err != nil {
				panic(err)
			}
		// you became your own target
		case returnNumber == 6:
			GlobalCDModelCost(ctx)
			// now no wife mode.
			getHideMsg := dict["hide_mode"][rand.Intn(len(dict["hide_mode"]))]
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(getHideMsg, "\n貌似没有任何反馈～"))
			generatePairKey := GenerateMD5(ctx.Event.UserID, ChooseAPerson, ctx.Event.GroupID)
			err := InsertUserGlobalMarryList(marryList, ctx.Event.GroupID, ctx.Event.UserID, ctx.Event.UserID, 6, generatePairKey)
			if err != nil {
				panic(err)
			}
		}
	})

	engine.OnRegex(`^(娶|嫁)\[CQ:at,qq=(\d+)\]`, zero.OnlyGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		choice := ctx.State["regex_matched"].([]string)[1]
		fiancee, _ := strconv.ParseInt(ctx.State["regex_matched"].([]string)[2], 10, 64)
		uid := ctx.Event.UserID
		// fast check
		if !CheckTheUserStatusAndDoRepeat(ctx) {
			return
		}
		if GlobalCDModelCostLeastReply(ctx) == 0 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("今天的机会已经使用完了哦～12小时后再来试试吧"))
			return
		}
		// check others.
		if uid == fiancee {
			switch rand.Intn(5) {
			case 1:
				GlobalCDModelCost(ctx)
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(ReplyMentMode("貌似Lucy故意添加了 --force 的命令，成功了(笑 ", ctx.Event.UserID, 1, ctx)))
				generatePairKey := GenerateMD5(ctx.Event.UserID, ctx.Event.UserID, ctx.Event.GroupID)
				err := InsertUserGlobalMarryList(marryList, ctx.Event.GroupID, ctx.Event.UserID, ctx.Event.UserID, 3, generatePairKey)
				if err != nil {
					panic(err)
				}
			default:
				GlobalCDModelCost(ctx)
				ctx.SendChain(message.Text("笨蛋！娶你自己干什么a"))
			}
			return
		}
		// However Lucy is only available to be married. LOL.
		if fiancee == ctx.Event.SelfID {
			// not work yet, so just the next path.
			if rand.Intn(100) > 90 {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("笨蛋！不准娶~ ama"))
				GlobalCDModelCost(ctx)
				return
			} else {
				// do it.
				GlobalCDModelCost(ctx)
				getSuccessMsg := dict["success"][rand.Intn(len(dict["success"]))]
				// normal mode. nothing happened.
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(ReplyMentMode(getSuccessMsg, fiancee, 1, ctx)))
				generatePairKey := GenerateMD5(ctx.Event.UserID, fiancee, ctx.Event.GroupID)
				err := InsertUserGlobalMarryList(marryList, ctx.Event.GroupID, ctx.Event.UserID, fiancee, 1, generatePairKey)
				if err != nil {
					fmt.Print(err)
					return
				}
				return
			}
		}
		switch {
		case choice == "娶":
			ResuitTheReferUserAndMakeIt(ctx, dict, uid, fiancee)
		case choice == "嫁":
			ResuitTheReferUserAndMakeIt(ctx, dict, fiancee, uid)
		}
	})

	engine.OnFullMatch("我要离婚", zero.OnlyToMe, zero.OnlyGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getStatusCode, _ := CheckTheUserIsTargetOrUser(marryList, ctx, ctx.Event.UserID)
		if getStatusCode == -1 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("貌似？没有对象的样子x"))
			return
		}
		if LeaveCDModelCostLeastReply(ctx) == 0 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("今天的次数已经用完了哦～或许可以试一下别的方式？"))
			return
		}
		LeaveCDModelCost(ctx)
		getlostSuccessedMsg := dict["lost_success"][rand.Intn(len(dict["lost_success"]))]
		getLostFailedMsg := dict["lost_failed"][rand.Intn(len(dict["lost_failed"]))]
		if rand.Intn(4) >= 2 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(getLostFailedMsg))
		} else {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(getlostSuccessedMsg))
			getPairKey := CheckThePairKey(marryList, ctx.Event.UserID, ctx.Event.GroupID)
			_ = RemoveUserGlobalMarryList(marryList, getPairKey, ctx.Event.GroupID)
		}
	})

	engine.OnRegex(`^试着骗(\[CQ:at,qq=(\d+)\]\s?|(\d+))做我的老婆`, zero.OnlyGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		fid := ctx.State["regex_matched"].([]string)
		fiancee, _ := strconv.ParseInt(fid[2]+fid[3], 10, 64)
		uid := ctx.Event.UserID
		if fiancee == uid {
			ctx.SendChain(message.Text("要骗别人哦~为什么要骗自己呢x"))
			return
		}
		// this case should other people existed.

	})

}
