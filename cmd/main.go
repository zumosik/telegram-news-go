package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/zumosik/telegram-news-go/internal/bot"
	"github.com/zumosik/telegram-news-go/internal/bot/middleware"
	"github.com/zumosik/telegram-news-go/internal/botkit"
	"github.com/zumosik/telegram-news-go/internal/config"
	"github.com/zumosik/telegram-news-go/internal/fetcher"
	"github.com/zumosik/telegram-news-go/internal/notifier"
	"github.com/zumosik/telegram-news-go/internal/storage/postgres"
)

func main() {
	botAPI, err := tgbotapi.NewBotAPI(config.Get().TelegramBotToken)
	if err != nil {
		log.Printf("failed to create bot: %s", err.Error())
		return
	}

	db, err := sqlx.Connect("postgres", config.Get().DatabaseDSN)
	if err != nil {
		log.Printf("failed to connect to db: %s", err.Error())
		return
	}
	defer db.Close()

	log.Println("Connected to db")

	var (
		articleStorage = postgres.NewArticleStorage(db)
		sourceStorage  = postgres.NewSourceStorage(db)

		fetcher = fetcher.New(
			articleStorage,
			sourceStorage,
			config.Get().FetchInterval,
			config.Get().FilterKeywords,
		)

		notifier = notifier.New(
			articleStorage,
			botAPI,
			config.Get().NotificationInterval,
			2*config.Get().FetchInterval,
			config.Get().TelegramChannelID,
		)
	)

	newsBot := botkit.New(botAPI)
	newsBot.RegisterCmdView("start", bot.ViewCmdStart()) // /start
	newsBot.RegisterCmdView("help", bot.ViewCmdHelp())   // /help
	newsBot.RegisterCmdView("add", middleware.AdminOnly(
		config.Get().TelegramChannelID,
		bot.ViewCmdSource(sourceStorage),
	)) // /add Name, URL
	newsBot.RegisterCmdView("list", middleware.AdminOnly(
		config.Get().TelegramChannelID,
		bot.ViewCmdListSources(sourceStorage),
	)) // /list
	newsBot.RegisterCmdView("delete", middleware.AdminOnly(
		config.Get().TelegramChannelID,
		bot.ViewCmdDelete(sourceStorage),
	)) // /delete 2

	// TODO: more admin commands

	log.Println("Registered commands")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Println("Starting notifier and fetcher")

	go func(ctx context.Context) {
		if err := fetcher.Fetch(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("failed to start fetcher: %w", err)
				return
			}

			log.Printf("fetcher stopped")
		}
	}(ctx)

	_ = notifier

	go func(ctx context.Context) {
		if err := notifier.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("failed to start notifier: %w", err)
				return
			}

			log.Printf("notifier stopped")
		}
	}(ctx)

	if err := newsBot.Run(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			log.Printf("failed to run bot: %v", err)
			return
		}

		log.Println("Bot stopped")
	}
}
