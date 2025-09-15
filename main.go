package main

import (
	"fmt"
	"go_todo/config"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
)

const (
	PL     string = "[+]"
	RS     string = " → "
	DS     string = " / "
	NS     string = "\n"
	PS     string = " "
	Suffix string = "~"
)

var (
	Emoji      []string = []string{"🌻", "🌶️", "🌵", "🚀", "👾", "🍎", "⚙️", "🎲", "🎯", "🏀", "⚽", "🎳", "♥️", "♠️", "♦️", "♣️"}
	SuffixLine string   = strings.Repeat(Suffix, 39)
)

// color: 1 red, 2 green, 3 yello, 4 blue, 5 purple, 6 blue
func p(color int, sep string, str ...any) {
	newStr := []any{}
	for index, v := range str {
		if index == 0 {
			newStr = append(newStr, v)
		} else {
			newStr = append(newStr, sep, v)
		}
	}

	suffixColor := "\033[3" + strconv.Itoa(color) + "m"
	fmt.Printf("%s%s%s", suffixColor, fmt.Sprint(newStr...), "\033[0m\n")
}

func connectWithTg(token string, url string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot.Debug = config.DEBUG
	p(3, " ~ ", PL, bot.Self.UserName, url)

	whUrl := url + "/" + token
	wh, _ := tgbotapi.NewWebhook(whUrl)
	wh.AllowedUpdates = []string{"message", "edited_channel_post", "callback_query"}
	_, err = bot.Request(wh)
	if err != nil {
		return nil, err
	}

	commandStart := tgbotapi.BotCommand{
		Command:     "start",
		Description: Emoji[3] + " Start bot",
	}

	commandHi := tgbotapi.BotCommand{
		Command:     "settings",
		Description: Emoji[6] + " The settings",
	}

	bc := tgbotapi.NewSetMyCommands(commandStart, commandHi)
	_, err = bot.Request(bc)
	if err != nil {
		return nil, err
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		return nil, err
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	return bot, nil
}

func setTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("it's ok, v" + config.VERSION))
}

func handleButton(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	// Извлечь данные обратного вызова
	data := callback.Data
	switch {
	case strings.HasPrefix(data, "paging_"):
		typeButton, u, p, err := config.GetCallbackPaging(data)
		if err != nil {
			return
		}

		msgArr, pagingBool, markup := config.GetTodoList(u, p, config.PAGE_SIZE, "paging")

		// Ответить на запрос обратного вызова
		callbackMess := tgbotapi.NewCallback(callback.ID, config.GetCallbackTitle(typeButton))
		bot.Request(callbackMess)

		// Опционально — отредактировать сообщение, чтобы отразить выбор
		edit := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, strings.Join(msgArr, NS+NS))
		if pagingBool {
			edit.ReplyMarkup = &markup
		}

		bot.Send(edit)

	case strings.HasPrefix(data, "sending_"):
		uid, boolValue, err := config.GetCallbackSending(data)
		if err != nil {
			return
		}

		user := config.SetUserSending(uid, boolValue)
		nameButton, valueButton, callbackButton := config.GetButtonSending(&user)

		// Ответить на запрос обратного вызова
		callbackMess := tgbotapi.NewCallback(callback.ID, callbackButton)
		bot.Request(callbackMess)

		// Опционально — отредактировать сообщение, чтобы отразить выбор
		markup := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(nameButton, valueButton),
			),
		)

		edit := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, callbackButton)
		edit.ReplyMarkup = &markup
		bot.Send(edit)
	}
}

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	p(2, " → ", PL, message.Chat.UserName, message.Chat.ID, message.Text)
	// ~~~ add user DB
	userName := message.From.UserName
	if message.Chat.Type == "group" {
		userName = message.Chat.Title
	}

	user := config.SetUser(message.Chat.ID, userName)
	// ~~~ end

	switch {
	case strings.HasPrefix(message.Text, "/start"):
		setStartCommand(bot, message, &user)

	case strings.HasPrefix(message.Text, "/settings"):
		setSettingsCommand(bot, message, &user)

	case message != nil && len(message.Text) > 0:
		setDefaultMessage(bot, message, &user)
	}
}

func setSettingsCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, user *config.User) {
	msgArr := []string{
		"📝 Cервис списка дел 👍",
		SuffixLine,
		fmt.Sprintf("🎯 Вы @%s", message.From.UserName),
		fmt.Sprintf("📌 Ваш id %d", message.From.ID),
		SuffixLine,
		fmt.Sprintf("🕜 %s", time.Now().Format("15:04 ~ 02.01.2006")),
		fmt.Sprintf("✉️ рассылка 👇 по часовому поясу %s", time.Now().Format("MST")),
		fmt.Sprintf("⏰ %s часы", "10,11,12,13,14,15,16,17,18,19"),
		SuffixLine,
		fmt.Sprintf("%s %s версия", Emoji[12], config.VERSION),
	}

	nameButton, valueButton, _ := config.GetButtonSending(user)

	msg := tgbotapi.NewMessage(message.Chat.ID, strings.Join(msgArr, NS))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(nameButton, valueButton),
		),
	)
	bot.Send(msg)
}

// command start
func setStartCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, user *config.User) {
	msgArr, pagingBool, markup := config.GetTodoList(user.ID, 1, config.PAGE_SIZE, "paging")

	msg := tgbotapi.NewMessage(message.Chat.ID, strings.Join(msgArr, NS+NS))
	if pagingBool {
		msg.ReplyMarkup = &markup
	}
	mes, _ := bot.Send(msg)

	pin := tgbotapi.PinChatMessageConfig{
		ChatID:              mes.Chat.ID,
		ChannelUsername:     mes.From.UserName,
		MessageID:           mes.MessageID,
		DisableNotification: true,
	}
	bot.Request(pin)
}

// default message
func setDefaultMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, user *config.User) {
	td := config.AddTodo(user, message.Text)
	msgArr := []string{
		fmt.Sprintf("Дело ✅%d успешно добавлено 👍", td.Num),
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, strings.Join(msgArr, NS))
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)

	// get PinnedMessage for update Todo
	ch := tgbotapi.ChatInfoConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: message.Chat.ID,
		},
	}
	chat, _ := bot.GetChat(ch)

	if chat.PinnedMessage != nil && chat.PinnedMessage.MessageID > 0 {
		msgArr, pagingBool, markup := config.GetTodoList(user.ID, 1, config.PAGE_SIZE, "paging")

		edit := tgbotapi.NewEditMessageText(message.Chat.ID, chat.PinnedMessage.MessageID, strings.Join(msgArr, NS+NS))
		if pagingBool {
			edit.ReplyMarkup = &markup
		}
		bot.Send(edit)
	}
}

// mailing
func sendMessageEveryone(bot *tgbotapi.BotAPI) {
	users := config.GetListUsers()

	for _, us := range users {
		msgArr := []string{
			fmt.Sprintf("%s %s %s → @%s", "🕜", time.Now().Format("15:04:05 ~ 02.01.2006"), Emoji[8], us.Name),
			SuffixLine,
		}
		msg := tgbotapi.NewMessage(us.Tid, strings.Join(msgArr, NS))
		bot.Send(msg)
	}
}

func setCroneStarted(bot *tgbotapi.BotAPI) *cron.Cron {
	c := cron.New(cron.WithSeconds())
	// A job to run every 15 seconds ~ "*/15 * * * * *"
	// A job to run every day at 07:30:00  "0 30 7 * * *"
	// job hour every day "0 0 10,11,12,13,14,15,16,17,18,19 * * *"
	c.AddFunc("0 0 10,11,12,13,14,15,16,17,18,19 * * *", func() {
		sendMessageEveryone(bot)
	})
	c.Start()
	p(5, " ~ ", PL, "Crone started", "🚀")

	return c
}

func main() {
	bot, err := connectWithTg(config.TOKEN, config.URL_BOT)
	if err != nil {
		log.Fatal(err)
	}

	c := setCroneStarted(bot)
	defer c.Stop()

	updates := bot.ListenForWebhook("/" + config.TOKEN)
	http.HandleFunc("/", setTest)
	go http.ListenAndServe(":8080", nil)

	for update := range updates {
		switch {
		// Handle messages
		case update.Message != nil:
			handleMessage(bot, update.Message)

		// Handle button clicks
		case update.CallbackQuery != nil:
			handleButton(bot, update.CallbackQuery)
		}
	}
}
