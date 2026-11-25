<?php
namespace App\Controllers;

use App\Repositories\CategoryRepository;

final class ReportController extends Controller {
    private CategoryRepository $categories;

    public function __construct() {
        parent::__construct();
        $this->categories = new CategoryRepository();
    }

    public function categories(): void {
        $rows = $this->categories->withProductCounts();
        $this->render('reports/categories', [
            'title' => 'Ҳисобот',
            'active' => 'reports',
            'rows' => $rows,
        ]);
    }
}