package postgres

import (
	"context"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/zumosik/telegram-news-go/internal/model"
)

func CreateTestSourcePostgresStorage(t *testing.T) *SourcePostgresStorage {
	godotenv.Load()
	db, err := sqlx.Connect("postgres", os.Getenv("POSTGRES_TEST_STORAGE_URL")) // poot .env in this folder if doesnt work
	if err != nil {
		t.Error(err)
		t.FailNow()
		return nil
	}
	return NewSourceStorage(db)
}

func TestSourcePostgresStorage_Sources(t *testing.T) {
	storage := CreateTestSourcePostgresStorage(t)
	ctx := context.Background()

	tests := []struct {
		name    string
		ctx     context.Context
		wantErr bool
	}{
		{
			name:    "ok",
			ctx:     ctx,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := storage.Sources(tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("SourcePostgresStorage.Sources() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSourcePostgresStorage_SourceByID(t *testing.T) {
	storage := CreateTestSourcePostgresStorage(t)
	ctx := context.Background()

	// creating user
	id, err := storage.Add(ctx, model.Source{
		Name:     "test source for founding",
		FeedURL:  "https://exmaple.org",
		Priority: 0,
	})

	if err != nil {
		t.FailNow()
	}

	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "error",
			args: args{
				ctx: ctx,
				id:  0,
			},
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				ctx: ctx,
				id:  id,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := storage.SourceByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("SourcePostgresStorage.SourceByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSourcePostgresStorage_Add(t *testing.T) {
	storage := CreateTestSourcePostgresStorage(t)
	ctx := context.Background()
	type args struct {
		ctx    context.Context
		source model.Source
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				ctx: ctx,
				source: model.Source{
					Name:     "new test source 1",
					FeedURL:  "https://aaa",
					Priority: 0,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := storage.Add(tt.args.ctx, tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("SourcePostgresStorage.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSourcePostgresStorage_Delete(t *testing.T) {
	storage := CreateTestSourcePostgresStorage(t)
	ctx := context.Background()

	// creating user
	id, err := storage.Add(ctx, model.Source{
		Name:     "test source for deleting",
		FeedURL:  "https://exmaple.org",
		Priority: 0,
	})

	if err != nil {
		t.FailNow()
	}

	type args struct {
		ctx context.Context
		id  int64
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
				id:  id,
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				ctx: ctx,
				id:  9999,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := storage.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("SourcePostgresStorage.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
