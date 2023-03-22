package handler

import (
	"fmt"

	"emby_library_manager/lib"
	tele "gopkg.in/telebot.v3"
)

const successMsg = "刷新成功！"

func RefreshItems(c tele.Context) error {
	ids := c.Args()
	var pollingIds []string
	text := ""

	if len(ids) == 0 {
		return c.Reply("未指定要刷新的项目！")
	}

	if ids[0] == "all" {
		if items, err := getLibraryItems(); err != "" {
			text = "Error: " + err + "\n获取全部 Library 失败！"
		} else {
			for i, item := range items {
				result := refresh(item.ItemId)
				text += fmt.Sprintf("*%d. %s*(`%s`)\n%s\n\n", i+1, item.Name, item.ItemId, result)
				if result == successMsg {
					pollingIds = append(pollingIds, item.ItemId)
				}
			}
		}
	} else {
		for i, id := range ids {
			result := refresh(id)
			text += fmt.Sprintf("*%d. %s*\n%s\n\n", i+1, id, result)
			if result == successMsg {
				pollingIds = append(pollingIds, id)
			}
		}
	}

	timerStart(pollingIds)
	return c.Reply(text, tele.ModeMarkdown)
}

func refresh(id string) string {
	resp, err := lib.Resty.R().
		EnableTrace().
		SetHeader("X-Emby-Token", EmbyToken).
		SetQueryParam("Recursive", "true").
		SetQueryParam("ImageRefreshMode", "Default").
		SetQueryParam("MetadataRefreshMode", "Default").
		SetQueryParam("ReplaceAllImages", "false").
		SetQueryParam("ReplaceAllMetadata", "false").
		Post(EmbyUrl + "/Items/" + id + "/Refresh")

	if err != nil {
		return "Error: 发送请求失败！"
	} else {
		if resp.StatusCode() != 204 {
			return fmt.Sprintf("Error: %d %s", resp.StatusCode(), resp.String())
		}
	}
	return successMsg
}
