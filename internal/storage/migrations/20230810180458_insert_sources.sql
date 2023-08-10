-- +goose Up
-- +goose StatementBegin
INSERT INTO sources (name, feed_url, created_at, priority) VALUES ('test source 1', 'https://test/source', '2023-08-01 10:10:10', 0);
INSERT INTO sources (name, feed_url, created_at, priority) VALUES ('test source 2', 'https://test/source/2', '2023-08-01 10:10:11', 0);
INSERT INTO sources (name, feed_url, created_at, priority) VALUES ('test source 3', 'https://test/source/3', '2023-08-01 10:10:12', 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM sources WHERE name = 'test source 1' AND feed_url = 'https://test/source';
DELETE FROM sources WHERE name = 'test source 2' AND feed_url = 'https://test/source/2';
DELETE FROM sources WHERE name = 'test source 3' AND feed_url = 'https://test/source/3';
-- +goose StatementEnd