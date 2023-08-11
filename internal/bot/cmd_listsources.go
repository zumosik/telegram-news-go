package bot

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
	"github.com/zumosik/telegram-news-go/internal/botkit"
	"github.com/zumosik/telegram-news-go/internal/model"
	"github.com/zumosik/telegram-news-go/internal/storage"
)

func ViewCmdListSources(repo storage.SourceStorage) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		sources, err := repo.Sources(ctx)
		if err != nil {
			return err
		}

		var (
			sourceInfos = lo.Map(sources, func(source model.Source, _ int) string {
				return formatSource(source)
			})
			msgText = fmt.Sprintf(listSourcesMsg, len(sources), strings.Join(sourceInfos, "\n\n"))
		)

		msg := tgbotapi.NewMessage(update.Message.From.ID, msgText)

		_, err = bot.Send(msg)
		return err
	}

}

func formatSource(source model.Source) string {
	const msg = "🌐 %s\nID: %d\nURL: %s"
	res := fmt.Sprintf(msg,
		source.Name, source.ID, source.FeedURL)
	return res

}
