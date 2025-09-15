package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GetActiveIcon(active bool) string {
	icon := "⛔️"
	if !active {
		icon = "✅"
	}

	return icon
}

func GetViewList(tds []Todo) []string {
	msgArr := []string{}
	loc, _ := time.LoadLocation("Europe/Moscow")
	for _, td := range tds {
		icon := GetActiveIcon(td.Active)
		timeItem := td.CreatedAt.In(loc).Format("15:04 02.01.2006")
		s := fmt.Sprintf("✓%d ~ %v ~ %s", td.Num, timeItem, icon) + "\n" + td.Description
		msgArr = append(msgArr, s)
	}

	return msgArr
}

func GetEmptyList() []string {
	return []string{
		"📝 Список дел ☝️",
		"🙏 Список пуст, добавьте задачу 👍",
		"💬 Чтобы добавить, напишите сообщение 👇",
	}
}

func GetButtonSending(user *User) (string, string, string) {
	nameButton := "❌ Выключить рассылку"
	valueButton := fmt.Sprintf("sending_false_%d", user.ID)
	callbackButton := "👍 Ваша рассылка включена."

	if !user.Sending {
		nameButton = "✅ Включить рассылку"
		valueButton = fmt.Sprintf("sending_true_%d", user.ID)
		callbackButton = "✋ Ваша рассылка выключена."
	}

	return nameButton, valueButton, callbackButton
}

func GetButtonPaging(count int, page int, pageSize int) (int, int, int) {
	pageCount := count / pageSize
	if (count % pageSize) > 0 {
		pageCount += 1
	}

	previous := page - 1
	nexts := page + 1

	if previous == 0 {
		previous = pageCount
	}

	if nexts > pageCount {
		nexts = 1
	}

	return pageCount, previous, nexts
}

func GetCallbackSending(data string) (uint, bool, error) {
	commandParams := strings.Split(data, "_")
	if len(commandParams) < 3 {
		return 0, false, errors.New("command is not array")
	}

	uid64, err := strconv.Atoi(commandParams[2])
	if err != nil {
		return 0, false, err
	}

	boolValue, err := strconv.ParseBool(commandParams[1])
	if err != nil {
		return 0, false, err
	}

	return uint(uid64), boolValue, err
}

func GetCallbackPaging(data string) (string, uint, int, error) {
	commandParams := strings.Split(data, "_")
	typeButton := commandParams[1]
	if len(commandParams) < 4 {
		return typeButton, 0, 0, errors.New("command is not array")
	}

	p, err := strconv.Atoi(commandParams[3])
	if err != nil {
		return typeButton, 0, 0, err
	}

	uid, err := strconv.Atoi(commandParams[2])
	if err != nil {
		return typeButton, 0, 0, err
	}

	u := uint(uid)

	return typeButton, u, p, nil
}

func GetCallbackTitle(typeButton string) string {
	var title string
	switch {
	case typeButton == "next":
		title = "предыдущее " + "👉"
	case typeButton == "previous":
		title = "👈" + " следующие"
	case typeButton == "page":
		title = "🔄⏳" + " обновление списка"
	}

	return title
}
