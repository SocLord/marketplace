-- ==============================================================
-- НАЧАЛЬНЫЕ ДАННЫЕ ДЛЯ МАРКЕТПЛЕЙСА (MARKETPLACE SEED DATA)
-- Microsoft SQL Server / SSMS
-- ==============================================================

USE MarketplaceDB;
GO

SET NOCOUNT ON;
GO

PRINT N'Инициализация демонстрационных данных...';
GO

-- 1. Пользователи (Users)
-- Пароль для всех тестовых аккаунтов: 123456
-- Хэш bcrypt: $2a$10$7v.G97.V2rP524r8L19e4.kL94/R08nJqB0uF9k0g4k5yNfV.9z1a
-- Роли: buyer, seller, moderator
IF NOT EXISTS (SELECT 1 FROM dbo.users WHERE email = 'buyer@market.tj')
BEGIN
    INSERT INTO dbo.users (name, email, password_hash, role, phone, avatar_url)
    VALUES 
    (N'Дилшод Раҳимов', 'buyer@market.tj', '$2a$10$B0u6R4G9bI3cR31.8c7bqeC038a39nFsq81YvG7M8z5.6N8jB5h9G', 'buyer', '+992 90 111 2233', 'https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?w=150'),
    (N'Парвиз Алиев (TechStore)', 'seller1@market.tj', '$2a$10$B0u6R4G9bI3cR31.8c7bqeC038a39nFsq81YvG7M8z5.6N8jB5h9G', 'seller', '+992 91 222 3344', 'https://images.unsplash.com/photo-1570295999919-56ceb5ecca61?w=150'),
    (N'Нигина Саидова (ModnoStyle)', 'seller2@market.tj', '$2a$10$B0u6R4G9bI3cR31.8c7bqeC038a39nFsq81YvG7M8z5.6N8jB5h9G', 'seller', '+992 92 333 4455', 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=150'),
    (N'Фаридун Умаров (Модератор)', 'moderator@market.tj', '$2a$10$B0u6R4G9bI3cR31.8c7bqeC038a39nFsq81YvG7M8z5.6N8jB5h9G', 'moderator', '+992 93 444 5566', 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150');
END
GO

PRINT N'Пользователи готовы.';
GO
