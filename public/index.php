<?php
declare(strict_types=1);

use App\Controllers\ProductController;
use App\Controllers\ReportController;
use App\Controllers\ShopController;
use App\Core\Router;

define('BASE_PATH', dirname(__DIR__));

// 1. Автозагрузка Composer (PSR-4)
if (is_file(BASE_PATH . '/vendor/autoload.php')) {
    require BASE_PATH . '/vendor/autoload.php';
} else {
    spl_autoload_register(function (string $class): void {
        if (str_starts_with($class, 'App\\')) {
            $file = BASE_PATH . '/src/' . str_replace('\\', '/', substr($class, 4)) . '.php';
            if (is_file($file)) {
                require $file;
            }
        }
    });
}

// 2. Конфигурация и обработка ошибок
$configFile = BASE_PATH . '/config/config.php';
if (!is_file($configFile)) {
    $configFile = BASE_PATH . '/config/config.example.php';
}
$config = require $configFile;

ini_set('display_errors', $config['debug'] ? '1' : '0');
error_reporting(E_ALL);

if (class_exists('App\Core\Database')) {
    App\Core\Database::configure($config['db']);
}

set_exception_handler(function (Throwable $e) use ($config): void {
    error_log((string)$e);
    http_response_code(500);
    echo $config['debug'] ? $e->getMessage() : 'Баъдтар кӯшиш кунед ё ба администратор муроҷиат намоед.';
});

// 3. Настройка маршрутов (Router)
$router = new Router();

$router->get('/', [ProductController::class, 'home']);
$router->get('/products', [ProductController::class, 'index']);
$router->get('/products/create', [ProductController::class, 'create']);
$router->post('/products/create', [ProductController::class, 'store']);
$router->get('/products/{id}', [ProductController::class, 'show']);
$router->post('/products/{id}/delete', [ProductController::class, 'delete']);

$router->get('/shops', [ShopController::class, 'index']);

$router->get('/reports/categories', [ReportController::class, 'categories']);

// 4. Запуск роутера
$router->dispatch($_SERVER['REQUEST_METHOD'], $_SERVER['REQUEST_URI']);