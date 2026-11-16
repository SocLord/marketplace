<template>
  <div id="marketplace-app">
    <Navbar @open-auth="showAuthModal = true" />
    <main class="app-main">
      <router-view />
    </main>
    <Footer />

    <!-- Модальное окно авторизации -->
    <AuthModal v-if="showAuthModal" @close="showAuthModal = false" />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import Navbar from './components/Navbar.vue'
import Footer from './components/Footer.vue'
import AuthModal from './components/AuthModal.vue'

const showAuthModal = ref(false)
</script>

<style scoped>
#marketplace-app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.app-main {
  flex: 1;
}
</style>
