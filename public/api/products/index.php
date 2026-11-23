<?php
declare(strict_types=1);
header('Content-Type: application/json; charset=utf-8');

$products = [
    [
        'id' => 1,
        'name' => 'Ноутбук Asus ZenBook 14 OLED',
        'category' => 'Электроника',
        'price' => 14500,
        'stock' => 10,
        'status' => 'active',
        'image' => 'https://placehold.co/60x60/0D6EFD/FFFFFF?text=Asus'
    ],
    [
        'id' => 2,
        'name' => 'Смартфон Samsung Galaxy S24',
        'category' => 'Электроника',
        'price' => 8900,
        'stock' => 5,
        'status' => 'active',
        'image' => 'https://placehold.co/60x60/198754/FFFFFF?text=Samsung'
    ],
    [
        'id' => 3,
        'name' => 'Наушники Sony WH-1000XM5',
        'category' => 'Электроника',
        'price' => 2500,
        'stock' => 0,
        'status' => 'pending',
        'image' => 'https://placehold.co/60x60/FFC107/212529?text=Sony'
    ],
    [
        'id' => 4,
        'name' => 'Клавиатура Logitech MX Keys',
        'category' => 'Электроника',
        'price' => 1200,
        'stock' => 20,
        'status' => 'active',
        'image' => 'https://placehold.co/60x60/6C757D/FFFFFF?text=Logitech'
    ],
    [
        'id' => 5,
        'name' => 'Монитор LG 27 UltraFine',
        'category' => 'Электроника',
        'price' => 3400,
        'stock' => 0,
        'status' => 'out',
        'image' => 'https://placehold.co/60x60/DC3545/FFFFFF?text=LG'
    ],
    [
        'id' => 6,
        'name' => 'Смарт-часы Xiaomi Band 8',
        'category' => 'Электроника',
        'price' => 650,
        'stock' => 15,
        'status' => 'active',
        'image' => 'https://placehold.co/60x60/0D6EFD/FFFFFF?text=Xiaomi'
    ],
    [
        'id' => 7,
        'name' => '<img src=x onerror="alert(\'XSS\')"> Тестовый товар',
        'category' => 'Тест',
        'price' => 100,
        'stock' => 1,
        'status' => 'active',
        'image' => 'https://placehold.co/60x60/DC3545/FFFFFF?text=XSS'
    ]
];

$q = mb_strtolower(trim((string)($_GET['q'] ?? '')));
if ($q !== '') {
    $products = array_values(array_filter(
        $products,
        fn(array $p): bool => str_contains(mb_strtolower($p['name']), $q)
    ));
}

echo json_encode($products, JSON_UNESCAPED_UNICODE);