package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"github.com/zumosik/telegram-news-go/internal/model"
)

type ArticlePostgresStorage struct {
	db *sqlx.DB
}

func NewArticleStorage(db *sqlx.DB) *ArticlePostgresStorage {
	return &ArticlePostgresStorage{db: db}
}

type dbArticle struct {
	ID             int64          `db:"id"`
	SourcePriority int64          `db:"priority"`
	SourceID       int64          `db:"source_id"`
	Title          string         `db:"title"`
	Link           string         `db:"link"`
	Summary        sql.NullString `db:"summary"`
	PublishedAt    time.Time      `db:"published_at"`
	PostedAt       sql.NullTime   `db:"posted_at"`
	CreatedAt      time.Time      `db:"created_at"`
}

func (s *ArticlePostgresStorage) Store(ctx context.Context, article model.Article) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(
		ctx,
		`INSERT INTO articles (source_id, title, link, summary, published_at)
	    				VALUES ($1, $2, $3, $4, $5)
	    				ON CONFLICT DO NOTHING;`,
		article.SourceID,
		article.Title,
		article.Link,
		article.Summary,
		article.PublishedAt,
	); err != nil {
		return err
	}

	return nil

}
func (s *ArticlePostgresStorage) AllNotPosted(ctx context.Context, since time.Time, limit uint64) ([]model.Article, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var articles []dbArticle
	query := "SELECT * FROM articles WHERE posted_at IS NULL ORDER BY published_at LIMIT $1"
	if err := conn.SelectContext(ctx, &articles, query, limit); err != nil {
		return nil, err
	}

	return lo.Map(articles, func(article dbArticle, _ int) model.Article {
		return model.Article{
			ID:          article.ID,
			SourceID:    article.SourceID,
			Title:       article.Title,
			Link:        article.Link,
			Summary:     article.Link,
			PublishedAt: article.PublishedAt,
			PostedAt:    article.PostedAt.Time,
			CreatedAt:   article.CreatedAt,
		}
	}), nil
}
func (s *ArticlePostgresStorage) MarkAsPosted(ctx context.Context, article model.Article) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(
		ctx,
		`UPDATE articles SET posted_at = $1::timestamp WHERE id = $2;`,
		time.Now().UTC().Format(time.RFC3339),
		article.ID,
	); err != nil {
		return err
	}

	return nil

}
