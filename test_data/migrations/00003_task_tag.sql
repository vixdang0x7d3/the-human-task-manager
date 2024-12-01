-- +goose Up
-- +goose StatementBegin
ALTER TABLE  tasks
ADD tags text[];
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tasks
DROP COLUMN tags;
-- +goose StatementEnd
