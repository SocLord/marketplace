<?php
namespace App\Controllers;

use App\Core\Validator;
use App\Repositories\ProductRepository;
use App\Repositories\CategoryRepository;
use App\Repositories\ShopRepository;
use \PDOException;

final class ProductController extends Controller {
    private ProductRepository $products;
    private CategoryRepository $categories;
    private ShopRepository $shops;

    public function __construct() {
        parent::__construct();
        $this->products = new ProductRepository();
        $this->categories = new CategoryRepository();
        $this->shops = new ShopRepository();
    }

    public function home(): void {
        $this->redirect('/products');
    }

    public function index(): void {
        $page = max(1, (int)($_GET['page'] ?? 1));
        $filters = [
            'q' => $_GET['q'] ?? '',
            'category_id' => $_GET['category_id'] ?? null,
        ];
        
        $result = $this->products->paginate($filters, $page, 10);
        
        $this->render('products/index', [
            'title' => 'Каталог',
            'active' => 'products',
            'items' => $result['items'],
            'total' => $result['total'],
            'page' => $result['page'],
            'pages' => $result['pages'],
            'filters' => $filters,
        ]);
    }

    public function show(int $id): void {
        $product = $this->products->find($id);
        if ($product === null) {
            $this->abort(404, 'Мол ёфт нашуд');
        }
        $this->render('products/show', [
            'title' => $product['name'],
            'active' => 'products',
            'product' => $product,
        ]);
    }

    public function create(): void {
        $this->render('products/form', [
            'title' => 'Иловаи мол',
            'active' => 'products',
            'errors' => [],
            'data' => [],
            'categories' => $this->categories->all(),
            'shops' => $this->shops->all(),
        ]);
    }

    public function store(): void {
        $v = new Validator($_POST);
        $v->required('name', 'Ном')->maxLength('name', 255, 'Ном')
          ->required('price', 'Нарх')
          ->required('category_id', 'Категория')
          ->required('shop_id', 'Мағоза')
          ->in('status', ['approved', 'pending', 'rejected'], 'Статус');

        if ($v->fails()) {
            $this->render('products/form', [
                'title' => 'Иловаи мол',
                'active' => 'products',
                'errors' => $v->errors(),
                'data' => $_POST,
                'categories' => $this->categories->all(),
                'shops' => $this->shops->all(),
            ], 422);
            return;
        }

        $insertData = [
            'name' => trim($_POST['name']),
            'slug' => $this->generateSlug($_POST['name']),
            'description' => trim($_POST['description'] ?? '') ?: '',
            'price' => (float)$_POST['price'],
            'old_price' => !empty($_POST['old_price']) ? (float)$_POST['old_price'] : 0,
            'stock_quantity' => (int)($_POST['stock_quantity'] ?? 0),
            'image_url' => trim($_POST['image_url'] ?? '') ?: '',
            'additional_images' => '',
            'status' => $_POST['status'] ?? 'pending',
            'moderation_comment' => '',
            'category_id' => (int)$_POST['category_id'],
            'shop_id' => (int)$_POST['shop_id'],
        ];

        $id = $this->products->create($insertData);
        $this->redirect('/products?msg=created');
    }

    public function delete(int $id): void {
        try {
            $this->products->delete($id);
            $this->redirect('/products?msg=deleted');
        } catch (PDOException $e) {
            if (strpos($e->getMessage(), '547') !== false) {
                $this->redirect('/products?msg=fk_error');
            }
            throw $e;
        }
    }

    private function generateSlug(string $text): string {
        $text = mb_strtolower($text, 'UTF-8');
        $text = preg_replace('/[^a-zа-яё0-9\s-]/u', '', $text);
        $text = preg_replace('/[\s-]+/', '-', $text);
        $slug = trim($text, '-');
        if ($slug === '') {
            $slug = 'product-' . time();
        }
        return mb_substr($slug, 0, 255);
    }
}