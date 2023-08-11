package notifier

import (
	"context"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/zumosik/telegram-news-go/internal/botkit/markup"
	"github.com/zumosik/telegram-news-go/internal/model"
	"github.com/zumosik/telegram-news-go/internal/storage"
)

type Notifier struct {
	articles     storage.ArticleStorage
	bot          *tgbotapi.BotAPI
	sendInterval time.Duration
	lookUpWindow time.Duration
	channelID    int64
}

func New(
	articles storage.ArticleStorage,
	bot *tgbotapi.BotAPI,
	sendInterval time.Duration,
	lookUpWindow time.Duration,
	channelID int64,
) *Notifier {
	return &Notifier{
		articles:     articles,
		bot:          bot,
		sendInterval: sendInterval,
		lookUpWindow: lookUpWindow,
		channelID:    channelID,
	}
}

func (n *Notifier) Start(ctx context.Context) error {
	ticker := time.NewTicker(n.sendInterval)
	defer ticker.Stop()

	if err := n.SelectAndSendArticle(ctx); err != nil {
		return err
	}

	for {
		select {
		case <-ticker.C:
			log.Println("send 1")
			if err := n.SelectAndSendArticle(ctx); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (n *Notifier) SelectAndSendArticle(ctx context.Context) error {
	topOneArticles, err := n.articles.AllNotPosted(ctx, time.Now().Add(-n.lookUpWindow), 1)
	if err != nil {
		log.Print(err)
		return err
	}

	if len(topOneArticles) <= 0 {
		log.Print("len = 0")
		return nil
	}

	article := topOneArticles[0]

	if err := n.sendArticle(article); err != nil {
		log.Print(err)
		return err
	}

	log.Println("send 2")

	return n.articles.MarkAsPosted(ctx, article)
}

func (n *Notifier) sendArticle(article model.Article) error {
	log.Println(article.Title)

	const msgFormat = "*%s*\n%s"

	msg := tgbotapi.NewMessage(n.channelID, fmt.Sprintf(msgFormat, markup.EscapeForMarkdown(article.Title), markup.EscapeForMarkdown(article.Link)))
	msg.ParseMode = tgbotapi.ModeMarkdownV2

	_, err := n.bot.Send(msg)
	if err != nil {
		return err
	}

	log.Println("send 3")

	return nil
}
