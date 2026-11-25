<?php
namespace App\Controllers;

use App\Repositories\ShopRepository;

final class ShopController extends Controller {
    private ShopRepository $shops;

    public function __construct() {
        parent::__construct();
        $this->shops = new ShopRepository();
    }

    public function index(): void {
        $this->render('shops/index', [
            'title' => 'Мағозаҳо',
            'active' => 'shops',
            'shops' => $this->shops->all(),
        ]);
    }
}