CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS "users" (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fullname TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT ''
);
