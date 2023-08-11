-- +goose Up
-- +goose StatementBegin
INSERT INTO sources (name, feed_url, priority) VALUES ('dev_to', 'https://dev.to/feed/tag/golang', 0);
INSERT INTO sources (name, feed_url, priority) VALUES ('hashnode', 'https://hashnode.com/n/golang/rss', 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM sources WHERE feed_url = 'https://dev.to/feed/tag/golang';
DELETE FROM sources WHERE feed_url = 'https://hashnode.com/n/golang/rss';
-- +goose StatementEnd