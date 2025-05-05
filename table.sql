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