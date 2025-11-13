-- Insert categories.
INSERT INTO categories (code, name) VALUES
    ('clothing', 'Clothing'),
    ('shoes', 'Shoes'),
    ('accessories', 'Accessories'),
    ('boots', 'Boots');

-- Assign products to categories
UPDATE products SET category_id = (SELECT id FROM categories WHERE code = 'clothing')
WHERE code IN ('PROD001', 'PROD004', 'PROD007');

UPDATE products SET category_id = (SELECT id FROM categories WHERE code = 'shoes')
WHERE code IN ('PROD002', 'PROD006');

UPDATE products SET category_id = (SELECT id FROM categories WHERE code = 'accessories')
WHERE code IN ('PROD003', 'PROD005', 'PROD008');

-- Assign PROD009 to boots category for discount testing (30% off boots)
-- This product has a variant with SKU 000003 (15% off), so it will have both discounts
UPDATE products SET category_id = (SELECT id FROM categories WHERE code = 'boots')
WHERE code = 'PROD009';
