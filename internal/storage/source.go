package storage

import (
	"context"

	"github.com/zumosik/telegram-news-go/internal/model"
)

type SourceStorage interface {
	Sources(ctx context.Context) ([]model.Source, error)
	SourceByID(ctx context.Context, id int64) (*model.Source, error)
	Add(ctx context.Context, source model.Source) (int64, error)
	Delete(ctx context.Context, id int64) error
}
