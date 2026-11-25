<div class="container py-5 text-center">
    <h1 class="display-1 text-muted"><?= $code ?? 404 ?></h1>
    <h2><?= $this->e($title ?? 'Хато') ?></h2>
    <p class="lead"><?= $this->e($message ?? '') ?></p>
    <a href="<?= $this->url('/products') ?>" class="btn btn-primary">Ба каталог баргаштан</a>
</div>