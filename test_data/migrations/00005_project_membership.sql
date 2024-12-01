-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_memberships (
    user_id uuid,
    project_id uuid,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, project_id)
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_membership;
-- +goose StatementEnd 
