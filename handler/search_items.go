package handler

import (
	"fmt"
	"strconv"
	"time"

	"emby_library_manager/lib"
	tele "gopkg.in/telebot.v3"
)

const countPerPage = 7

var (
	status = map[string]string{
		"Idle":   "空闲",
		"Active": "刷新中...",
		"Queued": "排队中...",
	}

	overallStatus = map[string]string{
		"Idle":       "空闲",
		"Running":    "运行中...",
		"Cancelling": "正在停止",
	}

	pageTurning = "items_page_turning"
)

type EmbyLibrary struct {
	Name               string                 `json:"Name"`
	Locations          []string               `json:"Locations"`
	CollectionType     string                 `json:"CollectionType"`
	LibraryOptions     map[string]interface{} `json:"LibraryOptions"`
	ItemId             string                 `json:"ItemId"`
	Guid               string                 `json:"Guid"`
	PrimaryImageItemId string                 `json:"PrimaryImageItemId"`
	RefreshProgress    float64                `json:"RefreshProgress"`
	RefreshStatus      string                 `json:"RefreshStatus"`
}

type TaskStatus struct {
	Name                      string  `json:"Name"`
	State                     string  `json:"State"`
	CurrentProgressPercentage float64 `json:"CurrentProgressPercentage"`
	Id                        string  `json:"Id"`
	LastExecutionResult       struct {
		StartTimeUtc time.Time `json:"StartTimeUtc"`
		EndTimeUtc   time.Time `json:"EndTimeUtc"`
		Status       string    `json:"Status"`
		Name         string    `json:"Name"`
		Key          string    `json:"Key"`
		Id           string    `json:"Id"`
	} `json:"LastExecutionResult"`
	Triggers []struct {
		Type           string `json:"Type"`
		TimeOfDayTicks int64  `json:"TimeOfDayTicks"`
	} `json:"Triggers"`
	Description string `json:"Description"`
	Category    string `json:"Category"`
	IsHidden    bool   `json:"IsHidden"`
	Key         string `json:"Key"`
}

func SearchItems(c tele.Context) error {
	// 获取最新媒体库列表
	itemList, err := getLibraryItems()
	if err != "" {
		return c.Send("Error: " + err)
	}

	text, markup := generateMsg(itemList, 0)
	return c.Send(text, tele.ModeMarkdown, markup)
}

func TurnPage(c tele.Context, p string) error {
	// 获取最新媒体库列表
	itemList, err1 := getLibraryItems()
	if err1 != "" {
		return c.Send("Error: " + err1)
	}

	pageNum, err2 := strconv.Atoi(p)
	if err2 != nil {
		return c.Send("Error: " + err2.Error())
	}

	text, markup := generateMsg(itemList, pageNum)
	return c.Edit(text, tele.ModeMarkdown, markup)
}

func getLibraryItems() (items []EmbyLibrary, error string) {
	if resp, err := lib.Resty.R().
		EnableTrace().
		SetHeader("X-Emby-Token", EmbyToken).
		SetResult(&items).
		Get(EmbyUrl + "/Library/VirtualFolders"); err != nil {
		error = "发送请求失败！"
	} else {
		if resp.StatusCode() != 200 {
			error = fmt.Sprintf("%d %s", resp.StatusCode(), resp.String())
		}
	}
	return
}

func generateMsg(l []EmbyLibrary, num int) (text string, markup *tele.ReplyMarkup) {
	// 获取总扫描状态
	text = getScanStatus()

	// 计算页码数据
	pageNum := len(l)/countPerPage + 1
	prev := strconv.Itoa(num - 1)
	next := strconv.Itoa(num + 1)

	// 构造按钮
	markup = &tele.ReplyMarkup{}
	BtnPrev := markup.Data("⬅", pageTurning, []string{prev}...)
	BtnCurrent := markup.Data("", "current")
	BtnNext := markup.Data("➡", pageTurning, []string{next}...)
	BtnCurrent.Text = fmt.Sprintf("%d/%d", num+1, pageNum)

	// 根据当前页码添加按钮
	if num == 0 {
		markup.Inline(
			markup.Row(BtnCurrent, BtnNext),
		)
	} else if num+1 == pageNum {
		markup.Inline(
			markup.Row(BtnPrev, BtnCurrent),
		)
	} else {
		markup.Inline(
			markup.Row(BtnPrev, BtnCurrent, BtnNext),
		)
	}

	// 计算首尾索引
	startIndex := num * countPerPage
	var endIndex int
	if (num+1)*countPerPage > len(l) {
		endIndex = len(l)
	} else {
		endIndex = (num + 1) * countPerPage
	}
	displayItems := l[startIndex:endIndex]

	// 构造消息内容
	for i, l := range displayItems {
		info := fmt.Sprintf("*%d. %s*(`%s`)\n当前状态：%s", i+startIndex+1, l.Name, l.ItemId, status[l.RefreshStatus])
		if l.RefreshStatus == "Active" {
			info += fmt.Sprintf("\n刷新进度：%.2f%%", l.RefreshProgress)
		}
		text += info + "\n\n"
	}
	return
}

func getScanStatus() (result string) {
	overallScanStatus := TaskStatus{}
	if resp, err := lib.Resty.R().
		EnableTrace().
		SetHeader("X-Emby-Token", EmbyToken).
		SetResult(&overallScanStatus).
		Get(EmbyUrl + "/ScheduledTasks/" + EmbyScanTaskId); err != nil {
		result = "总扫描状态请求失败！"
	} else {
		if resp.StatusCode() != 200 {
			result = fmt.Sprintf("Error: 总扫描状态获取失败！%d %s", resp.StatusCode(), resp.String())
		} else {
			result = "总扫描状态：" + overallStatus[overallScanStatus.State] + "\n"
			if overallScanStatus.State == "Running" {
				result += fmt.Sprintf("总进度：%.2f%%\n", overallScanStatus.CurrentProgressPercentage)
			}
			result += "\n"
		}
	}
	return
}
