ALTER TABLE products
ADD COLUMN category_id INTEGER NULL REFERENCES categories(id) ON DELETE SET NULL;

CREATE INDEX idx_products_category_id ON products(category_id);
