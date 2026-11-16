<template>
  <div class="catalog-page">
    <div class="container">
      <!-- Хлебные крошки и заголовок -->
      <div class="page-header">
        <h1 class="page-title">
          {{ selectedCategory ? selectedCategory.name : t('allProducts') }}
        </h1>
        <span class="results-count">{{ t('foundProductsCount') }} {{ products.length }}</span>
      </div>

      <div class="catalog-layout">
        <!-- Боковая панель фильтрации (Sidebar Filters) -->
        <aside class="filters-sidebar card">
          <div class="filters-header">
            <h3><SlidersHorizontal :size="18" /> {{ t('filtersTitle') }}</h3>
            <button @click="resetFilters" class="reset-btn">{{ t('clearFilters') }}</button>
          </div>

          <!-- Поиск по названию -->
          <div class="filter-group">
            <label class="filter-title">{{ t('textSearch') }}</label>
            <div class="filter-input-wrap">
              <input 
                type="text" 
                v-model="filters.q" 
                @keyup.enter="applyFilters"
                :placeholder="t('textSearch') + '...'" 
              />
              <button @click="applyFilters" class="search-mini-btn">
                <Search :size="14" />
              </button>
            </div>
          </div>

          <!-- Фильтр по категории -->
          <div class="filter-group">
            <label class="filter-title">{{ t('categoriesFilter') }}</label>
            <div class="categories-list">
              <label 
                class="category-radio" 
                :class="!filters.category && 'active'">
                <input 
                  type="radio" 
                  name="cat" 
                  value="" 
                  v-model="filters.category" 
                  @change="applyFilters" 
                />
                <span>{{ t('allSections') }}</span>
              </label>

              <label 
                v-for="cat in categories" 
                :key="cat.id" 
                class="category-radio"
                :class="filters.category === cat.slug && 'active'">
                <input 
                  type="radio" 
                  name="cat" 
                  :value="cat.slug" 
                  v-model="filters.category" 
                  @change="applyFilters" 
                />
                <span>{{ cat.name }}</span>
                <small>({{ cat.products_count || 0 }})</small>
              </label>
            </div>
          </div>

          <!-- Фильтр по продавцам / магазинам -->
          <div class="filter-group">
            <label class="filter-title">{{ t('sellersFilter') }}</label>
            <div class="shops-list">
              <label 
                class="category-radio" 
                :class="!filters.shop && 'active'">
                <input 
                  type="radio" 
                  name="shop" 
                  value="" 
                  v-model="filters.shop" 
                  @change="applyFilters" 
                />
                <span>{{ t('allShops') }}</span>
              </label>

              <label 
                v-for="sh in shops" 
                :key="sh.id" 
                class="category-radio"
                :class="filters.shop === sh.slug && 'active'">
                <input 
                  type="radio" 
                  name="shop" 
                  :value="sh.slug" 
                  v-model="filters.shop" 
                  @change="applyFilters" 
                />
                <span>{{ sh.name }}</span>
                <Star :size="12" class="shop-star" />
              </label>
            </div>
          </div>

          <!-- Фильтр по цене (диапазон цен в сомони) -->
          <div class="filter-group">
            <label class="filter-title">{{ t('priceFilter') }}</label>
            <div class="price-inputs">
              <input 
                type="number" 
                v-model.number="filters.min_price" 
                placeholder="от" 
                min="0" 
              />
              <span class="price-dash">—</span>
              <input 
                type="number" 
                v-model.number="filters.max_price" 
                placeholder="до" 
                min="0" 
              />
            </div>
            <button @click="applyFilters" class="btn btn-primary btn-sm apply-price-btn">
              {{ t('applyPrice') }}
            </button>
          </div>
        </aside>

        <!-- Основная область с товарами -->
        <main class="catalog-content">
          <!-- Сортировка и активные фильтры -->
          <div class="catalog-toolbar card">
            <div class="active-chips">
              <span v-if="filters.q" class="chip">
                {{ filters.q }}
                <X :size="12" @click="filters.q = ''; applyFilters()" />
              </span>
              <span v-if="filters.category" class="chip">
                {{ getCategoryName(filters.category) }}
                <X :size="12" @click="filters.category = ''; applyFilters()" />
              </span>
              <span v-if="filters.shop" class="chip">
                {{ getShopName(filters.shop) }}
                <X :size="12" @click="filters.shop = ''; applyFilters()" />
              </span>
              <span v-if="filters.min_price || filters.max_price" class="chip">
                {{ filters.min_price || 0 }} - {{ filters.max_price || '∞' }} с.
                <X :size="12" @click="filters.min_price = null; filters.max_price = null; applyFilters()" />
              </span>
            </div>

            <!-- Сортировка -->
            <div class="sort-selector">
              <label>{{ t('sortLabel') }}</label>
              <select v-model="filters.sort" @change="applyFilters">
                <option value="newest">{{ t('sortNewest') }}</option>
                <option value="popular">{{ t('sortPopular') }}</option>
                <option value="rating">{{ t('sortRating') }}</option>
                <option value="price_asc">{{ t('sortPriceAsc') }}</option>
                <option value="price_desc">{{ t('sortPriceDesc') }}</option>
              </select>
            </div>
          </div>

          <!-- Состояние загрузки -->
          <div v-if="loading" class="loading-box card">
            <span>Интизор шавед... / Загрузка товаров...</span>
          </div>

          <!-- Пустой результат -->
          <div v-else-if="products.length === 0" class="empty-box card">
            <PackageSearch :size="48" color="#f59e0b" />
            <h3>{{ t('noProductsFound') }}</h3>
            <button @click="resetFilters" class="btn btn-primary btn-sm">{{ t('clearFilters') }}</button>
          </div>

          <!-- Сетка товаров -->
          <div v-else class="catalog-grid">
            <ProductCard 
              v-for="product in products" 
              :key="product.id" 
              :product="product" 
            />
          </div>
        </main>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { SlidersHorizontal, Search, X, Star, PackageSearch } from 'lucide-vue-next'
import { api } from '../api'
import { t } from '../utils/i18n'
import ProductCard from '../components/ProductCard.vue'

const route = useRoute()
const router = useRouter()

const loading = ref(true)
const products = ref([])
const categories = ref([])
const shops = ref([])

const filters = reactive({
  q: route.query.q || '',
  category: route.query.category || '',
  shop: route.query.shop || '',
  min_price: route.query.min_price ? Number(route.query.min_price) : null,
  max_price: route.query.max_price ? Number(route.query.max_price) : null,
  sort: route.query.sort || 'newest'
})

const selectedCategory = computed(() => {
  if (!filters.category) return null
  return categories.value.find(c => c.slug === filters.category)
})

function getCategoryName(slug) {
  const cat = categories.value.find(c => c.slug === slug)
  return cat ? cat.name : slug
}

function getShopName(slug) {
  const sh = shops.value.find(s => s.slug === slug)
  return sh ? sh.name : slug
}

async function loadInitialData() {
  try {
    const [catsRes, shopsRes] = await Promise.all([
      api.categories.list(),
      api.shops.list()
    ])
    categories.value = catsRes.data || []
    shops.value = shopsRes.data || []
  } catch (e) {
    console.error(e)
  }
}

async function fetchProducts() {
  loading.value = true
  try {
    const params = {
      q: filters.q,
      category: filters.category,
      shop: filters.shop,
      min_price: filters.min_price,
      max_price: filters.max_price,
      sort: filters.sort
    }
    const res = await api.products.list(params)
    products.value = res.data || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function applyFilters() {
  const q = {}
  if (filters.q) q.q = filters.q
  if (filters.category) q.category = filters.category
  if (filters.shop) q.shop = filters.shop
  if (filters.min_price) q.min_price = filters.min_price
  if (filters.max_price) q.max_price = filters.max_price
  if (filters.sort) q.sort = filters.sort

  router.push({ path: '/catalog', query: q })
  fetchProducts()
}

function resetFilters() {
  filters.q = ''
  filters.category = ''
  filters.shop = ''
  filters.min_price = null
  filters.max_price = null
  filters.sort = 'newest'
  router.push('/catalog')
  fetchProducts()
}

watch(() => route.query, (newQ) => {
  filters.q = newQ.q || ''
  filters.category = newQ.category || ''
  filters.shop = newQ.shop || ''
  filters.min_price = newQ.min_price ? Number(newQ.min_price) : null
  filters.max_price = newQ.max_price ? Number(newQ.max_price) : null
  filters.sort = newQ.sort || 'newest'
  fetchProducts()
})

onMounted(async () => {
  await loadInitialData()
  fetchProducts()
})
</script>

<style scoped>
.catalog-page {
  padding: 24px 0 60px 0;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  margin-bottom: 20px;
}

.page-title {
  font-size: 26px;
  font-weight: 800;
  color: var(--text-main);
}

.results-count {
  font-size: 13px;
  color: var(--text-muted);
}

.catalog-layout {
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: 24px;
  align-items: start;
}

@media (max-width: 900px) {
  .catalog-layout {
    grid-template-columns: 1fr;
  }
}

/* Боковая панель фильтрации */
.filters-sidebar {
  padding: 20px;
  background: #ffffff;
}

.filters-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 18px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--border);
}

.filters-header h3 {
  font-size: 15px;
  font-weight: 800;
  display: flex;
  align-items: center;
  gap: 6px;
  color: #0f172a;
}

.reset-btn {
  font-size: 12px;
  color: #d97706;
  font-weight: 700;
}

.reset-btn:hover {
  text-decoration: underline;
}

.filter-group {
  margin-bottom: 22px;
}

.filter-title {
  font-size: 13px;
  font-weight: 800;
  color: var(--text-main);
  margin-bottom: 10px;
  display: block;
}

.filter-input-wrap {
  position: relative;
  display: flex;
}

.filter-input-wrap input {
  width: 100%;
  padding: 8px 32px 8px 10px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  font-size: 13px;
}

.filter-input-wrap input:focus {
  border-color: #f59e0b;
}

.search-mini-btn {
  position: absolute;
  right: 6px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-muted);
}

.categories-list, .shops-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
  max-height: 200px;
  overflow-y: auto;
}

.category-radio {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 8px;
  border-radius: 6px;
  font-size: 13px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.15s;
}

.category-radio:hover {
  background: var(--bg-page);
  color: #b45309;
}

.category-radio.active {
  background: #fef3c7;
  color: #92400e;
  font-weight: 700;
}

.category-radio small {
  margin-left: auto;
  font-size: 11px;
  color: var(--text-muted);
}

.shop-star {
  margin-left: auto;
  color: #f59e0b;
  fill: #f59e0b;
}

.price-inputs {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.price-inputs input {
  width: 100%;
  padding: 8px 10px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  font-size: 13px;
}

.price-inputs input:focus {
  border-color: #f59e0b;
}

.price-dash {
  color: var(--text-muted);
}

.apply-price-btn {
  width: 100%;
}

/* Панель инструментов каталога */
.catalog-toolbar {
  padding: 12px 18px;
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.active-chips {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: #fef3c7;
  color: #92400e;
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 700;
}

.chip svg {
  cursor: pointer;
}

.sort-selector {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  font-weight: 700;
}

.sort-selector select {
  padding: 6px 12px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  font-size: 13px;
  background: #ffffff;
  outline: none;
}

.sort-selector select:focus {
  border-color: #f59e0b;
}

/* Сетка товаров */
.catalog-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
}

@media (max-width: 1150px) {
  .catalog-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 600px) {
  .catalog-grid {
    grid-template-columns: 1fr;
  }
}

.loading-box, .empty-box {
  padding: 60px 20px;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.empty-box h3 {
  font-size: 18px;
  font-weight: 800;
}
</style>
