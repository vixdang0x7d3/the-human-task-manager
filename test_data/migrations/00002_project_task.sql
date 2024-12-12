-- +goose Up
-- +goose StatementBegin
CREATE TYPE task_state AS ENUM (
    'started',
    'waiting',
    'completed',
    'deleted'
);

CREATE TYPE task_priority AS ENUM (
    'H',
    'M',
    'L',
    'none'
);

CREATE TABLE projects (
    "id"        uuid PRIMARY KEY,
    "user_id"   uuid NOT NULL,
    "title"     varchar(64) NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE tasks (
    "id"           uuid PRIMARY KEY,
    "user_id"      uuid  NOT NULL,
    "project_id"   uuid,
    "completed_by" uuid,
    "description"  text NOT NULL ,
    "priority"     task_priority NOT NULL,
    "state"       task_state NOT NULL,
    "deadline"     timestamptz NOT NULL,
    "schedule"     timestamptz NOT NULL,
    "wait"         timestamptz NOT NULL,
    "create"       timestamptz NOT NULL,
    "end"          timestamptz NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (completed_by) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE task_status CASCADE;
DROP TYPE task_priority CASCADE;
DROP TABLE projects CASCADE;
DROP TABLE tasks CASCADE;
-- +goose StatementEnd
