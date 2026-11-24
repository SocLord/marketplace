<!doctype html>
<html lang="tg">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Ҳисобот — Бозори Онлайн</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
<nav class="navbar navbar-dark bg-primary mb-4">
    <div class="container">
        <a class="navbar-brand" href="/products">Бозори Онлайн</a>
    </div>
</nav>

<main class="container py-4">
    <nav aria-label="Занҷири навигатсионӣ">
        <ol class="breadcrumb">
            <li class="breadcrumb-item"><a href="/products">Каталог</a></li>
            <li class="breadcrumb-item active">Ҳисобот</li>
        </ol>
    </nav>

    <h1 class="mb-4">Ҳисобот: Миқдори молҳо аз рӯи категорияҳо</h1>
    
    <div class="card">
        <div class="card-body">
            <table class="table table-bordered table-hover mb-0">
                <thead class="table-light">
                    <tr>
                        <th>Категория</th>
                        <th class="text-end">Миқдори молҳо</th>
                    </tr>
                </thead>
                <tbody>
                    <?php foreach ($rows as $row): ?>
                    <tr>
                        <td><?= htmlspecialchars($row['name']) ?></td>
                        <td class="text-end"><?= (int)$row['total'] ?></td>
                    </tr>
                    <?php endforeach; ?>
                </tbody>
            </table>
        </div>
    </div>

    <div class="mt-4">
        <a href="/products" class="btn btn-secondary">← Ба каталог баргаштан</a>
    </div>
</main>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>
</body>
</html>