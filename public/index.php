<?php
define('BASE_PATH', dirname(__DIR__));

// 1. Подключаем ядро
require BASE_PATH . '/src/Core/Database.php';
require BASE_PATH . '/src/Core/Validator.php';

// 2. Подключаем репозитории (нужно минимум 3 по методичке!)
require BASE_PATH . '/src/Repositories/ProductRepository.php';
require BASE_PATH . '/src/Repositories/CategoryRepository.php';
require BASE_PATH . '/src/Repositories/ShopRepository.php';

// 3. Подключаем контроллеры
require BASE_PATH . '/src/Controllers/ProductController.php';
require BASE_PATH . '/src/Controllers/ReportController.php';

// 4. Загрузка конфигурации
$configFile = BASE_PATH . '/config/config.php';
if (!is_file($configFile)) {
    $configFile = BASE_PATH . '/config/config.example.php';
}
$config = require $configFile;

ini_set('display_errors', $config['debug'] ? '1' : '0');
error_reporting(E_ALL);

Database::configure($config['db']);

// 5. Обработка глобальных ошибок (страница 500)
set_exception_handler(function(Throwable $e) use ($config): void {
    error_log((string)$e);
    http_response_code(500);
    echo $config['debug'] ? $e->getMessage() : 'Баъдтар кӯшиш кунед ё ба администратор муроҷиат намоед.';
});

// 6. Роутинг (Front Controller)
$uri = parse_url($_SERVER['REQUEST_URI'], PHP_URL_PATH);
$method = $_SERVER['REQUEST_METHOD'];

$productController = new ProductController();
$reportController = new ReportController();

if ($uri === '/products' && $method === 'GET') {
    $productController->index();
} elseif ($uri === '/products/create' && $method === 'GET') {
    $productController->create();
} elseif ($uri === '/products/create' && $method === 'POST') {
    $productController->store();
} elseif (preg_match('#^/products/(\d+)/delete$#', $uri, $matches) && $method === 'POST') {
    $productController->delete((int)$matches[1]);
} elseif ($uri === '/reports/categories' && $method === 'GET') {
    $reportController->categories();
} else {
    http_response_code(404);
    echo '<h1>404</h1><p>Саҳифа ёфт нашуд.</p><a href="/products">Ба каталог баргаштан</a>';
}