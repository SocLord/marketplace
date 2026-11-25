<!doctype html>
<html lang="tg">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title><?= $this->e($title ?? 'Бозори Онлайн') ?></title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
    <?= $this->partial('nav', ['active' => $active ?? '']) ?>
    
    <main class="container py-4">
        <?= $content ?>
    </main>

    <footer class="app-footer container py-4 text-center">
        <p class="mb-0">© <?= date('Y') ?> Бозори Онлайн · Душанбе, Тоҷикистон</p>
    </footer>

    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>
</body>
</html>