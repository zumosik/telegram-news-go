package bot

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zumosik/telegram-news-go/internal/botkit"
	"github.com/zumosik/telegram-news-go/internal/model"
	"github.com/zumosik/telegram-news-go/internal/storage"
)

func ViewCmdSource(repo storage.SourceStorage) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		argsText := update.Message.CommandArguments() // hello, https://aaa
		args := strings.Split(argsText, ", ")         // ["hello", "https://aaa"]
		if len(args) < 2 {
			_, err := bot.Send(tgbotapi.NewMessage(update.Message.From.ID, invalidDataMsg))
			return err
		}

		id, err := repo.Add(ctx, model.Source{
			Name:     args[0],
			FeedURL:  args[1],
			Priority: 0,
		})

		if err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(update.Message.From.ID, fmt.Sprintf(sourceIsAddedMsg, id))
		msg.ParseMode = tgbotapi.ModeMarkdownV2

		_, err = bot.Send(msg)
		return err

	}
}
