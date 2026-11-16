<template>
  <div class="shop-page">
    <div class="container">
      <div v-if="loading" class="loading-state">
        <span>{{ currentLang === 'ru' ? 'Загрузка магазина...' : 'Боргирии мағоза...' }}</span>
      </div>

      <div v-else-if="!shop" class="not-found-box card">
        <h2>{{ currentLang === 'ru' ? 'Магазин не найден' : 'Мағоза ёфт нашуд' }}</h2>
        <router-link to="/shops" class="btn btn-primary">{{ currentLang === 'ru' ? 'К списку магазинов' : 'Ба рӯйхати мағозаҳо' }}</router-link>
      </div>

      <div v-else>
        <!-- Баннер и профиль магазина -->
        <div class="shop-hero-card card">
          <div class="shop-banner">
            <img :src="shop.banner_url || 'https://images.unsplash.com/photo-1518770660439-4636190af475?w=1200'" alt="Banner" class="banner-img" />
          </div>

          <div class="shop-profile-bar">
            <div class="shop-logo-wrap">
              <img :src="shop.logo_url || 'https://images.unsplash.com/photo-1550009158-9ebf69173e03?w=200'" alt="Logo" class="profile-logo" />
            </div>

            <div class="shop-title-meta">
              <div class="name-row">
                <h1 class="shop-name">{{ shop.name }}</h1>
                <span class="verified-pill">
                  <CheckCheck :size="14" />
                  <span>{{ currentLang === 'ru' ? 'Официальный магазин' : 'Мағозаи расмӣ' }}</span>
                </span>
              </div>
              <p class="shop-description">{{ shop.description }}</p>
              <div class="shop-contacts">
                <span v-if="shop.address">📍 {{ shop.address }}</span>
                <span v-if="shop.phone">📞 {{ shop.phone }}</span>
                <span class="shop-rating-pill">★ {{ shop.rating ? shop.rating.toFixed(1) : '5.0' }} {{ currentLang === 'ru' ? 'рейтинг' : 'рейтинг' }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Каталог товаров магазина -->
        <div class="shop-products-section">
          <div class="section-header">
            <div>
              <h2 class="section-title">{{ currentLang === 'ru' ? `Все товары ${shop.name}` : `Ҳамаи молҳои ${shop.name}` }}</h2>
              <span class="products-count">{{ currentLang === 'ru' ? 'Количество товаров' : 'Шумораи молҳо' }}: {{ (shop.products || []).length }} {{ currentLang === 'ru' ? 'шт.' : 'дона' }}</span>
            </div>
          </div>

          <div v-if="(shop.products || []).length === 0" class="empty-shop-box card">
            <Package :size="40" color="#94a3b8" />
            <p>{{ currentLang === 'ru' ? 'В этом магазине сейчас нет активных товаров.' : 'Дар ин мағоза айни замон моли фаъол нашр нашудааст.' }}</p>
          </div>

          <div v-else class="products-grid">
            <ProductCard 
              v-for="product in shop.products" 
              :key="product.id" 
              :product="{ ...product, shop: shop }" 
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { CheckCheck, Package } from 'lucide-vue-next'
import { api } from '../api'
import { currentLang } from '../utils/i18n'
import ProductCard from '../components/ProductCard.vue'

const route = useRoute()
const loading = ref(true)
const shop = ref(null)

async function loadShop() {
  loading.value = true
  try {
    const idOrSlug = route.params.idOrSlug
    const res = await api.shops.get(idOrSlug)
    shop.value = res.data
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadShop()
})
</script>

<style scoped>
.shop-page {
  padding: 24px 0 60px 0;
}

.shop-hero-card {
  overflow: hidden;
  margin-bottom: 36px;
  background: #ffffff;
}

.shop-banner {
  height: 220px;
  background: #e2e8f0;
  overflow: hidden;
}

.banner-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.shop-profile-bar {
  padding: 0 32px 28px 32px;
  display: flex;
  gap: 24px;
  position: relative;
}

@media (max-width: 700px) {
  .shop-profile-bar {
    flex-direction: column;
    padding: 0 20px 20px 20px;
    align-items: center;
    text-align: center;
  }
}

.shop-logo-wrap {
  margin-top: -46px;
  width: 92px;
  height: 92px;
  border-radius: 20px;
  border: 4px solid #ffffff;
  background: #ffffff;
  box-shadow: var(--shadow-md);
  overflow: hidden;
  flex-shrink: 0;
}

.profile-logo {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.shop-title-meta {
  padding-top: 14px;
  flex: 1;
}

.name-row {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.shop-name {
  font-size: 26px;
  font-weight: 800;
  color: var(--text-main);
}

.verified-pill {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  background: #fef3c7;
  color: #92400e;
  font-size: 12px;
  font-weight: 700;
  padding: 3px 10px;
  border-radius: 20px;
}

.shop-description {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.5;
  margin: 8px 0 12px 0;
}

.shop-contacts {
  display: flex;
  gap: 16px;
  font-size: 12px;
  color: var(--text-muted);
  flex-wrap: wrap;
}

.shop-rating-pill {
  background: #fffbeb;
  color: #b45309;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 4px;
}

.section-header {
  margin-bottom: 22px;
}

.section-title {
  font-size: 22px;
  font-weight: 800;
}

.products-count {
  font-size: 13px;
  color: var(--text-muted);
}

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
  }
}

@media (max-width: 480px) {
  .products-grid {
    grid-template-columns: 1fr;
  }
}

.empty-shop-box {
  padding: 60px;
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
