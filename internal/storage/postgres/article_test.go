package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/zumosik/telegram-news-go/internal/model"
)

func CreateTestArticlePostgresStorage(t *testing.T) *ArticlePostgresStorage {
	godotenv.Load()
	db, err := sqlx.Connect("postgres", os.Getenv("POSTGRES_TEST_STORAGE_URL")) // poot .env in this folder if doesnt work
	if err != nil {
		t.Error(err)
		t.FailNow()
		return nil
	}
	return NewArticleStorage(db)
}

func TestArticlePostgresStorage_Store(t *testing.T) {
	store := CreateTestArticlePostgresStorage(t)
	ctx := context.Background()
	type args struct {
		ctx     context.Context
		article model.Article
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				ctx: ctx,
				article: model.Article{
					SourceID:    4,
					Title:       "test",
					Summary:     "test summary",
					Link:        "https://",
					PublishedAt: time.Now(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := store.Store(tt.args.ctx, tt.args.article); (err != nil) != tt.wantErr {
				t.Errorf("ArticlePostgresStorage.Store() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TODO: add other tests
