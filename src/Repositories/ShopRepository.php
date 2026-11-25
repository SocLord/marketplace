<?php
namespace App\Repositories;

use App\Core\Database;
use PDO;
final class ShopRepository {
    private PDO $db;

    public function __construct(?PDO $db = null) {
        $this->db = $db ?? Database::connection();
    }

    public function all(): array {
        return $this->db->query('SELECT * FROM shops ORDER BY name')->fetchAll();
    }

    public function ids(): array {
        return array_column($this->db->query('SELECT id FROM shops')->fetchAll(), 'id');
    }
}