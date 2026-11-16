<template>
  <div class="product-detail-page">
    <div class="container">
      <div v-if="loading" class="loading-state">
        <span>{{ currentLang === 'ru' ? 'Загрузка товара...' : 'Боргирии маълумоти мол...' }}</span>
      </div>

      <div v-else-if="!product" class="not-found-box card">
        <h2>{{ currentLang === 'ru' ? 'Товар не найден' : 'Мол ёфт нашуд' }}</h2>
        <p>{{ currentLang === 'ru' ? 'Возможно товар был удален или еще проходит модерацию.' : 'Эҳтимол ин мол нест карда шудааст ё ҳанӯз аз ҷониби модератор тасдиқ нагардидааст.' }}</p>
        <router-link to="/catalog" class="btn btn-primary">{{ t('goToCatalog') }}</router-link>
      </div>

      <div v-else class="product-layout">
        <!-- Левая колонка: Галерея фото -->
        <div class="gallery-card card">
          <div class="main-image-wrap">
            <img :src="activeImage" :alt="product.name" class="main-img" />
            <span v-if="discountPercent" class="badge-discount-tag">
              -{{ discountPercent }}%
            </span>
          </div>
          <!-- Миниатюры (дополнительные фото) -->
          <div v-if="allImages.length > 1" class="thumbnails-row">
            <div 
              v-for="(img, idx) in allImages" 
              :key="idx"
              :class="['thumbnail-item', activeImage === img && 'active']"
              @click="activeImage = img">
              <img :src="img" :alt="`Photo ${idx+1}`" />
            </div>
          </div>
        </div>

        <!-- Центральная колонка: Основные детали товара -->
        <div class="details-card card">
          <div class="detail-header">
            <span class="category-crumb" v-if="product.category">
              {{ product.category.name }}
            </span>
            <h1 class="product-title">{{ product.name }}</h1>
            
            <div class="rating-badge-row">
              <div class="stars">
                <Star :size="16" class="star-icon" />
                <span class="rating-val">{{ product.rating ? product.rating.toFixed(1) : '5.0' }}</span>
              </div>
              <span class="reviews-link" @click="scrollToReviews">
                {{ product.reviews_count || 0 }} {{ t('reviewsCountSuffix') }}
              </span>
              <span class="sales-tag">
                {{ t('soldCount') }} {{ product.sales_count || 0 }} шт.
              </span>
            </div>
          </div>

          <!-- Нарх ва тахфиф -->
          <div class="price-section">
            <div class="price-main">
              <span class="current-price">{{ formatPrice(product.price) }} <small>{{ t('currency') }}</small></span>
              <span v-if="product.old_price && product.old_price > product.price" class="old-price">
                {{ formatPrice(product.old_price) }} с.
              </span>
            </div>
            <span class="stock-status" :class="product.stock_quantity > 0 ? 'in-stock' : 'out-stock'">
              <CheckCircle2 v-if="product.stock_quantity > 0" :size="15" />
              <AlertCircle v-else :size="15" />
              <span>{{ product.stock_quantity > 0 ? `${t('inStockCount')} (${product.stock_quantity} шт.)` : t('outOfStock') }}</span>
            </span>
          </div>

          <!-- Блок выбора количества и покупки -->
          <div class="purchase-box">
            <div class="quantity-control">
              <button @click="changeQty(-1)" :disabled="buyQty <= 1">-</button>
              <input type="number" v-model.number="buyQty" min="1" :max="product.stock_quantity" />
              <button @click="changeQty(1)" :disabled="buyQty >= product.stock_quantity">+</button>
            </div>

            <button 
              @click="handleAddToCart" 
              class="btn btn-primary btn-lg add-btn"
              :disabled="product.stock_quantity <= 0">
              <ShoppingCart :size="18" />
              <span>{{ addedSuccess ? t('inCart') : t('addToCart') }}</span>
            </button>
          </div>

          <!-- Тавсифи мол (Description) -->
          <div class="description-section">
            <h3>{{ t('productDesc') }}</h3>
            <p class="description-text">{{ product.description || 'Описание отсутствует.' }}</p>
          </div>
        </div>

        <!-- Правая колонка: Продавец и условия доставки -->
        <div class="sidebar-col">
          <!-- Карточка магазина -->
          <div class="shop-card card" v-if="product.shop">
            <div class="shop-header">
              <img :src="product.shop.logo_url || 'https://images.unsplash.com/photo-1550009158-9ebf69173e03?w=100'" class="shop-logo" alt="Shop logo" />
              <div class="shop-info">
                <span class="shop-title">{{ product.shop.name }}</span>
                <span class="shop-verified">
                  <CheckCheck :size="13" />
                  <span>{{ t('verifiedShop') }}</span>
                </span>
              </div>
            </div>
            <p class="shop-desc">{{ product.shop.description }}</p>
            <div class="shop-meta">
              <span v-if="product.shop.address">📍 {{ product.shop.address }}</span>
              <span v-if="product.shop.phone">📞 {{ product.shop.phone }}</span>
            </div>
            <router-link :to="`/shop/${product.shop.slug || product.shop.id}`" class="btn btn-secondary btn-sm visit-shop-btn">
              <Store :size="15" />
              <span>{{ t('visitShop') }}</span>
            </router-link>
          </div>

          <!-- Гарантии маркетплейса -->
          <div class="guarantee-box card">
            <div class="guarantee-item">
              <Truck :size="18" color="#f59e0b" />
              <div>
                <strong>{{ t('featDeliveryTitle') }}</strong>
                <p>1-3 дня в любой город</p>
              </div>
            </div>
            <div class="guarantee-item">
              <RotateCcw :size="18" color="#f59e0b" />
              <div>
                <strong>Легкий возврат</strong>
                <p>14 дней на обмен или возврат</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Раздел: Отзывы и оценки -->
      <section class="reviews-section" id="reviews-section">
        <div class="section-title-row">
          <h2>{{ t('reviewsTitle') }}</h2>
          <span class="reviews-total-badge">{{ reviews.length }} {{ t('reviewsCountSuffix') }}</span>
        </div>

        <div class="reviews-layout">
          <!-- Форма добавления нового отзыва -->
          <div class="add-review-card card">
            <h3>{{ t('writeReview') }}</h3>
            <p class="review-prompt">{{ currentLang === 'ru' ? 'Если вы покупали этот товар, поставьте оценку:' : 'Агар шумо ин молро харидорӣ намуда бошед, ба он баҳо диҳед:' }}</p>

            <form @submit.prevent="submitReview">
              <!-- Звезды оценки (1-5) -->
              <div class="form-group rating-select">
                <label>{{ t('yourRating') }}</label>
                <div class="interactive-stars">
                  <Star 
                    v-for="s in 5" 
                    :key="s" 
                    :size="26"
                    :class="['star-pick', s <= newReview.rating && 'active']"
                    @click="newReview.rating = s" 
                  />
                  <span class="selected-rating-text">{{ newReview.rating }} / 5 ★</span>
                </div>
              </div>

              <!-- Текст отзыва -->
              <div class="form-group">
                <label>{{ currentLang === 'ru' ? 'Текст отзыва' : 'Матни тақриз' }}</label>
                <textarea 
                  v-model="newReview.comment" 
                  rows="4" 
                  required
                  :placeholder="currentLang === 'ru' ? 'Опишите качество товара, скорость доставки и впечатления...' : 'Сифати мол, мувофиқати он ба тавсиф ва таассуроти худро муфассал баён намоед...'">
                </textarea>
              </div>

              <button type="submit" class="btn btn-primary" :disabled="submittingReview">
                {{ submittingReview ? 'Сохранение...' : t('sendReview') }}
              </button>
            </form>
          </div>

          <!-- Список отзывов -->
          <div class="reviews-list-col">
            <div v-if="reviews.length === 0" class="no-reviews card">
              <MessageSquare :size="36" color="#94a3b8" />
              <p>{{ currentLang === 'ru' ? 'Отзывов пока нет. Будьте первым, кто оставит оценку!' : 'Ҳанӯз ягон тақриз барои ин мол навишта нашудааст. Аввалин шуда баҳо диҳед!' }}</p>
            </div>

            <div v-else class="reviews-list">
              <div v-for="rev in reviews" :key="rev.id" class="review-item card">
                <div class="review-header">
                  <div class="author-info">
                    <img 
                      :src="rev.user?.avatar_url || 'https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?w=100'" 
                      class="author-avatar" 
                      alt="Avatar" 
                    />
                    <div>
                      <strong class="author-name">{{ rev.user?.name || 'Покупатель' }}</strong>
                      <span class="review-date">{{ formatDate(rev.created_at) }}</span>
                    </div>
                  </div>
                  <!-- Оценка звездами -->
                  <div class="review-stars">
                    <Star 
                      v-for="s in 5" 
                      :key="s" 
                      :size="14"
                      :class="['review-star', s <= rev.rating && 'filled']" 
                    />
                  </div>
                </div>
                <p class="review-text">{{ rev.comment }}</p>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { 
  Star, ShoppingCart, CheckCircle2, AlertCircle, 
  Store, CheckCheck, Truck, RotateCcw, MessageSquare 
} from 'lucide-vue-next'
import { api, addToCart, getStoredUser } from '../api'
import { t, currentLang } from '../utils/i18n'

const route = useRoute()

const loading = ref(true)
const product = ref(null)
const reviews = ref([])
const activeImage = ref('')
const buyQty = ref(1)
const addedSuccess = ref(false)

const submittingReview = ref(false)
const newReview = ref({
  rating: 5,
  comment: ''
})

const allImages = computed(() => {
  if (!product.value) return []
  const imgs = [product.value.image_url]
  if (product.value.additional_images) {
    try {
      const extra = JSON.parse(product.value.additional_images)
      if (Array.isArray(extra)) imgs.push(...extra)
    } catch (e) {}
  }
  return [...new Set(imgs)]
})

const discountPercent = computed(() => {
  if (product.value?.old_price && product.value.old_price > product.value.price) {
    return Math.round(((product.value.old_price - product.value.price) / product.value.old_price) * 100)
  }
  return null
})

function formatPrice(val) {
  if (!val) return '0'
  return new Intl.NumberFormat('ru-RU').format(val)
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleDateString(currentLang.value === 'ru' ? 'ru-RU' : 'tg-TJ', { day: 'numeric', month: 'long', year: 'numeric' })
}

function changeQty(diff) {
  const next = buyQty.value + diff
  if (next >= 1 && next <= (product.value?.stock_quantity || 1)) {
    buyQty.value = next
  }
}

function handleAddToCart() {
  if (!product.value) return
  addToCart(product.value, buyQty.value)
  addedSuccess.value = true
  setTimeout(() => {
    addedSuccess.value = false
  }, 2000)
}

function scrollToReviews() {
  const el = document.getElementById('reviews-section')
  if (el) el.scrollIntoView({ behavior: 'smooth' })
}

async function loadProduct() {
  loading.value = true
  try {
    const idOrSlug = route.params.idOrSlug
    const res = await api.products.get(idOrSlug)
    product.value = res.data
    activeImage.value = res.data.image_url
    reviews.value = res.data.reviews || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function submitReview() {
  const user = getStoredUser()
  if (!user) {
    alert(currentLang.value === 'ru' ? 'Пожалуйста, авторизуйтесь чтобы оставить отзыв.' : 'Илтимос, аввал ба система ворид шавед то тақриз гузошта тавонед.')
    return
  }

  if (!newReview.value.comment.trim()) {
    alert(currentLang.value === 'ru' ? 'Введите текст отзыва' : 'Матни тақризро ворид намоед')
    return
  }

  submittingReview.value = true
  try {
    await api.reviews.create({
      product_id: product.value.id,
      rating: newReview.value.rating,
      comment: newReview.value.comment.trim()
    })
    newReview.value.comment = ''
    newReview.value.rating = 5
    await loadProduct()
    alert(currentLang.value === 'ru' ? 'Спасибо! Ваш отзыв опубликован.' : 'Ташаккур! Фикру мулоҳизаи шумо сабт шуд.')
  } catch (err) {
    alert('Ошибка: ' + err.message)
  } finally {
    submittingReview.value = false
  }
}

onMounted(() => {
  loadProduct()
})
</script>

<style scoped>
.product-detail-page {
  padding: 24px 0 60px 0;
}

.product-layout {
  display: grid;
  grid-template-columns: 420px 1fr 300px;
  gap: 24px;
  align-items: start;
}

@media (max-width: 1100px) {
  .product-layout {
    grid-template-columns: 1fr 1fr;
  }
  .sidebar-col {
    grid-column: span 2;
  }
}

@media (max-width: 768px) {
  .product-layout {
    grid-template-columns: 1fr;
  }
  .sidebar-col {
    grid-column: span 1;
  }
}

/* Галерея */
.gallery-card {
  padding: 16px;
  background: #ffffff;
}

.main-image-wrap {
  position: relative;
  width: 100%;
  padding-top: 100%;
  background: #f8fafc;
  border-radius: var(--radius-md);
  overflow: hidden;
}

.main-img {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.badge-discount-tag {
  position: absolute;
  bottom: 12px;
  left: 12px;
  background: #ef4444;
  color: #ffffff;
  font-size: 13px;
  font-weight: 800;
  padding: 4px 10px;
  border-radius: 8px;
}

.thumbnails-row {
  display: flex;
  gap: 10px;
  margin-top: 14px;
}

.thumbnail-item {
  width: 60px;
  height: 60px;
  border-radius: var(--radius-sm);
  border: 2px solid transparent;
  cursor: pointer;
  overflow: hidden;
  transition: all 0.15s;
}

.thumbnail-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.thumbnail-item.active {
  border-color: #f59e0b;
}

/* Основные детали */
.details-card {
  padding: 28px;
  background: #ffffff;
}

.category-crumb {
  font-size: 12px;
  font-weight: 800;
  color: #b45309;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.product-title {
  font-size: 24px;
  font-weight: 800;
  color: var(--text-main);
  line-height: 1.3;
  margin: 8px 0 14px 0;
}

.rating-badge-row {
  display: flex;
  align-items: center;
  gap: 14px;
  flex-wrap: wrap;
  padding-bottom: 18px;
  border-bottom: 1px solid var(--border);
}

.stars {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #d97706;
  font-weight: 700;
  font-size: 15px;
}

.star-icon {
  fill: #f59e0b;
  color: #f59e0b;
}

.reviews-link {
  color: var(--text-secondary);
  font-size: 13px;
  text-decoration: underline;
  cursor: pointer;
}

.sales-tag {
  font-size: 12px;
  color: var(--text-muted);
  background: var(--bg-page);
  padding: 3px 8px;
  border-radius: 4px;
}

/* Цена */
.price-section {
  padding: 20px 0;
  border-bottom: 1px solid var(--border);
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.price-main {
  display: flex;
  align-items: baseline;
  gap: 12px;
}

.current-price {
  font-size: 28px;
  font-weight: 800;
  color: #0f172a;
}

.current-price small {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-secondary);
}

.old-price {
  font-size: 16px;
  color: var(--text-muted);
  text-decoration: line-through;
}

.stock-status {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 700;
}

.in-stock { color: var(--success); }
.out-stock { color: var(--danger); }

/* Блок покупки */
.purchase-box {
  display: flex;
  gap: 14px;
  padding: 20px 0;
  border-bottom: 1px solid var(--border);
}

.quantity-control {
  display: flex;
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  overflow: hidden;
  height: 48px;
}

.quantity-control button {
  width: 44px;
  background: var(--bg-page);
  font-size: 18px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.15s;
}

.quantity-control button:hover:not(:disabled) {
  background: #fef3c7;
}

.quantity-control input {
  width: 50px;
  text-align: center;
  border: none;
  outline: none;
  font-size: 16px;
  font-weight: 700;
}

.add-btn {
  flex: 1;
}

/* Тавсиф */
.description-section {
  padding-top: 22px;
}

.description-section h3 {
  font-size: 16px;
  font-weight: 800;
  margin-bottom: 10px;
}

.description-text {
  font-size: 14px;
  line-height: 1.7;
  color: var(--text-secondary);
}

/* Правая колонка / Магазин */
.shop-card {
  padding: 20px;
  background: #ffffff;
  margin-bottom: 16px;
}

.shop-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.shop-logo {
  width: 46px;
  height: 46px;
  border-radius: 12px;
  object-fit: cover;
  border: 1px solid var(--border);
}

.shop-title {
  font-size: 15px;
  font-weight: 800;
  display: block;
}

.shop-verified {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 11px;
  color: #b45309;
  font-weight: 700;
}

.shop-desc {
  font-size: 12px;
  color: var(--text-secondary);
  line-height: 1.5;
  margin-bottom: 12px;
}

.shop-meta {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 11px;
  color: var(--text-muted);
  margin-bottom: 16px;
}

.visit-shop-btn {
  width: 100%;
}

.guarantee-box {
  padding: 20px;
  background: #ffffff;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.guarantee-item {
  display: flex;
  gap: 12px;
}

.guarantee-item strong {
  font-size: 13px;
  display: block;
}

.guarantee-item p {
  font-size: 11px;
  color: var(--text-muted);
}

/* Секция отзывов */
.reviews-section {
  margin-top: 60px;
}

.section-title-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 24px;
}

.section-title-row h2 {
  font-size: 22px;
  font-weight: 800;
}

.reviews-total-badge {
  background: #fef3c7;
  color: #b45309;
  font-size: 12px;
  font-weight: 800;
  padding: 3px 10px;
  border-radius: 20px;
}

.reviews-layout {
  display: grid;
  grid-template-columns: 380px 1fr;
  gap: 24px;
  align-items: start;
}

@media (max-width: 850px) {
  .reviews-layout {
    grid-template-columns: 1fr;
  }
}

.add-review-card {
  padding: 24px;
  background: #ffffff;
}

.add-review-card h3 {
  font-size: 16px;
  font-weight: 800;
  margin-bottom: 6px;
}

.review-prompt {
  font-size: 12px;
  color: var(--text-muted);
  margin-bottom: 18px;
}

.interactive-stars {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 6px;
}

.star-pick {
  color: #cbd5e1;
  cursor: pointer;
  transition: transform 0.15s, color 0.15s;
}

.star-pick:hover {
  transform: scale(1.15);
}

.star-pick.active {
  color: #f59e0b;
  fill: #f59e0b;
}

.selected-rating-text {
  font-size: 13px;
  font-weight: 800;
  color: #b45309;
  margin-left: 8px;
}

.add-review-card textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  font-size: 13px;
  outline: none;
  resize: vertical;
}

.add-review-card textarea:focus {
  border-color: #f59e0b;
}

.reviews-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.review-item {
  padding: 18px 20px;
  background: #ffffff;
}

.review-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.author-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.author-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  object-fit: cover;
}

.author-name {
  font-size: 13px;
  font-weight: 700;
  display: block;
}

.review-date {
  font-size: 11px;
  color: var(--text-muted);
}

.review-stars {
  display: flex;
  gap: 2px;
}

.review-star {
  color: #e2e8f0;
}

.review-star.filled {
  color: #f59e0b;
  fill: #f59e0b;
}

.review-text {
  font-size: 13px;
  line-height: 1.6;
  color: var(--text-main);
}

.no-reviews {
  padding: 40px;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  color: var(--text-muted);
}

.not-found-box {
  padding: 60px;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14px;
}
</style>
