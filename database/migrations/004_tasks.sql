CREATE TABLE tasks (
    id           BIGINT      GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id      UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title        TEXT        NOT NULL,
    description  TEXT        NOT NULL DEFAULT '',
    deadline     TIMESTAMPTZ,
    status       TEXT        NOT NULL DEFAULT 'planned'
                             CHECK (status IN ('planned', 'in_progress', 'done')),
    subject_name TEXT        NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX tasks_user_status_deadline_idx ON tasks (user_id, status, deadline);
