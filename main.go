package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	Bot "emby_library_manager/bot"
	"emby_library_manager/handler"
	_ "github.com/joho/godotenv/autoload"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

var (
	adminList []int64
)

func init() {
	for _, v := range strings.Split(os.Getenv("BOT_ADMIN_LIST"), ",") {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			fmt.Println("id 解析失败，请检查 .env 的 BOT_ADMIN_LIST >> ", err)
		} else {
			adminList = append(adminList, id)
		}
	}
}

func main() {
	// 全局中间件
	Bot.B.Use(middleware.Recover())
	Bot.B.Use(middleware.AutoRespond())

	admin := Bot.B.Group()
	admin.Use(middleware.Whitelist(adminList...))

	// BOT 功能
	admin.Handle("/chat", handler.Chat)
	admin.Handle("/items", handler.SearchItems)
	admin.Handle("/refresh", handler.RefreshItems)
	admin.Handle("/cancel", handler.CancelScan)
	admin.Handle(tele.OnCallback, handler.OnCallback)

	Bot.B.Start()
}
