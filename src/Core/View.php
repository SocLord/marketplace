<?php
namespace App\Core;

final class View {
    private string $viewsPath;

    public function __construct(string $viewsPath) {
        $this->viewsPath = rtrim($viewsPath, '/');
    }

    public function render(string $template, array $data = [], ?string $layout = 'layout'): string {
        $content = $this->renderFile($template, $data);
        if ($layout === null) {
            return $content;
        }
        return $this->renderFile($layout, $data + ['content' => $content]);
    }

    private function renderFile(string $template, array $data): string {
        $file = $this->viewsPath . '/' . $template . '.php';
        if (!is_file($file)) {
            throw new \RuntimeException("Шаблон ёфт нашуд: {$template}");
        }
        extract($data, EXTR_SKIP);
        ob_start();
        require $file;
        return (string) ob_get_clean();
    }

    /** Экранирование для защиты от XSS */
    public function e(mixed $value): string {
        return htmlspecialchars((string) $value, ENT_QUOTES | ENT_SUBSTITUTE, 'UTF-8');
    }

    /** Ссылка с учётом basePath */
    public function url(string $path = '/'): string {
        return Router::basePath() . $path;
    }

    /** Фрагмент из views/partials/ */
    public function partial(string $name, array $data = []): string {
        return $this->render('partials/' . $name, $data, null);
    }

    /** Форматирование даты */
    public function date(?string $iso): string {
        if (!$iso) return '';
        [$year, $month, $day] = explode('-', $iso);
        return "{$day}.{$month}.{$year}";
    }
}