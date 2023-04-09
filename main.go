package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	helpCmd    = "Just use /dashboard command."
	unknownCmd = "Unknown command"
	exPrefix   = "Курс USD to KZT: "
	diffToday  = "Изменения за сегодня: "
	mxSuffix   = "by Moscow Exchange"
	timenow    = "Дата: "
)

var (
	bot      *tgbotapi.BotAPI
	BotToken = "6257743064:AAEtG6mPn8mLot_tDEPhsuV9aJMwvDGy5pw"
)

type User struct {
	User *tgbotapi.User `json:"user"`
	Date time.Time      `json:"date"`
}

func main() {
	var err error
	bot, err = tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Panic(err)
	}

	log.Debugf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		msg := tgbotapi.MessageConfig{}
		if update.Message != nil {
			if !update.Message.IsCommand() {
				continue
			}

			msg = messageByCommand(update.Message.Chat.ID, update.Message.Command())
		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				log.Error(err)
			}
		}

		msg.ParseMode = tgbotapi.ModeMarkdown
		if _, err := bot.Send(msg); err != nil {
			log.Error(err)
		}
	}
}

func messageByCommand(chatId int64, command string) (m tgbotapi.MessageConfig) {
	m.ChatID = chatId
	
	client := NewClient()
	rate, difference, err := client.GetRate(USDKZT)
	if err != nil {
		panic(err)
	}

	dt := time.Now()

	switch command {
	case "start", "dashboard":
		m.Text = fmt.Sprintf("%s%.2f\n%s%.2f\n%s%s\n*%s*",
			exPrefix, rate, diffToday, difference, timenow, dt.Format("2006-01-02 3:4:5 pm"), mxSuffix)
	case "help":
		m.Text = helpCmd
	default:
		m.Text = unknownCmd
	}

	return
}
