package wife

import (
	"fmt"
	"github.com/FloatTech/floatbox/binary"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"hash/crc64"
	"strconv"
)

var engine = control.Register("wife", &ctrl.Options[*zero.Ctx]{
	DisableOnDefault:  false,
	Help:              "Hi NekoPachi!\n说明书: https://lucy.impart.icu",
	PrivateDataFolder: "score",
})

// usage :
/*
TODO:
1. get ArrayList;
2. make pair key, to check it quickly.
3. cdCheckSum
*/

func init() {
	engine.OnFullMatch("娶群友").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getHandleUser := ctx.Event.UserID
		getTargetUser := getUserListAndChooseOne(ctx)
		// find the key is existed?

		// first check the HandleUser.
		err := userlistDatabase.Find(strconv.FormatInt(ctx.Event.GroupID, 10), &UserGroup{}, "where handleid or targatid is "+strconv.FormatInt(getHandleUser, 10))
		if err != nil {
			// none. go next?

		}
		// err
		// TODO: Handle If One of them have.
		// As if it done, pause.

		// both of them is none,go next.
		id := int64(crc64.Checksum(binary.StringToBytes(fmt.Sprintf("%d_%d_%d", ctx.Event.GroupID, getHandleUser, getTargetUser)), crc64.MakeTable(crc64.ISO)))
		_ = JointListAdd(userlistDatabase, getHandleUser, ctx.Event.GroupID)
		_ = InsertUserMarriedMember(getHandleUser, getTargetUser, ctx, userlistDatabase, id)
		// key done.

	})

}
