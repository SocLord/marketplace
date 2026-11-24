<!doctype html>
<html lang="tg">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Иловаи мол — Бозори Онлайн</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
<nav class="navbar navbar-dark bg-primary mb-4">
    <div class="container">
        <a class="navbar-brand" href="/products">Бозори Онлайн</a>
    </div>
</nav>

<main class="container">
    <nav aria-label="breadcrumb">
        <ol class="breadcrumb">
            <li class="breadcrumb-item"><a href="/products">Каталог</a></li>
            <li class="breadcrumb-item active">Иловаи мол</li>
        </ol>
    </nav>

    <h1 class="mb-4">Илова кардани моли нав</h1>

    <form class="card card-body" method="post" action="/products/create" novalidate>
        <div class="mb-3">
            <label class="form-label" for="name">Номи мол *</label>
            <input class="form-control <?= isset($errors['name']) ? 'is-invalid' : '' ?>" 
                   type="text" id="name" name="name" required maxlength="255" 
                   value="<?= htmlspecialchars($data['name'] ?? '') ?>">
            <?php if (isset($errors['name'])): ?>
                <div class="invalid-feedback"><?= htmlspecialchars($errors['name']) ?></div>
            <?php endif; ?>
        </div>

        <div class="mb-3">
            <label class="form-label" for="description">Тавсиф</label>
            <textarea class="form-control" id="description" name="description" rows="3"><?= htmlspecialchars($data['description'] ?? '') ?></textarea>
        </div>

        <div class="row g-3 mb-3">
            <div class="col-md-4">
                <label class="form-label" for="price">Нарх (сомонӣ) *</label>
                <input class="form-control <?= isset($errors['price']) ? 'is-invalid' : '' ?>" 
                       type="number" step="0.01" id="price" name="price" required 
                       value="<?= htmlspecialchars($data['price'] ?? '') ?>">
                <?php if (isset($errors['price'])): ?>
                    <div class="invalid-feedback"><?= htmlspecialchars($errors['price']) ?></div>
                <?php endif; ?>
            </div>
            <div class="col-md-4">
                <label class="form-label" for="old_price">Нархи кӯҳна (ихтиёрӣ)</label>
                <input class="form-control" type="number" step="0.01" id="old_price" name="old_price" 
                       value="<?= htmlspecialchars($data['old_price'] ?? '') ?>">
            </div>
            <div class="col-md-4">
                <label class="form-label" for="stock_quantity">Миқдор дар анбор *</label>
                <input class="form-control" type="number" id="stock_quantity" name="stock_quantity" required 
                       value="<?= htmlspecialchars($data['stock_quantity'] ?? '0') ?>">
            </div>
        </div>

        <div class="row g-3 mb-4">
            <div class="col-md-4">
                <label class="form-label" for="category_id">Категория *</label>
                <select class="form-select <?= isset($errors['category_id']) ? 'is-invalid' : '' ?>" id="category_id" name="category_id" required>
                    <option value="">Интихоб кунед...</option>
                    <?php foreach ($categories as $cat): ?>
                        <option value="<?= $cat['id'] ?>" <?= ($data['category_id'] ?? '') == $cat['id'] ? 'selected' : '' ?>>
                            <?= htmlspecialchars($cat['name']) ?>
                        </option>
                    <?php endforeach; ?>
                </select>
                <?php if (isset($errors['category_id'])): ?>
                    <div class="invalid-feedback"><?= htmlspecialchars($errors['category_id']) ?></div>
                <?php endif; ?>
            </div>

            <div class="col-md-4">
                <label class="form-label" for="shop_id">Мағоза *</label>
                <select class="form-select <?= isset($errors['shop_id']) ? 'is-invalid' : '' ?>" id="shop_id" name="shop_id" required>
                    <option value="">Интихоб кунед...</option>
                    <?php foreach ($shops as $shop): ?>
                        <option value="<?= $shop['id'] ?>" <?= ($data['shop_id'] ?? '') == $shop['id'] ? 'selected' : '' ?>>
                            <?= htmlspecialchars($shop['name']) ?>
                        </option>
                    <?php endforeach; ?>
                </select>
                <?php if (isset($errors['shop_id'])): ?>
                    <div class="invalid-feedback"><?= htmlspecialchars($errors['shop_id']) ?></div>
                <?php endif; ?>
            </div>

            <div class="col-md-4">
                <label class="form-label" for="status">Статус</label>
                <select class="form-select" id="status" name="status">
                    <option value="pending" <?= ($data['status'] ?? 'pending') === 'pending' ? 'selected' : '' ?>>На модератсия (pending)</option>
                    <option value="approved" <?= ($data['status'] ?? '') === 'approved' ? 'selected' : '' ?>>Тасдиқ шуд (approved)</option>
                    <option value="rejected" <?= ($data['status'] ?? '') === 'rejected' ? 'selected' : '' ?>>Рад шуд (rejected)</option>
                </select>
            </div>
        </div>

        <div class="d-flex gap-2">
            <button type="submit" class="btn btn-primary">Захира кардан</button>
            <a href="/products" class="btn btn-outline-secondary">Бекор кардан</a>
        </div>
    </form>
</main>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>
</body>
</html>