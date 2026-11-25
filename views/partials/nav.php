<nav class="navbar navbar-dark bg-primary mb-4">
    <div class="container">
        <a class="navbar-brand" href="<?= $this->url('/products') ?>">Бозори Онлайн</a>
        <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#mainNav">
            <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="mainNav">
            <ul class="navbar-nav me-auto">
                <li class="nav-item">
                    <a class="nav-link <?= $active === 'products' ? 'active' : '' ?>" 
                       href="<?= $this->url('/products') ?>">Каталог</a>
                </li>
                <li class="nav-item">
                    <a class="nav-link <?= $active === 'reports' ? 'active' : '' ?>" 
                       href="<?= $this->url('/reports/categories') ?>">Ҳисобот</a>
                </li>
            </ul>
            <a class="btn btn-light btn-sm" href="<?= $this->url('/products/create') ?>">+ Илова кардан</a>
        </div>
    </div>
</nav>