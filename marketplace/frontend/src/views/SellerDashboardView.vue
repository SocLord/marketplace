<template>
  <div class="seller-dashboard-page">
    <div class="container">
      <!-- Шапка кабинета продавца -->
      <div class="dashboard-header">
        <div class="header-left">
          <span class="badge badge-seller-tag">{{ t('sellerTitle') }}</span>
          <h1 class="page-title">{{ shopInfo?.name || t('sellerTitle') }}</h1>
          <p class="page-desc">{{ t('sellerSubtitle') }}</p>
        </div>

        <button @click="openAddProductModal" class="btn btn-primary btn-add-prod">
          <Plus :size="18" />
          <span>{{ t('addNewProduct') }}</span>
        </button>
      </div>

      <!-- Вкладки навигации кабинета -->
      <div class="dashboard-nav">
        <button 
          :class="['nav-tab', activeTab === 'stats' && 'active']"
          @click="activeTab = 'stats'">
          <BarChart3 :size="18" />
          <span>{{ t('tabStats') }}</span>
        </button>

        <button 
          :class="['nav-tab', activeTab === 'products' && 'active']"
          @click="activeTab = 'products'">
          <Package :size="18" />
          <span>{{ t('tabProducts') }} ({{ products.length }})</span>
        </button>

        <button 
          :class="['nav-tab', activeTab === 'orders' && 'active']"
          @click="activeTab = 'orders'">
          <ShoppingBag :size="18" />
          <span>{{ t('tabOrders') }} ({{ orders.length }})</span>
        </button>

        <button 
          :class="['nav-tab', activeTab === 'settings' && 'active']"
          @click="activeTab = 'settings'">
          <Settings :size="18" />
          <span>{{ t('tabSettings') }}</span>
        </button>
      </div>

      <!-- ==================== ВКЛАДКА 1: ОМОРИ ФУРӮШ АЗ РӮИ МОҲҲО ==================== -->
      <div v-if="activeTab === 'stats'" class="tab-content">
        <!-- Карточки показателей (KPI Cards) -->
        <div class="kpi-grid">
          <div class="kpi-card card">
            <div class="kpi-icon icon-yellow">
              <TrendingUp :size="24" />
            </div>
            <div class="kpi-data">
              <span class="kpi-label">{{ t('totalRevenue') }}</span>
              <span class="kpi-value">{{ formatPrice(stats?.total_revenue) }} <small>с.</small></span>
              <span class="kpi-hint">сомони</span>
            </div>
          </div>

          <div class="kpi-card card">
            <div class="kpi-icon icon-amber">
              <Calendar :size="24" />
            </div>
            <div class="kpi-data">
              <span class="kpi-label">{{ t('monthRevenue') }}</span>
              <span class="kpi-value">{{ formatPrice(stats?.current_month_revenue) }} <small>с.</small></span>
              <span class="kpi-hint">сомони</span>
            </div>
          </div>

          <div class="kpi-card card">
            <div class="kpi-icon icon-dark">
              <ShoppingBag :size="24" />
            </div>
            <div class="kpi-data">
              <span class="kpi-label">{{ t('tabOrders') }}</span>
              <span class="kpi-value">{{ stats?.total_orders_count || 0 }}</span>
              <span class="kpi-hint">заказов магазина</span>
            </div>
          </div>

          <div class="kpi-card card">
            <div class="kpi-icon icon-emerald">
              <PackageCheck :size="24" />
            </div>
            <div class="kpi-data">
              <span class="kpi-label">{{ t('tabProducts') }}</span>
              <span class="kpi-value">{{ stats?.total_products_count || 0 }}</span>
              <span class="kpi-hint" v-if="stats?.pending_products">({{ stats.pending_products }} на проверке)</span>
              <span class="kpi-hint" v-else>активных</span>
            </div>
          </div>
        </div>

        <!-- Визуальная диаграмма: График продаж по месяцам (Желтые стильные столбцы) -->
        <div class="chart-card card">
          <div class="chart-header">
            <div>
              <h3>{{ t('chartTitle') }}</h3>
              <p>{{ t('chartSub') }}</p>
            </div>
            <span class="badge badge-processing">KPI Analytics</span>
          </div>

          <div class="monthly-bars-chart">
            <div 
              v-for="m in (stats?.monthly_stats || [])" 
              :key="m.monthName" 
              class="bar-column">
              <div class="bar-value-label">{{ formatPrice(m.revenue) }} с.</div>
              <div class="bar-track">
                <div 
                  class="bar-fill" 
                  :style="{ height: getBarHeight(m.revenue) }">
                </div>
              </div>
              <div class="bar-name">{{ m.monthName }}</div>
              <div class="bar-orders-badge">{{ m.ordersCount }} зак.</div>
            </div>
          </div>
        </div>

        <!-- Сводная таблица по месяцам -->
        <div class="table-card card">
          <div class="table-header">
            <h3>Ҳисоботи моҳона / Ежемесячная сводка</h3>
          </div>
          <table class="data-table">
            <thead>
              <tr>
                <th>Период</th>
                <th>Заказы</th>
                <th>Продано единиц</th>
                <th>Сумма продаж</th>
                <th>Статус</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="m in (stats?.monthly_stats || [])" :key="m.monthName">
                <td><strong>{{ m.monthName }}</strong></td>
                <td>{{ m.ordersCount }}</td>
                <td>{{ m.itemsSold }} шт.</td>
                <td><strong class="text-primary">{{ formatPrice(m.revenue) }} сомони</strong></td>
                <td><span class="badge badge-approved">Подтверждено</span></td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Топ продаваемых товаров этого продавца -->
        <div class="top-products-card card" v-if="stats?.top_products?.length">
          <h3>Топ продаваемых товаров магазина</h3>
          <div class="top-list">
            <div v-for="(p, idx) in stats.top_products" :key="p.product_id" class="top-item">
              <span class="top-rank">#{{ idx + 1 }}</span>
              <img :src="p.image_url" :alt="p.product_name" class="top-img" />
              <div class="top-info">
                <strong>{{ p.product_name }}</strong>
                <span>Продано: {{ p.units_sold }} шт.</span>
              </div>
              <div class="top-earn">
                <strong>{{ formatPrice(p.total_earned) }} с.</strong>
                <small>доход</small>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- ==================== ВКЛАДКА 2: МОЛҲОИ МАН (CRUD + СТАТУС МОДЕРАЦИИ) ==================== -->
      <div v-if="activeTab === 'products'" class="tab-content">
        <div class="table-card card">
          <div class="table-header-flex">
            <div>
              <h3>{{ t('tabProducts') }}</h3>
              <p>Товары проходят модерацию перед публикацией в каталоге</p>
            </div>
            <button @click="openAddProductModal" class="btn btn-primary btn-sm">
              <Plus :size="15" />
              <span>{{ t('addNewProduct') }}</span>
            </button>
          </div>

          <div v-if="products.length === 0" class="empty-state">
            <p>У вас пока нет добавленных товаров.</p>
          </div>

          <table v-else class="data-table">
            <thead>
              <tr>
                <th>Фото</th>
                <th>Название товара</th>
                <th>Категория</th>
                <th>Цена</th>
                <th>Остаток</th>
                <th>{{ t('moderationCol') }}</th>
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="prod in products" :key="prod.id">
                <td>
                  <img :src="prod.image_url" :alt="prod.name" class="table-thumb" />
                </td>
                <td>
                  <strong>{{ prod.name }}</strong>
                  <div class="prod-slug">ID: #{{ prod.id }} | Рейтинг: ★ {{ prod.rating || 5.0 }}</div>
                </td>
                <td>{{ prod.category?.name || '—' }}</td>
                <td>
                  <strong class="text-primary">{{ formatPrice(prod.price) }} с.</strong>
                  <div v-if="prod.old_price" class="old-price">{{ formatPrice(prod.old_price) }} с.</div>
                </td>
                <td>{{ prod.stock_quantity }} шт.</td>
                <td>
                  <div class="moderation-status-cell">
                    <span :class="['badge', getProductStatusBadge(prod.status)]">
                      {{ getProductStatusLabel(prod.status) }}
                    </span>
                    <small v-if="prod.moderation_comment" class="mod-comment">
                      {{ prod.moderation_comment }}
                    </small>
                  </div>
                </td>
                <td>
                  <div class="action-buttons">
                    <button @click="openEditProductModal(prod)" class="btn-icon edit" title="Редактировать">
                      <Pencil :size="15" />
                    </button>
                    <button @click="deleteProduct(prod.id)" class="btn-icon delete" title="Удалить">
                      <Trash2 :size="15" />
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- ==================== ВКЛАДКА 3: ФАРМОИШҲОИ МАҒОЗА (ТОЛЬКО СВОЯ ЧАСТЬ!) ==================== -->
      <div v-if="activeTab === 'orders'" class="tab-content">
        <div class="seller-orders-section card">
          <div class="table-header">
            <div>
              <h3>{{ t('tabOrders') }}</h3>
              <p class="strict-isolation-badge">
                {{ t('privacyNotice') }}
              </p>
            </div>
          </div>

          <div v-if="orders.length === 0" class="empty-state">
            <p>Заказов для вашего магазина пока нет.</p>
          </div>

          <div v-else class="seller-orders-list">
            <div v-for="subOrder in orders" :key="subOrder.id" class="seller-order-item">
              <div class="so-header">
                <div>
                  <span class="so-number">Посылка №{{ subOrder.id }} (Общий заказ: #{{ subOrder.order?.order_number }})</span>
                  <span class="so-date">{{ formatDate(subOrder.created_at) }}</span>
                </div>
                <div class="so-status-action">
                  <span :class="['badge', getSubOrderStatusBadge(subOrder.status)]">
                    {{ getSubOrderStatusLabel(subOrder.status) }}
                  </span>
                  <!-- Смена статуса продавцом -->
                  <select 
                    :value="subOrder.status" 
                    @change="updateStatus(subOrder.id, $event.target.value)" 
                    class="status-select">
                    <option value="pending">Ожидает (Pending)</option>
                    <option value="processing">В сборке (Processing)</option>
                    <option value="shipped">Отправлен (Shipped)</option>
                    <option value="delivered">Доставлен (Delivered)</option>
                    <option value="cancelled">Отменен (Cancelled)</option>
                  </select>
                </div>
              </div>

              <!-- Данные доставки для курьера продавца -->
              <div class="so-shipping-info">
                <span>📍 <strong>Адрес:</strong> {{ subOrder.order?.shipping_address }}</span>
                <span>📞 <strong>Клиент:</strong> {{ subOrder.order?.recipient_name }} ({{ subOrder.order?.recipient_phone }})</span>
                <span v-if="subOrder.shipping_tracking">📦 <strong>Трек-номер:</strong> {{ subOrder.shipping_tracking }}</span>
              </div>

              <!-- Товары в заказе -->
              <div class="so-items">
                <div v-for="it in subOrder.items" :key="it.id" class="so-item-row">
                  <img :src="it.image_url" :alt="it.product_name" class="so-img" />
                  <div class="so-item-info">
                    <strong>{{ it.product_name }}</strong>
                    <span>{{ formatPrice(it.product_price) }} с. × {{ it.quantity }} шт.</span>
                  </div>
                  <strong class="so-item-total">{{ formatPrice(it.total_price) }} с.</strong>
                </div>
              </div>

              <div class="so-footer">
                <span>Итого по вашей части:</span>
                <strong>{{ formatPrice(subOrder.subtotal_amount) }} сомони</strong>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- ==================== ВКЛАДКА 4: ТАНЗИМОТИ МАҒОЗА ==================== -->
      <div v-if="activeTab === 'settings'" class="tab-content">
        <div class="settings-card card">
          <h3>{{ t('tabSettings') }}</h3>
          <p>Эта информация отображается на публичной странице вашего магазина для покупателей.</p>

          <form @submit.prevent="saveShopSettings" class="settings-form">
            <div class="form-group">
              <label>Название магазина *</label>
              <input type="text" v-model="shopForm.name" required />
            </div>

            <div class="form-group">
              <label>Описание магазина</label>
              <textarea v-model="shopForm.description" rows="3"></textarea>
            </div>

            <div class="form-grid">
              <div class="form-group">
                <label>Телефон</label>
                <input type="tel" v-model="shopForm.phone" />
              </div>

              <div class="form-group">
                <label>Адрес</label>
                <input type="text" v-model="shopForm.address" />
              </div>
            </div>

            <div class="form-group">
              <label>URL логотипа</label>
              <input type="url" v-model="shopForm.logo_url" />
            </div>

            <div class="form-group">
              <label>URL баннера магазина</label>
              <input type="url" v-model="shopForm.banner_url" />
            </div>

            <button type="submit" class="btn btn-primary" :disabled="savingSettings">
              {{ savingSettings ? 'Сохранение...' : 'Сохранить изменения' }}
            </button>
          </form>
        </div>
      </div>

      <!-- Модальное окно добавления / редактирования товара -->
      <div v-if="showProductModal" class="modal-backdrop" @click.self="showProductModal = false">
        <div class="modal-content product-form-modal">
          <div class="modal-header">
            <h3>{{ isEditing ? 'Редактировать товар' : 'Добавить новый товар' }}</h3>
            <button @click="showProductModal = false" class="close-btn"><X :size="20" /></button>
          </div>
          <p class="modal-subtitle">После добавления или редактирования товар отправляется на проверку модератору.</p>

          <form @submit.prevent="submitProductForm" class="prod-modal-form">
            <div class="form-group">
              <label>Название товара *</label>
              <input type="text" v-model="prodForm.name" required placeholder="Например: Смартфон Samsung Galaxy S24" />
            </div>

            <div class="form-grid">
              <div class="form-group">
                <label>Категория *</label>
                <select v-model.number="prodForm.category_id" required>
                  <option value="" disabled>Выберите категорию</option>
                  <option v-for="c in categories" :key="c.id" :value="c.id">{{ c.name }}</option>
                </select>
              </div>

              <div class="form-group">
                <label>Количество на складе *</label>
                <input type="number" v-model.number="prodForm.stock_quantity" min="1" required />
              </div>
            </div>

            <div class="form-grid">
              <div class="form-group">
                <label>Цена продажи (сомони) *</label>
                <input type="number" step="0.01" v-model.number="prodForm.price" required />
              </div>

              <div class="form-group">
                <label>Старая цена (для скидки)</label>
                <input type="number" step="0.01" v-model.number="prodForm.old_price" />
              </div>
            </div>

            <div class="form-group">
              <label>URL фото товара *</label>
              <input type="url" v-model="prodForm.image_url" required placeholder="https://images.unsplash.com/..." />
            </div>

            <div class="form-group">
              <label>Описание товара *</label>
              <textarea v-model="prodForm.description" rows="4" required placeholder="Характеристики, гарантия, свойства..."></textarea>
            </div>

            <button type="submit" class="btn btn-primary btn-lg" :disabled="savingProduct">
              {{ savingProduct ? 'Сохранение...' : (isEditing ? 'Сохранить изменения и отправить на проверку' : 'Добавить товар и отправить на модерацию') }}
            </button>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { 
  BarChart3, Package, ShoppingBag, Settings, Plus, 
  TrendingUp, Calendar, PackageCheck, Pencil, Trash2, X 
} from 'lucide-vue-next'
import { api } from '../api'
import { t } from '../utils/i18n'

const activeTab = ref('stats')
const stats = ref(null)
const products = ref([])
const orders = ref([])
const categories = ref([])
const shopInfo = ref(null)

const showProductModal = ref(false)
const isEditing = ref(false)
const savingProduct = ref(false)
const savingSettings = ref(false)

const prodForm = reactive({
  id: null,
  name: '',
  category_id: null,
  price: 0,
  old_price: null,
  stock_quantity: 10,
  image_url: '',
  description: ''
})

const shopForm = reactive({
  name: '',
  description: '',
  phone: '',
  address: '',
  logo_url: '',
  banner_url: ''
})

function formatPrice(val) {
  if (!val) return '0'
  return new Intl.NumberFormat('ru-RU').format(val)
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleDateString('ru-RU', { day: 'numeric', month: 'short', year: 'numeric' })
}

function getBarHeight(rev) {
  const max = Math.max(...(stats.value?.monthly_stats || []).map(m => m.revenue), 1000)
  const pct = Math.max((rev / max) * 100, 8)
  return `${pct}%`
}

function getProductStatusLabel(st) {
  switch (st) {
    case 'approved': return 'Опубликован'
    case 'pending': return 'На проверке'
    case 'rejected': return 'Отклонен'
    default: return st
  }
}

function getProductStatusBadge(st) {
  switch (st) {
    case 'approved': return 'badge-approved'
    case 'pending': return 'badge-pending'
    case 'rejected': return 'badge-rejected'
    default: return 'badge-pending'
  }
}

function getSubOrderStatusLabel(st) {
  switch (st) {
    case 'pending': return 'Ожидает'
    case 'processing': return 'В сборке'
    case 'shipped': return 'Отправлен'
    case 'delivered': return 'Доставлен'
    case 'cancelled': return 'Отменен'
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

async function loadAll() {
  try {
    const [statsRes, prodsRes, ordersRes, catsRes, shopRes] = await Promise.all([
      api.seller.getMonthlyStats(),
      api.seller.getProducts(),
      api.seller.getOrders(),
      api.categories.list(),
      api.seller.getMyShop()
    ])

    stats.value = statsRes.data
    products.value = prodsRes.data || []
    orders.value = ordersRes.data || []
    categories.value = catsRes.data || []
    shopInfo.value = shopRes.data

    if (shopRes.data) {
      shopForm.name = shopRes.data.name || ''
      shopForm.description = shopRes.data.description || ''
      shopForm.phone = shopRes.data.phone || ''
      shopForm.address = shopRes.data.address || ''
      shopForm.logo_url = shopRes.data.logo_url || ''
      shopForm.banner_url = shopRes.data.banner_url || ''
    }
  } catch (e) {
    console.error(e)
  }
}

function openAddProductModal() {
  isEditing.value = false
  prodForm.id = null
  prodForm.name = ''
  prodForm.category_id = categories.value[0]?.id || null
  prodForm.price = 100
  prodForm.old_price = null
  prodForm.stock_quantity = 15
  prodForm.image_url = 'https://images.unsplash.com/photo-1526170375885-4d8ecf77b99f?w=600'
  prodForm.description = ''
  showProductModal.value = true
}

function openEditProductModal(p) {
  isEditing.value = true
  prodForm.id = p.id
  prodForm.name = p.name
  prodForm.category_id = p.category_id
  prodForm.price = p.price
  prodForm.old_price = p.old_price
  prodForm.stock_quantity = p.stock_quantity
  prodForm.image_url = p.image_url
  prodForm.description = p.description
  showProductModal.value = true
}

async function submitProductForm() {
  savingProduct.value = true
  try {
    if (isEditing.value) {
      await api.seller.updateProduct(prodForm.id, prodForm)
      alert('Товар успешно обновлен и отправлен на модерацию.')
    } else {
      await api.seller.createProduct(prodForm)
      alert('Товар успешно добавлен и отправлен на модерацию.')
    }
    showProductModal.value = false
    await loadAll()
  } catch (err) {
    alert('Ошибка: ' + err.message)
  } finally {
    savingProduct.value = false
  }
}

async function deleteProduct(id) {
  if (confirm('Вы уверены, что хотите удалить этот товар?')) {
    try {
      await api.seller.deleteProduct(id)
      alert('Товар успешно удален.')
      await loadAll()
    } catch (e) {
      alert('Ошибка при удалении: ' + e.message)
    }
  }
}

async function updateStatus(subOrderId, newStatus) {
  try {
    await api.seller.updateOrderStatus(subOrderId, { status: newStatus })
    alert(`Статус посылки обновлен: ${newStatus}`)
    await loadAll()
  } catch (e) {
    alert('Ошибка обновления: ' + e.message)
  }
}

async function saveShopSettings() {
  savingSettings.value = true
  try {
    await api.seller.updateMyShop(shopForm)
    alert('Настройки магазина успешно сохранены!')
    await loadAll()
  } catch (e) {
    alert('Ошибка: ' + e.message)
  } finally {
    savingSettings.value = false
  }
}

onMounted(() => {
  loadAll()
})
</script>

<style scoped>
.seller-dashboard-page {
  padding: 24px 0 60px 0;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;
  margin-bottom: 24px;
}

.badge-seller-tag {
  background: #fef3c7;
  color: #b45309;
}

.page-title {
  font-size: 26px;
  font-weight: 800;
  margin: 6px 0 4px 0;
  color: var(--text-main);
}

.page-desc {
  font-size: 13px;
  color: var(--text-muted);
}

.dashboard-nav {
  display: flex;
  gap: 10px;
  background: var(--bg-card);
  padding: 6px;
  border-radius: var(--radius-lg);
  border: 1px solid var(--border);
  margin-bottom: 28px;
  overflow-x: auto;
}

.nav-tab {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 18px;
  border-radius: var(--radius-md);
  font-size: 13px;
  font-weight: 700;
  color: var(--text-secondary);
  transition: all 0.2s;
  white-space: nowrap;
}

.nav-tab:hover {
  background: var(--bg-page);
  color: #b45309;
}

.nav-tab.active {
  background: #f59e0b;
  color: #0f172a;
  box-shadow: 0 4px 12px rgba(245, 158, 11, 0.35);
}

/* KPI */
.kpi-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 28px;
}

@media (max-width: 1000px) {
  .kpi-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 550px) {
  .kpi-grid {
    grid-template-columns: 1fr;
  }
}

.kpi-card {
  padding: 20px;
  background: #ffffff;
  display: flex;
  align-items: center;
  gap: 16px;
}

.kpi-icon {
  width: 52px;
  height: 52px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.icon-yellow { background: #fef3c7; color: #b45309; }
.icon-amber { background: #fffbeb; color: #d97706; }
.icon-dark { background: #0f172a; color: #fbbf24; }
.icon-emerald { background: #ecfdf5; color: #059669; }

.kpi-data {
  display: flex;
  flex-direction: column;
}

.kpi-label {
  font-size: 12px;
  color: var(--text-muted);
  font-weight: 700;
}

.kpi-value {
  font-size: 22px;
  font-weight: 800;
  color: var(--text-main);
  line-height: 1.2;
  margin: 2px 0;
}

.kpi-value small {
  font-size: 13px;
  font-weight: 600;
}

.kpi-hint {
  font-size: 11px;
  color: var(--text-secondary);
}

/* График */
.chart-card {
  padding: 24px;
  background: #ffffff;
  margin-bottom: 28px;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.chart-header h3 {
  font-size: 16px;
  font-weight: 800;
}

.chart-header p {
  font-size: 12px;
  color: var(--text-muted);
}

.monthly-bars-chart {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 16px;
  height: 240px;
  align-items: end;
  padding: 20px 10px 0 10px;
  border-bottom: 2px solid var(--border);
}

.bar-column {
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 100%;
  justify-content: flex-end;
}

.bar-value-label {
  font-size: 11px;
  font-weight: 800;
  color: #b45309;
  margin-bottom: 6px;
}

.bar-track {
  width: 40px;
  height: 150px;
  background: #fef3c7;
  border-radius: 8px 8px 0 0;
  display: flex;
  align-items: flex-end;
  overflow: hidden;
}

.bar-fill {
  width: 100%;
  background: linear-gradient(180deg, #f59e0b 0%, #d97706 100%);
  border-radius: 6px 6px 0 0;
  transition: height 0.4s ease;
}

.bar-name {
  font-size: 11px;
  font-weight: 700;
  color: var(--text-main);
  margin-top: 8px;
  text-align: center;
}

.bar-orders-badge {
  font-size: 10px;
  color: var(--text-muted);
}

/* Таблицы */
.table-card {
  padding: 24px;
  background: #ffffff;
  margin-bottom: 28px;
}

.table-header, .table-header-flex {
  margin-bottom: 18px;
}

.table-header-flex {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.table-header h3, .table-header-flex h3 {
  font-size: 16px;
  font-weight: 800;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}

.data-table th {
  text-align: left;
  padding: 12px 14px;
  background: var(--bg-page);
  color: var(--text-secondary);
  font-weight: 700;
  border-bottom: 1px solid var(--border);
}

.data-table td {
  padding: 14px;
  border-bottom: 1px solid var(--border-light);
}

.table-thumb {
  width: 44px;
  height: 44px;
  border-radius: 6px;
  object-fit: cover;
}

.prod-slug {
  font-size: 11px;
  color: var(--text-muted);
}

.moderation-status-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.mod-comment {
  font-size: 11px;
  color: var(--text-muted);
}

.action-buttons {
  display: flex;
  gap: 8px;
}

.btn-icon {
  padding: 6px;
  border-radius: 6px;
  border: 1px solid var(--border);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s;
}

.btn-icon.edit:hover {
  background: #fef3c7;
  color: #b45309;
  border-color: #f59e0b;
}

.btn-icon.delete:hover {
  background: #fee2e2;
  color: var(--danger);
  border-color: var(--danger);
}

/* Заказы продавца */
.seller-orders-section {
  padding: 24px;
  background: #ffffff;
}

.strict-isolation-badge {
  font-size: 12px;
  color: #92400e;
  background: #fef3c7;
  display: inline-block;
  padding: 4px 10px;
  border-radius: 6px;
  margin-top: 6px;
  font-weight: 600;
}

.seller-orders-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-top: 18px;
}

.seller-order-item {
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  overflow: hidden;
  background: #ffffff;
}

.so-header {
  padding: 12px 16px;
  background: #f8fafc;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
  border-bottom: 1px solid var(--border-light);
}

.so-number {
  font-size: 14px;
  font-weight: 800;
  display: block;
}

.so-date {
  font-size: 11px;
  color: var(--text-muted);
}

.so-status-action {
  display: flex;
  align-items: center;
  gap: 10px;
}

.status-select {
  padding: 5px 10px;
  border: 1px solid var(--border);
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  background: #ffffff;
  outline: none;
}

.so-shipping-info {
  padding: 10px 16px;
  background: #fafafa;
  font-size: 12px;
  color: var(--text-secondary);
  display: flex;
  gap: 18px;
  flex-wrap: wrap;
  border-bottom: 1px solid var(--border-light);
}

.so-items {
  padding: 0 16px;
}

.so-item-row {
  display: grid;
  grid-template-columns: 44px 1fr auto;
  gap: 12px;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid var(--border-light);
}

.so-item-row:last-child {
  border-bottom: none;
}

.so-img {
  width: 44px;
  height: 44px;
  border-radius: 6px;
  object-fit: cover;
}

.so-item-info {
  display: flex;
  flex-direction: column;
  font-size: 12px;
}

.so-item-total {
  font-size: 13px;
  font-weight: 700;
  color: var(--text-main);
}

.so-footer {
  padding: 10px 16px;
  background: #fafafa;
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  font-size: 13px;
}

/* Топ товаров */
.top-products-card {
  padding: 24px;
  background: #ffffff;
  margin-bottom: 28px;
}

.top-products-card h3 {
  font-size: 16px;
  font-weight: 800;
  margin-bottom: 16px;
}

.top-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.top-item {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 10px 14px;
  background: var(--bg-page);
  border-radius: var(--radius-md);
}

.top-rank {
  font-size: 14px;
  font-weight: 800;
  color: #b45309;
  width: 24px;
}

.top-img {
  width: 42px;
  height: 42px;
  border-radius: 6px;
  object-fit: cover;
}

.top-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  font-size: 13px;
}

.top-info span {
  font-size: 11px;
  color: var(--text-muted);
}

.top-earn {
  text-align: right;
  display: flex;
  flex-direction: column;
}

.top-earn strong {
  font-size: 14px;
  color: #b45309;
}

.top-earn small {
  font-size: 10px;
  color: var(--text-muted);
}

/* Настройки */
.settings-card {
  padding: 28px;
  background: #ffffff;
  max-width: 680px;
}

.settings-card h3 {
  font-size: 18px;
  font-weight: 800;
  margin-bottom: 4px;
}

.settings-card p {
  font-size: 13px;
  color: var(--text-muted);
  margin-bottom: 22px;
}

.settings-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group label {
  font-size: 12px;
  font-weight: 700;
}

.form-group input, .form-group textarea, .form-group select {
  padding: 10px 12px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  font-size: 13px;
  outline: none;
}

.form-group input:focus, .form-group textarea:focus, .form-group select:focus {
  border-color: #f59e0b;
}

/* Модалка товара */
.product-form-modal {
  max-width: 580px;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.modal-header h3 {
  font-size: 18px;
  font-weight: 800;
}

.modal-subtitle {
  font-size: 12px;
  color: var(--text-muted);
  margin-bottom: 18px;
}

.prod-modal-form {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.empty-state {
  padding: 40px;
  text-align: center;
  color: var(--text-muted);
}
</style>
