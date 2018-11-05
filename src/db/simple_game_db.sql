DROP TABLE users CASCADE;

CREATE TABLE IF NOT EXISTS users(
  email CITEXT PRIMARY KEY,        -- Обязательное поле
  name CITEXT DEFAULT '',
  last_name CITEXT DEFAULT '',
  age SMALLINT NOT NULL,              -- Обязательное поле
  nick CITEXT NOT NULL UNIQUE,     -- Обязательное поле
  password TEXT NOT NULL,          -- Обязательное поле
  score SMALLINT NOT NULL DEFAULT 0
)