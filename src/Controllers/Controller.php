<?php
namespace App\Controllers;

use App\Core\View;

abstract class Controller {
    protected View $view;

    public function __construct() {
        $this->view = new View(BASE_PATH . '/views');
    }

    protected function render(string $template, array $data = [], int $status = 200, string $layout = 'layout'): void {
        http_response_code($status);
        echo $this->view->render($template, $data, $layout);
    }

    protected function redirect(string $path): never {
        header('Location: ' . $this->view->url($path));
        exit;
    }

    protected function json(array $data, int $status = 200): never {
        http_response_code($status);
        header('Content-Type: application/json; charset=utf-8');
        echo json_encode($data, JSON_UNESCAPED_UNICODE | JSON_UNESCAPED_SLASHES);
        exit;
    }

    protected function abort(int $code, string $message = ''): never {
        http_response_code($code);
        echo $this->view->render('errors/error', [
            'title' => "Хатои {$code}",
            'code' => $code,
            'message' => $message,
        ], null);
        exit;
    }
}