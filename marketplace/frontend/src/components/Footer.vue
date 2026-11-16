<template>
  <footer class="marketplace-footer">
    <div class="container footer-content">
      <div class="footer-col">
        <div class="footer-logo">
          <div class="logo-box">
            <ShoppingBag :size="20" color="#0f172a" />
          </div>
          <span class="logo-title">БОЗОРГОҲ</span>
        </div>
        <p class="footer-desc">
          {{ currentLang === 'ru' 
            ? 'Современная платформа электронной коммерции с поддержкой продажи товаров несколькими независимыми продавцами (Multi-Vendor Marketplace).' 
            : 'Платформаи муосири тиҷорати электронӣ бо имконияти фурӯши молҳо аз ҷониби якчанд фурӯшандагон (Multi-Vendor Marketplace).' 
          }}
        </p>
        <div class="footer-badges">
          <span class="tech-pill">Golang 1.24</span>
          <span class="tech-pill">Microsoft SQL Server</span>
          <span class="tech-pill">Vue 3 + Vite</span>
        </div>
      </div>

      <div class="footer-col">
        <h4>{{ currentLang === 'ru' ? 'Покупателям' : 'Барои харидорон' }}</h4>
        <ul>
          <li><router-link to="/catalog">{{ t('catalog') }}</router-link></li>
          <li><router-link to="/shops">{{ t('shops') }}</router-link></li>
          <li><router-link to="/cart">{{ t('cart') }}</router-link></li>
          <li><router-link to="/orders">{{ t('myOrders') }}</router-link></li>
        </ul>
      </div>

      <div class="footer-col">
        <h4>{{ currentLang === 'ru' ? 'Продавцам' : 'Барои фурӯшандагон' }}</h4>
        <ul>
          <li><router-link to="/seller/dashboard">{{ t('sellerCabinet') }}</router-link></li>
          <li><router-link to="/seller/dashboard">{{ t('tabStats') }}</router-link></li>
          <li><router-link to="/seller/dashboard">{{ t('tabProducts') }}</router-link></li>
          <li><router-link to="/seller/dashboard">{{ t('tabOrders') }}</router-link></li>
        </ul>
      </div>

      <div class="footer-col">
        <h4>{{ currentLang === 'ru' ? 'Безопасность и правила' : 'Амният ва санҷиш' }}</h4>
        <ul>
          <li><router-link to="/moderator">{{ t('moderationPanel') }}</router-link></li>
          <li><a href="#rules">{{ currentLang === 'ru' ? 'Правила маркетплейса' : 'Қоидаҳои маркетплейс' }}</a></li>
          <li><a href="#guarantee">{{ currentLang === 'ru' ? 'Гарантия возврата средств' : 'Кафолати бозпасгирии маблағ' }}</a></li>
          <li><a href="#support">{{ currentLang === 'ru' ? 'Техподдержка 24/7' : 'Дастгирии техникӣ 24/7' }}</a></li>
        </ul>
      </div>
    </div>

    <div class="footer-bottom">
      <div class="container bottom-inner">
        <span>© 2026 Бозоргоҳ (Marketplace TJ). {{ currentLang === 'ru' ? 'Все права защищены.' : 'Ҳамаи ҳуқуқҳо маҳфузанд.' }}</span>
        <span>Пойгоҳи додаҳо: Microsoft SQL Server (SSMS) | MarketplaceDB</span>
      </div>
    </div>
  </footer>
</template>

<script setup>
import { ShoppingBag } from 'lucide-vue-next'
import { t, currentLang } from '../utils/i18n'
</script>

<style scoped>
.marketplace-footer {
  background: #ffffff;
  border-top: 1px solid var(--border);
  margin-top: 60px;
  padding-top: 48px;
}

.footer-content {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1fr;
  gap: 36px;
  padding-bottom: 40px;
}

@media (max-width: 900px) {
  .footer-content {
    grid-template-columns: 1fr 1fr;
  }
}

@media (max-width: 600px) {
  .footer-content {
    grid-template-columns: 1fr;
  }
}

.footer-logo {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
}

.logo-box {
  width: 34px;
  height: 34px;
  background: #f59e0b;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-title {
  font-size: 19px;
  font-weight: 800;
  color: #0f172a;
}

.footer-desc {
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.6;
  margin-bottom: 16px;
}

.footer-badges {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.tech-pill {
  font-size: 11px;
  font-weight: 700;
  padding: 4px 8px;
  background: #fef3c7;
  border: 1px solid #fde68a;
  border-radius: 6px;
  color: #92400e;
}

.footer-col h4 {
  font-size: 14px;
  font-weight: 800;
  color: var(--text-main);
  margin-bottom: 14px;
}

.footer-col ul {
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.footer-col a {
  font-size: 13px;
  color: var(--text-secondary);
  transition: color 0.15s;
}

.footer-col a:hover {
  color: #d97706;
}

.footer-bottom {
  border-top: 1px solid var(--border-light);
  padding: 18px 0;
  background: var(--bg-page);
  font-size: 12px;
  color: var(--text-muted);
}

.bottom-inner {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
}
</style>
