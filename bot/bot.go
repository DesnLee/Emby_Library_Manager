package Bot

import (
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	tele "gopkg.in/telebot.v3"
)

var B *tele.Bot

func init() {
	botToken := os.Getenv("BOT_TOKEN")

	pref := tele.Settings{
		Token:  botToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	// 初始化 BOT
	if B == nil {
		bot, err := tele.NewBot(pref)
		if err != nil {
			log.Fatal(err)
			return
		}
		B = bot
	}
}
