<?php
namespace App\Core;

final class Database {
    private static ?\PDO $pdo = null;
    private static array $config = [];

    public static function configure(array $config): void {
        self::$config = $config;
    }

    public static function connection(): \PDO {
        if (self::$pdo === null) {
            $c = self::$config;
            $dsn = sprintf('sqlsrv:Server=%s,%d;Database=%s', 
                $c['server'], $c['port'], $c['database']);
            
            self::$pdo = new \PDO($dsn, $c['user'], $c['pass'], [
                \PDO::ATTR_ERRMODE => \PDO::ERRMODE_EXCEPTION,
                \PDO::ATTR_DEFAULT_FETCH_MODE => \PDO::FETCH_ASSOC,
            ]);
        }
        return self::$pdo;
    }
}