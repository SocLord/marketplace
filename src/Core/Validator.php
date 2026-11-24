<?php
final class Validator {
    private array $data;
    private array $errors = [];

    public function __construct(array $data) {
        $this->data = $data;
    }

    public function required(string $field, string $label): self {
        if (empty(trim((string)$this->data[$field]))) {
            $this->errors[$field] = "Майдони «{$label}» ҳатмист.";
        }
        return $this;
    }

    public function maxLength(string $field, int $max, string $label): self {
        if (isset($this->data[$field]) && mb_strlen((string)$this->data[$field]) > $max) {
            $this->errors[$field] = "«{$label}» на зиёда аз {$max} аломат.";
        }
        return $this;
    }

    // НОВЫЙ МЕТОД: проверка, что значение есть в списке допустимых
    public function in(string $field, array $allowed, string $label): self {
        if (isset($this->data[$field]) && !in_array($this->data[$field], $allowed, true)) {
            $this->errors[$field] = "Қимати «{$label}» нодуруст аст.";
        }
        return $this;
    }

    // НОВЫЙ МЕТОД: добавить произвольную ошибку (для кастомных правил)
    public function addError(string $field, string $message): self {
        $this->errors[$field] = $message;
        return $this;
    }

    public function errors(): array {
        return $this->errors;
    }

    public function fails(): bool {
        return !empty($this->errors);
    }
}