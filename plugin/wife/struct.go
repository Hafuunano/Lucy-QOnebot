package wife

import (
	zero "github.com/wdvxdr1123/ZeroBot"
	"math/rand"
)

func init() {

}

func getUserListAndChooseOne(ctx *zero.Ctx) int64 {
	getThisGroupList := ctx.GetGroupMemberList(ctx.Event.GroupID).Array()
	getThisUserID := getThisGroupList[rand.Intn(len(getThisGroupList))].Get("user_id").Int()
	return getThisUserID
}
