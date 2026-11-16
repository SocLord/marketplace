-- ==============================================================
-- БАЗА ДАННЫХ ДЛЯ МАРКЕТПЛЕЙСА (MARKETPLACE)
-- Microsoft SQL Server / SSMS
-- ==============================================================

USE master;
GO

IF NOT EXISTS (SELECT name FROM sys.databases WHERE name = N'MarketplaceDB')
BEGIN
    CREATE DATABASE MarketplaceDB COLLATE Cyrillic_General_CI_AS;
END
GO

USE MarketplaceDB;
GO

-- 1. Таблица пользователей (Users)
-- Нақшҳо: buyer (харидор), seller (фурӯшанда), moderator (модератор), admin (администратор)
IF OBJECT_ID('dbo.order_items', 'U') IS NOT NULL DROP TABLE dbo.order_items;
IF OBJECT_ID('dbo.sub_orders', 'U') IS NOT NULL DROP TABLE dbo.sub_orders;
IF OBJECT_ID('dbo.orders', 'U') IS NOT NULL DROP TABLE dbo.orders;
IF OBJECT_ID('dbo.reviews', 'U') IS NOT NULL DROP TABLE dbo.reviews;
IF OBJECT_ID('dbo.products', 'U') IS NOT NULL DROP TABLE dbo.products;
IF OBJECT_ID('dbo.categories', 'U') IS NOT NULL DROP TABLE dbo.categories;
IF OBJECT_ID('dbo.shops', 'U') IS NOT NULL DROP TABLE dbo.shops;
IF OBJECT_ID('dbo.users', 'U') IS NOT NULL DROP TABLE dbo.users;
GO

CREATE TABLE dbo.users (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    name NVARCHAR(150) NOT NULL,
    email NVARCHAR(150) NOT NULL UNIQUE,
    password_hash NVARCHAR(255) NOT NULL,
    role NVARCHAR(50) NOT NULL DEFAULT 'buyer', -- 'buyer', 'seller', 'moderator', 'admin'
    phone NVARCHAR(50) NULL,
    avatar_url NVARCHAR(500) NULL,
    created_at DATETIME2 NOT NULL DEFAULT SYSUTCDATETIME(),
    updated_at DATETIME2 NOT NULL DEFAULT SYSUTCDATETIME()
);
GO

-- 2. Таблица магазинов продавцов (Shops)
-- Саҳифаи алоҳидаи мағозаи фурӯшанда
CREATE TABLE dbo.shops (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    user_id BIGINT NOT NULL,
    name NVARCHAR(200) NOT NULL,
    slug NVARCHAR(200) NOT NULL UNIQUE,
    description NVARCHAR(MAX) NULL,
    logo_url NVARCHAR(500) NULL,
    banner_url NVARCHAR(500) NULL,
    phone NVARCHAR(50) NULL,
    address NVARCHAR(255) NULL,
    rating FLOAT NOT NULL DEFAULT 5.0,
    is_verified BIT NOT NULL DEFAULT 1,
    created_at DATETIME2 NOT NULL DEFAULT SYSUTCDATETIME(),
    updated_at DATETIME2 NOT NULL DEFAULT SYSUTCDATETIME(),
    CONSTRAINT FK_Shops_Users FOREIGN KEY (user_id) REFERENCES dbo.users(id) ON DELETE CASCADE
);
GO

-- 3. Таблица категорий товаров (Categories)
-- Феҳристи молҳо бо категорияҳо
CREATE TABLE dbo.categories (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    name NVARCHAR(150) NOT NULL,
    slug NVARCHAR(150) NOT NULL UNIQUE,
    icon NVARCHAR(100) NULL,
    image_url NVARCHAR(500) NULL,
    description NVARCHAR(MAX) NULL,
    sort_order INT NOT NULL DEFAULT 0,
    created_at DATETIME2 NOT NULL DEFAULT SYSUTCDATETIME()
);
GO

-- 4. Таблица товаров (Products)
-- Статус: 'pending' (дар санҷиш), 'approved' (нашр шуд), 'rejected' (рад шуд)
CREATE TABLE dbo.products (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    shop_id BIGINT NOT NULL,
    category_id BIGINT NOT NULL,
    name NVARCHAR(300) NOT NULL,
    slug NVARCHAR(300) NOT NULL,
    description NVARCHAR(MAX) NULL,
    price DECIMAL(14,2) NOT NULL DEFAULT 0.00,
    old_price DECIMAL(14,2) NULL,
    stock_quantity INT NOT NULL DEFAULT 0,
    image_url NVARCHAR(500) NOT NULL,
    additional_images NVARCHAR(MAX) NULL, -- JSON массив картинок
    status NVARCHAR(50) NOT NULL DEFAULT 'pending', -- 'pending', 'approved', 'rejected'
    moderation_comment NVARCHAR(MAX) NULL,
    rating FLOAT NOT NULL DEFAULT 0.0,
    reviews_count INT NOT NULL DEFAULT 0,
    sales_count INT NOT NULL DEFAULT 0,
    created_at DATETIME2 NOT NULL DEFAULT SYSUTCDATETIME(),
    updated_at DATETIME2 NOT NULL DEFAULT SYSUTCDATETIME(),
    CONSTRAINT FK_Products_Shops FOREIGN KEY (shop_id) REFERENCES dbo.shops(id) ON DELETE CASCADE,
    CONSTRAINT FK_Products_Categories FOREIGN KEY (category_id) REFERENCES dbo.categories(id)
);
GO

-- 5. Таблица общих заказов (Orders)
CREATE TABLE dbo.orders (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    order_number NVARCHAR(60) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    total_amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
    status NVARCHAR(50) NOT NULL DEFAULT 'processing', -- 'processing', 'partially_shipped', 'completed', 'cancelled'
    recipient_name NVARCHAR(150) NOT NULL,
    recipient_phone NVARCHAR(50) NOT NULL,
    shipping_address NVARCHAR(500) NOT NULL,
    payment_method NVARCHAR(50) NOT NULL DEFAULT 'cash_on_delivery', -- 'card', 'cash_on_delivery', 'wallet'
    notes NVARCHAR(MAX) NULL,
    created_at DATETIME2 NOT NULL DEFAULT SYSUTCDATETIME(),
    updated_at DATETIME2 NOT NULL DEFAULT SYSUTCDATETIME(),
    CONSTRAINT FK_Orders_Users FOREIGN KEY (user_id) REFERENCES dbo.users(id)
);
GO

-- 6. Таблица подзаказов продавцов (SubOrders)
-- ТАЛАБОТИ МАХСУС: Фармоише, ки молҳои якчанд фурӯшандаро дорад, ба қисмҳо тақсим мешавад
-- ва ҳар фурӯшанда танҳо қисми худро мебинад!
CREATE TABLE dbo.sub_orders (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    order_id BIGINT NOT NULL,
    shop_id BIGINT NOT NULL,
    subtotal_amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
    status NVARCHAR(50) NOT NULL DEFAULT 'pending', -- 'pending', 'processing', 'shipped', 'delivered', 'cancelled'
    shipping_tracking NVARCHAR(100) NULL,
    created_at DATETIME2 NOT NULL DEFAULT SYSUTCDATETIME(),
    updated_at DATETIME2 NOT NULL DEFAULT SYSUTCDATETIME(),
    CONSTRAINT FK_SubOrders_Orders FOREIGN KEY (order_id) REFERENCES dbo.orders(id) ON DELETE CASCADE,
    CONSTRAINT FK_SubOrders_Shops FOREIGN KEY (shop_id) REFERENCES dbo.shops(id)
);
GO

-- 7. Позиции заказа (Order Items)
CREATE TABLE dbo.order_items (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    order_id BIGINT NOT NULL,
    sub_order_id BIGINT NOT NULL,
    shop_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    product_name NVARCHAR(300) NOT NULL,
    product_price DECIMAL(14,2) NOT NULL DEFAULT 0.00,
    quantity INT NOT NULL DEFAULT 1,
    total_price DECIMAL(14,2) NOT NULL DEFAULT 0.00,
    image_url NVARCHAR(500) NULL,
    created_at DATETIME2 NOT NULL DEFAULT SYSUTCDATETIME(),
    CONSTRAINT FK_OrderItems_Orders FOREIGN KEY (order_id) REFERENCES dbo.orders(id) ON DELETE NO ACTION,
    CONSTRAINT FK_OrderItems_SubOrders FOREIGN KEY (sub_order_id) REFERENCES dbo.sub_orders(id) ON DELETE CASCADE,
    CONSTRAINT FK_OrderItems_Shops FOREIGN KEY (shop_id) REFERENCES dbo.shops(id),
    CONSTRAINT FK_OrderItems_Products FOREIGN KEY (product_id) REFERENCES dbo.products(id)
);
GO

-- 8. Отзывы и оценки (Reviews)
-- Фикру мулоҳизаҳо ва баҳо ба мол
CREATE TABLE dbo.reviews (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    product_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    order_id BIGINT NULL,
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment NVARCHAR(MAX) NOT NULL,
    created_at DATETIME2 NOT NULL DEFAULT SYSUTCDATETIME(),
    updated_at DATETIME2 NOT NULL DEFAULT SYSUTCDATETIME(),
    CONSTRAINT FK_Reviews_Products FOREIGN KEY (product_id) REFERENCES dbo.products(id) ON DELETE CASCADE,
    CONSTRAINT FK_Reviews_Users FOREIGN KEY (user_id) REFERENCES dbo.users(id)
);
GO

-- Создание полезных индексов для быстродействия поиска, фильтрации и статистики
CREATE INDEX IX_Products_Shop ON dbo.products(shop_id);
CREATE INDEX IX_Products_Category ON dbo.products(category_id);
CREATE INDEX IX_Products_Status ON dbo.products(status);
CREATE INDEX IX_Products_Price ON dbo.products(price);
CREATE INDEX IX_SubOrders_Shop ON dbo.sub_orders(shop_id);
CREATE INDEX IX_SubOrders_Order ON dbo.sub_orders(order_id);
CREATE INDEX IX_OrderItems_SubOrder ON dbo.order_items(sub_order_id);
CREATE INDEX IX_OrderItems_Shop ON dbo.order_items(shop_id);
CREATE INDEX IX_Reviews_Product ON dbo.reviews(product_id);
GO
