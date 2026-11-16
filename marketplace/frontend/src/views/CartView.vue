<template>
  <div class="cart-page">
    <div class="container">
      <div class="cart-header">
        <h1 class="cart-title">{{ t('cartTitle') }}</h1>
        <button v-if="cartItems.length > 0" @click="handleClearCart" class="clear-cart-btn">
          <Trash2 :size="15" />
          <span>{{ t('clearCartText') }}</span>
        </button>
      </div>

      <!-- Если корзина пуста -->
      <div v-if="cartItems.length === 0" class="empty-cart-card card">
        <div class="empty-icon-wrap">
          <ShoppingCart :size="52" color="#94a3b8" />
        </div>
        <h2>{{ t('emptyCartText') }}</h2>
        <p>{{ currentLang === 'ru' ? 'Вы можете перейти в каталог и добавить интересующие товары.' : 'Шумо метавонед ба каталог гузашта, молҳои дилхоҳатонро ба сабад илова намоед.' }}</p>
        <router-link to="/catalog" class="btn btn-primary btn-lg">
          {{ t('goToCatalog') }}
        </router-link>
      </div>

      <!-- Если в корзине есть товары -->
      <div v-else class="cart-layout">
        <!-- Левая колонка: Товары, сгруппированные по продавцам / магазинам -->
        <div class="cart-items-column">
          <div class="multi-vendor-notice card">
            <Info :size="18" color="#d97706" />
            <p>
              <strong>{{ currentLang === 'ru' ? 'Разделение заказа:' : 'Тақсими фармоиш:' }}</strong> 
              {{ currentLang === 'ru' 
                 ? `В вашей корзине товары из ${Object.keys(groupedByShop).length} магазинов. Заказ разделен на части, каждый магазин отправляет посылку независимо.` 
                 : `Дар сабади шумо молҳо аз ${Object.keys(groupedByShop).length} мағоза мавҷуданд. Фармоиш ба қисмҳо тақсим шуда, ҳар як мағоза бастаи худро мустақилона мерасонад.` }}
            </p>
          </div>

          <!-- Группы по магазинам (Как в WB / Ozon) -->
          <div 
            v-for="(group, shopId) in groupedByShop" 
            :key="shopId" 
            class="shop-group-card card">
            
            <div class="shop-group-header">
              <Store :size="18" class="shop-icon" />
              <span class="shop-title">{{ t('packageFrom') }} <strong>{{ group.shopName }}</strong></span>
              <span class="shop-items-badge">{{ group.items.length }} {{ currentLang === 'ru' ? 'тов.' : 'мол' }}</span>
            </div>

            <div class="items-list">
              <div 
                v-for="item in group.items" 
                :key="item.product_id" 
                class="cart-item-row">
                
                <img :src="item.image_url" :alt="item.name" class="item-img" />

                <div class="item-details">
                  <router-link :to="`/product/${item.product_id}`" class="item-name">
                    {{ item.name }}
                  </router-link>
                  <div class="item-price">
                    <span class="current">{{ formatPrice(item.price) }} {{ t('currency') }}</span>
                    <span v-if="item.old_price" class="old">{{ formatPrice(item.old_price) }} {{ currentLang === 'ru' ? 'с.' : 'с.' }}</span>
                  </div>
                </div>

                <!-- Управление количеством -->
                <div class="item-qty-control">
                  <button @click="updateQty(item.product_id, item.quantity - 1)">-</button>
                  <span>{{ item.quantity }}</span>
                  <button @click="updateQty(item.product_id, item.quantity + 1)" :disabled="item.quantity >= item.stock_quantity">+</button>
                </div>

                <!-- Сумма за позицию -->
                <div class="item-subtotal">
                  <span>{{ formatPrice(item.price * item.quantity) }} {{ currentLang === 'ru' ? 'сом.' : 'с.' }}</span>
                </div>

                <!-- Удалить -->
                <button @click="removeItem(item.product_id)" class="remove-btn" :title="currentLang === 'ru' ? 'Удалить' : 'Нест кардан'">
                  <Trash2 :size="16" />
                </button>
              </div>
            </div>

            <!-- Итог по магазину -->
            <div class="shop-group-footer">
              <span>{{ currentLang === 'ru' ? 'Посылка магазина ' + group.shopName + ':' : 'Ҷамъи мағозаи ' + group.shopName + ':' }}</span>
              <strong>{{ formatPrice(group.subtotal) }} {{ t('currency') }}</strong>
            </div>
          </div>
        </div>

        <!-- Правая колонка: Общий расчет и кнопка перехода к оформлению -->
        <div class="cart-summary-column">
          <div class="summary-card card">
            <h3 class="summary-title">{{ t('checkoutSummary') }}</h3>

            <div class="summary-rows">
              <div class="summary-row">
                <span>{{ t('totalItems') }}</span>
                <span>{{ totalItemsCount }} {{ currentLang === 'ru' ? 'шт.' : 'дона' }}</span>
              </div>
              <div class="summary-row">
                <span>{{ t('packagesCount') }}</span>
                <span>{{ Object.keys(groupedByShop).length }} {{ currentLang === 'ru' ? 'посылок' : 'қисм' }}</span>
              </div>
              <div class="summary-row">
                <span>{{ t('delivery') }}</span>
                <span class="text-free">{{ t('freeDelivery') }}</span>
              </div>
            </div>

            <div class="summary-divider"></div>

            <div class="grand-total-row">
              <span>{{ t('totalToPay') }}</span>
              <span class="grand-total-price">{{ formatPrice(grandTotal) }} <small>{{ t('currency') }}</small></span>
            </div>

            <router-link to="/checkout" class="btn btn-primary btn-lg checkout-btn">
              <span>{{ t('proceedCheckout') }}</span>
              <ArrowRight :size="18" />
            </router-link>

            <div class="cart-perks">
              <div class="perk">
                <Check :size="14" color="#10b981" />
                <span>{{ currentLang === 'ru' ? 'Оплата наличными или картой при получении' : 'Пардохти нақдӣ ё бо корт ҳангоми гирифтан' }}</span>
              </div>
              <div class="perk">
                <Check :size="14" color="#10b981" />
                <span>{{ currentLang === 'ru' ? 'Гарантия возврата и качества товаров' : 'Кафолати сифат ва имконияти бозгардонидан' }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { 
  ShoppingCart, Trash2, Info, Store, 
  ArrowRight, Check 
} from 'lucide-vue-next'
import { getCart, updateCartQuantity, removeFromCart, clearCart } from '../api'
import { t, currentLang } from '../utils/i18n'

const cartItems = ref([])

function loadCart() {
  cartItems.value = getCart()
}

// Группировка товаров по магазинам (Специальное требование!)
const groupedByShop = computed(() => {
  const groups = {}
  for (const item of cartItems.value) {
    const sId = item.shop_id || 0
    if (!groups[sId]) {
      groups[sId] = {
        shopId: sId,
        shopName: item.shop_name || (currentLang.value === 'ru' ? `Магазин #${sId}` : `Мағозаи #${sId}`),
        items: [],
        subtotal: 0
      }
    }
    groups[sId].items.push(item)
    groups[sId].subtotal += item.price * item.quantity
  }
  return groups
})

const totalItemsCount = computed(() => {
  return cartItems.value.reduce((sum, it) => sum + it.quantity, 0)
})

const grandTotal = computed(() => {
  return cartItems.value.reduce((sum, it) => sum + (it.price * it.quantity), 0)
})

function formatPrice(val) {
  if (!val) return '0'
  return new Intl.NumberFormat('ru-RU').format(val)
}

function updateQty(productId, qty) {
  updateCartQuantity(productId, qty)
  loadCart()
}

function removeItem(productId) {
  removeFromCart(productId)
  loadCart()
}

function handleClearCart() {
  const msg = currentLang.value === 'ru' ? 'Очистить все товары из корзины?' : 'Оё мехоҳед ҳамаи молҳоро аз сабад тоза кунед?'
  if (confirm(msg)) {
    clearCart()
    loadCart()
  }
}

onMounted(() => {
  loadCart()
})
</script>

<style scoped>
.cart-page {
  padding: 24px 0 60px 0;
}

.cart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.cart-title {
  font-size: 26px;
  font-weight: 800;
}

.clear-cart-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--danger);
  font-size: 13px;
  font-weight: 600;
}

.clear-cart-btn:hover {
  text-decoration: underline;
}

.empty-cart-card {
  padding: 60px 20px;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14px;
}

.empty-icon-wrap {
  width: 90px;
  height: 90px;
  background: var(--bg-page);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.empty-cart-card h2 {
  font-size: 20px;
  font-weight: 800;
}

.empty-cart-card p {
  color: var(--text-muted);
  font-size: 14px;
  max-width: 420px;
}

/* Макет корзины */
.cart-layout {
  display: grid;
  grid-template-columns: 1fr 360px;
  gap: 24px;
  align-items: start;
}

@media (max-width: 950px) {
  .cart-layout {
    grid-template-columns: 1fr;
  }
}

.multi-vendor-notice {
  padding: 16px 20px;
  background: #fffbeb;
  border: 1px solid rgba(245, 158, 11, 0.3);
  border-radius: 12px;
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
  font-size: 13px;
  color: #78350f;
}

/* Группа магазина */
.shop-group-card {
  background: #ffffff;
  margin-bottom: 20px;
  overflow: hidden;
}

.shop-group-header {
  padding: 14px 20px;
  background: var(--bg-page);
  border-bottom: 1px solid var(--border);
  display: flex;
  align-items: center;
  gap: 10px;
}

.shop-icon {
  color: var(--primary);
}

.shop-title {
  font-size: 14px;
  color: var(--text-main);
}

.shop-items-badge {
  margin-left: auto;
  font-size: 11px;
  font-weight: 700;
  background: #fef3c7;
  color: #92400e;
  padding: 2px 8px;
  border-radius: 10px;
}

.items-list {
  padding: 0 20px;
}

.cart-item-row {
  display: grid;
  grid-template-columns: 64px 1fr auto auto auto;
  gap: 16px;
  align-items: center;
  padding: 18px 0;
  border-bottom: 1px solid var(--border-light);
}

.item-img {
  width: 64px;
  height: 64px;
  border-radius: 8px;
  object-fit: cover;
}

.item-details {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.item-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-main);
  line-height: 1.3;
}

.item-name:hover {
  color: var(--primary-hover);
}

.item-price {
  display: flex;
  gap: 8px;
  font-size: 13px;
}

.item-price .current {
  font-weight: 700;
  color: #0f172a;
}

.item-price .old {
  color: var(--text-muted);
  text-decoration: line-through;
  font-size: 12px;
}

.item-qty-control {
  display: flex;
  align-items: center;
  border: 1px solid var(--border);
  border-radius: 6px;
  overflow: hidden;
}

.item-qty-control button {
  width: 30px;
  height: 30px;
  background: var(--bg-page);
  font-weight: 700;
  transition: background 0.15s;
}

.item-qty-control button:hover:not(:disabled) {
  background: #fef3c7;
}

.item-qty-control span {
  width: 34px;
  text-align: center;
  font-size: 13px;
  font-weight: 700;
}

.item-subtotal {
  font-size: 14px;
  font-weight: 800;
  color: #0f172a;
  min-width: 80px;
  text-align: right;
}

.remove-btn {
  color: var(--text-muted);
  padding: 6px;
  border-radius: 6px;
  transition: all 0.15s;
}

.remove-btn:hover {
  color: var(--danger);
  background: #fee2e2;
}

.shop-group-footer {
  padding: 12px 20px;
  background: #fafafa;
  border-top: 1px solid var(--border-light);
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  font-size: 13px;
}

/* Блок итогов */
.summary-card {
  padding: 24px;
  background: #ffffff;
  position: sticky;
  top: 120px;
}

.summary-title {
  font-size: 18px;
  font-weight: 800;
  margin-bottom: 18px;
}

.summary-rows {
  display: flex;
  flex-direction: column;
  gap: 12px;
  font-size: 13px;
  color: var(--text-secondary);
}

.summary-row {
  display: flex;
  justify-content: space-between;
}

.text-free {
  color: var(--success);
  font-weight: 700;
}

.summary-divider {
  height: 1px;
  background: var(--border);
  margin: 18px 0;
}

.grand-total-row {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  margin-bottom: 22px;
}

.grand-total-price {
  font-size: 24px;
  font-weight: 800;
  color: var(--text-main);
}

.grand-total-price small {
  font-size: 14px;
  font-weight: 600;
}

.checkout-btn {
  width: 100%;
}

.cart-perks {
  margin-top: 20px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.perk {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: var(--text-muted);
}
</style>
