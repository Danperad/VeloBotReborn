package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"tgBotVeloBot/Controllers"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("7222803242:AAGJsA3X9etgUr8cyd2GP5FkMkO-7ZQK1Ks")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Получите канал апдейтов
	updates := bot.GetUpdatesChan(u)

	// Обрабатывайте апдейты
	for update := range updates {
		if update.Message != nil { // Если сообщение не пустое
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "start":
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ну привет")
					bot.Send(msg)
					go Controllers.CreateUser(update)
				default:
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда")
					bot.Send(msg)
				}
			}
		}
	}
}
