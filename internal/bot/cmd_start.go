package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zumosik/telegram-news-go/internal/botkit"
)

func ViewCmdStart() botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) (err error) {
		_, err = bot.Send(tgbotapi.NewMessage(update.FromChat().ID, helloMsg))
		return
	}
}
