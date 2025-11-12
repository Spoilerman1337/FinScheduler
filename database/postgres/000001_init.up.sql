CREATE TABLE items (
   id UUID PRIMARY KEY,
   name TEXT NOT NULL UNIQUE,
   price NUMERIC(16, 2) NOT NULL DEFAULT 0 CHECK (price >= 0),
   description TEXT NULL,
   is_active BOOLEAN NOT NULL DEFAULT FALSE,
   created_at TIMESTAMP NOT NULL DEFAULT now(),
   updated_at TIMESTAMP NULL
);