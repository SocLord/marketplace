<?php
final class ReportController {
    private CategoryRepository $categories;

    public function __construct() {
        $this->categories = new CategoryRepository();
    }

    public function categories(): void {

        $rows = $this->categories->withProductCounts();
        
        
        require BASE_PATH . '/views/reports/categories.php';
    }
}