<!doctype html>
<html lang="tg">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Каталог — Бозори Онлайн</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
<nav class="navbar navbar-dark bg-primary mb-4">
    <div class="container">
        <a class="navbar-brand" href="/products">Бозори Онлайн</a>
        <a class="btn btn-light btn-sm" href="/products/create">+ Илова кардани мол</a>
    </div>
</nav>

<main class="container">
    <?php if (isset($_GET['msg'])): ?>
        <div class="alert <?= $_GET['msg'] === 'fk_error' ? 'alert-warning' : 'alert-success' ?> alert-dismissible fade show">
            <?php
            $messages = [
                'created' => 'Мол бомуваффақият илова карда шуд!',
                'deleted' => 'Мол нест карда шуд!',
                'fk_error' => 'Ин молро нест кардан мумкин нест, зеро ба он дар ҷадвалҳои дигар истинодҳо ҳастанд.',
            ];
            echo $messages[$_GET['msg']] ?? 'Амал иҷро шуд.';
            ?>
            <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
        </div>
    <?php endif; ?>

    <h1 class="mb-4">Каталоги молҳо</h1>

    <form method="get" action="/products" class="card card-body mb-4">
        <div class="row g-3">
            <div class="col-md-5">
                <label class="form-label" for="q">Ҷустуҷӯ</label>
                <input type="search" class="form-control" id="q" name="q" 
                       value="<?= htmlspecialchars($filters['q'] ?? '') ?>" placeholder="Номи мол...">
            </div>
            <div class="col-md-4">
                <label class="form-label" for="category_id">Категория</label>
                <select class="form-select" id="category_id" name="category_id">
                    <option value="">Ҳама категорияҳо</option>
                    <option value="1" <?= ($filters['category_id'] ?? '') == 1 ? 'selected' : '' ?>>Электроника</option>
                    <option value="2" <?= ($filters['category_id'] ?? '') == 2 ? 'selected' : '' ?>>Одежда</option>
                    <option value="3" <?= ($filters['category_id'] ?? '') == 3 ? 'selected' : '' ?>>Бытовая техника</option>
                </select>
            </div>
            <div class="col-md-3 d-flex align-items-end">
                <button type="submit" class="btn btn-primary w-100">Ёфтан</button>
            </div>
        </div>
    </form>

    <div class="table-responsive">
        <table class="table table-hover align-middle">
            <thead class="table-light">
                <tr>
                    <th>ID</th>
                    <th>Ном</th>
                    <th>Категория</th>
                    <th>Мағоза</th>
                    <th>Нарх</th>
                    <th>Дар анбор</th>
                    <th>Статус</th>
                    <th class="text-end">Амалҳо</th>
                </tr>
            </thead>
            <tbody>
                <?php foreach ($items as $product): ?>
                <tr>
                    <td><?= $product['id'] ?></td>
                    <td><?= htmlspecialchars($product['name']) ?></td>
                    <td><?= htmlspecialchars($product['category_name']) ?></td>
                    <td><?= htmlspecialchars($product['shop_name']) ?></td>
                    <td><?= number_format($product['price'], 2, ',', ' ') ?> с.</td>
                    <td><?= (int)$product['stock_quantity'] ?></td>
                    <td>
                        <span class="badge <?= $product['status'] === 'approved' ? 'bg-success' : ($product['status'] === 'pending' ? 'bg-warning text-dark' : 'bg-danger') ?>">
                            <?= htmlspecialchars($product['status']) ?>
                        </span>
                    </td>
                    <td class="text-end">
                        <form method="post" action="/products/<?= $product['id'] ?>/delete" style="display:inline" onsubmit="return confirm('Оё шумо мутмаин ҳастед, ки ин молро нест кардан мехоҳед?')">
                            <button type="submit" class="btn btn-sm btn-outline-danger">Нест кардан</button>
                        </form>
                    </td>
                </tr>
                <?php endforeach; ?>
                
                <?php if (empty($items)): ?>
                <tr>
                    <td colspan="8" class="text-center py-4 text-muted">Ягон мол ёфт нашуд.</td>
                </tr>
                <?php endif; ?>
            </tbody>
        </table>
    </div>

    <?php if ($pages > 1): ?>
    <nav aria-label="Саҳифаҳо">
        <ul class="pagination justify-content-center">
            <?php for ($i = 1; $i <= $pages; $i++): ?>
                <li class="page-item <?= $i === $page ? 'active' : '' ?>">
                    <a class="page-link" href="/products?page=<?= $i ?>&q=<?= urlencode($filters['q'] ?? '') ?>&category_id=<?= urlencode($filters['category_id'] ?? '') ?>">
                        <?= $i ?>
                    </a>
                </li>
            <?php endfor; ?>
        </ul>
    </nav>
    <?php endif; ?>
    
    <div class="text-muted small mt-3">Ҳамагӣ: <?= $total ?> мол. Саҳифа <?= $page ?> аз <?= $pages ?>.</div>
</main>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>
</body>
</html>