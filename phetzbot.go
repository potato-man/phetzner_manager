package main

import (
   "context"
   "fmt"
   "log"
   "os"
   "time"

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
      return c.Send("Извините, у вас нет доступа.")
   }

   serversButton := tele.InlineButton{
      Unique: "servers_btn",
      Text:   "Список серверов",
   }

   inlineKeys := [][]tele.InlineButton{
      {serversButton},
   }

   return c.Send("Привет! Что вы хотите сделать?", &tele.ReplyMarkup{InlineKeyboard: inlineKeys})
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

      return c.Edit("Что вы хотите сделать?", &tele.ReplyMarkup{InlineKeyboard: inlineKeys})
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

func main() {
   telegramToken := os.Getenv("TELEGRAM_TOKEN")
   hcloudToken := os.Getenv("HCLOUD_TOKEN")
   allowedUserID := int64({TELEGRAM_ID}) // Замените на свой ID

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

   for i := 0; i < 10; i++ { // Ограничение на 10 серверов, вы можете изменить это значение
      bot.Handle(&tele.InlineButton{Unique: fmt.Sprintf("server_%d", i)}, botHandler.HandleServerActions(i))
   }

   bot.Handle(&tele.InlineButton{Unique: "power_on"}, botHandler.HandlePowerOn)
   bot.Handle(&tele.InlineButton{Unique: "power_off"}, botHandler.HandlePowerOff)

   // Запуск бота
   bot.Start()
}
