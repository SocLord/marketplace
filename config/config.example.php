<?php

return [
    'db' => [
        'driver' => 'sqlsrv',
        'server' => getenv('DB_SERVER') ?: 'localhost',
        'port' => (int)(getenv('DB_PORT') ?: 14330),
        'database' => getenv('DB_NAME') ?: 'MarketplaceDB',
        'user' => getenv('DB_USER') ?: 'Shahzod',
        'pass' => getenv('DB_PASS') ?: '1012',
    ],
    'debug' => (getenv('APP_DEBUG') ?: '1') === '1',
];