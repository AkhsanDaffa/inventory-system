ALTER TABLE products
DROP CONSTRAINT fk_products_category;

ALTER TABLE products
DROP COLUMN category_id;