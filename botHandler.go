package main

import (
	"context"
	"fmt"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	tele "gopkg.in/telebot.v3"
)

type BotHandler struct {
	bot            *tele.Bot
	allowedUserID  int64
	hcloudClient   *hcloud.Client
	selectedServer *hcloud.Server
}

func NewBotHandler(bot *tele.Bot, allowedUserID int64, hcloudClient *hcloud.Client) *BotHandler {
	return &BotHandler{
		bot:           bot,
		allowedUserID: allowedUserID,
		hcloudClient:  hcloudClient,
	}
}

func (bh *BotHandler) CheckAuthorization(c tele.Context) bool {
	return c.Sender().ID == bh.allowedUserID
}

func (bh *BotHandler) HandleStart(c tele.Context) error {
	if !bh.CheckAuthorization(c) {
		return c.Send("Sorry, you don't have access right")
	}

	serversButton := tele.InlineButton{
		Unique: "servers_btn",
		Text:   "Список серверов",
	}

	inlineKeys := [][]tele.InlineButton{
		{serversButton},
	}

	return c.Send("Hello my lord! What do you want to do?", &tele.ReplyMarkup{InlineKeyboard: inlineKeys})
}

func (bh *BotHandler) HandleServerList(c tele.Context) error {
	servers, err := bh.hcloudClient.Server.All(context.Background())
	if err != nil {
		return c.Send(fmt.Sprintf("Ошибка при получении списка серверов: %s", err))
	}

	var serverButtons []tele.InlineButton
	for i, server := range servers {
		serverButtons = append(serverButtons, tele.InlineButton{
			Unique: fmt.Sprintf("server_%d", i),
			Text:   server.Name,
		})
	}

	inlineKeys := make([][]tele.InlineButton, len(serverButtons))
	for i, btn := range serverButtons {
		inlineKeys[i] = []tele.InlineButton{btn}
	}

	return c.Edit("Выберите сервер:", &tele.ReplyMarkup{InlineKeyboard: inlineKeys})
}

func (bh *BotHandler) HandleServerActions(index int) func(tele.Context) error {
	return func(c tele.Context) error {
		servers, err := bh.hcloudClient.Server.All(context.Background())
		if err != nil {
			return c.Send(fmt.Sprintf("Ошибка при получении списка серверов: %s", err))
		}

		if index < 0 || index >= len(servers) {
			return c.Send("Некорректный индекс сервера.")
		}

		bh.selectedServer = servers[index]

		inlineKeys := [][]tele.InlineButton{
			{
				{Unique: "power_on", Text: "Включить сервер"},
				{Unique: "power_off", Text: "Выключить сервер"},
			},
		}

		return c.Edit("What do you want to do?", &tele.ReplyMarkup{InlineKeyboard: inlineKeys})
	}
}

func (bh *BotHandler) HandlePowerOn(c tele.Context) error {
	_, _, err := bh.hcloudClient.Server.Poweron(context.Background(), bh.selectedServer)
	if err != nil {
		return c.Send(fmt.Sprintf("Ошибка при включении сервера: %s", err))
	}

	return c.Edit("Сервер запускается.")
}

func (bh *BotHandler) HandlePowerOff(c tele.Context) error {
	_, _, err := bh.hcloudClient.Server.Shutdown(context.Background(), bh.selectedServer)
	if err != nil {
		return c.Send(fmt.Sprintf("Ошибка при выключении сервера: %s", err))
	}

	return c.Edit("Сервер выключается.")
}
