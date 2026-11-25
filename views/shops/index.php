<?php
/**
 * @var App\Core\View $this
 * @var array $shops
 */
?>
<h1 class="mb-4">Мағозаҳо</h1>

<div class="row g-3">
    <?php foreach ($shops as $shop): ?>
    <div class="col-12 col-sm-6 col-lg-4">
        <div class="card h-100">
            <div class="card-body">
                <h2 class="h5"><?= $this->e($shop['name']) ?></h2>
                <p class="text-muted mb-0"><?= $this->e($shop['description'] ?? '') ?></p>
            </div>
        </div>
    </div>
    <?php endforeach; ?>
</div>