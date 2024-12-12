-- +goose Up
-- +goose StatementBegin
CREATE TYPE membership_role AS ENUM (
    'owner',
    'invited',
    'requested',
    'member'
);

CREATE TABLE project_memberships (
    user_id uuid,
    project_id uuid,
    role membership_role NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, project_id)
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE membership_role CASCADE;
DROP TABLE project_memberships CASCADE;
-- +goose StatementEnd 
