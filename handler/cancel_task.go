package handler

import (
	"fmt"

	"emby_library_manager/lib"
	tele "gopkg.in/telebot.v3"
)

func CancelScan(c tele.Context) error {
	tasks, err1 := getTaskList()
	if err1 != "" {
		return c.Reply(err1)
	}

	id := ""
	for _, task := range tasks {
		if task.Name == "Scan media library" {
			id = task.Id
			break
		}
	}

	if id == "" {
		return c.Reply("未找到扫描任务！")
	}

	if resp, err := lib.Resty.R().
		EnableTrace().
		SetHeader("X-Emby-Token", EmbyToken).
		Post(EmbyUrl + "/ScheduledTasks/Running/" + id + "/Delete"); err != nil {
		return c.Reply("Error: 请求取消任务失败！")
	} else {
		if resp.StatusCode() != 204 {
			return c.Reply(fmt.Sprintf("Error: %d %s", resp.StatusCode(), resp.String()))
		}
	}

	return c.Reply("取消扫描任务成功！可使用 /items 查看当前状态")
}

func getTaskList() (tasks []TaskStatus, error string) {
	resp, err := lib.Resty.R().
		EnableTrace().
		SetResult(&tasks).
		SetHeader("X-Emby-Token", EmbyToken).
		Get(EmbyUrl + "/ScheduledTasks")

	if err != nil {
		error = "Error: 请求 Task 列表失败！"
	} else {
		if resp.StatusCode() != 200 {
			error = fmt.Sprintf("Error: %d 获取 Task 列表失败！ %s", resp.StatusCode(), resp.String())
		}
	}
	return
}
