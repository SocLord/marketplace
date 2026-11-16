<template>
  <div class="shops-list-page">
    <div class="container">
      <div class="page-header">
        <h1 class="page-title">{{ currentLang === 'ru' ? 'Магазины-партнёры маркетплейса' : 'Мағозаҳои шарики маркетплейс' }}</h1>
        <p class="page-subtitle">{{ currentLang === 'ru' ? 'Официальные проверенные магазины продавцов с гарантией подлинности' : 'Мағозаҳои расмӣ ва санҷидашудаи фурӯшандагон бо кафолати аслӣ' }}</p>
      </div>

      <div v-if="loading" class="loading-state">
        <span>{{ currentLang === 'ru' ? 'Загрузка магазинов...' : 'Боргирии мағозаҳо...' }}</span>
      </div>

      <div v-else class="shops-grid">
        <div v-for="shop in shops" :key="shop.id" class="shop-card card">
          <div class="shop-banner">
            <img :src="shop.banner_url || 'https://images.unsplash.com/photo-1518770660439-4636190af475?w=600'" :alt="shop.name" />
          </div>

          <div class="shop-body">
            <div class="shop-logo-wrap">
              <img :src="shop.logo_url || 'https://images.unsplash.com/photo-1550009158-9ebf69173e03?w=150'" :alt="shop.name" />
            </div>

            <div class="shop-info">
              <div class="shop-name-row">
                <h3 class="shop-name">{{ shop.name }}</h3>
                <span class="rating-badge">★ {{ shop.rating ? shop.rating.toFixed(1) : '5.0' }}</span>
              </div>
              <p class="shop-desc">{{ shop.description }}</p>

              <div class="shop-meta">
                <span v-if="shop.address">📍 {{ shop.address }}</span>
                <span v-if="shop.phone">📞 {{ shop.phone }}</span>
              </div>

              <router-link :to="`/shop/${shop.slug || shop.id}`" class="btn btn-primary btn-sm visit-btn">
                <span>{{ currentLang === 'ru' ? 'Перейти в магазин' : 'Ба мағоза гузаштан' }}</span>
                <ArrowRight :size="14" />
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ArrowRight } from 'lucide-vue-next'
import { api } from '../api'
import { currentLang } from '../utils/i18n'

const loading = ref(true)
const shops = ref([])

async function loadShops() {
  loading.value = true
  try {
    const res = await api.shops.list()
    shops.value = res.data || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadShops()
})
</script>

<style scoped>
.shops-list-page {
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

.shops-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
}

@media (max-width: 950px) {
  .shops-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 600px) {
  .shops-grid {
    grid-template-columns: 1fr;
  }
}

.shop-card {
  background: #ffffff;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.shop-banner {
  height: 120px;
  background: #e2e8f0;
  overflow: hidden;
}

.shop-banner img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.shop-body {
  padding: 0 20px 20px 20px;
  position: relative;
  display: flex;
  flex-direction: column;
  flex: 1;
}

.shop-logo-wrap {
  margin-top: -30px;
  width: 60px;
  height: 60px;
  border-radius: 14px;
  border: 3px solid #ffffff;
  box-shadow: var(--shadow-md);
  overflow: hidden;
  background: #ffffff;
  margin-bottom: 12px;
}

.shop-logo-wrap img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.shop-name-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.shop-name {
  font-size: 16px;
  font-weight: 800;
  color: var(--text-main);
}

.rating-badge {
  background: #fffbeb;
  color: #b45309;
  font-size: 12px;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 4px;
}

.shop-desc {
  font-size: 12px;
  color: var(--text-secondary);
  line-height: 1.5;
  margin-bottom: 12px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.shop-meta {
  display: flex;
  flex-direction: column;
  gap: 3px;
  font-size: 11px;
  color: var(--text-muted);
  margin-bottom: 16px;
  margin-top: auto;
}

.visit-btn {
  width: 100%;
}
</style>
