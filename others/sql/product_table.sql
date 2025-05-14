CREATE TABLE product_category (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
)

CREATE TABLE product (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price NUMERIC,
    stock INTEGER,
    category_id INTEGER not NULL,
    CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES product_category(id) ON DELETE CASCADE
)

-- Tambahkan extension pg_trgm untuk fuzzy matching
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Tambahkan index untuk mempercepat pencarian trigram
CREATE INDEX product_name_trgm_idx ON product USING GIN (name gin_trgm_ops);