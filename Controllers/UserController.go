package Controllers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"tgBotVeloBot/database"
	"tgBotVeloBot/db_models"
)

func CreateUser(update tgbotapi.Update) {
	if update.Message == nil {
		log.Print("Кто то пытается сломать бота")
	}
	var user db_models.User
	user.UserId = strconv.FormatInt(update.Message.From.ID, 10)
	user.UserName = update.Message.From.FirstName
	var db = database.Connect()
	db.Create(&user)
}
