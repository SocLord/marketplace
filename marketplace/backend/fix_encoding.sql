UPDATE products 
SET name = N'Ноутбук Asus ZenBook 14 OLED Ultra 7',
    description = N'Ноутбуки ултрабуки борик бо экрани 2.8K OLED, протсессори пуриқтидори Intel Core Ultra 7, 32GB RAM ва 1TB SSD. Идеалӣ барои кор, барномасозӣ ва дизайн.'
WHERE id = 1;

IF NOT EXISTS (SELECT 1 FROM products WHERE name LIKE N'%Trek Marlin 7%')
BEGIN
    INSERT INTO products (
        shop_id, category_id, name, slug, description, price, old_price, stock_quantity, 
        image_url, additional_images, status, moderation_comment, rating, reviews_count, sales_count, created_at, updated_at
    ) VALUES (
        1, 5, 
        N'Велосипеди кӯҳии касбии Trek Marlin 7', 
        N'trek-marlin-7-bike', 
        N'Велосипеди пуриқтидор бо тормозҳои гидравликии дискӣ, амортизатори танзимшаванда ва чархҳои 29 дюйма.', 
        5800.00, 6500.00, 7, 
        N'https://images.unsplash.com/photo-1485965120184-e220f721d03e?w=600', 
        N'', 
        N'pending', 
        N'Интизори санҷиши модератор', 
        5.0, 6, 12, GETDATE(), GETDATE()
    );
END
ELSE
BEGIN
    UPDATE products 
    SET name = N'Велосипеди кӯҳии касбии Trek Marlin 7',
        description = N'Велосипеди пуриқтидор бо тормозҳои гидравликии дискӣ, амортизатори танзимшаванда ва чархҳои 29 дюйма.',
        status = N'pending'
    WHERE name LIKE N'%Trek Marlin 7%';
END
