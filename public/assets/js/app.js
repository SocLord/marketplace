// ===== УТИЛИТЫ =====
function debounce(fn, delay = 300) {
    let timerId;
    return (...args) => {
        clearTimeout(timerId);
        timerId = setTimeout(() => fn(...args), delay);
    };
}

function formatPrice(price) {
    if (!price) return '';
    return price.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ' ') + ' с.';
}

// ===== ВАЛИДАЦИЯ ФОРМЫ =====
const FIELD_MESSAGES = {
    customer_name: { valueMissing: 'Номи муштариро ворид кунед.' },
    phone: {
        valueMissing: 'Рақами телефонро ворид кунед.',
        patternMismatch: 'Формати телефон: +992 XX XXX XX XX'
    },
    address: { valueMissing: 'Суроғаро ворид кунед.' },
    product_id: { valueMissing: 'Маҳсулотро интихоб кunед.' }
};

function messageFor(field) {
    const v = field.validity;
    if (v.valid) return '';
    if (v.customError) return field.validationMessage;
    const custom = FIELD_MESSAGES[field.name] ?? {};
    if (v.valueMissing) return custom.valueMissing ?? 'Ин майдонро пур кунед.';
    if (v.patternMismatch) return custom.patternMismatch ?? 'Формати нодуруст.';
    if (v.tooLong) return `На зиёда аз ${field.maxLength} аломат.`;
    return 'Қимати майдонро санҷед.';
}

function showFieldState(field) {
    const message = messageFor(field);
    const feedback = document.getElementById(`${field.id}-error`);
    field.classList.toggle('is-invalid', message !== '');
    if (message) field.setAttribute('aria-invalid', 'true');
    else field.removeAttribute('aria-invalid');
    if (feedback) feedback.textContent = message;
    return message === '';
}

function checkPhone(form) {
    const phone = form.elements.phone;
    phone.setCustomValidity('');
    if (phone.value && !phone.value.startsWith('+992')) {
        phone.setCustomValidity('Телефон бояд бо +992 оғоз шавад.');
    }
}

function initOrderForm() {
    const form = document.getElementById('order-form');
    if (!form) return;

    const fields = () => form.querySelectorAll('input, select, textarea');

    const refresh = (field) => {
        checkPhone(form);
        showFieldState(field);
    };

    form.addEventListener('focusout', (event) => {
        if (event.target.matches('input, select, textarea')) {
            refresh(event.target);
        }
    });

    form.addEventListener('input', (event) => {
        if (event.target.classList.contains('is-invalid')) {
            refresh(event.target);
        }
    });

    form.addEventListener('submit', (event) => {
        checkPhone(form);
        let firstInvalid = null;
        fields().forEach((field) => {
            if (!showFieldState(field) && !firstInvalid) {
                firstInvalid = field;
            }
        });
        if (firstInvalid) {
            event.preventDefault();
            firstInvalid.focus();
        }
    });
}

// ===== ЖИВОЙ ПОИСК =====
const STATUS_LABELS = { active: 'Дар анбор', pending: 'На модератсия', out: 'Том шуд' };
const STATUS_CLASS = { active: 'status-active', pending: 'status-pending', out: 'status-out' };

function createCell(text, className = '') {
    const td = document.createElement('td');
    td.textContent = text ?? ''; // ЗАЩИТА ОТ XSS
    if (className) td.className = className;
    return td;
}

function createProductRow(product) {
    const tr = document.createElement('tr');

    const photoCell = document.createElement('td');
    const img = document.createElement('img');
    img.src = product.image || 'https://placehold.co/60x60?text=No+Photo';
    img.alt = product.name;
    img.className = 'product-thumb-sm';
    photoCell.append(img);

    const nameCell = document.createElement('td');
    const link = document.createElement('a');
    link.href = `product.html?id=${encodeURIComponent(product.id)}`;
    link.textContent = product.name; // ЗАЩИТА ОТ XSS
    nameCell.append(link);

    const statusCell = document.createElement('td');
    const badge = document.createElement('span');
    badge.className = `badge status ${STATUS_CLASS[product.status] ?? ''}`;
    badge.textContent = STATUS_LABELS[product.status] ?? product.status;
    statusCell.append(badge);

    tr.append(
        photoCell,
        nameCell,
        createCell(product.category),
        createCell(formatPrice(product.price)),
        createCell(String(product.stock), 'd-none d-md-table-cell'),
        statusCell
    );
    return tr;
}

function renderProducts(tbody, products) {
    tbody.replaceChildren();
    if (products.length === 0) {
        const tr = document.createElement('tr');
        const td = document.createElement('td');
        td.colSpan = 6;
        td.textContent = 'Маҳсулот ёфт нашуд.';
        td.className = 'text-center text-muted py-4';
        tr.append(td);
        tbody.append(tr);
        return;
    }
    products.forEach((product) => tbody.append(createProductRow(product)));
}

function initProductSearch() {
    const input = document.getElementById('product-search');
    const tbody = document.getElementById('products-tbody');
    if (!input || !tbody) return;

    const statusBox = document.getElementById('product-search-status');
    const setStatus = (text) => { if (statusBox) statusBox.textContent = text; };
    let controller = null;

    async function search(query) {
        controller?.abort();
        controller = new AbortController();
        tbody.setAttribute('aria-busy', 'true');
        setStatus('Ҷустуҷӯ...');

        try {
            const response = await fetch(`/api/products?q=${encodeURIComponent(query)}`, {
                headers: { Accept: 'application/json' },
                signal: controller.signal
            });
            if (!response.ok) throw new Error(`HTTP ${response.status}`);
            const products = await response.json();
            renderProducts(tbody, products);
            setStatus(products.length ? `Ёфт шуд: ${products.length}` : 'Маҳсулот ёфт нашуд.');
        } catch (error) {
            if (error.name === 'AbortError') return;
            console.error('Хатои ҷустуҷӯ:', error);
            setStatus('Ҷустуҷӯро иҷро кардан муяссар нашуд.');
        } finally {
            tbody.removeAttribute('aria-busy');
        }
    }

    const debouncedSearch = debounce(search, 300);
    input.addEventListener('input', () => debouncedSearch(input.value.trim()));
}

// ===== МОДАЛЬНОЕ ОКНО =====
function initDeleteConfirm() {
    const modalElement = document.getElementById('confirm-delete-modal');
    if (!modalElement || !window.bootstrap) return;

    const form = document.getElementById('confirm-delete-form');
    const nameElement = document.getElementById('confirm-delete-name');
    const submitButton = form.querySelector('[type="submit"]');
    const modal = window.bootstrap.Modal.getOrCreateInstance(modalElement);
    let lastTrigger = null;

    document.addEventListener('click', (event) => {
        const button = event.target.closest('[data-delete-url]');
        if (!button) return;

        lastTrigger = button;
        form.action = button.dataset.deleteUrl;
        nameElement.textContent = button.dataset.productName ?? '';
        submitButton.disabled = false;
        modal.show();
    });

    modalElement.addEventListener('hidden.bs.modal', () => {
        if (lastTrigger?.isConnected) lastTrigger.focus();
    });

    form.addEventListener('submit', () => {
        submitButton.disabled = true;
    });
}

// ===== ЗАПУСК =====
initOrderForm();
initProductSearch();
initDeleteConfirm();