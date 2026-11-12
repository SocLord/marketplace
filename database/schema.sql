-- database/schema.sql — физическая модель marketplace
-- MS SQL Server 2019+ с поддержкой UTF-8
-- Перезапускаемый скрипт

IF DB_ID('marketplace') IS NOT NULL
BEGIN
    ALTER DATABASE marketplace SET SINGLE_USER WITH ROLLBACK IMMEDIATE;
    DROP DATABASE marketplace;
END
GO

CREATE DATABASE marketplace
COLLATE Latin1_General_100_CI_AS_SC_UTF8;
GO

USE marketplace;
GO

-- Таблица пользователей
CREATE TABLE users (
    id              BIGINT IDENTITY(1,1) PRIMARY KEY,
    name            NVARCHAR(150) NOT NULL,
    email           NVARCHAR(150) NOT NULL,
    password_hash   NVARCHAR(255) NOT NULL,
    role            NVARCHAR(50) NOT NULL DEFAULT 'buyer',
    phone           NVARCHAR(50) NULL,
    avatar_url      NVARCHAR(500) NULL,
    created_at      DATETIME2 NOT NULL DEFAULT GETDATE(),
    updated_at      DATETIME2 NOT NULL DEFAULT GETDATE(),
    CONSTRAINT uq_users_email UNIQUE (email)
);
GO

-- Таблица магазинов
CREATE TABLE shops (
    id              BIGINT IDENTITY(1,1) PRIMARY KEY,
    user_id         BIGINT NOT NULL,
    name            NVARCHAR(200) NOT NULL,
    slug            NVARCHAR(200) NOT NULL,
    description     NVARCHAR(MAX) NULL,
    logo_url        NVARCHAR(500) NULL,
    banner_url      NVARCHAR(500) NULL,
    phone           NVARCHAR(50) NULL,
    address         NVARCHAR(255) NULL,
    rating          DECIMAL(3,2) NOT NULL DEFAULT 5.0,
    is_verified     BIT NOT NULL DEFAULT 1,
    created_at      DATETIME2 NOT NULL DEFAULT GETDATE(),
    updated_at      DATETIME2 NOT NULL DEFAULT GETDATE(),
    CONSTRAINT uq_shops_slug UNIQUE (slug),
    CONSTRAINT fk_shops_user FOREIGN KEY (user_id) 
        REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);
GO

-- Таблица категорий
CREATE TABLE categories (
    id              BIGINT IDENTITY(1,1) PRIMARY KEY,
    name            NVARCHAR(150) NOT NULL,
    slug            NVARCHAR(150) NOT NULL,
    icon            NVARCHAR(100) NULL,
    image_url       NVARCHAR(500) NULL,
    description     NVARCHAR(MAX) NULL,
    sort_order      INT NOT NULL DEFAULT 0,
    created_at      DATETIME2 NOT NULL DEFAULT GETDATE(),
    CONSTRAINT uq_categories_slug UNIQUE (slug)
);
GO

-- Таблица товаров
CREATE TABLE products (
    id                  BIGINT IDENTITY(1,1) PRIMARY KEY,
    shop_id             BIGINT NOT NULL,
    category_id         BIGINT NOT NULL,
    name                NVARCHAR(300) NOT NULL,
    slug                NVARCHAR(300) NOT NULL,
    description         NVARCHAR(MAX) NULL,
    price               DECIMAL(14,2) NOT NULL DEFAULT 0,
    old_price           DECIMAL(14,2) NULL,
    stock_quantity      INT NOT NULL DEFAULT 0,
    image_url           NVARCHAR(500) NOT NULL,
    additional_images   NVARCHAR(MAX) NULL,
    status              NVARCHAR(50) NOT NULL DEFAULT 'pending',
    moderation_comment  NVARCHAR(MAX) NULL,
    rating              DECIMAL(3,2) NOT NULL DEFAULT 0.0,
    reviews_count       INT NOT NULL DEFAULT 0,
    sales_count         INT NOT NULL DEFAULT 0,
    created_at          DATETIME2 NOT NULL DEFAULT GETDATE(),
    updated_at          DATETIME2 NOT NULL DEFAULT GETDATE(),
    CONSTRAINT fk_products_shop FOREIGN KEY (shop_id) 
        REFERENCES shops(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_products_category FOREIGN KEY (category_id) 
        REFERENCES categories(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT chk_products_price CHECK (price >= 0),
    INDEX idx_products_name (name),
    INDEX idx_products_shop (shop_id)
);
GO

-- Таблица заказов
CREATE TABLE orders (
    id                  BIGINT IDENTITY(1,1) PRIMARY KEY,
    order_number        NVARCHAR(60) NOT NULL,
    user_id             BIGINT NOT NULL,
    total_amount        DECIMAL(14,2) NOT NULL DEFAULT 0,
    status              NVARCHAR(50) NOT NULL DEFAULT 'processing',
    recipient_name      NVARCHAR(150) NOT NULL,
    recipient_phone     NVARCHAR(50) NOT NULL,
    shipping_address    NVARCHAR(500) NOT NULL,
    payment_method      NVARCHAR(50) NOT NULL DEFAULT 'cash_on_delivery',
    notes               NVARCHAR(MAX) NULL,
    created_at          DATETIME2 NOT NULL DEFAULT GETDATE(),
    updated_at          DATETIME2 NOT NULL DEFAULT GETDATE(),
    CONSTRAINT uq_orders_number UNIQUE (order_number),
    CONSTRAINT fk_orders_user FOREIGN KEY (user_id) 
        REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    INDEX idx_orders_user (user_id)
);
GO

-- Таблица подзаказов
CREATE TABLE sub_orders (
    id                  BIGINT IDENTITY(1,1) PRIMARY KEY,
    order_id            BIGINT NOT NULL,
    shop_id             BIGINT NOT NULL,
    subtotal_amount     DECIMAL(14,2) NOT NULL DEFAULT 0,
    status              NVARCHAR(50) NOT NULL DEFAULT 'pending',
    shipping_tracking   NVARCHAR(100) NULL,
    created_at          DATETIME2 NOT NULL DEFAULT GETDATE(),
    updated_at          DATETIME2 NOT NULL DEFAULT GETDATE(),
    CONSTRAINT fk_sub_orders_order FOREIGN KEY (order_id) 
        REFERENCES orders(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_sub_orders_shop FOREIGN KEY (shop_id) 
        REFERENCES shops(id) ON DELETE RESTRICT ON UPDATE CASCADE
);
GO

-- Таблица позиций заказа
CREATE TABLE order_items (
    id              BIGINT IDENTITY(1,1) PRIMARY KEY,
    order_id        BIGINT NOT NULL,
    sub_order_id    BIGINT NOT NULL,
    shop_id         BIGINT NOT NULL,
    product_id      BIGINT NOT NULL,
    product_name    NVARCHAR(300) NOT NULL,
    product_price   DECIMAL(14,2) NOT NULL,
    quantity        INT NOT NULL DEFAULT 1,
    total_price     DECIMAL(14,2) NOT NULL,
    image_url       NVARCHAR(500) NULL,
    created_at      DATETIME2 NOT NULL DEFAULT GETDATE(),
    CONSTRAINT fk_order_items_order FOREIGN KEY (order_id) 
        REFERENCES orders(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_order_items_sub_order FOREIGN KEY (sub_order_id) 
        REFERENCES sub_orders(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_order_items_shop FOREIGN KEY (shop_id) 
        REFERENCES shops(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_order_items_product FOREIGN KEY (product_id) 
        REFERENCES products(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT chk_order_items_quantity CHECK (quantity > 0)
);
GO

-- Таблица отзывов
CREATE TABLE reviews (
    id              BIGINT IDENTITY(1,1) PRIMARY KEY,
    product_id      BIGINT NOT NULL,
    user_id         BIGINT NOT NULL,
    order_id        BIGINT NULL,
    rating          INT NOT NULL,
    comment         NVARCHAR(MAX) NOT NULL,
    created_at      DATETIME2 NOT NULL DEFAULT GETDATE(),
    updated_at      DATETIME2 NOT NULL DEFAULT GETDATE(),
    CONSTRAINT fk_reviews_product FOREIGN KEY (product_id) 
        REFERENCES products(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_reviews_user FOREIGN KEY (user_id) 
        REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_reviews_order FOREIGN KEY (order_id) 
        REFERENCES orders(id) ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT chk_reviews_rating CHECK (rating >= 1 AND rating <= 5)
);
GO

PRINT '✅ База данных marketplace успешно создана!';
GO