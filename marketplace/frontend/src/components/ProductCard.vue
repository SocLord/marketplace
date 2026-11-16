<template>
  <div class="product-card card">
    <!-- Акс ва бейҷи тахфиф -->
    <div class="image-wrapper">
      <router-link :to="`/product/${product.id}`">
        <img :src="product.image_url" :alt="product.name" class="product-img" loading="lazy" />
      </router-link>
      <span v-if="discountPercent" class="badge-discount-tag">
        -{{ discountPercent }}%
      </span>
      <span v-if="product.stock_quantity <= 0" class="badge-out-of-stock">
        {{ t('outOfStock') }}
      </span>
    </div>

    <!-- Маълумот дар бораи мол -->
    <div class="card-body">
      <!-- Нарх (ҷорӣ ва пешина) -->
      <div class="price-row">
        <span class="current-price">{{ formatPrice(product.price) }} <small>{{ t('currency') }}</small></span>
        <span v-if="product.old_price && product.old_price > product.price" class="old-price">
          {{ formatPrice(product.old_price) }} с.
        </span>
      </div>

      <!-- Номи мол -->
      <router-link :to="`/product/${product.id}`" class="product-title" :title="product.name">
        {{ product.name }}
      </router-link>

      <!-- Рейтинг ва тақризҳо -->
      <div class="rating-row">
        <div class="stars">
          <Star :size="13" class="star-icon" />
          <span class="rating-value">{{ product.rating ? product.rating.toFixed(1) : '5.0' }}</span>
        </div>
        <span class="reviews-count">({{ product.reviews_count || 0 }} {{ t('reviewsCountSuffix') }})</span>
      </div>

      <!-- Фурӯшанда / Мағоза -->
      <div class="shop-badge-row" v-if="product.shop">
        <Store :size="12" class="shop-icon" />
        <router-link :to="`/shop/${product.shop.slug || product.shop.id}`" class="shop-name">
          {{ product.shop.name }}
        </router-link>
      </div>

      <!-- Тугмаи харид / Ба сабад -->
      <div class="card-action">
        <button 
          @click.stop="handleAddToCart" 
          :class="['add-to-cart-btn', isAdded && 'added']"
          :disabled="product.stock_quantity <= 0">
          <ShoppingCart :size="16" />
          <span>{{ isAdded ? t('inCart') : t('addToCart') }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { Star, Store, ShoppingCart } from 'lucide-vue-next'
import { addToCart } from '../api'
import { t } from '../utils/i18n'

const props = defineProps({
  product: {
    type: Object,
    required: true
  }
})

const isAdded = ref(false)

const discountPercent = computed(() => {
  if (props.product.old_price && props.product.old_price > props.product.price) {
    return Math.round(((props.product.old_price - props.product.price) / props.product.old_price) * 100)
  }
  return null
})

function formatPrice(val) {
  if (!val) return '0'
  return new Intl.NumberFormat('ru-RU').format(val)
}

function handleAddToCart() {
  addToCart(props.product, 1)
  isAdded.value = true
  setTimeout(() => {
    isAdded.value = false
  }, 1800)
}
</script>

<style scoped>
.product-card {
  display: flex;
  flex-direction: column;
  overflow: hidden;
  height: 100%;
  position: relative;
  background: #ffffff;
  border-radius: var(--radius-lg);
  border: 1px solid var(--border);
  transition: all 0.25s ease;
}

.product-card:hover {
  border-color: #f59e0b;
  box-shadow: 0 10px 24px -4px rgba(245, 158, 11, 0.25);
  transform: translateY(-2px);
}

.image-wrapper {
  position: relative;
  width: 100%;
  padding-top: 100%;
  background: #f8fafc;
  overflow: hidden;
}

.product-img {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.product-card:hover .product-img {
  transform: scale(1.05);
}

.badge-discount-tag {
  position: absolute;
  bottom: 8px;
  left: 8px;
  background: #ef4444;
  color: #ffffff;
  font-size: 11px;
  font-weight: 800;
  padding: 3px 8px;
  border-radius: 6px;
  box-shadow: 0 2px 6px rgba(239, 68, 68, 0.3);
}

.badge-out-of-stock {
  position: absolute;
  top: 8px;
  right: 8px;
  background: rgba(15, 23, 42, 0.75);
  color: #ffffff;
  font-size: 10px;
  font-weight: 700;
  padding: 3px 8px;
  border-radius: 4px;
}

.card-body {
  padding: 14px;
  display: flex;
  flex-direction: column;
  flex: 1;
}

.price-row {
  display: flex;
  align-items: baseline;
  gap: 8px;
  margin-bottom: 6px;
}

.current-price {
  font-size: 18px;
  font-weight: 800;
  color: #0f172a;
}

.current-price small {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-secondary);
}

.old-price {
  font-size: 12px;
  color: var(--text-muted);
  text-decoration: line-through;
}

.product-title {
  font-size: 13px;
  font-weight: 700;
  color: var(--text-main);
  line-height: 1.35;
  margin-bottom: 8px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  height: 35px;
  transition: color 0.15s;
}

.product-title:hover {
  color: #d97706;
}

.rating-row {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  margin-bottom: 6px;
}

.stars {
  display: flex;
  align-items: center;
  gap: 3px;
  color: #d97706;
  font-weight: 700;
}

.star-icon {
  fill: #f59e0b;
  color: #f59e0b;
}

.reviews-count {
  color: var(--text-muted);
  font-size: 11px;
}

.shop-badge-row {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 11px;
  color: var(--text-secondary);
  margin-bottom: 12px;
}

.shop-icon {
  color: #d97706;
}

.shop-name {
  color: var(--text-secondary);
  font-weight: 600;
}

.shop-name:hover {
  color: #b45309;
  text-decoration: underline;
}

.card-action {
  margin-top: auto;
}

.add-to-cart-btn {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 9px 12px;
  background: #fef3c7;
  color: #92400e;
  border: 1px solid #fde68a;
  border-radius: var(--radius-md);
  font-size: 13px;
  font-weight: 800;
  transition: all 0.2s;
}

.add-to-cart-btn:hover:not(:disabled) {
  background: #f59e0b;
  color: #0f172a;
  border-color: #f59e0b;
  box-shadow: 0 4px 12px rgba(245, 158, 11, 0.35);
}

.add-to-cart-btn.added {
  background: var(--success);
  color: #ffffff;
  border-color: var(--success);
}

.add-to-cart-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
