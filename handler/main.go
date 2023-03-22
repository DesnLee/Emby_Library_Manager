package handler

import (
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

var (
	EmbyUrl        string
	EmbyToken      string
	EmbyScanTaskId string
	NotifyChatId   int64
)

func init() {
	EmbyUrl = os.Getenv("EMBY_URL")
	EmbyToken = os.Getenv("EMBY_TOKEN")
	EmbyScanTaskId = os.Getenv("EMBY_SCAN_TASK_ID")

	id, err := strconv.ParseInt(os.Getenv("BOT_SCAN_COMPLETE_NOTIFY_CHAT_ID"), 10, 64)
	if err != nil {
		panic(err)
	} else {
		NotifyChatId = id
	}
}
