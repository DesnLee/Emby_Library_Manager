package handler

import (
	"fmt"

	tele "gopkg.in/telebot.v3"
)

func Chat(c tele.Context) error {
	chat := c.Chat()
	text := fmt.Sprintf("Chat ID: `%v`", chat.ID)
	return c.Reply(text, tele.ModeMarkdown)
}
