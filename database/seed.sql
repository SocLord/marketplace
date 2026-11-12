-- database/seed.sql — тестовые данные для marketplace
-- Пароль для всех: 123456 (хеш bcrypt)

USE marketplace;
GO

-- Пользователи (все роли)
INSERT INTO users (name, email, password_hash, role, phone, avatar_url) VALUES
(N'Дилшод Раҳимов', N'buyer@market.tj', 
 N'$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 
 N'buyer', N'+992 90 111 2233', 
 N'https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?w=150'),
(N'Парвиз Алиев', N'seller1@market.tj', 
 N'$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 
 N'seller', N'+992 91 222 3344', 
 N'https://images.unsplash.com/photo-1570295999919-56ceb5ecca61?w=150'),
(N'Нигина Саидова', N'seller2@market.tj', 
 N'$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 
 N'seller', N'+992 92 333 4455', 
 N'https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=150'),
(N'Фаридун Умаров', N'moderator@market.tj', 
 N'$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 
 N'moderator', N'+992 93 444 5566', 
 N'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150'),
(N'Администратор Систем', N'admin@market.tj', 
 N'$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 
 N'admin', N'+992 94 555 6677', 
 N'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150');
GO

-- Магазины
INSERT INTO shops (user_id, name, slug, description, logo_url, banner_url, phone, address, rating, is_verified) VALUES
(2, N'TechStore', N'techstore', 
 N'Мағозаи расмии техникаи рақамӣ.',
 N'https://images.unsplash.com/photo-1550009158-9ebf69173e03?w=200',
 N'https://images.unsplash.com/photo-1518770660439-4636190af475?w=1200',
 N'+992 91 222 3344', N'ш. Душанбе, хиёбони Рӯдакӣ, 45', 4.9, 1),
(3, N'ModnoStyle', N'modnostyle', 
 N'Либосҳои замонавӣ.',
 N'https://images.unsplash.com/photo-1441986300917-64674bd600d8?w=200',
 N'https://images.unsplash.com/photo-1441984904996-e0b6ba687e04?w=1200',
 N'+992 92 333 4455', N'ш. Хуҷанд, кӯчаи И. Сомонӣ, 12', 4.8, 1);
GO

-- Категории (6 штук)
INSERT INTO categories (name, slug, icon, image_url, description, sort_order) VALUES
(N'Электроника', N'electronics', N'Laptop', 
 N'https://images.unsplash.com/photo-1498049794561-7780e7231661?w=400',
 N'Смартфонҳо, ноутбукҳо, гӯшмонакҳо', 1),
(N'Либос ва пойафзол', N'clothing', N'Shirt', 
 N'https://images.unsplash.com/photo-1489987707025-afc232f7ea0f?w=400',
 N'Либоси мардона, занона', 2),
(N'Хона ва ошхона', N'home-kitchen', N'Home', 
 N'https://images.unsplash.com/photo-1513694203232-719a280e022f?w=400',
 N'Зарфҳо, нассоҷӣ, рӯзгор', 3),
(N'Зебоӣ ва саломатӣ', N'beauty', N'Sparkles', 
 N'https://images.unsplash.com/photo-1522335789203-aabd1fc54bc9?w=400',
 N'Пардоз, атр, маҳсулоти нигоҳубин', 4),
(N'Варзиш ва сайёҳӣ', N'sports', N'Activity', 
 N'https://images.unsplash.com/photo-1517838277536-f5f99be501cd?w=400',
 N'Анҷомҳои варзишӣ, хаймаҳо', 5),
(N'Китобҳо', N'books', N'BookOpen', 
 N'https://images.unsplash.com/photo-1495446815901-a7297e633e8d?w=400',
 N'Адабиёти бадеӣ, илмӣ', 6);
GO

-- Товары (12 штук — больше 10, как требует методичка)
INSERT INTO products (shop_id, category_id, name, slug, description, price, old_price, stock_quantity, image_url, status, rating, reviews_count, sales_count) VALUES
(1, 1, N'Ноутбук Asus ZenBook 14 OLED', N'asus-zenbook-14', 
 N'Ноутбуки ултрабуки борик бо экрани 2.8K OLED.', 
 14500.00, 16200.00, 15, N'https://images.unsplash.com/photo-1541807084-5c52b6b3adef?w=600', 
 N'approved', 4.9, 8, 24),
(1, 1, N'Смартфон Apple iPhone 15 Pro 256GB', N'iphone-15-pro', 
 N'Смартфони флагманӣ бо корпуси титанӣ.', 
 13800.00, 14900.00, 12, N'https://images.unsplash.com/photo-1510557880182-3d4d3cba35a5?w=600', 
 N'approved', 5.0, 14, 35),
(1, 1, N'Гӯшмонакҳои бесими Sony WH-1000XM5', N'sony-wh-1000xm5', 
 N'Гӯшмонакҳои беҳтарин бо Noise Cancelling.', 
 3450.00, 3900.00, 28, N'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=600', 
 N'approved', 4.8, 9, 19),
(1, 1, N'Смарт-соат Samsung Galaxy Watch 6', N'samsung-watch-6', 
 N'Соати ҳушманд бо мониторинги хоб ва ЭКГ.', 
 2890.00, 3300.00, 20, N'https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=600', 
 N'approved', 4.7, 6, 15),
(1, 1, N'Планшет Apple iPad Air M2', N'ipad-air-m2', 
 N'Планшети навтарин бо чипи Apple M2.', 
 7900.00, 8500.00, 10, N'https://images.unsplash.com/photo-1544244015-0df4b3ffc6b0?w=600', 
 N'pending', 0.0, 0, 0),
(2, 2, N'Курткаи чармии мардона Classic Biker', N'leather-jacket-biker', 
 N'Куртка аз чарми табиӣ.', 
 1250.00, 1600.00, 22, N'https://images.unsplash.com/photo-1520975916090-3105956dac38?w=600', 
 N'approved', 4.9, 11, 42),
(2, 2, N'Кросссовкаҳои варзишии Nike Air Max', N'nike-air-max', 
 N'Пойафзоли варзишии бароҳат.', 
 980.00, 1250.00, 30, N'https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=600', 
 N'approved', 5.0, 15, 58),
(2, 2, N'Либоси шевои занонаи шомгоҳӣ', N'emerald-silk-dress', 
 N'Либоси абрешими нафис.', 
 850.00, 1100.00, 18, N'https://images.unsplash.com/photo-1515372039744-b8f02a3ae446?w=600', 
 N'approved', 4.8, 7, 26),
(2, 2, N'Костюми классикии мардона Slim Fit', N'mens-suit-slim', 
 N'Костюми босифати 2-қисма.', 
 1950.00, 2400.00, 14, N'https://images.unsplash.com/photo-1594938298603-c8148c4dae35?w=600', 
 N'approved', 4.9, 5, 18),
(2, 2, N'Пальтои гарми зимистонаи кашемирӣ', N'cashmere-winter-coat', 
 N'Пальтои дароз бо матои кашемири 100%.', 
 1600.00, 2100.00, 8, N'https://images.unsplash.com/photo-1539533018447-63fcce667823?w=600', 
 N'pending', 0.0, 0, 0),
(2, 3, N'Маҷмӯи зарфҳои қаҳвапазии керамикӣ', N'ceramic-coffee-set', 
 N'Маҷмӯи қаҳванушӣ аз сафоли табиӣ.', 
 420.00, 550.00, 25, N'https://images.unsplash.com/photo-1514432324607-a09d9b4aefdd?w=600', 
 N'approved', 4.7, 4, 16),
(1, 5, N'Велосипеди кӯҳии касбии Trek Marlin 7', N'trek-marlin-7', 
 N'Велосипеди пуриқтидор бо тормозҳои гидравликии.', 
 5800.00, 6500.00, 7, N'https://images.unsplash.com/photo-1485965120184-e220f721d03e?w=600', 
 N'approved', 5.0, 6, 12);
GO

-- Заказы
INSERT INTO orders (order_number, user_id, total_amount, status, recipient_name, recipient_phone, shipping_address, payment_method, notes) VALUES
(N'ORD-2026-9001', 1, 4430.00, N'processing', N'Дилшод Раҳимов', N'+992 90 111 2233', 
 N'ш. Душанбе, н. Сино, кӯчаи Борбад 48', N'cash_on_delivery', N'Илтимос пеш аз омадан занг занед'),
(N'ORD-2026-9002', 1, 15750.00, N'completed', N'Дилшод Раҳимов', N'+992 90 111 2233', 
 N'ш. Душанбе, н. И. Сомонӣ', N'card', NULL),
(N'ORD-2026-9003', 1, 16690.00, N'completed', N'Дилшод Раҳимов', N'+992 90 111 2233', 
 N'ш. Душанбе', N'card', NULL);
GO

-- Подзаказы
INSERT INTO sub_orders (order_id, shop_id, subtotal_amount, status, shipping_tracking) VALUES
(1, 1, 3450.00, N'processing', N'TRACK-TJ-TECH-101'),
(1, 2, 980.00, N'shipped', N'TRACK-TJ-MOD-202'),
(2, 1, 14500.00, N'delivered', N'TRACK-OCT-01'),
(2, 2, 1250.00, N'delivered', N'TRACK-OCT-02'),
(3, 1, 13800.00, N'delivered', N'TRACK-SEP-01'),
(3, 2, 2890.00, N'delivered', N'TRACK-SEP-02');
GO

-- Позиции заказа
INSERT INTO order_items (order_id, sub_order_id, shop_id, product_id, product_name, product_price, quantity, total_price, image_url) VALUES
(1, 1, 1, 3, N'Гӯшмонакҳои бесими Sony WH-1000XM5', 3450.00, 1, 3450.00, 
 N'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=600'),
(1, 2, 2, 7, N'Кросссовкаҳои варзишии Nike Air Max', 980.00, 1, 980.00, 
 N'https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=600'),
(2, 3, 1, 1, N'Ноутбук Asus ZenBook 14 OLED', 14500.00, 1, 14500.00, 
 N'https://images.unsplash.com/photo-1541807084-5c52b6b3adef?w=600'),
(2, 4, 2, 6, N'Курткаи чармии мардона Classic Biker', 1250.00, 1, 1250.00, 
 N'https://images.unsplash.com/photo-1520975916090-3105956dac38?w=600'),
(3, 5, 1, 2, N'Смартфон Apple iPhone 15 Pro 256GB', 13800.00, 1, 13800.00, 
 N'https://images.unsplash.com/photo-1510557880182-3d4d3cba35a5?w=600'),
(3, 6, 2, 9, N'Костюми классикии мардона Slim Fit', 1950.00, 1, 1950.00, 
 N'https://images.unsplash.com/photo-1594938298603-c8148c4dae35?w=600');
GO

-- Отзывы
INSERT INTO reviews (product_id, user_id, order_id, rating, comment) VALUES
(1, 1, 2, 5, N'Ноутбуки бисёр олиҷаноб! Экрани OLED чашмро хаста намекунад.'),
(2, 1, 3, 5, N'Сифати аслӣ! Телефон дар ҳақиқат аслии Apple аст.'),
(6, 1, 2, 5, N'Куртка хеле хуб дӯхта шудааст. Андозааш дақиқ рост омад.');
GO

PRINT '✅ Тестовые данные успешно загружены!';
GO