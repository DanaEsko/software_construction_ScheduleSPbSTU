CREATE TABLE schedule_cache (
    group_id   INTEGER     NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    week_start DATE        NOT NULL,
    payload    JSONB       NOT NULL,
    fetched_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (group_id, week_start)
);
