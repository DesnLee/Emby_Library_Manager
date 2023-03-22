package handler

import (
	"fmt"
	"strconv"
	"time"

	Bot "emby_library_manager/bot"
	tele "gopkg.in/telebot.v3"
)

type ChatID int64

func (i ChatID) Recipient() string {
	return strconv.FormatInt(int64(i), 10)
}

func timerStart(ids []string) {
	for _, id := range ids {
		// 创建一个定时器每 30 秒调用接口轮询，检查接口返回内容，如果已完成则停止轮询，销毁定时器，如果未完成则继续轮询
		ticker := time.NewTicker(30 * time.Second)
		go func(t *time.Ticker, i string) {
			for {
				select {
				case <-t.C:
					itemList, err := getLibraryItems()
					if err != "" {
						fmt.Println("Error: 发送轮询请求失败！")
					} else {
						for _, v := range itemList {
							if v.ItemId == i {
								if v.RefreshStatus == "Idle" {
									fmt.Println(i, " 刷新完成，停止轮询...")
									t.Stop()
									text := fmt.Sprintf("%s 刷新完成！", v.Name)
									recipient := ChatID(NotifyChatId)
									if _, err := Bot.B.Send(recipient, text, tele.ModeMarkdown); err != nil {
										return
									}
								} else if v.RefreshStatus == "Active" {
									fmt.Println(i, "正在刷新")
								} else if v.RefreshStatus == "Queued" {
									fmt.Println(i, "排队中")
								}
							}
						}
					}
				}
			}
		}(ticker, id)
	}
}
