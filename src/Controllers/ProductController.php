<?php
final class ProductController {
    private ProductRepository $products;
    private CategoryRepository $categories;
    private ShopRepository $shops;

    public function __construct() {
        $this->products = new ProductRepository();
        $this->categories = new CategoryRepository();
        $this->shops = new ShopRepository();
    }

    public function index(): void {
        $page = max(1, (int)($_GET['page'] ?? 1));
        $filters = [
            'q' => $_GET['q'] ?? '',
            'category_id' => $_GET['category_id'] ?? null,
        ];
        
        $result = $this->products->paginate($filters, $page, 10);
        
        $items = $result['items'];
        $total = $result['total'];
        $page = $result['page'];
        $pages = $result['pages'];
        
        require BASE_PATH . '/views/products/index.php';
    }

    public function create(): void {
        $errors = [];
        $data = [];
        $categories = $this->categories->all();
        $shops = $this->shops->all();
        require BASE_PATH . '/views/products/form.php';
    }

public function store(): void {
    $v = new Validator($_POST);
    $v->required('name', 'Ном')->maxLength('name', 255, 'Ном')
      ->required('price', 'Нарх')
      ->required('category_id', 'Категория')
      ->required('shop_id', 'Мағоза')
      ->in('status', ['approved', 'pending', 'rejected'], 'Статус');

    if ($v->fails()) {
        http_response_code(422);
        $errors = $v->errors();
        $data = $_POST;
        $categories = $this->categories->all();
        $shops = $this->shops->all();
        require BASE_PATH . '/views/products/form.php';
        return;
    }

    // Собираем данные для INSERT
    // ВАЖНО: для NOT NULL строковых полей используем '' вместо null!
    $insertData = [
        'name' => trim($_POST['name']),
        'slug' => $this->generateSlug($_POST['name']),
        'description' => trim($_POST['description'] ?? '') ?: '',  // '' вместо null
        'price' => (float)$_POST['price'],
        'old_price' => !empty($_POST['old_price']) ? (float)$_POST['old_price'] : 0,  // 0 вместо null
        'stock_quantity' => (int)($_POST['stock_quantity'] ?? 0),
        'image_url' => trim($_POST['image_url'] ?? '') ?: '',  // '' вместо null
        'additional_images' => '',  // '' вместо null
        'status' => $_POST['status'] ?? 'pending',
        'moderation_comment' => '',  // '' вместо null
        'category_id' => (int)$_POST['category_id'],
        'shop_id' => (int)$_POST['shop_id'],
    ];

    $id = $this->products->create($insertData);
    header('Location: /products?msg=created');
    exit;
}
    public function delete(int $id): void {
        try {
            $this->products->delete($id);
            header('Location: /products?msg=deleted');
            exit;
        } catch (PDOException $e) {
            // 547 = нарушение внешнего ключа в SQL Server
            if (strpos($e->getMessage(), '547') !== false) {
                header('Location: /products?msg=fk_error');
                exit;
            }
            throw $e;
        }
    }

    private function generateSlug(string $text): string {
        $text = mb_strtolower($text, 'UTF-8');
        $text = preg_replace('/[^a-zа-яё0-9\s-]/u', '', $text);
        $text = preg_replace('/[\s-]+/', '-', $text);
        return trim($text, '-') ?: 'product-' . time();
    }
}