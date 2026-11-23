CREATE TABLE notes (
    id           BIGINT      GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id      UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title        TEXT        NOT NULL,
    body         TEXT        NOT NULL DEFAULT '',
    subject_name TEXT        NOT NULL DEFAULT '',
    lesson_date  DATE,
    lesson_start TIME,
    tsv          TSVECTOR GENERATED ALWAYS AS (
                     to_tsvector('russian', coalesce(title, '') || ' ' || coalesce(body, ''))
                 ) STORED,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX notes_user_updated_idx ON notes (user_id, updated_at DESC);
CREATE INDEX notes_user_subject_idx ON notes (user_id, subject_name);
CREATE INDEX notes_tsv_idx ON notes USING GIN (tsv);
