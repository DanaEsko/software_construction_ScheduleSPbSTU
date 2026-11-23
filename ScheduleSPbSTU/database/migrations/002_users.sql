-- Пользователи приложения.
CREATE TABLE users (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    email         TEXT        NOT NULL UNIQUE,   -- UNIQUE: второй пользователь с тем же e-mail получит ошибку (ответ 409)
    password_hash TEXT        NOT NULL,          -- хеш (bcrypt), пароль в открытом виде не хранится
    group_id      INTEGER     REFERENCES groups(id), -- может быть пустым: группу выбирают после регистрации
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Ускоряет выборки пользователей по группе.
CREATE INDEX users_group_idx ON users (group_id);
