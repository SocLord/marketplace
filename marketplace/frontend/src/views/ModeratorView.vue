<template>
  <div class="moderator-page">
    <div class="container">
      <div class="moderator-header">
        <div>
          <span class="badge badge-processing">{{ currentLang === 'ru' ? 'Панель проверки' : 'Панели санҷиш' }}</span>
          <h1 class="page-title">{{ t('modTitle') }}</h1>
          <p class="page-desc">{{ t('modSub') }}</p>
        </div>
      </div>

      <!-- Вкладки навигации модератора -->
      <div class="moderator-nav">
        <button 
          :class="['nav-tab', activeTab === 'pending' && 'active']"
          @click="activeTab = 'pending'">
          <ShieldAlert :size="18" />
          <span>{{ t('tabPending') }} ({{ pendingProducts.length }})</span>
        </button>

        <button 
          :class="['nav-tab', activeTab === 'all' && 'active']"
          @click="activeTab = 'all'">
          <Package :size="18" />
          <span>{{ t('tabAllProds') }}</span>
        </button>

        <button 
          :class="['nav-tab', activeTab === 'categories' && 'active']"
          @click="activeTab = 'categories'">
          <Grid :size="18" />
          <span>{{ t('tabCategories') }}</span>
        </button>

        <button 
          :class="['nav-tab', activeTab === 'stats' && 'active']"
          @click="activeTab = 'stats'">
          <BarChart3 :size="18" />
          <span>{{ t('tabPlatformStats') }}</span>
        </button>
      </div>

      <!-- ==================== ВКЛАДКА 1: МОЛҲО ДАР НАВБАТИ САНҶИШ (PENDING) ==================== -->
      <div v-if="activeTab === 'pending'" class="tab-content">
        <div v-if="pendingProducts.length === 0" class="empty-pending card">
          <CheckCircle2 :size="52" color="#10b981" />
          <h3>{{ currentLang === 'ru' ? 'Все товары проверены!' : 'Ҳамаи молҳо санҷида шудаанд!' }}</h3>
          <p>{{ currentLang === 'ru' ? 'Нет товаров в очереди на модерацию. Все поступившие товары обработаны.' : 'Дар навбат ягон моли нав барои тасдиқ мавҷуд нест. Ҳамаи молҳои воридшуда коркард шудаанд.' }}</p>
        </div>

        <div v-else class="pending-grid">
          <div v-for="prod in pendingProducts" :key="prod.id" class="pending-card card">
            <div class="pending-img-col">
              <img :src="prod.image_url" :alt="prod.name" class="pending-img" />
              <span class="badge badge-pending status-pill">{{ currentLang === 'ru' ? 'На проверке' : 'Дар санҷиш' }}</span>
            </div>

            <div class="pending-info-col">
              <div class="pending-meta">
                <span class="shop-badge">🏪 {{ prod.shop?.name || (currentLang === 'ru' ? `Магазин #${prod.shop_id}` : `Мағозаи #${prod.shop_id}`) }}</span>
                <span class="cat-badge">📁 {{ prod.category?.name }}</span>
                <span class="time-badge">🕒 {{ formatDate(prod.created_at) }}</span>
              </div>

              <h2 class="pending-title">{{ prod.name }}</h2>

              <div class="pending-pricing">
                <span class="price-val">{{ formatPrice(prod.price) }} {{ t('currency') }}</span>
                <span v-if="prod.old_price" class="old-price-val">{{ formatPrice(prod.old_price) }} {{ currentLang === 'ru' ? 'сом.' : 'с.' }}</span>
                <span class="stock-val">{{ currentLang === 'ru' ? 'Остаток' : 'Бақия' }}: {{ prod.stock_quantity }} {{ currentLang === 'ru' ? 'шт.' : 'дона' }}</span>
              </div>

              <div class="pending-desc-box">
                <strong>{{ currentLang === 'ru' ? 'Описание от продавца:' : 'Тавсифи пешниҳодшудаи фурӯшанда:' }}</strong>
                <p>{{ prod.description }}</p>
              </div>

              <!-- Блок принятия решения модератором -->
              <div class="decision-box">
                <div class="comment-input-wrap">
                  <input 
                    type="text" 
                    v-model="reviewComments[prod.id]" 
                    :placeholder="currentLang === 'ru' ? 'Причина отклонения или комментарий (необязательно)...' : 'Сабаби радкунӣ ё эзоҳ барои фурӯшанда (ихтиёрӣ)...'" 
                  />
                </div>

                <div class="decision-buttons">
                  <button 
                    @click="handleReview(prod.id, 'approve')" 
                    class="btn btn-primary"
                    :disabled="submittingId === prod.id">
                    <Check :size="16" />
                    <span>{{ t('approveBtn') }}</span>
                  </button>

                  <button 
                    @click="handleReview(prod.id, 'reject')" 
                    class="btn btn-danger"
                    :disabled="submittingId === prod.id">
                    <X :size="16" />
                    <span>{{ t('rejectBtn') }}</span>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- ==================== ВКЛАДКА 2: ВСЕ ТОВАРЫ ==================== -->
      <div v-if="activeTab === 'all'" class="tab-content">
        <div class="table-card card">
          <div class="filter-bar-flex">
            <h3>{{ currentLang === 'ru' ? 'Все товары на платформе' : 'Ҳамаи молҳои воридшуда ба пойгоҳи додаҳо' }}</h3>
            <div class="status-tabs">
              <button 
                :class="['filter-btn', statusFilter === 'all' && 'active']"
                @click="filterAllProducts('all')">
                {{ currentLang === 'ru' ? 'Все' : 'Ҳамааш' }}
              </button>
              <button 
                :class="['filter-btn', statusFilter === 'approved' && 'active']"
                @click="filterAllProducts('approved')">
                {{ currentLang === 'ru' ? 'Опубликованы (Approved)' : 'Нашршуда (Approved)' }}
              </button>
              <button 
                :class="['filter-btn', statusFilter === 'pending' && 'active']"
                @click="filterAllProducts('pending')">
                {{ currentLang === 'ru' ? 'На проверке (Pending)' : 'Дар санҷиш (Pending)' }}
              </button>
              <button 
                :class="['filter-btn', statusFilter === 'rejected' && 'active']"
                @click="filterAllProducts('rejected')">
                {{ currentLang === 'ru' ? 'Отклонены (Rejected)' : 'Радшуда (Rejected)' }}
              </button>
            </div>
          </div>

          <table class="data-table">
            <thead>
              <tr>
                <th>{{ currentLang === 'ru' ? 'Фото' : 'Акс' }}</th>
                <th>{{ currentLang === 'ru' ? 'Название' : 'Номи мол' }}</th>
                <th>{{ currentLang === 'ru' ? 'Магазин' : 'Мағоза' }}</th>
                <th>{{ currentLang === 'ru' ? 'Категория' : 'Категория' }}</th>
                <th>{{ currentLang === 'ru' ? 'Цена' : 'Нарх' }}</th>
                <th>{{ currentLang === 'ru' ? 'Статус' : 'Ҳолат' }}</th>
                <th>{{ currentLang === 'ru' ? 'Комментарий модератора' : 'Шарҳи модератор' }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="p in allProducts" :key="p.id">
                <td><img :src="p.image_url" :alt="p.name" class="table-thumb" /></td>
                <td>
                  <strong>{{ p.name }}</strong>
                  <div class="prod-slug">ID: #{{ p.id }}</div>
                </td>
                <td>{{ p.shop?.name || '—' }}</td>
                <td>{{ p.category?.name || '—' }}</td>
                <td><strong>{{ formatPrice(p.price) }} {{ currentLang === 'ru' ? 'сом.' : 'с.' }}</strong></td>
                <td>
                  <span :class="['badge', getProductStatusBadge(p.status)]">
                    {{ getProductStatusLabel(p.status) }}
                  </span>
                </td>
                <td><small>{{ p.moderation_comment || '—' }}</small></td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- ==================== ВКЛАДКА 3: КАТЕГОРИИ ==================== -->
      <div v-if="activeTab === 'categories'" class="tab-content">
        <div class="categories-admin-grid">
          <!-- Форма добавления категории -->
          <div class="cat-form-card card">
            <h3>{{ currentLang === 'ru' ? 'Добавить категорию' : 'Иловаи категорияи нав' }}</h3>
            <form @submit.prevent="submitCategory" class="cat-form">
              <div class="form-group">
                <label>{{ currentLang === 'ru' ? 'Название категории *' : 'Номи категория *' }}</label>
                <input type="text" v-model="newCat.name" required :placeholder="currentLang === 'ru' ? 'Например: Игрушки и дети' : 'Масалан: Бозичаҳо ва кӯдакон'" />
              </div>

              <div class="form-group">
                <label>{{ currentLang === 'ru' ? 'Краткое описание' : 'Тавсифи кӯтоҳ' }}</label>
                <input type="text" v-model="newCat.description" :placeholder="currentLang === 'ru' ? 'Детские товары...' : 'Молҳои бачагона...'" />
              </div>

              <div class="form-group">
                <label>{{ currentLang === 'ru' ? 'Ссылка на изображение (URL)' : 'Нишонии расм (Image URL)' }}</label>
                <input type="url" v-model="newCat.image_url" placeholder="https://images.unsplash.com/..." />
              </div>

              <button type="submit" class="btn btn-primary">
                <Plus :size="16" />
                <span>{{ currentLang === 'ru' ? 'Добавить категорию' : 'Илова кардани категория' }}</span>
              </button>
            </form>
          </div>

          <!-- Список существующих категорий -->
          <div class="cat-list-card card">
            <h3>{{ currentLang === 'ru' ? 'Активные категории каталога' : 'Категорияҳои фаъоли маркетплейс' }}</h3>
            <div class="cats-table-wrap">
              <table class="data-table">
                <thead>
                  <tr>
                    <th>{{ currentLang === 'ru' ? 'Фото' : 'Расм' }}</th>
                    <th>{{ currentLang === 'ru' ? 'Название' : 'Ном' }}</th>
                    <th>{{ currentLang === 'ru' ? 'Кол-во товаров' : 'Шумораи молҳо' }}</th>
                    <th>{{ currentLang === 'ru' ? 'Действия' : 'Амалҳо' }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="cat in categories" :key="cat.id">
                    <td><img :src="cat.image_url" :alt="cat.name" class="cat-thumb" /></td>
                    <td>
                      <strong>{{ cat.name }}</strong>
                      <div class="cat-slug">/{{ cat.slug }}</div>
                    </td>
                    <td>{{ cat.products_count || 0 }} {{ currentLang === 'ru' ? 'тов.' : 'мол' }}</td>
                    <td>
                      <button @click="deleteCat(cat.id)" class="btn-icon delete" :title="currentLang === 'ru' ? 'Удалить' : 'Нест кардан'">
                        <Trash2 :size="14" />
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

      <!-- ==================== ВКЛАДКА 4: СТАТИСТИКА ПЛАТФОРМЫ ==================== -->
      <div v-if="activeTab === 'stats'" class="tab-content">
        <div class="platform-stats-grid">
          <div class="stat-card card">
            <span class="stat-num">{{ modStats?.total_turnover ? formatPrice(modStats.total_turnover) : 0 }} {{ currentLang === 'ru' ? 'сом.' : 'с.' }}</span>
            <span class="stat-label">{{ currentLang === 'ru' ? 'Общий финансовый оборот маркетплейса' : 'Гардиши умумии молиявии маркетплейс' }}</span>
          </div>

          <div class="stat-card card">
            <span class="stat-num">{{ modStats?.total_orders || 0 }}</span>
            <span class="stat-label">{{ currentLang === 'ru' ? 'Всего заказов' : 'Ҳамагӣ фармоишҳо' }}</span>
          </div>

          <div class="stat-card card">
            <span class="stat-num">{{ modStats?.total_products || 0 }}</span>
            <span class="stat-label">{{ currentLang === 'ru' ? 'Всего зарегистрированных товаров' : 'Ҳамагӣ молҳои ба қайд гирифташуда' }}</span>
          </div>

          <div class="stat-card card">
            <span class="stat-num">{{ modStats?.total_shops || 0 }}</span>
            <span class="stat-label">{{ currentLang === 'ru' ? 'Активных магазинов' : 'Мағозаҳои фаъол' }}</span>
          </div>

          <div class="stat-card card">
            <span class="stat-num">{{ modStats?.total_sellers || 0 }}</span>
            <span class="stat-label">{{ currentLang === 'ru' ? 'Зарегистрированных продавцов' : 'Фурӯшандагони ба қайд гирифташуда' }}</span>
          </div>

          <div class="stat-card card">
            <span class="stat-num">{{ modStats?.total_buyers || 0 }}</span>
            <span class="stat-label">{{ currentLang === 'ru' ? 'Активных покупателей' : 'Харидорони фаъол' }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { 
  ShieldAlert, Package, Grid, BarChart3, 
  CheckCircle2, Check, X, Plus, Trash2 
} from 'lucide-vue-next'
import { api } from '../api'
import { t, currentLang } from '../utils/i18n'

const activeTab = ref('pending')
const pendingProducts = ref([])
const allProducts = ref([])
const categories = ref([])
const modStats = ref(null)

const statusFilter = ref('all')
const reviewComments = reactive({})
const submittingId = ref(null)

const newCat = reactive({
  name: '',
  description: '',
  image_url: 'https://images.unsplash.com/photo-1513694203232-719a280e022f?w=400'
})

function formatPrice(val) {
  if (!val) return '0'
  return new Intl.NumberFormat('ru-RU').format(val)
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  const locale = currentLang.value === 'ru' ? 'ru-RU' : 'tg-TJ'
  return d.toLocaleDateString(locale, { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })
}

function getProductStatusLabel(st) {
  if (currentLang.value === 'ru') {
    switch (st) {
      case 'approved': return 'Опубликован'
      case 'pending': return 'На проверке'
      case 'rejected': return 'Отклонён'
      default: return st
    }
  }
  switch (st) {
    case 'approved': return 'Нашр шуд'
    case 'pending': return 'Дар санҷиш'
    case 'rejected': return 'Рад шуд'
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

async function loadData() {
  try {
    const [pendRes, catsRes, statsRes] = await Promise.all([
      api.moderator.getPendingProducts(),
      api.categories.list(),
      api.moderator.getStats()
    ])
    pendingProducts.value = pendRes.data || []
    categories.value = catsRes.data || []
    modStats.value = statsRes.data

    await filterAllProducts(statusFilter.value)
  } catch (e) {
    console.error(e)
  }
}

async function filterAllProducts(status) {
  statusFilter.value = status
  try {
    const res = await api.moderator.getAllProducts(status)
    allProducts.value = res.data || []
  } catch (e) {
    console.error(e)
  }
}

async function handleReview(productId, action) {
  submittingId.value = productId
  const comment = reviewComments[productId] || ''
  try {
    await api.moderator.reviewProduct(productId, { action, comment })
    const msg = currentLang.value === 'ru'
      ? (action === 'approve' ? 'Товар одобрен и опубликован!' : 'Товар отклонён.')
      : (action === 'approve' ? 'Мол тасдиқ шуд ва нашр гардид!' : 'Мол рад карда шуд.')
    alert(msg)
    delete reviewComments[productId]
    await loadData()
  } catch (err) {
    alert((currentLang.value === 'ru' ? 'Ошибка: ' : 'Хатогӣ: ') + err.message)
  } finally {
    submittingId.value = null
  }
}

async function submitCategory() {
  try {
    await api.moderator.createCategory(newCat)
    alert(currentLang.value === 'ru' ? 'Категория добавлена!' : 'Категория илова шуд!')
    newCat.name = ''
    newCat.description = ''
    await loadData()
  } catch (e) {
    alert((currentLang.value === 'ru' ? 'Ошибка: ' : 'Хатогӣ: ') + e.message)
  }
}

async function deleteCat(id) {
  const msg = currentLang.value === 'ru' ? 'Удалить эту категорию?' : 'Оё мехоҳед ин категорияро нест кунед?'
  if (confirm(msg)) {
    try {
      await api.moderator.deleteCategory(id)
      await loadData()
    } catch (e) {
      alert((currentLang.value === 'ru' ? 'Ошибка: ' : 'Хатогӣ: ') + e.message)
    }
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.moderator-page {
  padding: 24px 0 60px 0;
}

.moderator-header {
  margin-bottom: 24px;
}

.page-title {
  font-size: 26px;
  font-weight: 800;
  margin: 6px 0 4px 0;
}

.page-desc {
  font-size: 13px;
  color: var(--text-muted);
}

.moderator-nav {
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
  color: var(--primary);
}

.nav-tab.active {
  background: #0f172a;
  color: #fbbf24;
  box-shadow: 0 4px 10px rgba(15, 23, 42, 0.25);
}

/* Карточки на проверке */
.empty-pending {
  padding: 60px 20px;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.pending-grid {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.pending-card {
  padding: 24px;
  background: #ffffff;
  display: grid;
  grid-template-columns: 240px 1fr;
  gap: 24px;
}

@media (max-width: 768px) {
  .pending-card {
    grid-template-columns: 1fr;
  }
}

.pending-img-col {
  position: relative;
  width: 100%;
  height: 240px;
  border-radius: var(--radius-md);
  overflow: hidden;
  background: #f1f5f9;
}

.pending-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.status-pill {
  position: absolute;
  top: 10px;
  left: 10px;
}

.pending-info-col {
  display: flex;
  flex-direction: column;
}

.pending-meta {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: var(--text-secondary);
  margin-bottom: 8px;
  flex-wrap: wrap;
}

.pending-title {
  font-size: 20px;
  font-weight: 800;
  color: var(--text-main);
  margin-bottom: 12px;
}

.pending-pricing {
  display: flex;
  align-items: baseline;
  gap: 14px;
  margin-bottom: 14px;
}

.price-val {
  font-size: 22px;
  font-weight: 800;
  color: #0f172a;
}

.old-price-val {
  font-size: 14px;
  color: var(--text-muted);
  text-decoration: line-through;
}

.stock-val {
  font-size: 12px;
  color: var(--text-secondary);
  background: var(--bg-page);
  padding: 3px 8px;
  border-radius: 4px;
}

.pending-desc-box {
  background: var(--bg-page);
  padding: 12px 16px;
  border-radius: var(--radius-md);
  font-size: 13px;
  color: var(--text-secondary);
  margin-bottom: 20px;
}

.pending-desc-box strong {
  display: block;
  font-size: 12px;
  margin-bottom: 4px;
  color: var(--text-main);
}

.decision-box {
  margin-top: auto;
  display: flex;
  flex-direction: column;
  gap: 12px;
  background: #fffbeb;
  padding: 16px;
  border-radius: var(--radius-md);
  border: 1px solid rgba(245, 158, 11, 0.2);
}

.comment-input-wrap input {
  width: 100%;
  padding: 9px 12px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  font-size: 13px;
}

.decision-buttons {
  display: flex;
  gap: 12px;
}

/* Все товары */
.table-card {
  padding: 24px;
  background: #ffffff;
}

.filter-bar-flex {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  flex-wrap: wrap;
  gap: 12px;
}

.filter-bar-flex h3 {
  font-size: 16px;
  font-weight: 800;
}

.status-tabs {
  display: flex;
  gap: 8px;
}

.filter-btn {
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  border: 1px solid var(--border);
  background: var(--bg-page);
  color: var(--text-secondary);
}

.filter-btn.active {
  background: #0f172a;
  color: #fbbf24;
  border-color: #0f172a;
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

/* Категории */
.categories-admin-grid {
  display: grid;
  grid-template-columns: 360px 1fr;
  gap: 24px;
}

@media (max-width: 850px) {
  .categories-admin-grid {
    grid-template-columns: 1fr;
  }
}

.cat-form-card, .cat-list-card {
  padding: 24px;
  background: #ffffff;
}

.cat-form-card h3, .cat-list-card h3 {
  font-size: 16px;
  font-weight: 800;
  margin-bottom: 18px;
}

.cat-form {
  display: flex;
  flex-direction: column;
  gap: 14px;
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

.form-group input {
  padding: 9px 12px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  font-size: 13px;
}

.cat-thumb {
  width: 36px;
  height: 36px;
  border-radius: 6px;
  object-fit: cover;
}

.cat-slug {
  font-size: 11px;
  color: var(--text-muted);
}

.btn-icon.delete {
  padding: 6px;
  border-radius: 6px;
  color: var(--danger);
  border: 1px solid #fee2e2;
  background: #fef2f2;
}

/* Статистика платформы */
.platform-stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
}

@media (max-width: 900px) {
  .platform-stats-grid {
    grid-template-columns: 1fr;
  }
}

.stat-card {
  padding: 30px;
  text-align: center;
  background: #ffffff;
}

.stat-num {
  font-size: 32px;
  font-weight: 800;
  color: #d97706;
  display: block;
  margin-bottom: 6px;
}

.stat-label {
  font-size: 13px;
  color: var(--text-muted);
  font-weight: 600;
}
</style>
