CREATE EXTENSION IF NOT EXISTS citext;
CREATE TABLE IF NOT EXISTS users(
  email CITEXT PRIMARY KEY,        -- Обязательное поле
  nick CITEXT NOT NULL UNIQUE,     -- Обязательное поле
  password TEXT NOT NULL,          -- Обязательное поле
  score SMALLINT NOT NULL DEFAULT 0
)
