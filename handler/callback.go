package handler

import (
	"strings"

	tele "gopkg.in/telebot.v3"
)

func OnCallback(c tele.Context) error {
	// 获取回调数据
	q := strings.Split(c.Data(), "|")
	a := strings.ReplaceAll(q[0], "\f", "")
	d := q[1:]

	switch a {
	case pageTurning:
		return TurnPage(c, d[0])
	}

	return c.Respond()
}
