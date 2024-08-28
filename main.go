package main

import (
	"VeloBotReborn/controllers"
	"VeloBotReborn/utils"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/conversation"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	FatalLog *log.Logger
)

var timeStart time.Time

func main() {
	InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	FatalLog = log.New(os.Stderr, "FATAL\t", log.Ldate|log.Ltime|log.Lshortfile)
	apiKey := os.Getenv("TELEGRAM_API_KEY")
	if apiKey == "" {
		FatalLog.Fatalln("TELEGRAM_API_KEY environment variable not set")
	}
	b, err := gotgbot.NewBot(apiKey, &gotgbot.BotOpts{})
	if err != nil {
		FatalLog.Fatalln(err)
	}

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			ErrorLog.Println(err)
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})

	updater := ext.NewUpdater(dispatcher, &ext.UpdaterOpts{})

	c := &client{}
	dispatcher.AddHandler(handlers.NewCommand("start", startHandler))
	dispatcher.AddHandler(handlers.NewMessage(myResultsMessage, myResultsHandler))
	dispatcher.AddHandler(handlers.NewMessage(totalResultsMessage, totalResultHandler))
	dispatcher.AddHandler(handlers.NewMessage(workTimeMessage, workTime))
	cancel := handlers.NewMessage(cancelMessage, cancelHandler)
	dispatcher.AddHandler(handlers.NewConversation(
		[]ext.Handler{handlers.NewMessage(writeResultMessage, writeResultHandler)},
		map[string][]ext.Handler{
			SPEED:    {handlers.NewMessage(noCommands, c.getSpeedHandler)},
			DISTANCE: {handlers.NewMessage(noCommands, c.getDistanceHandler)},
		},
		&handlers.ConversationOpts{
			Exits:        []ext.Handler{cancel},
			StateStorage: conversation.NewInMemoryStorage(conversation.KeyStrategySenderAndChat),
			AllowReEntry: true,
		},
	))
	dispatcher.AddHandler(cancel)
	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 60,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 60,
			},
		},
	})
	if err != nil {
		FatalLog.Fatalln(err)
	}

	InfoLog.Printf("%s has been started...\n", b.User.Username)
	timeStart = time.Now()
	updater.Idle()
}

const (
	SPEED    = "speed"
	DISTANCE = "distance"
)

const (
	CANCEL       = "Отменить"
	WRITERESULT  = "👀 Записать результат"
	TOTALRESULTS = "✋Общие результаты"
	MYRESULTS    = "👑Мои результаты"
	WORKTIME     = "❓ Сколько я работаю без сбоев"
)

func noCommands(msg *gotgbot.Message) bool {
	return message.Text(msg) && !message.Command(msg)
}

func writeResultMessage(msg *gotgbot.Message) bool {
	return !message.Command(msg) && message.Text(msg) && msg.Text == WRITERESULT
}

func totalResultsMessage(msg *gotgbot.Message) bool {
	return !message.Command(msg) && message.Text(msg) && msg.Text == TOTALRESULTS
}

func myResultsMessage(msg *gotgbot.Message) bool {
	return !message.Command(msg) && message.Text(msg) && msg.Text == MYRESULTS
}

func workTimeMessage(msg *gotgbot.Message) bool {
	return !message.Command(msg) && message.Text(msg) && msg.Text == WORKTIME
}

func cancelMessage(msg *gotgbot.Message) bool {
	return !message.Command(msg) && message.Text(msg) && msg.Text == CANCEL
}

func createReplyMarkup() *gotgbot.ReplyKeyboardMarkup {
	keyboard := gotgbot.ReplyKeyboardMarkup{}
	keyboard.IsPersistent = true
	keyboard.Selective = false
	keyboard.ResizeKeyboard = true
	rows := [][]gotgbot.KeyboardButton{
		{gotgbot.KeyboardButton{Text: WRITERESULT}, gotgbot.KeyboardButton{Text: MYRESULTS}},
		{gotgbot.KeyboardButton{Text: TOTALRESULTS}, gotgbot.KeyboardButton{Text: WORKTIME}},
	}
	keyboard.Keyboard = rows
	return &keyboard
}

func createCancelReplyMarkup() *gotgbot.ReplyKeyboardMarkup {
	keyboard := gotgbot.ReplyKeyboardMarkup{}
	keyboard.IsPersistent = true
	keyboard.Selective = false
	keyboard.ResizeKeyboard = true
	keyboard.Keyboard = [][]gotgbot.KeyboardButton{}
	var rowButtons []gotgbot.KeyboardButton
	rowButtons = append(rowButtons, gotgbot.KeyboardButton{Text: CANCEL})
	keyboard.Keyboard = append(keyboard.Keyboard, rowButtons)
	return &keyboard
}

func startHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	keyboard := createReplyMarkup()
	_, err := b.SendMessage(ctx.EffectiveChat.Id, "Вас приветствует шайтан бот", &gotgbot.SendMessageOpts{ParseMode: "HTML", ReplyMarkup: keyboard})
	if err != nil {
		return err
	}
	err = controllers.CreateUser(ctx.EffectiveMessage)
	if err != nil {
		return err
	}
	return nil
}

func writeResultHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := b.SendMessage(ctx.EffectiveChat.Id, "Ага молдоец, а теперь максимальную скорость введи", &gotgbot.SendMessageOpts{ParseMode: "HTML", ReplyMarkup: createCancelReplyMarkup()})
	if err != nil {
		return err
	}
	return handlers.NextConversationState(SPEED)
}

func myResultsHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	sendMessage, err := controllers.GetResultsForUser(ctx.EffectiveMessage.From.Id)
	if err != nil {
		return err
	}
	_, err = b.SendMessage(ctx.EffectiveChat.Id, sendMessage, &gotgbot.SendMessageOpts{ParseMode: "HTML"})
	if err != nil {
		return err
	}
	return nil
}

func totalResultHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	sendMessage, err := controllers.GetAllResults()
	if err != nil {
		return err
	}
	_, err = b.SendMessage(ctx.EffectiveChat.Id, sendMessage, &gotgbot.SendMessageOpts{ParseMode: "HTML"})
	if err != nil {
		return err
	}
	return nil
}

func workTime(b *gotgbot.Bot, ctx *ext.Context) error {
	timeNow := time.Now()
	_, err := b.SendMessage(ctx.EffectiveChat.Id, fmt.Sprintf("%d секунд", uint64(timeNow.Sub(timeStart).Seconds())), &gotgbot.SendMessageOpts{ParseMode: "HTML"})
	if err != nil {
		return err
	}
	return nil
}

func cancelHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	keyboard := createReplyMarkup()
	_, err := ctx.EffectiveMessage.Reply(b, "Ну и ладно", &gotgbot.SendMessageOpts{ParseMode: "HTML", ReplyMarkup: keyboard})
	if err != nil {
		return fmt.Errorf("failed to send cancelHandler message: %w", err)
	}
	return handlers.EndConversation()
}

func (c *client) getSpeedHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	speed, err := strconv.ParseFloat(ctx.Message.Text, 64)
	if err != nil {
		_, err = b.SendMessage(ctx.EffectiveChat.Id, "Жулик, не ломай, а теперь все сначала", &gotgbot.SendMessageOpts{ParseMode: "HTML"})
		return handlers.NextConversationState(SPEED)
	}
	c.setUserData(ctx, "speed", utils.ToFixedPrecision(speed, 2))
	_, err = b.SendMessage(ctx.EffectiveChat.Id, "Отлично! Теперь введи дистанцию:", &gotgbot.SendMessageOpts{ParseMode: "HTML"})
	return handlers.NextConversationState(DISTANCE)
}

func (c *client) getDistanceHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	distance, err := strconv.ParseFloat(ctx.Message.Text, 64)
	if err != nil {
		_, err = b.SendMessage(ctx.EffectiveChat.Id, "Жулик, не ломай, а теперь все сначала", &gotgbot.SendMessageOpts{ParseMode: "HTML"})
		return handlers.NextConversationState(DISTANCE)
	}
	keyboard := createReplyMarkup()
	speed, ok := c.getUserData(ctx, "speed")
	if !ok {
		_, err = b.SendMessage(ctx.EffectiveChat.Id, "Ну молодец, все сломал", &gotgbot.SendMessageOpts{ParseMode: "HTML", ReplyMarkup: keyboard})
		return handlers.EndConversation()
	}
	user, err := controllers.AddResult(ctx.Message.From.Id, utils.ToFixedPrecision(speed, 2), utils.ToFixedPrecision(distance, 2))
	if err != nil {
		_, err = b.SendMessage(ctx.EffectiveChat.Id, "Ну молодец, все сломал", &gotgbot.SendMessageOpts{ParseMode: "HTML", ReplyMarkup: keyboard})
		return handlers.EndConversation()
	}
	_, err = b.SendMessage(ctx.EffectiveChat.Id, fmt.Sprintf("Записал для %s, %.2f, %.2f", user.UserName, speed, distance), &gotgbot.SendMessageOpts{ParseMode: "HTML", ReplyMarkup: keyboard})
	return handlers.EndConversation()
}
