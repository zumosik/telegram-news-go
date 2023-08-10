package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"github.com/zumosik/telegram-news-go/internal/model"
)

var (
	errSourceNotFound = errors.New("source was not found")
)

type SourcePostgresStorage struct {
	db *sqlx.DB
}

type dbSource struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	FeedURL   string    `db:"feed_url"`
	CreatedAt time.Time `db:"created_at"`
	Priority  int       `db:"priority"`
}

func NewSourceStorage(db *sqlx.DB) *SourcePostgresStorage {
	return &SourcePostgresStorage{
		db: db,
	}
}

func (s *SourcePostgresStorage) Sources(ctx context.Context) ([]model.Source, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var sources []dbSource
	if err := conn.SelectContext(ctx, &sources, `SELECT * FROM sources`); err != nil {
		return nil, err
	}

	return lo.Map(sources, func(source dbSource, _ int) model.Source { return model.Source(source) }), nil
}

func (s *SourcePostgresStorage) SourceByID(ctx context.Context, id int64) (*model.Source, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var source dbSource
	if err := conn.GetContext(ctx, &source, "SELECT * FROM sources WHERE id = $1", id); err != nil {
		return nil, err
	}

	return (*model.Source)(&source), nil
}

func (s *SourcePostgresStorage) Add(ctx context.Context, source model.Source) (int64, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	query := "INSERT INTO sources (name, feed_url, created_at, priority) VALUES ($1, $2, $3, $4) RETURNING id"
	var id int64

	if err := conn.QueryRowContext(ctx, query, source.Name, source.FeedURL, source.CreatedAt, source.Priority).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *SourcePostgresStorage) Delete(ctx context.Context, id int64) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	source, err := s.SourceByID(ctx, id)
	if err != nil {
		return err
	}

	if source == nil {
		return errSourceNotFound
	}

	if _, err := conn.ExecContext(ctx, "DELETE FROM sources WHERE id = $1", id); err != nil {
		return err
	}

	return nil
}
