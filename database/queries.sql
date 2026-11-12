-- database/queries.sql — контрольные запросы для marketplace
-- Каждый запрос проверяет одно требование из ТЗ (лаб. 3)

USE marketplace;
GO

-- 1. JOIN трёх таблиц: заказы с пользователями и подзаказами
-- Проверяет: ФТ-06 (Оформление заказа)
SELECT 
    o.order_number AS [Номер заказа],
    u.name AS [Покупатель],
    o.total_amount AS [Сумма],
    o.status AS [Статус],
    s.name AS [Магазин],
    so.subtotal_amount AS [Сумма подзаказа]
FROM orders o
JOIN users u ON u.id = o.user_id
JOIN sub_orders so ON so.order_id = o.id
JOIN shops s ON s.id = so.shop_id
ORDER BY o.created_at DESC;
GO

-- 2. LIKE — поиск товаров по названию (кириллица)
-- Проверяет: ФТ-04 (Поиск товаров)
SELECT 
    id, name AS [Название товара], price AS [Цена], status AS [Статус]
FROM products
WHERE name LIKE N'%iPhone%' OR name LIKE N'%Ноутбук%';
GO

-- 3. GROUP BY + COUNT — статистика товаров по категориям
-- Проверяет: ФТ-10 (Управление справочниками)
SELECT 
    c.name AS [Категория],
    COUNT(p.id) AS [Количество товаров],
    AVG(p.rating) AS [Средний рейтинг],
    SUM(p.sales_count) AS [Всего продаж]
FROM categories c
LEFT JOIN products p ON p.category_id = c.id AND p.status = 'approved'
GROUP BY c.name
ORDER BY [Количество товаров] DESC;
GO

-- 4. LEFT JOIN — магазины с количеством товаров
-- Проверяет: ФТ-02 (Создание карточки товара)
SELECT 
    s.name AS [Магазин],
    s.slug AS [URL],
    COUNT(p.id) AS [Всего товаров],
    SUM(CASE WHEN p.status = 'approved' THEN 1 ELSE 0 END) AS [Одобренных],
    SUM(CASE WHEN p.status = 'pending' THEN 1 ELSE 0 END) AS [На модерации]
FROM shops s
LEFT JOIN products p ON p.shop_id = s.id
GROUP BY s.name, s.slug
ORDER BY [Всего товаров] DESC;
GO

-- 5. Фильтр по статусу — товары на модерации
-- Проверяет: ФТ-03 (Модерация товаров)
SELECT 
    p.id, p.name AS [Название], s.name AS [Магазин], 
    c.name AS [Категория], p.price AS [Цена], p.created_at AS [Дата]
FROM products p
JOIN shops s ON s.id = p.shop_id
JOIN categories c ON c.id = p.category_id
WHERE p.status = 'pending'
ORDER BY p.created_at ASC;
GO

-- 6. Пагинация (OFFSET-FETCH) — вторая страница товаров
-- Проверяет: ФТ-04 (Поиск с пагинацией)
SELECT 
    id, name AS [Название], price AS [Цена], 
    rating AS [Рейтинг], sales_count AS [Продажи]
FROM products
WHERE status = 'approved'
ORDER BY sales_count DESC, rating DESC
OFFSET 5 ROWS
FETCH NEXT 5 ROWS ONLY;
GO

PRINT '✅ Все запросы выполнены успешно!';
GO