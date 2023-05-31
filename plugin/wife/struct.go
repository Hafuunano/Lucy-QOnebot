package wife

import (
	zero "github.com/wdvxdr1123/ZeroBot"
	"math/rand"
)

func init() {

}

/*
func getArrayListAndMakeSure(ctx *zero.Ctx) (groupIdStatus int, target int64) {
	// make sure every user will be chosen.
	/*
		method:
		1. status 200 --> everything is ok:P
		2. status 403 --> user cannot get this user, for it has others.
		{ TODO: check the user may have others list. };
		3. status 404 --> user didnt join this group. (setted by users.)
		4. status 401 --> unauth (blacklist)
		5. status 1 --> you are yourself :p

	getThisGroupList := ctx.GetGroupMemberList(ctx.Event.GroupID).Array()
	getThisUserID := getThisGroupList[rand.Intn(len(getThisGroupList))].Get("user_id").Int()
	currentHandlerUser := ctx.Event.UserID
	switch {
	case getThisUserID == currentHandlerUser:
		return 1, currentHandlerUser
	}
}
*/

func getUserListAndChooseOne(ctx *zero.Ctx) int64 {
	getThisGroupList := ctx.GetGroupMemberList(ctx.Event.GroupID).Array()
	getThisUserID := getThisGroupList[rand.Intn(len(getThisGroupList))].Get("user_id").Int()
	return getThisUserID
}
