<template>
  <div class="checkout-page">
    <div class="container">
      <div class="checkout-header">
        <h1 class="page-title">{{ t('checkoutPageTitle') }}</h1>
        <router-link to="/cart" class="back-link">
          ← {{ currentLang === 'ru' ? 'Вернуться в корзину' : 'Баргаштан ба сабад' }}
        </router-link>
      </div>

      <!-- Если корзина пуста -->
      <div v-if="cartItems.length === 0" class="empty-box card">
        <p>{{ currentLang === 'ru' ? 'В вашей корзине пока нет товаров.' : 'Дар сабади шумо ҳанӯз ягон мол нест.' }}</p>
        <router-link to="/catalog" class="btn btn-primary">{{ t('goToCatalog') }}</router-link>
      </div>

      <div v-else class="checkout-layout">
        <!-- Левая колонка: Форма данных получателя и доставки -->
        <div class="checkout-form-col">
          <form @submit.prevent="submitOrder" class="checkout-form">
            <!-- 1. Данные получателя -->
            <div class="form-section card">
              <h3 class="section-heading">{{ t('recipientData') }}</h3>
              
              <div class="form-grid">
                <div class="form-group">
                  <label>{{ t('recipientNameLabel') }}</label>
                  <input 
                    type="text" 
                    v-model="orderForm.recipient_name" 
                    required 
                    :placeholder="currentLang === 'ru' ? 'Например: Дилшод Рахимов' : 'Масалан: Дилшод Раҳимов'" 
                  />
                </div>

                <div class="form-group">
                  <label>{{ t('recipientPhoneLabel') }}</label>
                  <input 
                    type="tel" 
                    v-model="orderForm.recipient_phone" 
                    required 
                    placeholder="+992 90 111 2233" 
                  />
                </div>
              </div>
            </div>

            <!-- 2. Адрес доставки -->
            <div class="form-section card">
              <h3 class="section-heading">{{ t('shippingAddressHeading') }}</h3>

              <div class="form-group">
                <label>{{ t('shippingAddressLabel') }}</label>
                <textarea 
                  v-model="orderForm.shipping_address" 
                  rows="3" 
                  required
                  :placeholder="currentLang === 'ru' ? 'Например: г. Душанбе, р-н Сино, ул. Борбад 48, кв. 12' : 'Масалан: ш. Душанбе, н. Сино, кӯчаи Борбад 48, хонаи 12'">
                </textarea>
              </div>

              <div class="form-group">
                <label>{{ t('notesLabel') }}</label>
                <input 
                  type="text" 
                  v-model="orderForm.notes" 
                  :placeholder="currentLang === 'ru' ? 'Позвонить перед доставкой, код домофона 45...' : 'Пеш аз омадан занг занед, коди домофон 45 ва ғайра.'" 
                />
              </div>
            </div>

            <!-- 3. Способ оплаты -->
            <div class="form-section card">
              <h3 class="section-heading">{{ t('paymentHeading') }}</h3>

              <div class="payment-methods">
                <label class="payment-card" :class="orderForm.payment_method === 'cash_on_delivery' && 'selected'">
                  <input type="radio" value="cash_on_delivery" v-model="orderForm.payment_method" />
                  <Banknote :size="20" />
                  <div>
                    <strong>{{ t('cashPayment') }}</strong>
                    <p>{{ currentLang === 'ru' ? 'Оплата курьеру после проверки товара' : 'Маблағро ба курер пас аз санҷидани мол супоред' }}</p>
                  </div>
                </label>

                <label class="payment-card" :class="orderForm.payment_method === 'card' && 'selected'">
                  <input type="radio" value="card" v-model="orderForm.payment_method" />
                  <CreditCard :size="20" />
                  <div>
                    <strong>{{ t('cardPayment') }}</strong>
                    <p>{{ currentLang === 'ru' ? 'Мгновенная оплата банковской картой' : 'Пардохти фаврӣ ва бехатар тариқи корти бонкӣ' }}</p>
                  </div>
                </label>
              </div>
            </div>

            <button type="submit" class="btn btn-primary btn-lg submit-order-btn" :disabled="submitting">
              {{ submitting 
                 ? (currentLang === 'ru' ? 'Оформление и разделение заказа...' : 'Дар ҳоли сабт ва тақсими фармоиш...') 
                 : t('confirmOrderBtn') }}
            </button>
          </form>
        </div>

        <!-- Правая колонка: Просмотр состава и разделения по продавцам -->
        <div class="checkout-summary-col">
          <div class="summary-card card">
            <h3>{{ currentLang === 'ru' ? 'Состав заказа' : 'Таркиби фармоиш' }}</h3>
            <p class="split-info-text">
              ✨ {{ currentLang === 'ru'
                 ? `Ваш заказ разделен на ${Object.keys(groupedByShop).length} отдельные посылки. Каждый продавец видит только свои товары!`
                 : `Молҳои шумо ба ${Object.keys(groupedByShop).length} бастаи алоҳида тақсим шуда, ҳар фурӯшанда танҳо қисми худро мебинад.` }}
            </p>

            <div class="split-packages-list">
              <div 
                v-for="(grp, sId) in groupedByShop" 
                :key="sId" 
                class="split-package-item">
                <div class="pkg-header">
                  <Store :size="14" color="#d97706" />
                  <strong>{{ grp.shopName }}</strong>
                </div>
                <div class="pkg-items">
                  <div v-for="it in grp.items" :key="it.product_id" class="pkg-item-row">
                    <span>{{ it.name }} × {{ it.quantity }}</span>
                    <span>{{ formatPrice(it.price * it.quantity) }} {{ currentLang === 'ru' ? 'сом.' : 'с.' }}</span>
                  </div>
                </div>
                <div class="pkg-footer">
                  <span>{{ currentLang === 'ru' ? 'Итого по магазину:' : 'Ҷамъи баста:' }}</span>
                  <strong>{{ formatPrice(grp.subtotal) }} {{ t('currency') }}</strong>
                </div>
              </div>
            </div>

            <div class="summary-divider"></div>

            <div class="total-row">
              <span>{{ t('totalToPay') }}</span>
              <strong class="total-sum">{{ formatPrice(grandTotal) }} {{ t('currency') }}</strong>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Banknote, CreditCard, Store } from 'lucide-vue-next'
import { getCart, clearCart, api, getStoredUser } from '../api'
import { t, currentLang } from '../utils/i18n'

const router = useRouter()
const cartItems = ref([])
const submitting = ref(false)

const user = getStoredUser()

const orderForm = reactive({
  recipient_name: user?.name || '',
  recipient_phone: user?.phone || '',
  shipping_address: 'ш. Душанбе, н. Сино, кӯчаи Борбад 48',
  payment_method: 'cash_on_delivery',
  notes: ''
})

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

const grandTotal = computed(() => {
  return cartItems.value.reduce((sum, it) => sum + (it.price * it.quantity), 0)
})

function formatPrice(val) {
  if (!val) return '0'
  return new Intl.NumberFormat('ru-RU').format(val)
}

async function submitOrder() {
  if (!user) {
    const msg = currentLang.value === 'ru'
      ? 'Пожалуйста, сначала войдите в систему (или нажмите кнопку Покупатель в верхней панели демо).'
      : 'Илтимос аввал ба система ворид шавед (ё тугмаи Харидорро аз панели боло пахш намоед).'
    alert(msg)
    return
  }

  submitting.value = true
  try {
    const payload = {
      items: cartItems.value.map(it => ({
        product_id: it.product_id,
        quantity: it.quantity
      })),
      recipient_name: orderForm.recipient_name,
      recipient_phone: orderForm.recipient_phone,
      shipping_address: orderForm.shipping_address,
      payment_method: orderForm.payment_method,
      notes: orderForm.notes
    }

    const res = await api.orders.create(payload)
    clearCart()
    const successMsg = currentLang.value === 'ru'
      ? `🎉 Поздравляем! Заказ №${res.data.order_number} успешно оформлен и разделен по продавцам!`
      : `🎉 Табрик! Фармоиши №${res.data.order_number} бомуваффақият сабт шуд ва ба қисмҳои фурӯшандагон тақсим гардид!`
    alert(successMsg)
    router.push('/orders')
  } catch (err) {
    alert((currentLang.value === 'ru' ? 'Ошибка оформления: ' : 'Хатогии сабти фармоиш: ') + err.message)
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  cartItems.value = getCart()
})
</script>

<style scoped>
.checkout-page {
  padding: 24px 0 60px 0;
}

.checkout-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-title {
  font-size: 26px;
  font-weight: 800;
}

.back-link {
  color: var(--primary);
  font-weight: 600;
  font-size: 13px;
}

.checkout-layout {
  display: grid;
  grid-template-columns: 1fr 380px;
  gap: 24px;
  align-items: start;
}

@media (max-width: 900px) {
  .checkout-layout {
    grid-template-columns: 1fr;
  }
}

.checkout-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-section {
  padding: 24px;
  background: #ffffff;
}

.section-heading {
  font-size: 16px;
  font-weight: 800;
  margin-bottom: 18px;
  padding-bottom: 10px;
  border-bottom: 1px solid var(--border);
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

@media (max-width: 600px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 14px;
}

.form-group label {
  font-size: 13px;
  font-weight: 700;
  color: var(--text-main);
}

.form-group input, .form-group textarea {
  padding: 10px 14px;
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  font-size: 14px;
  outline: none;
}

.form-group input:focus, .form-group textarea:focus {
  border-color: var(--primary);
}

.payment-methods {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.payment-card {
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  padding: 14px;
  display: flex;
  align-items: center;
  gap: 14px;
  cursor: pointer;
  transition: all 0.15s;
}

.payment-card input {
  display: none;
}

.payment-card.selected {
  border-color: var(--primary);
  background: var(--primary-light);
  color: var(--primary);
}

.payment-card strong {
  font-size: 13px;
  display: block;
}

.payment-card p {
  font-size: 11px;
  color: var(--text-muted);
}

.submit-order-btn {
  width: 100%;
  padding: 14px;
}

/* Правая колонка */
.summary-card {
  padding: 24px;
  background: #ffffff;
}

.summary-card h3 {
  font-size: 16px;
  font-weight: 800;
  margin-bottom: 8px;
}

.split-info-text {
  font-size: 12px;
  color: var(--text-secondary);
  line-height: 1.4;
  margin-bottom: 16px;
  background: var(--primary-light);
  padding: 8px 12px;
  border-radius: 8px;
}

.split-packages-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.split-package-item {
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  padding: 12px;
  background: var(--bg-page);
}

.pkg-header {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: var(--primary);
  margin-bottom: 8px;
}

.pkg-items {
  display: flex;
  flex-direction: column;
  gap: 6px;
  font-size: 12px;
  color: var(--text-secondary);
}

.pkg-item-row {
  display: flex;
  justify-content: space-between;
}

.pkg-footer {
  display: flex;
  justify-content: space-between;
  border-top: 1px dashed var(--border);
  margin-top: 8px;
  padding-top: 6px;
  font-size: 12px;
}

.summary-divider {
  height: 1px;
  background: var(--border);
  margin: 18px 0;
}

.total-row {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
}

.total-sum {
  font-size: 22px;
  font-weight: 800;
  color: var(--text-main);
}

.empty-box {
  padding: 40px;
  text-align: center;
}
</style>
