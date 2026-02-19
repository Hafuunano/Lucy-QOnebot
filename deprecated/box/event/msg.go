package event

import (
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"time"
)

func WaitForNextMessage(ctx *zero.Ctx) message.Message {
	t := time.NewTimer(30 * time.Second)
	defer t.Stop()
	r, cancel := ctx.FutureEvent("message", ctx.CheckSession()).Repeat()
	defer cancel()
	select {
	case e := <-r:
		return e.Event.Message
	case <-t.C:
		return nil
	}
}
