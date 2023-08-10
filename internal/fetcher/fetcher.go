package fetcher

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/zumosik/telegram-news-go/internal/model"
	"github.com/zumosik/telegram-news-go/internal/source"
	"github.com/zumosik/telegram-news-go/internal/storage"

	"github.com/tomakado/containers/set"
)

type Source interface {
	ID() int64
	Name() string
	Fetch(ctx context.Context) ([]model.Item, error)
}

type Fetcher struct {
	articles storage.ArticleStorage
	sources  storage.SourceStorage

	fetchInterval  time.Duration
	filterKeywords []string
}

func New(aStorage storage.ArticleStorage, sStorage storage.SourceStorage, fetchInterval time.Duration, filterKeywords []string) *Fetcher {
	return &Fetcher{
		articles:       aStorage,
		sources:        sStorage,
		fetchInterval:  fetchInterval,
		filterKeywords: filterKeywords,
	}
}

func (f *Fetcher) Start(ctx context.Context) error {
	ticker := time.NewTicker(f.fetchInterval)
	defer ticker.Stop()

	if err := f.Fetch(ctx); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			if err := f.Fetch(ctx); err != nil {
				return err
			}
		}
	}
}

func (f *Fetcher) Fetch(ctx context.Context) error {
	sourcer, err := f.sources.Sources(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, src := range sourcer {
		wg.Add(1)

		rssSource := source.NewRSSSourceFromModel(src)

		go func(source Source) {
			defer wg.Done()

			items, err := source.Fetch(ctx)
			if err != nil {
				log.Printf("can't fetch items from source: %s", err.Error())
				return
			}

			if err := f.processItems(ctx, source, items); err != nil {
				log.Printf("can't process items from source: %s", err.Error())
				return
			}

		}(rssSource)
	}

	wg.Wait()

	return nil
}

func (f *Fetcher) processItems(ctx context.Context, source Source, items []model.Item) error {
	for _, item := range items {
		item.Date = item.Date.UTC()

		if f.itemShouldBeSkipped(item) {
			continue
		}

		if err := f.articles.Store(ctx, model.Article{
			SourceID: source.ID(),
			Title:    item.Title,
			Link:     item.Link,
			Summary:  item.Summary,
			PostedAt: item.Date,
		}); err != nil {
			return err
		}
	}

	return nil
}

// itemShouldBeSkipped checks if item keywords contains filterKeywords, or have filterKeywords in title.
func (f *Fetcher) itemShouldBeSkipped(item model.Item) bool {
	categoriesSet := set.New(item.Categories...)

	for _, keyword := range f.filterKeywords {
		titleContainsKeyword := strings.Contains(strings.ToLower(item.Title), keyword)
		if categoriesSet.Contains(keyword) || titleContainsKeyword {
			return true
		}
	}

	return false
}
