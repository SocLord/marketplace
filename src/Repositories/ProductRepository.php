<?php
namespace App\Repositories;

use App\Core\Database;
use \PDO;

final class ProductRepository {
    public const FIELDS = [
        'name', 'slug', 'description', 'price', 'old_price', 
        'stock_quantity', 'image_url', 'additional_images', 
        'status', 'moderation_comment', 'category_id', 'shop_id'
    ];
    
    private const SELECT = '
        SELECT p.*, c.name AS category_name, s.name AS shop_name 
        FROM products p 
        JOIN categories c ON c.id = p.category_id 
        JOIN shops s ON s.id = p.shop_id
    ';

    private \PDO $db;

    public function __construct(?\PDO $db = null) {
        $this->db = $db ?? Database::connection();
    }

    private function buildWhere(array $filters): array {
        $conditions = [];
        $params = [];
        
        $q = trim((string)($filters['q'] ?? ''));
        if ($q !== '') {
            $conditions[] = "p.name LIKE :q";
            $params[':q'] = '%' . addcslashes($q, '%_\\') . '%';
        }
        
        if (!empty($filters['category_id'])) {
            $conditions[] = 'p.category_id = :category_id';
            $params[':category_id'] = (int)$filters['category_id'];
        }

        $where = $conditions ? 'WHERE ' . implode(' AND ', $conditions) : '';
        return [$where, $params];
    }

    public function paginate(array $filters, int $page, int $perPage = 10): array {
        [$where, $params] = $this->buildWhere($filters);

        $stmt = $this->db->prepare("SELECT COUNT(*) FROM products p {$where}");
        $stmt->execute($params);
        $total = (int)$stmt->fetchColumn();

        $pages = max(1, (int)ceil($total / $perPage));
        $page = min(max(1, $page), $pages);

        $sql = self::SELECT . "{$where} ORDER BY p.id DESC OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY";
        $stmt = $this->db->prepare($sql);
        
        foreach ($params as $name => $value) {
            $stmt->bindValue($name, $value);
        }

        $stmt->bindValue(':offset', (int)(($page - 1) * $perPage), \PDO::PARAM_INT);
        $stmt->bindValue(':limit', (int)$perPage, \PDO::PARAM_INT);
        $stmt->execute();

        return [
            'items' => $stmt->fetchAll(),
            'total' => $total,
            'page' => $page,
            'pages' => $pages
        ];
    }

    public function find(int $id): ?array {
        $stmt = $this->db->prepare(self::SELECT . ' WHERE p.id = :id');
        $stmt->execute([':id' => $id]);
        $row = $stmt->fetch();
        return $row ?: null;
    }

    public function create(array $data): int {
        $columns = implode(',', self::FIELDS);
        $placeholders = ':' . implode(',:', self::FIELDS);
        $stmt = $this->db->prepare("INSERT INTO products ({$columns}) VALUES ({$placeholders})");
        
        $execData = [];
        foreach (self::FIELDS as $field) {
            $execData[':' . $field] = $data[$field] ?? null;
        }
        
        $stmt->execute($execData);
        return (int)$this->db->lastInsertId();
    }

    public function delete(int $id): void {
        $stmt = $this->db->prepare('DELETE FROM products WHERE id = :id');
        $stmt->execute([':id' => $id]);
    }
}