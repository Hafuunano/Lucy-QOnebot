// Package main, Lucy, are you here?
package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/HafuuNano/Lucy-QOnebot/core"
	"github.com/joho/godotenv"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
)

func main() {
	_ = godotenv.Load() // load .env from current directory if present
	core.Init()
	nickNames := os.Getenv("NICK_NAMES")
	nickNameList := strings.Split(nickNames, ",")
	for i := range nickNameList {
		nickNameList[i] = strings.TrimSpace(nickNameList[i])
	}

	commandPrefix := os.Getenv("COMMAND_PREFIX")
	superUsersStr := os.Getenv("SUPER_USERS")
	var superUsers []int64
	for p := range strings.SplitSeq(superUsersStr, ",") {
		n, _ := strconv.ParseInt(strings.TrimSpace(p), 10, 64)
		if n != 0 {
			superUsers = append(superUsers, n)
		}
	}

	wsURL := os.Getenv("WS_URL")
	wsToken := os.Getenv("WS_TOKEN")

	zero.RunAndBlock(&zero.Config{
		NickName:      nickNameList,
		CommandPrefix: commandPrefix,
		SuperUsers:    superUsers,
		Driver: []zero.Driver{
			driver.NewWebSocketServer(16, wsURL, wsToken),
		},
	}, nil)
}
