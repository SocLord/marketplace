# UI-kit: Маркетплейс marketplace

**Figma файл:** [[ВСТАВЬ ССЫЛКУ ИЗ FIGMA](https://www.figma.com/design/RkrOs8xLNOm7H66iIU2lpV/marketplace?node-id=1-5&t=yQkSUb4p6uqN0pJ0-1)]  
**Прототип:** [[ВСТАВЬ ССЫЛКУ НА ПРОТОТИП](https://www.figma.com/design/RkrOs8xLNOm7H66iIU2lpV/marketplace?node-id=0-1&t=yQkSUb4p6uqN0pJ0-1)]

## 1. Сетка (Layout Grid)
- **Desktop (1440px):** 12 колонок, ширина колонки 88px, gutter 24px. Ширина контента 1320px (как `.container` в Bootstrap 5.3).
- **Mobile (390px):** 4 колонки, margin 16px, gutter 16px.

## 2. Цвета (Color Styles)
| Стиль | HEX | Переменная Bootstrap | Где используется |
|---|---|---|---|
| `color/primary` | `#0D6EFD` | `--bs-primary` | Основные кнопки, ссылки, цены |
| `color/success` | `#198754` | `--bs-success` | Статус "В наличии", "Одобрено" |
| `color/warning` | `#FFC107` | `--bs-warning` | Статус "На модерации" (текст тёмный) |
| `color/danger` | `#DC3545` | `--bs-danger` | Ошибки, кнопка "Удалить" |
| `color/text` | `#212529` | `--bs-body-color` | Основной текст |
| `color/bg` | `#F8F9FA` | `--bs-light` | Фон страницы |

## 3. Типографика
Шрифт: **Noto Sans** (полная кириллица, включая ғ, ӣ, ӯ, ҳ, қ).

| Стиль | Размер/межстрочный | Начертание |
|---|---|---|
| `text/h1` | 32/40 | SemiBold |
| `text/h2` | 24/32 | SemiBold |
| `text/body` | 16/24 | Regular |
| `text/small` | 14/20 | Regular |

## 4. Компоненты
- **Button:** Варианты `Primary`, `Secondary`, `Danger`. Состояния `Default`, `Hover`. Высота 44px.
- **Input:** Состояния `Default`, `Focus`, `Error`. Высота 44px, с подписью и сообщением об ошибке.
- **ProductCard:** Карточка товара на Auto Layout (аналог CSS Flexbox). При длинном названии карточка растягивается автоматически.

## 5. Экраны (Mockups)
![Каталог товаров](screens/products.png)
![Карточка товара](screens/product-detail.png)