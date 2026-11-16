<template>
  <div class="home-page">
    <div class="container">
      <!-- Промо-баннер (Hero Slider / Banner в современном желто-графитовом стиле) -->
      <section class="hero-section">
        <div class="hero-banner">
          <div class="hero-content">
            <span class="hero-badge">{{ t('heroBadge') }}</span>
            <h1 class="hero-title">{{ t('heroTitle') }}</h1>
            <p class="hero-desc">{{ t('heroDesc') }}</p>
            <div class="hero-cta">
              <router-link to="/catalog" class="btn btn-primary btn-lg">
                <ShoppingBag :size="18" />
                <span>{{ t('goToCatalog') }}</span>
              </router-link>
              <router-link to="/shops" class="btn btn-secondary btn-lg hero-outline">
                <Store :size="18" />
                <span>{{ t('partnerShops') }}</span>
              </router-link>
            </div>
          </div>
          <div class="hero-image">
            <img 
              src="https://images.unsplash.com/photo-1607082348824-0a96f2a4b9da?w=800" 
              alt="Marketplace shopping" 
              class="hero-img" 
            />
          </div>
        </div>
      </section>

      <!-- Блок категорий (Категорияҳои асосӣ) -->
      <section class="section categories-section">
        <div class="section-header">
          <div>
            <h2 class="section-title">{{ t('popularCategories') }}</h2>
            <p class="section-subtitle">{{ t('popularCategoriesSub') }}</p>
          </div>
          <router-link to="/catalog" class="view-all-link">
            <span>{{ t('viewAll') }}</span>
            <ArrowRight :size="16" />
          </router-link>
        </div>

        <div class="categories-grid">
          <router-link 
            v-for="cat in categories" 
            :key="cat.id" 
            :to="`/catalog?category=${cat.slug}`"
            class="category-card">
            <div class="category-img-wrapper">
              <img :src="cat.image_url" :alt="cat.name" class="category-img" />
            </div>
            <div class="category-info">
              <h3 class="category-name">{{ cat.name }}</h3>
              <span class="category-count">{{ cat.products_count || 0 }} {{ t('allProducts').toLowerCase() }}</span>
            </div>
          </router-link>
        </div>
      </section>

      <!-- Хити фурӯш (Популярные товары) -->
      <section class="section">
        <div class="section-header">
          <div>
            <h2 class="section-title">{{ t('hitsTitle') }}</h2>
            <p class="section-subtitle">{{ t('hitsSub') }}</p>
          </div>
          <router-link to="/catalog?sort=popular" class="view-all-link">
            <span>{{ t('viewAll') }}</span>
            <ArrowRight :size="16" />
          </router-link>
        </div>

        <div v-if="loading" class="loading-state">
          <span>Боргирии молҳо...</span>
        </div>

        <div v-else class="products-grid">
          <ProductCard 
            v-for="product in popularProducts" 
            :key="product.id" 
            :product="product" 
          />
        </div>
      </section>

      <!-- Промо-баннер для продавцов (Стильный желтый фон с темным текстом) -->
      <section class="seller-promo-banner card">
        <div class="seller-promo-text">
          <span class="badge badge-seller-promo">{{ t('sellerBannerBadge') }}</span>
          <h2>{{ t('sellerBannerTitle') }}</h2>
          <p>{{ t('sellerBannerDesc') }}</p>
          <router-link to="/seller/dashboard" class="btn btn-dark btn-lg">
            <Store :size="18" />
            <span>{{ t('sellerBannerBtn') }}</span>
          </router-link>
        </div>
        <div class="seller-promo-img">
          <img 
            src="https://images.unsplash.com/photo-1556742049-0a67c5574f73?w=600" 
            alt="Seller on marketplace" 
          />
        </div>
      </section>

      <!-- Молҳои навтарин (Новинки каталога) -->
      <section class="section">
        <div class="section-header">
          <div>
            <h2 class="section-title">{{ t('newArrivalsTitle') }}</h2>
            <p class="section-subtitle">{{ t('newArrivalsSub') }}</p>
          </div>
          <router-link to="/catalog?sort=newest" class="view-all-link">
            <span>{{ t('viewAll') }}</span>
            <ArrowRight :size="16" />
          </router-link>
        </div>

        <div class="products-grid">
          <ProductCard 
            v-for="product in latestProducts" 
            :key="product.id" 
            :product="product" 
          />
        </div>
      </section>

      <!-- Афзалиятҳои маркетплейс -->
      <section class="features-section">
        <div class="feature-card card">
          <div class="feature-icon icon-yellow">
            <ShieldCheck :size="28" />
          </div>
          <h3>{{ t('featModTitle') }}</h3>
          <p>{{ t('featModDesc') }}</p>
        </div>

        <div class="feature-card card">
          <div class="feature-icon icon-amber">
            <Split :size="28" />
          </div>
          <h3>{{ t('featSplitTitle') }}</h3>
          <p>{{ t('featSplitDesc') }}</p>
        </div>

        <div class="feature-card card">
          <div class="feature-icon icon-yellow">
            <BarChart2 :size="28" />
          </div>
          <h3>{{ t('featStatsTitle') }}</h3>
          <p>{{ t('featStatsDesc') }}</p>
        </div>

        <div class="feature-card card">
          <div class="feature-icon icon-dark">
            <Truck :size="28" />
          </div>
          <h3>{{ t('featDeliveryTitle') }}</h3>
          <p>{{ t('featDeliveryDesc') }}</p>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ShoppingBag, Store, ArrowRight, ShieldCheck, Split, BarChart2, Truck } from 'lucide-vue-next'
import { api } from '../api'
import { t } from '../utils/i18n'
import ProductCard from '../components/ProductCard.vue'

const loading = ref(true)
const categories = ref([])
const popularProducts = ref([])
const latestProducts = ref([])

async function loadData() {
  loading.value = true
  try {
    const [catsRes, popRes, newRes] = await Promise.all([
      api.categories.list(),
      api.products.list({ sort: 'popular' }),
      api.products.list({ sort: 'newest' })
    ])
    categories.value = (catsRes.data || []).slice(0, 6)
    popularProducts.value = (popRes.data || []).slice(0, 8)
    latestProducts.value = (newRes.data || []).slice(0, 8)
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.home-page {
  padding: 24px 0 60px 0;
}

/* Баннер Hero в благородном графитово-желтом стиле */
.hero-section {
  margin-bottom: 40px;
}

.hero-banner {
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 60%, #334155 100%);
  border-radius: var(--radius-xl);
  overflow: hidden;
  display: grid;
  grid-template-columns: 1.2fr 1fr;
  align-items: center;
  color: #ffffff;
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.25);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

@media (max-width: 900px) {
  .hero-banner {
    grid-template-columns: 1fr;
  }
  .hero-image {
    display: none;
  }
}

.hero-content {
  padding: 48px;
}

.hero-badge {
  display: inline-block;
  background: rgba(245, 158, 11, 0.2);
  backdrop-filter: blur(10px);
  color: #fbbf24;
  padding: 6px 14px;
  border-radius: 30px;
  font-size: 13px;
  font-weight: 800;
  margin-bottom: 18px;
  border: 1px solid rgba(245, 158, 11, 0.4);
}

.hero-title {
  font-size: 38px;
  font-weight: 800;
  line-height: 1.2;
  margin-bottom: 16px;
  letter-spacing: -0.5px;
  color: #ffffff;
}

.hero-desc {
  font-size: 15px;
  line-height: 1.6;
  color: #cbd5e1;
  margin-bottom: 28px;
  max-width: 520px;
}

.hero-cta {
  display: flex;
  gap: 14px;
  flex-wrap: wrap;
}

.hero-outline {
  background: rgba(255, 255, 255, 0.1) !important;
  color: #ffffff !important;
  border-color: rgba(255, 255, 255, 0.3) !important;
}

.hero-outline:hover {
  background: rgba(255, 255, 255, 0.25) !important;
}

.hero-image {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: flex-end;
}

.hero-img {
  width: 100%;
  height: 380px;
  object-fit: cover;
}

/* Секции */
.section {
  margin-bottom: 48px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 22px;
}

.section-title {
  font-size: 22px;
  font-weight: 800;
  color: var(--text-main);
  letter-spacing: -0.3px;
}

.section-subtitle {
  font-size: 13px;
  color: var(--text-muted);
  margin-top: 4px;
}

.view-all-link {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  font-weight: 700;
  color: #d97706;
  transition: gap 0.2s;
}

.view-all-link:hover {
  gap: 10px;
  color: #b45309;
}

/* Сетка категорий */
.categories-grid {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 16px;
}

@media (max-width: 1100px) {
  .categories-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 600px) {
  .categories-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

.category-card {
  background: #ffffff;
  border-radius: var(--radius-lg);
  border: 1px solid var(--border);
  overflow: hidden;
  transition: all 0.2s ease;
  text-align: center;
  display: flex;
  flex-direction: column;
}

.category-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 10px 20px -4px rgba(245, 158, 11, 0.2);
  border-color: #f59e0b;
}

.category-img-wrapper {
  height: 120px;
  overflow: hidden;
  background: #f1f5f9;
}

.category-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}

.category-card:hover .category-img {
  transform: scale(1.08);
}

.category-info {
  padding: 12px 10px;
}

.category-name {
  font-size: 13px;
  font-weight: 700;
  color: var(--text-main);
  margin-bottom: 2px;
}

.category-count {
  font-size: 11px;
  color: var(--text-muted);
}

/* Сетка товаров */
.products-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
}

@media (max-width: 1100px) {
  .products-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 768px) {
  .products-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
  }
}

@media (max-width: 480px) {
  .products-grid {
    grid-template-columns: 1fr;
  }
}

/* Баннер для продавцов в золотисто-желтом градиенте */
.seller-promo-banner {
  margin: 40px 0;
  padding: 36px 40px;
  background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
  border: 1px solid #fcd34d;
  display: grid;
  grid-template-columns: 1.5fr 1fr;
  align-items: center;
  gap: 30px;
}

@media (max-width: 800px) {
  .seller-promo-banner {
    grid-template-columns: 1fr;
    padding: 24px;
  }
  .seller-promo-img {
    display: none;
  }
}

.badge-seller-promo {
  background: #f59e0b;
  color: #0f172a;
  font-weight: 800;
}

.seller-promo-text h2 {
  font-size: 26px;
  font-weight: 800;
  color: #78350f;
  margin: 12px 0 10px 0;
}

.seller-promo-text p {
  font-size: 14px;
  color: #92400e;
  line-height: 1.6;
  margin-bottom: 22px;
}

.seller-promo-img img {
  width: 100%;
  height: 220px;
  border-radius: var(--radius-lg);
  object-fit: cover;
  box-shadow: var(--shadow-md);
}

/* Преимущества */
.features-section {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-top: 50px;
}

@media (max-width: 1000px) {
  .features-section {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 550px) {
  .features-section {
    grid-template-columns: 1fr;
  }
}

.feature-card {
  padding: 24px;
  text-align: center;
}

.feature-icon {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px auto;
}

.icon-yellow { background: #fef3c7; color: #b45309; }
.icon-amber { background: #fffbeb; color: #d97706; }
.icon-dark { background: #f1f5f9; color: #0f172a; }

.feature-card h3 {
  font-size: 15px;
  font-weight: 800;
  margin-bottom: 8px;
}

.feature-card p {
  font-size: 12px;
  color: var(--text-secondary);
  line-height: 1.5;
}

.loading-state {
  text-align: center;
  padding: 40px;
  color: var(--text-muted);
}
</style>
