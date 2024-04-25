package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	tele "gopkg.in/telebot.v3"
)

func main() {
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	hcloudToken := os.Getenv("HCLOUD_TOKEN")
	allowedUserID := int64({TELEGRAM_ID})

	client := hcloud.NewClient(hcloud.WithToken(hcloudToken))

	pref := tele.Settings{
		Token:  telegramToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal("Ошибка при создании бота:", err)
		return
	}

	botHandler := NewBotHandler(bot, allowedUserID, client)

	bot.Handle("/start", botHandler.HandleStart)
	bot.Handle(&tele.InlineButton{Unique: "servers_btn"}, botHandler.HandleServerList)

	for i := 0; i < 10; i++ { // limitation 10 servers
		bot.Handle(&tele.InlineButton{Unique: fmt.Sprintf("server_%d", i)}, botHandler.HandleServerActions(i))
	}

	bot.Handle(&tele.InlineButton{Unique: "power_on"}, botHandler.HandlePowerOn)
	bot.Handle(&tele.InlineButton{Unique: "power_off"}, botHandler.HandlePowerOff)

	bot.Start()
}
