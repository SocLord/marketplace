<template>
  <header class="marketplace-header">
    <!-- Верхняя информационная полоса: статус / доставка + Переключатель языка -->
    <div class="top-info-bar">
      <div class="container top-inner">
        <div class="top-features">
          <span class="feature-item">
            <Sparkles :size="13" color="#fbbf24" />
            <span>{{ currentLang === 'ru' ? 'Онлайн-маркетплейс товаров в Таджикистане' : 'Бозоргоҳи калонтарини молҳои Тоҷикистон' }}</span>
          </span>
          <span class="feature-item hidden-mobile">🚀 {{ currentLang === 'ru' ? 'Быстрая доставка' : 'Расонидани зуд' }}</span>
          <span class="feature-item hidden-mobile">🔒 {{ currentLang === 'ru' ? 'Гарантия подлинности' : 'Кафолати сифат ва аслӣ' }}</span>
        </div>

        <!-- Переключатель языка (TG / RU) -->
        <div class="lang-switch-box">
          <Languages :size="14" />
          <button 
            :class="['lang-btn', currentLang === 'ru' && 'active']" 
            @click="changeLang('ru')">
            РУС
          </button>
          <span class="lang-div">|</span>
          <button 
            :class="['lang-btn', currentLang === 'tg' && 'active']" 
            @click="changeLang('tg')">
            ТОҶ
          </button>
        </div>
      </div>
    </div>

    <!-- Асосии шапка (Main Header) -->
    <div class="main-header">
      <div class="container header-container">
        <!-- Логотип -->
        <router-link to="/" class="logo">
          <div class="logo-icon">
            <ShoppingBag :size="24" color="#0f172a" />
          </div>
          <div class="logo-text">
            <span class="logo-title">БОЗОРГОҲ</span>
            <span class="logo-subtitle">MARKETPLACE</span>
          </div>
        </router-link>

        <!-- Тугмаи Каталог (желтая акцентная кнопка) -->
        <router-link to="/catalog" class="catalog-btn">
          <LayoutGrid :size="18" />
          <span>{{ t('catalog') }}</span>
        </router-link>

        <!-- Строка поиска -->
        <div class="search-box">
          <input 
            type="text" 
            v-model="searchQuery" 
            @keyup.enter="handleSearch"
            :placeholder="t('searchPlaceholder')" 
          />
          <button @click="handleSearch" class="search-btn" aria-label="Ҷустуҷӯ">
            <Search :size="18" />
          </button>
        </div>

        <!-- Навигационные кнопки -->
        <div class="header-actions">
          <!-- Мағозаҳо -->
          <router-link to="/shops" class="action-btn">
            <Store :size="20" />
            <span>{{ t('shops') }}</span>
          </router-link>

          <!-- Кабинеты по ролям -->
          <template v-if="currentUser">
            <!-- Если продавец -->
            <router-link v-if="currentUser.role === 'seller'" to="/seller/dashboard" class="action-btn highlight-seller">
              <BarChart3 :size="20" />
              <span>{{ t('sellerCabinet') }}</span>
            </router-link>

            <!-- Если модератор -->
            <router-link v-if="currentUser.role === 'moderator'" to="/moderator" class="action-btn highlight-mod">
              <ShieldCheck :size="20" />
              <span>{{ t('moderationPanel') }}</span>
            </router-link>

            <!-- Если покупатель: Фармоишҳои ман -->
            <router-link v-if="currentUser.role === 'buyer'" to="/orders" class="action-btn">
              <PackageCheck :size="20" />
              <span>{{ t('myOrders') }}</span>
            </router-link>
          </template>

          <!-- Корзина -->
          <router-link to="/cart" class="action-btn cart-btn">
            <div class="cart-icon-wrapper">
              <ShoppingCart :size="20" />
              <span v-if="cartCount > 0" class="cart-badge">{{ cartCount }}</span>
            </div>
            <span>{{ t('cart') }}</span>
          </router-link>

          <!-- Профиль / Вход -->
          <div class="user-profile-menu">
            <template v-if="currentUser">
              <div class="user-chip" @click="showUserMenu = !showUserMenu">
                <img :src="currentUser.avatar_url || 'https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?w=100'" class="user-avatar" alt="Avatar" />
                <span class="user-name">{{ currentUser.name.split(' ')[0] }}</span>
                <ChevronDown :size="14" />
              </div>

              <!-- Выпадающее меню -->
              <div v-if="showUserMenu" class="dropdown-menu">
                <div class="dropdown-header">
                  <strong>{{ currentUser.name }}</strong>
                  <span class="role-tag">{{ getRoleLabel(currentUser.role) }}</span>
                </div>
                <div class="dropdown-divider"></div>

                <router-link v-if="currentUser.role === 'seller'" to="/seller/dashboard" class="dropdown-item" @click="showUserMenu = false">
                  <BarChart3 :size="16" />
                  <span>{{ t('tabStats') }}</span>
                </router-link>
                <router-link v-if="currentUser.role === 'moderator'" to="/moderator" class="dropdown-item" @click="showUserMenu = false">
                  <ShieldCheck :size="16" />
                  <span>{{ t('moderationPanel') }}</span>
                </router-link>
                <router-link to="/orders" class="dropdown-item" @click="showUserMenu = false">
                  <Package :size="16" />
                  <span>{{ t('myOrders') }}</span>
                </router-link>

                <div class="dropdown-divider"></div>
                <button @click="logout" class="dropdown-item text-danger">
                  <LogOut :size="16" />
                  <span>{{ t('logout') }}</span>
                </button>
              </div>
            </template>

            <button v-else @click="$emit('open-auth')" class="login-btn">
              <User :size="18" />
              <span>{{ t('login') }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Полоса быстрых категорий -->
    <div class="category-strip">
      <div class="container strip-container">
        <router-link to="/catalog" class="strip-item all-link">
          <Grid :size="14" />
          <span>{{ t('allProducts') }}</span>
        </router-link>
        <router-link 
          v-for="cat in categories" 
          :key="cat.id" 
          :to="`/catalog?category=${cat.slug}`"
          class="strip-item">
          <span>{{ cat.name }}</span>
        </router-link>
      </div>
    </div>
  </header>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { 
  ShoppingBag, Search, LayoutGrid, ShoppingCart, 
  Store, User, ShieldCheck, BarChart3, PackageCheck, 
  Package, ChevronDown, LogOut, Grid, Sparkles, Languages 
} from 'lucide-vue-next'
import { api, getStoredUser, setStoredUser, setToken, getCart } from '../api'
import { t, currentLang, setLanguage } from '../utils/i18n'

const router = useRouter()
const currentUser = ref(getStoredUser())
const cartCount = ref(0)
const searchQuery = ref('')
const showUserMenu = ref(false)
const categories = ref([])

function changeLang(lang) {
  setLanguage(lang)
}

function updateCartCount() {
  const cart = getCart()
  cartCount.value = cart.reduce((sum, item) => sum + item.quantity, 0)
}

function handleSearch() {
  if (searchQuery.value.trim()) {
    router.push(`/catalog?q=${encodeURIComponent(searchQuery.value.trim())}`)
  } else {
    router.push('/catalog')
  }
}

async function quickLogin(email, password) {
  try {
    const res = await api.auth.login({ email, password })
    setToken(res.data.token)
    setStoredUser(res.data.user)
    currentUser.value = res.data.user
    showUserMenu.value = false
    
    if (res.data.user.role === 'seller') {
      router.push('/seller/dashboard')
    } else if (res.data.user.role === 'moderator') {
      router.push('/moderator')
    } else {
      router.push('/')
    }
  } catch (e) {
    alert('Хатогии воридшавӣ: ' + e.message)
  }
}

function logout() {
  api.auth.logout()
  currentUser.value = null
  showUserMenu.value = false
  router.push('/')
}

function getRoleLabel(role) {
  if (currentLang.value === 'ru') {
    switch (role) {
      case 'seller': return 'Продавец'
      case 'moderator': return 'Модератор'
      case 'admin': return 'Администратор'
      default: return 'Покупатель'
    }
  }
  switch (role) {
    case 'seller': return 'Фурӯшанда'
    case 'moderator': return 'Модератор'
    case 'admin': return 'Администратор'
    default: return 'Харидор'
  }
}

async function loadCategories() {
  try {
    const res = await api.categories.list()
    categories.value = res.data || []
  } catch (e) {
    console.error(e)
  }
}

function handleUserUpdate() {
  currentUser.value = getStoredUser()
}

onMounted(() => {
  updateCartCount()
  loadCategories()
  window.addEventListener('cart-updated', updateCartCount)
  window.addEventListener('user-updated', handleUserUpdate)
})

onUnmounted(() => {
  window.removeEventListener('cart-updated', updateCartCount)
  window.removeEventListener('user-updated', handleUserUpdate)
})
</script>

<style scoped>
.marketplace-header {
  background: #ffffff;
  border-bottom: 1px solid var(--border);
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.04);
}

/* Верхняя полоса */
.top-info-bar {
  background: #0f172a;
  color: #e2e8f0;
  font-size: 12px;
  padding: 8px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.top-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
}

.top-features {
  display: flex;
  align-items: center;
  gap: 20px;
  font-size: 12px;
  color: #94a3b8;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.feature-item:first-child {
  color: #f1f5f9;
  font-weight: 600;
}

@media (max-width: 768px) {
  .hidden-mobile {
    display: none;
  }
}

/* Переключатель языка */
.lang-switch-box {
  display: flex;
  align-items: center;
  gap: 6px;
  background: rgba(255, 255, 255, 0.12);
  padding: 3px 10px;
  border-radius: var(--radius-full);
  color: #fbbf24;
}

.lang-btn {
  color: #94a3b8;
  font-size: 11px;
  font-weight: 700;
  padding: 2px 4px;
  border-radius: 4px;
  transition: all 0.15s;
}

.lang-btn:hover {
  color: #ffffff;
}

.lang-btn.active {
  color: #0f172a;
  background: #fbbf24;
  padding: 2px 6px;
}

.lang-div {
  color: rgba(255, 255, 255, 0.3);
  font-size: 10px;
}

/* Главная шапка */
.main-header {
  padding: 14px 0;
}

.header-container {
  display: flex;
  align-items: center;
  gap: 20px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  text-decoration: none;
}

.logo-icon {
  width: 42px;
  height: 42px;
  background: linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(245, 158, 11, 0.35);
}

.logo-text {
  display: flex;
  flex-direction: column;
}

.logo-title {
  font-size: 19px;
  font-weight: 800;
  letter-spacing: 0.5px;
  color: #0f172a;
  line-height: 1.1;
}

.logo-subtitle {
  font-size: 10px;
  font-weight: 800;
  letter-spacing: 1.5px;
  color: #d97706;
}

/* Кнопка каталога (Желтая кнопка с темным текстом в стиле современных маркетплейсов) */
.catalog-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #f59e0b;
  color: #0f172a;
  padding: 11px 20px;
  border-radius: var(--radius-md);
  font-weight: 800;
  font-size: 14px;
  box-shadow: 0 2px 8px rgba(245, 158, 11, 0.3);
  transition: all 0.2s ease;
}

.catalog-btn:hover {
  background: #d97706;
  color: #ffffff;
  transform: translateY(-1px);
}

/* Поиск */
.search-box {
  flex: 1;
  display: flex;
  position: relative;
}

.search-box input {
  width: 100%;
  padding: 11px 45px 11px 16px;
  border: 2px solid var(--border);
  border-radius: var(--radius-md);
  font-size: 14px;
  outline: none;
  transition: all 0.2s ease;
  background: #f8fafc;
}

.search-box input:focus {
  border-color: #f59e0b;
  background: #ffffff;
  box-shadow: 0 0 0 3px rgba(245, 158, 11, 0.2);
}

.search-btn {
  position: absolute;
  right: 6px;
  top: 50%;
  transform: translateY(-50%);
  width: 34px;
  height: 34px;
  background: #f59e0b;
  color: #0f172a;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.2s;
}

.search-btn:hover {
  background: #d97706;
  color: #ffffff;
}

/* Действия */
.header-actions {
  display: flex;
  align-items: center;
  gap: 14px;
}

.action-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 3px;
  color: var(--text-secondary);
  font-size: 11px;
  font-weight: 700;
  padding: 6px 8px;
  border-radius: 8px;
  transition: all 0.2s;
}

.action-btn:hover {
  color: #b45309;
  background: #fffbeb;
}

.highlight-seller {
  color: #b45309;
  background: #fef3c7;
}

.highlight-mod {
  color: #6d28d9;
  background: #f5f3ff;
}

.cart-icon-wrapper {
  position: relative;
}

.cart-badge {
  position: absolute;
  top: -8px;
  right: -10px;
  background: #ef4444;
  color: #ffffff;
  font-size: 11px;
  font-weight: 800;
  padding: 2px 6px;
  border-radius: 10px;
}

.login-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: #fffbeb;
  border: 1px solid #fde68a;
  border-radius: var(--radius-md);
  font-size: 13px;
  font-weight: 700;
  color: #92400e;
  transition: all 0.2s;
}

.login-btn:hover {
  background: #fef3c7;
  border-color: #f59e0b;
  color: #0f172a;
}

.user-profile-menu {
  position: relative;
}

.user-chip {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 10px 4px 4px;
  background: #fef3c7;
  border: 1px solid #fde68a;
  border-radius: 30px;
  cursor: pointer;
  transition: all 0.2s;
}

.user-chip:hover {
  background: #fde68a;
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  object-fit: cover;
}

.user-name {
  font-size: 13px;
  font-weight: 700;
  color: #0f172a;
}

/* Выпадающее меню */
.dropdown-menu {
  position: absolute;
  right: 0;
  top: calc(100% + 8px);
  width: 230px;
  background: #ffffff;
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-xl);
  padding: 8px;
  z-index: 200;
}

.dropdown-header {
  padding: 8px 12px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.role-tag {
  font-size: 11px;
  background: #fef3c7;
  color: #b45309;
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 700;
  width: fit-content;
}

.dropdown-divider {
  height: 1px;
  background: var(--border);
  margin: 6px 0;
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  font-size: 13px;
  font-weight: 600;
  border-radius: 6px;
  width: 100%;
  color: var(--text-main);
  text-align: left;
  transition: background 0.15s;
}

.dropdown-item:hover {
  background: #fffbeb;
  color: #b45309;
}

.text-danger {
  color: var(--danger) !important;
}

/* Полоса категорий */
.category-strip {
  background: #ffffff;
  border-top: 1px solid var(--border-light);
  overflow-x: auto;
  white-space: nowrap;
}

.strip-container {
  display: flex;
  align-items: center;
  gap: 18px;
  padding: 8px 20px;
}

.strip-item {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-secondary);
  transition: color 0.15s;
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.strip-item:hover {
  color: #d97706;
}

.all-link {
  color: #b45309;
  font-weight: 800;
}
</style>
