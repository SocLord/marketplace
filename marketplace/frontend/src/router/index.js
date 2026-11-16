import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import CatalogView from '../views/CatalogView.vue'
import ProductDetailView from '../views/ProductDetailView.vue'
import ShopView from '../views/ShopView.vue'
import ShopsListView from '../views/ShopsListView.vue'
import CartView from '../views/CartView.vue'
import CheckoutView from '../views/CheckoutView.vue'
import BuyerOrdersView from '../views/BuyerOrdersView.vue'
import SellerDashboardView from '../views/SellerDashboardView.vue'
import ModeratorView from '../views/ModeratorView.vue'

const routes = [
  {
    path: '/',
    name: 'home',
    component: HomeView
  },
  {
    path: '/catalog',
    name: 'catalog',
    component: CatalogView
  },
  {
    path: '/product/:idOrSlug',
    name: 'product-detail',
    component: ProductDetailView
  },
  {
    path: '/shops',
    name: 'shops',
    component: ShopsListView
  },
  {
    path: '/shop/:idOrSlug',
    name: 'shop-view',
    component: ShopView
  },
  {
    path: '/cart',
    name: 'cart',
    component: CartView
  },
  {
    path: '/checkout',
    name: 'checkout',
    component: CheckoutView
  },
  {
    path: '/orders',
    name: 'orders',
    component: BuyerOrdersView
  },
  {
    path: '/seller/dashboard',
    name: 'seller-dashboard',
    component: SellerDashboardView
  },
  {
    path: '/moderator',
    name: 'moderator',
    component: ModeratorView
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

export default router
