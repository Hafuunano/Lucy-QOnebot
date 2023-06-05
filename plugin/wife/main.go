package wife

import (
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
)

var engine = control.Register("wife", &ctrl.Options[*zero.Ctx]{
	DisableOnDefault:  false,
	Help:              "Hi NekoPachi!\n说明书: https://lucy.impart.icu",
	PrivateDataFolder: "score",
})

/*
StatusID:

Type 1: Normal Mode, nothing happened.

Type 2: Cannot be the Target, Target became initative, so reverse.

(However the target and the initative should be in their position, DO NOT CHANGE. )

Type 3: Something is wrong, you are Target == initative Person.

Type 4: Removed.
(When User get others person. || IF REMARRIED, CHANGE IT TO TYPE1.) || (Be check more time to reduce to err.)

Type 5: NTR Mode

(Tips: NTR means changed their pairkey & TargetID || UserID, need to do some changes. ) ||
(Attempt to do once more every person.)

Type 6: No wife Mod?
Fake - Invisibile person here.
(Lucy Hides this and shows it in the next time if a person uses NTR,
shows nothing, and Lucy will make it for joke. LMAO)
*/
func init() {
	engine.OnFullMatch("娶群友").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		uid := ctx.Event.UserID
		gid := ctx.Event.GroupID
		getStatusCode := CheckTheUserIsTargetOrUser(marryList, ctx) // 判断这个user是否已经和别人在一起了
		switch {
		case getStatusCode == 0:
			// case target mode (0)
			return
		case getStatusCode == 1:
			// case user mode (1)
			return
		default:
		}
		// ok , go next.
		ChooseAPerson := GetUserListAndChooseOne(ctx)
		getStatusIfBannned := CheckTheUserIsInBlackListOrGroupList(uid, ChooseAPerson, gid)
		if getStatusIfBannned {
			// blocked.
			return
		}
		// go next. do something colorful.
		returnNumber := GetSomeRanDomChoiceProps(ctx)
		switch {
		case returnNumber == 1:
		// normal mode.
		case returnNumber == 2:
		// reverse Target Mode
		case returnNumber == 3:
		// you became your own target
		case returnNumber == 6:
			// now no wife mode.
		}
	})

}
