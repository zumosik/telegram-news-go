package bot

import (
	"context"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zumosik/telegram-news-go/internal/botkit"
	"github.com/zumosik/telegram-news-go/internal/storage"
)

func ViewCmdDelete(repo storage.SourceStorage) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		arg := update.Message.CommandArguments()
		arg = strings.TrimSpace(arg)
		if arg == "" {
			_, err := bot.Send(tgbotapi.NewMessage(update.Message.From.ID, invalidDataMsg))
			return err
		}

		id, err := strconv.Atoi(arg)
		if err != nil {
			return err
		}
		err = repo.Delete(ctx, int64(id))
		if err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, sorceIsDeleted)
		msg.ParseMode = tgbotapi.ModeMarkdownV2

		_, err = bot.Send(msg)
		return err

	}
}
