package storage

import (
	"context"
	"time"

	"github.com/zumosik/telegram-news-go/internal/model"
)

type ArticleStorage interface {
	Store(ctx context.Context, article model.Article) error
	AllNotPosted(ctx context.Context, since time.Time, limit uint64) ([]model.Article, error)
	MarkAsPosted(ctx context.Context, article model.Article) error
}
