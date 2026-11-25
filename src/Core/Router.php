<?php
namespace App\Core;

final class Router {
    private array $routes = [];

    public function get(string $path, array $handler): self {
        $this->addRoute('GET', $path, $handler);
        return $this;
    }

    public function post(string $path, array $handler): self {
        $this->addRoute('POST', $path, $handler);
        return $this;
    }

    private function addRoute(string $method, string $path, array $handler): void {
        $this->routes[] = [
            'method' => $method,
            'regex' => $this->toRegex($path),
            'handler' => $handler,
        ];
    }

    /**
     * '/products/{id}/delete' → '#^/products/(?P<id>\d+)/delete$#'
     */
    private function toRegex(string $path): string {
        $pattern = preg_replace('#\{(\w+)\}#', '(?P<$1>\d+)', $path);
        return '#^' . $pattern . '$#';
    }

    public static function basePath(): string {
        $dir = str_replace('\\', '/', dirname($_SERVER['SCRIPT_NAME'] ?? '/'));
        return rtrim($dir, '/');
    }

    private function normalizePath(string $uri): string {
        $path = parse_url($uri, PHP_URL_PATH) ?? '/';
        $base = self::basePath();
        if ($base !== '' && str_starts_with($path, $base)) {
            $path = substr($path, strlen($base));
        }
        return rtrim($path, '/') ?: '/';
    }

    public function dispatch(string $method, string $uri): void {
        $path = $this->normalizePath($uri);
        $allowedMethods = [];

        foreach ($this->routes as $route) {
            if (!preg_match($route['regex'], $path, $matches)) {
                continue;
            }
            if ($route['method'] !== $method) {
                $allowedMethods[] = $route['method'];
                continue;
            }

            // Только именованные параметры: ['id' => '15'] → [15]
            $params = array_filter($matches, 'is_string', ARRAY_FILTER_USE_KEY);
            $params = array_map('intval', array_values($params));

            [$class, $action] = $route['handler'];
            (new $class())->$action(...$params);
            return;
        }

        if ($allowedMethods !== []) {
            http_response_code(405);
            header('Allow: ' . implode(',', array_unique($allowedMethods)));
            echo '405 Method Not Allowed';
            return;
        }

        http_response_code(404);
        echo (new View(BASE_PATH . '/views'))->render('errors/error', [
            'title' => 'Саҳифа ёфт нашуд',
            'code' => 404,
            'message' => 'Суроғаро санҷед ё ба саҳифаи асосӣ баргардед.',
        ]);
    }
}