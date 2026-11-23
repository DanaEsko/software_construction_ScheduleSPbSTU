-- Справочник учебных групп СПбПУ.
-- id совпадает с идентификатором группы в ruz.spbstu.ru (например, 45608),
-- поэтому собственная нумерация и таблица соответствия не нужны.
CREATE TABLE groups (
    id           INTEGER     PRIMARY KEY,
    name         TEXT        NOT NULL,
    faculty_id   INTEGER     NOT NULL,
    faculty_name TEXT        NOT NULL,
    synced_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Ускоряет запрос «группы института» (GET /faculties/{id}/groups).
CREATE INDEX groups_faculty_idx ON groups (faculty_id);
