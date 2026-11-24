<?php
final class ReportController {
    private CategoryRepository $categories;

    public function __construct() {
        $this->categories = new CategoryRepository();
    }

    public function categories(): void {
        // Получаем данные из репозитория (там уже есть GROUP BY)
        $rows = $this->categories->withProductCounts();
        
        // Подключаем шаблон для отображения
        require BASE_PATH . '/views/reports/categories.php';
    }
}