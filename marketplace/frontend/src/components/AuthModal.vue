<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="modal-content auth-modal">
      <button class="close-btn" @click="$emit('close')">
        <X :size="20" />
      </button>

      <!-- Заголовок -->
      <div class="auth-header">
        <div class="auth-icon">
          <ShoppingBag :size="28" color="#0f172a" />
        </div>
        <h2>{{ currentLang === 'ru' ? 'Добро пожаловать в Бозоргоҳ' : 'Хуш омадед ба Бозоргоҳ' }}</h2>
        <p>{{ currentLang === 'ru' ? 'Онлайн-маркетплейс товаров от лучших продавцов' : 'Маркетплейси ягонаи молҳо ва фурӯшандагони Тоҷикистон' }}</p>
      </div>

      <!-- Вкладки: Вход / Регистрация -->
      <div class="tabs">
        <button 
          :class="['tab-btn', activeTab === 'login' && 'active']"
          @click="activeTab = 'login'">
          {{ currentLang === 'ru' ? 'Вход' : 'Воридшавӣ' }}
        </button>
        <button 
          :class="['tab-btn', activeTab === 'register' && 'active']"
          @click="activeTab = 'register'">
          {{ currentLang === 'ru' ? 'Регистрация' : 'Сабти ном' }}
        </button>
      </div>

      <!-- Сообщение об ошибке -->
      <div v-if="errorMsg" class="alert-error">
        {{ errorMsg }}
      </div>

      <!-- Форма Входа -->
      <form v-if="activeTab === 'login'" @submit.prevent="handleLogin" class="auth-form">
        <div class="form-group">
          <label>{{ currentLang === 'ru' ? 'Электронная почта (Email)' : 'Почтаи электронӣ (Email)' }}</label>
          <input 
            type="email" 
            v-model="loginForm.email" 
            required 
            placeholder="buyer@market.tj" 
          />
        </div>

        <div class="form-group">
          <label>{{ currentLang === 'ru' ? 'Пароль' : 'Парол' }}</label>
          <input 
            type="password" 
            v-model="loginForm.password" 
            required 
            placeholder="123456" 
          />
        </div>

        <button type="submit" class="btn btn-primary btn-submit" :disabled="loading">
          {{ loading ? (currentLang === 'ru' ? 'Вход...' : 'Интизор шавед...') : (currentLang === 'ru' ? 'Войти в аккаунт' : 'Ворид шудан') }}
        </button>
      </form>

      <!-- Форма Регистрации -->
      <form v-else @submit.prevent="handleRegister" class="auth-form">
        <div class="form-group">
          <label>{{ currentLang === 'ru' ? 'Ваше ФИО' : 'Ному насаб' }}</label>
          <input 
            type="text" 
            v-model="regForm.name" 
            required 
            placeholder="Дильшод Рахимов" 
          />
        </div>

        <div class="form-group">
          <label>Email</label>
          <input 
            type="email" 
            v-model="regForm.email" 
            required 
            placeholder="name@market.tj" 
          />
        </div>

        <div class="form-group">
          <label>{{ currentLang === 'ru' ? 'Номер телефона' : 'Рақами телефон' }}</label>
          <input 
            type="tel" 
            v-model="regForm.phone" 
            placeholder="+992 90 000 0000" 
          />
        </div>

        <div class="form-group">
          <label>{{ currentLang === 'ru' ? 'Ваша роль:' : 'Нақши шумо дар маркетплейс:' }}</label>
          <div class="role-selector">
            <label class="role-card" :class="regForm.role === 'buyer' && 'selected'">
              <input type="radio" value="buyer" v-model="regForm.role" />
              <User :size="20" />
              <div>
                <strong>{{ currentLang === 'ru' ? 'Покупатель' : 'Харидор' }}</strong>
                <p>{{ currentLang === 'ru' ? 'Поиск и покупка товаров' : 'Мехоҳам молҳоро ҷустуҷӯ ва харидорӣ намоям' }}</p>
              </div>
            </label>

            <label class="role-card" :class="regForm.role === 'seller' && 'selected'">
              <input type="radio" value="seller" v-model="regForm.role" />
              <Store :size="20" />
              <div>
                <strong>{{ currentLang === 'ru' ? 'Продавец' : 'Фурӯшанда' }}</strong>
                <p>{{ currentLang === 'ru' ? 'Управление магазином и товарами' : 'Мехоҳам мағоза кушоям ва мол фурӯшам' }}</p>
              </div>
            </label>
          </div>
        </div>

        <div v-if="regForm.role === 'seller'" class="form-group">
          <label>{{ currentLang === 'ru' ? 'Название вашего магазина' : 'Номи мағозаи шумо' }}</label>
          <input 
            type="text" 
            v-model="regForm.shop_name" 
            placeholder="SuperMart Dushanbe" 
          />
        </div>

        <div class="form-group">
          <label>{{ currentLang === 'ru' ? 'Пароль' : 'Парол' }}</label>
          <input 
            type="password" 
            v-model="regForm.password" 
            required 
            placeholder="Мин. 6 символов" 
          />
        </div>

        <button type="submit" class="btn btn-primary btn-submit" :disabled="loading">
          {{ loading ? (currentLang === 'ru' ? 'Регистрация...' : 'Сабт мешавад...') : (currentLang === 'ru' ? 'Зарегистрироваться' : 'Сабти ном ва кушодани ҳисоб') }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { X, ShoppingBag, User, Store, ShieldCheck } from 'lucide-vue-next'
import { api, setToken, setStoredUser } from '../api'
import { t, currentLang } from '../utils/i18n'

const emit = defineEmits(['close'])
const router = useRouter()

const activeTab = ref('login')
const loading = ref(false)
const errorMsg = ref('')

const loginForm = ref({
  email: '',
  password: ''
})

const regForm = ref({
  name: '',
  email: '',
  phone: '',
  role: 'buyer',
  shop_name: '',
  password: ''
})

function fillQuick(email, pass) {
  activeTab.value = 'login'
  loginForm.value.email = email
  loginForm.value.password = pass
  handleLogin()
}

async function handleLogin() {
  loading.value = true
  errorMsg.value = ''
  try {
    const res = await api.auth.login(loginForm.value)
    setToken(res.data.token)
    setStoredUser(res.data.user)
    window.dispatchEvent(new Event('user-updated'))
    emit('close')

    if (res.data.user.role === 'seller') {
      router.push('/seller/dashboard')
    } else if (res.data.user.role === 'moderator') {
      router.push('/moderator')
    }
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    loading.value = false
  }
}

async function handleRegister() {
  loading.value = true
  errorMsg.value = ''
  try {
    const res = await api.auth.register(regForm.value)
    setToken(res.data.token)
    setStoredUser(res.data.user)
    window.dispatchEvent(new Event('user-updated'))
    emit('close')

    if (res.data.user.role === 'seller') {
      router.push('/seller/dashboard')
    }
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-modal {
  position: relative;
  max-width: 480px;
}

.close-btn {
  position: absolute;
  top: 18px;
  right: 18px;
  color: var(--text-muted);
  padding: 4px;
  border-radius: 6px;
  transition: all 0.15s;
}

.close-btn:hover {
  background: var(--bg-page);
  color: var(--text-main);
}

.auth-header {
  text-align: center;
  margin-bottom: 20px;
}

.auth-icon {
  width: 52px;
  height: 52px;
  background: #fef3c7;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 12px auto;
}

.auth-header h2 {
  font-size: 20px;
  font-weight: 800;
  color: var(--text-main);
}

.auth-header p {
  font-size: 13px;
  color: var(--text-secondary);
  margin-top: 4px;
}

.tabs {
  display: flex;
  background: var(--bg-page);
  padding: 4px;
  border-radius: var(--radius-md);
  margin-bottom: 20px;
}

.tab-btn {
  flex: 1;
  padding: 8px 12px;
  font-size: 13px;
  font-weight: 700;
  border-radius: 8px;
  color: var(--text-secondary);
  transition: all 0.2s;
}

.tab-btn.active {
  background: #f59e0b;
  color: #0f172a;
  box-shadow: var(--shadow-sm);
}

.alert-error {
  background: #fee2e2;
  color: #b91c1c;
  padding: 10px 14px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  margin-bottom: 16px;
}

.auth-form {
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
  color: var(--text-main);
}

.form-group input {
  padding: 10px 14px;
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  font-size: 14px;
  outline: none;
  transition: border-color 0.15s;
}

.form-group input:focus {
  border-color: #f59e0b;
}

.role-selector {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.role-card {
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  padding: 12px;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  gap: 8px;
  transition: all 0.15s;
}

.role-card input {
  display: none;
}

.role-card.selected {
  border-color: #f59e0b;
  background: #fef3c7;
  color: #92400e;
}

.role-card strong {
  font-size: 13px;
  display: block;
}

.role-card p {
  font-size: 11px;
  color: var(--text-muted);
  line-height: 1.25;
}

.btn-submit {
  width: 100%;
  margin-top: 6px;
  padding: 12px;
}

.quick-demo-section {
  margin-top: 20px;
}

.demo-divider {
  text-align: center;
  position: relative;
  margin-bottom: 12px;
}

.demo-divider::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  width: 100%;
  height: 1px;
  background: var(--border);
}

.demo-divider span {
  position: relative;
  background: #ffffff;
  padding: 0 10px;
  font-size: 11px;
  color: var(--text-muted);
  font-weight: 600;
}

.demo-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}

.demo-pick {
  display: flex;
  align-items: center;
  gap: 6px;
  background: var(--bg-page);
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 8px 10px;
  font-size: 11px;
  font-weight: 700;
  color: var(--text-secondary);
  transition: all 0.15s;
}

.demo-pick:hover {
  background: #fef3c7;
  border-color: #f59e0b;
  color: #0f172a;
}
</style>
