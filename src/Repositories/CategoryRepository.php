<?php
namespace App\Repositories;

use App\Core\Database;
use \PDO;

final class CategoryRepository {
    private \PDO $db;

    public function __construct(?\PDO $db = null) {
        $this->db = $db ?? Database::connection();
    }

    public function all(): array {
        return $this->db->query('SELECT * FROM categories ORDER BY name')->fetchAll();
    }

    public function ids(): array {
        return array_column($this->db->query('SELECT id FROM categories')->fetchAll(), 'id');
    }

    public function withProductCounts(): array {
        $sql = 'SELECT c.name, COUNT(p.id) AS total
                FROM categories c
                LEFT JOIN products p ON p.category_id = c.id
                GROUP BY c.id, c.name
                ORDER BY total DESC';
        return $this->db->query($sql)->fetchAll();
    }
}