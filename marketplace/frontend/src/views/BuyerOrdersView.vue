<template>
  <div class="orders-page">
    <div class="container">
      <div class="page-header">
        <h1 class="page-title">{{ currentLang === 'ru' ? 'Мои заказы' : 'Фармоишҳои ман' }}</h1>
        <p class="page-subtitle">{{ currentLang === 'ru' ? 'История покупок и отслеживание доставки посылок от магазинов' : 'Таърихи харидҳо ва пайгирии ҳолати расонидани бастаҳои мағозаҳо' }}</p>
      </div>

      <div v-if="loading" class="loading-state">
        <span>{{ currentLang === 'ru' ? 'Загрузка заказов...' : 'Боргирии фармоишҳо...' }}</span>
      </div>

      <div v-else-if="orders.length === 0" class="empty-orders-card card">
        <Package :size="48" color="#94a3b8" />
        <h3>{{ currentLang === 'ru' ? 'У вас пока нет заказов' : 'Шумо ҳанӯз фармоиш надоред' }}</h3>
        <p>{{ currentLang === 'ru' ? 'Когда вы купите товары на маркетплейсе, они появятся на этой странице.' : 'Вақте ки шумо дар маркетплейс мол мехаред, онҳо дар ин саҳифа нишон дода мешаванд.' }}</p>
        <router-link to="/catalog" class="btn btn-primary">{{ t('goToCatalog') }}</router-link>
      </div>

      <div v-else class="orders-list">
        <div v-for="order in orders" :key="order.id" class="order-card card">
          <!-- Шапка заказа -->
          <div class="order-card-header">
            <div class="order-meta">
              <span class="order-number">{{ currentLang === 'ru' ? 'Заказ' : 'Фармоиш' }}: #{{ order.order_number }}</span>
              <span class="order-date">{{ formatDate(order.created_at) }}</span>
            </div>
            <div class="order-sum-status">
              <span class="order-sum">{{ formatPrice(order.total_amount) }} {{ t('currency') }}</span>
              <span :class="['badge', getOrderStatusBadge(order.status)]">
                {{ getOrderStatusLabel(order.status) }}
              </span>
            </div>
          </div>

          <!-- Адрес доставки -->
          <div class="order-address-bar">
            <span>📍 <strong>{{ currentLang === 'ru' ? 'Адрес' : 'Суроға' }}:</strong> {{ order.shipping_address }}</span>
            <span>📞 <strong>{{ currentLang === 'ru' ? 'Получатель' : 'Гиранда' }}:</strong> {{ order.recipient_name }} ({{ order.recipient_phone }})</span>
          </div>

          <!-- Разделенные подзаказы продавцов (Специальное требование!) -->
          <div class="sub-orders-list">
            <div 
              v-for="sub in order.sub_orders" 
              :key="sub.id" 
              class="sub-order-card">
              
              <!-- Шапка подзаказа конкретного магазина -->
              <div class="sub-order-header">
                <div class="sub-shop-info">
                  <Store :size="16" color="#d97706" />
                  <span class="sub-shop-name">{{ currentLang === 'ru' ? 'Магазин' : 'Мағозаи' }}: <strong>{{ sub.shop?.name || (currentLang === 'ru' ? `Магазин #${sub.shop_id}` : `Мағозаи #${sub.shop_id}`) }}</strong></span>
                </div>
                <div class="sub-tracking-info">
                  <span v-if="sub.shipping_tracking" class="tracking-pill">
                    📦 {{ currentLang === 'ru' ? 'Трек-код' : 'Трек-код' }}: {{ sub.shipping_tracking }}
                  </span>
                  <span :class="['badge', getSubOrderStatusBadge(sub.status)]">
                    {{ getSubOrderStatusLabel(sub.status) }}
                  </span>
                </div>
              </div>

              <!-- Товары в этой посылке -->
              <div class="sub-items-table">
                <div v-for="item in sub.items" :key="item.id" class="sub-item-row">
                  <img :src="item.image_url" :alt="item.product_name" class="sub-item-img" />
                  <div class="sub-item-title-col">
                    <router-link :to="`/product/${item.product_id}`" class="sub-item-title">
                      {{ item.product_name }}
                    </router-link>
                    <span class="sub-item-price-unit">{{ formatPrice(item.product_price) }} {{ currentLang === 'ru' ? 'сом.' : 'с.' }} × {{ item.quantity }} {{ currentLang === 'ru' ? 'шт.' : 'дона' }}</span>
                  </div>
                  <div class="sub-item-total">
                    {{ formatPrice(item.total_price) }} {{ t('currency') }}
                  </div>
                  <router-link :to="`/product/${item.product_id}#reviews-section`" class="btn btn-secondary btn-sm review-btn">
                    <Star :size="13" />
                    <span>{{ currentLang === 'ru' ? 'Оставить отзыв' : 'Тақриз навиштан' }}</span>
                  </router-link>
                </div>
              </div>

              <!-- Итог посылки -->
              <div class="sub-order-footer">
                <span>{{ currentLang === 'ru' ? 'Сумма этой посылки:' : 'Маблағи ин қисм:' }}</span>
                <strong>{{ formatPrice(sub.subtotal_amount) }} {{ t('currency') }}</strong>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Package, Store, Star } from 'lucide-vue-next'
import { api } from '../api'
import { t, currentLang } from '../utils/i18n'

const loading = ref(true)
const orders = ref([])

function formatPrice(val) {
  if (!val) return '0'
  return new Intl.NumberFormat('ru-RU').format(val)
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  const locale = currentLang.value === 'ru' ? 'ru-RU' : 'tg-TJ'
  return d.toLocaleDateString(locale, { day: 'numeric', month: 'long', year: 'numeric', hour: '2-digit', minute: '2-digit' })
}

function getOrderStatusLabel(st) {
  if (currentLang.value === 'ru') {
    switch (st) {
      case 'completed': return 'Выполнен (Всё доставлено)'
      case 'processing': return 'В обработке и сборке'
      case 'cancelled': return 'Отменён'
      default: return 'Выполняется'
    }
  }
  switch (st) {
    case 'completed': return 'Иҷро шуд (Ҳамагӣ расонида шуд)'
    case 'processing': return 'Дар ҳоли коркард ва ҷамъоварӣ'
    case 'cancelled': return 'Бекор шуд'
    default: return 'Дар иҷро'
  }
}

function getOrderStatusBadge(st) {
  switch (st) {
    case 'completed': return 'badge-delivered'
    case 'cancelled': return 'badge-rejected'
    default: return 'badge-processing'
  }
}

function getSubOrderStatusLabel(st) {
  if (currentLang.value === 'ru') {
    switch (st) {
      case 'pending': return 'Ожидает подтверждения продавца'
      case 'processing': return 'Собирается'
      case 'shipped': return 'Отправлено (В пути)'
      case 'delivered': return 'Доставлено ✓'
      case 'cancelled': return 'Отменено'
      default: return st
    }
  }
  switch (st) {
    case 'pending': return 'Интизори тасдиқи фурӯшанда'
    case 'processing': return 'Дар ҳоли ҷамъоварӣ'
    case 'shipped': return 'Ирсол шуд (Дар роҳ)'
    case 'delivered': return 'Расонида шуд ✓'
    case 'cancelled': return 'Бекор шуд'
    default: return st
  }
}

function getSubOrderStatusBadge(st) {
  switch (st) {
    case 'pending': return 'badge-pending'
    case 'processing': return 'badge-processing'
    case 'shipped': return 'badge-shipped'
    case 'delivered': return 'badge-delivered'
    case 'cancelled': return 'badge-rejected'
    default: return 'badge-pending'
  }
}

async function loadOrders() {
  loading.value = true
  try {
    const res = await api.orders.getMyOrders()
    orders.value = res.data || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadOrders()
})
</script>

<style scoped>
.orders-page {
  padding: 24px 0 60px 0;
}

.page-header {
  margin-bottom: 24px;
}

.page-title {
  font-size: 26px;
  font-weight: 800;
}

.page-subtitle {
  font-size: 13px;
  color: var(--text-muted);
  margin-top: 4px;
}

.orders-list {
  display: flex;
  flex-direction: column;
  gap: 22px;
}

.order-card {
  background: #ffffff;
  overflow: hidden;
}

.order-card-header {
  padding: 16px 20px;
  background: var(--bg-page);
  border-bottom: 1px solid var(--border);
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.order-meta {
  display: flex;
  align-items: center;
  gap: 12px;
}

.order-number {
  font-size: 15px;
  font-weight: 800;
  color: var(--text-main);
}

.order-date {
  font-size: 12px;
  color: var(--text-muted);
}

.order-sum-status {
  display: flex;
  align-items: center;
  gap: 14px;
}

.order-sum {
  font-size: 16px;
  font-weight: 800;
  color: #0f172a;
}

.order-address-bar {
  padding: 10px 20px;
  background: #fafafa;
  border-bottom: 1px solid var(--border-light);
  font-size: 12px;
  color: var(--text-secondary);
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
}

/* Подзаказы */
.sub-orders-list {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.sub-order-card {
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  overflow: hidden;
  background: #ffffff;
}

.sub-order-header {
  padding: 12px 16px;
  background: #fffbeb;
  border-bottom: 1px solid rgba(245, 158, 11, 0.15);
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
}

.sub-shop-info {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}

.sub-tracking-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.tracking-pill {
  font-size: 11px;
  font-weight: 700;
  background: #fef3c7;
  color: #92400e;
  padding: 3px 8px;
  border-radius: 4px;
}

.sub-items-table {
  padding: 0 16px;
}

.sub-item-row {
  display: grid;
  grid-template-columns: 48px 1fr auto auto;
  gap: 14px;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid var(--border-light);
}

.sub-item-row:last-child {
  border-bottom: none;
}

.sub-item-img {
  width: 48px;
  height: 48px;
  border-radius: 6px;
  object-fit: cover;
}

.sub-item-title-col {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.sub-item-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-main);
}

.sub-item-title:hover {
  color: var(--primary-hover);
}

.sub-item-price-unit {
  font-size: 11px;
  color: var(--text-muted);
}

.sub-item-total {
  font-size: 13px;
  font-weight: 700;
  color: var(--text-main);
}

.review-btn {
  font-size: 11px;
  padding: 4px 10px;
}

.sub-order-footer {
  padding: 10px 16px;
  background: #fafafa;
  border-top: 1px solid var(--border-light);
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  font-size: 12px;
}

.empty-orders-card {
  padding: 60px 20px;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14px;
}

.empty-orders-card h3 {
  font-size: 18px;
  font-weight: 800;
}

.empty-orders-card p {
  color: var(--text-muted);
  font-size: 13px;
  max-width: 400px;
}
</style>
