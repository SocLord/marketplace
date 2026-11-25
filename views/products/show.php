<?php
/**
 * @var App\Core\View $this
 * @var array $product
 */
?>
<nav aria-label="breadcrumb">
    <ol class="breadcrumb">
        <li class="breadcrumb-item"><a href="<?= $this->url('/products') ?>">Каталог</a></li>
        <li class="breadcrumb-item active"><?= $this->e($product['name']) ?></li>
    </ol>
</nav>

<div class="card">
    <div class="card-body">
        <h1 class="h3 mb-3"><?= $this->e($product['name']) ?></h1>

        <?php if (!empty($product['image_url'])): ?>
            <img src="<?= $this->e($product['image_url']) ?>" alt="<?= $this->e($product['name']) ?>"
                 class="rounded mb-3" style="max-width:300px">
        <?php endif; ?>

        <dl class="row mb-4">
            <dt class="col-sm-3">Категория</dt>
            <dd class="col-sm-9"><?= $this->e($product['category_name']) ?></dd>
            
            <dt class="col-sm-3">Мағоза</dt>
            <dd class="col-sm-9"><?= $this->e($product['shop_name']) ?></dd>
            
            <dt class="col-sm-3">Нарх</dt>
            <dd class="col-sm-9"><?= number_format((float)$product['price'], 2, ',', ' ') ?> с.</dd>
            
            <dt class="col-sm-3">Дар анбор</dt>
            <dd class="col-sm-9"><?= (int)$product['stock_quantity'] ?></dd>
            
            <dt class="col-sm-3">Статус</dt>
            <dd class="col-sm-9"><?= $this->e($product['status']) ?></dd>
        </dl>

        <?php if (!empty($product['description'])): ?>
            <h2 class="h5">Тавсиф</h2>
            <p><?= nl2br($this->e($product['description'])) ?></p>
        <?php endif; ?>

        <div class="d-flex gap-2">
            <a href="<?= $this->url('/products') ?>" class="btn btn-secondary">← Ба каталог</a>
            <form method="post" action="<?= $this->url('/products/' . (int)$product['id'] . '/delete') ?>"
                  onsubmit="return confirm('Нест кардан?')">
                <button type="submit" class="btn btn-outline-danger">Нест кардан</button>
            </form>
        </div>
    </div>
</div>