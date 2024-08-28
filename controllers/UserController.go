package controllers

import (
	"VeloBotReborn/dbModels"
	"VeloBotReborn/repositories"
	"VeloBotReborn/utils"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"sort"
	"strconv"
)

var medals = map[int]string{
	0: "🥇 1 место:",
	1: "🥈 2 место:",
	2: "🥉 3 место:",
}

func CreateUser(message *gotgbot.Message) error {
	var user dbModels.User
	user.UserId = strconv.FormatInt(message.From.Id, 10)
	user.UserName = message.From.FirstName
	return repositories.InsertUser(&user)
}

func AddResult(id int64, speed float64, distance float64) (*dbModels.User, error) {
	userId := strconv.FormatInt(id, 10)
	return repositories.AddResult(userId, speed, distance)
}

func GetResultsForUser(id int64) (string, error) {
	userId := strconv.FormatInt(id, 10)
	results, err := repositories.GetResults(userId)
	if err != nil {
		return "", err
	}
	message := "Результаты:\n\n"
	message += "Скорость:\n"
	maxI := 3
	if len(results) < 3 {
		maxI = len(results)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].MaxSpeed > results[j].MaxSpeed
	})
	for i := 0; i < maxI; i++ {
		message += fmt.Sprintf("%s %.2f км\\ч\n", medals[i], utils.ToFixedPrecision(results[i].MaxSpeed, 2))
	}
	message += "\nДистанция:\n"
	sort.Slice(results, func(i, j int) bool {
		return results[i].Distance > results[j].Distance
	})
	for i := 0; i < maxI; i++ {
		message += fmt.Sprintf("%s %.2f км\n", medals[i], utils.ToFixedPrecision(results[i].Distance, 2))
	}
	message += fmt.Sprintf("\nВсего: %.2f км", func() float64 {
		var sum float64 = 0
		for _, num := range results {
			sum += num.Distance
		}
		return utils.ToFixedPrecision(sum, 2)
	}())
	return message, nil
}

func GetAllResults() (string, error) {
	results, err := repositories.GetAllResults()
	if err != nil {
		return "", err
	}
	message := "Результаты:\n\n"
	message += "Скорость:\n"
	sort.Slice(results, func(i, j int) bool {
		return results[i].MaxSpeed > results[j].MaxSpeed
	})
	maxI := 3
	if len(results) < 3 {
		maxI = len(results)
	}
	for i := 0; i < maxI; i++ {
		message += fmt.Sprintf("%s %s %.2f км\\ч\n", medals[i], results[i].User.UserName, utils.ToFixedPrecision(results[i].MaxSpeed, 2))
	}
	message += "\nДистанция:\n"
	sort.Slice(results, func(i, j int) bool {
		return results[i].Distance > results[j].Distance
	})
	for i := 0; i < maxI; i++ {
		message += fmt.Sprintf("%s %s %.2f км\n", medals[i], results[i].User.UserName, utils.ToFixedPrecision(results[i].Distance, 2))
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].SumDistance > results[j].SumDistance
	})
	message += "\nВсего:\n"
	for i := 0; i < maxI; i++ {
		message += fmt.Sprintf("%s %s %.2f км\n", medals[i], results[i].User.UserName, utils.ToFixedPrecision(results[i].SumDistance, 2))
	}
	return message, nil
}
