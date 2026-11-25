<?php
/**
 * @var App\Core\View $this
 * @var array $rows
 */
?>

<nav aria-label="breadcrumb">
    <ol class="breadcrumb">
        <li class="breadcrumb-item"><a href="<?= $this->url('/products') ?>">Каталог</a></li>
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
                    <td><?= $this->e($row['name']) ?></td>
                    <td class="text-end"><?= (int)$row['total'] ?></td>
                </tr>
                <?php endforeach; ?>
            </tbody>
        </table>
    </div>
</div>

<div class="mt-4">
    <a href="<?= $this->url('/products') ?>" class="btn btn-secondary">← Ба каталог баргаштан</a>
</div>