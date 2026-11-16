const API_BASE = '/api'

export function getToken() {
  return localStorage.getItem('market_token') || ''
}

export function setToken(token) {
  if (token) {
    localStorage.setItem('market_token', token)
  } else {
    localStorage.removeItem('market_token')
  }
}

export function getStoredUser() {
  const u = localStorage.getItem('market_user')
  try {
    return u ? JSON.parse(u) : null
  } catch (e) {
    return null
  }
}

export function setStoredUser(user) {
  if (user) {
    localStorage.setItem('market_user', JSON.stringify(user))
  } else {
    localStorage.removeItem('market_user')
  }
}


export function getCart() {
  try {
    return JSON.parse(localStorage.getItem('market_cart') || '[]')
  } catch (e) {
    return []
  }
}

export function saveCart(items) {
  localStorage.setItem('market_cart', JSON.stringify(items))
  window.dispatchEvent(new Event('cart-updated'))
}

export function addToCart(product, quantity = 1) {
  const cart = getCart()
  const existing = cart.find(item => item.product_id === product.id)
  if (existing) {
    existing.quantity += quantity
  } else {
    cart.push({
      product_id: product.id,
      name: product.name,
      price: product.price,
      old_price: product.old_price,
      image_url: product.image_url,
      shop_id: product.shop_id,
      shop_name: product.shop?.name || `Мағозаи #${product.shop_id}`,
      stock_quantity: product.stock_quantity,
      quantity: quantity
    })
  }
  saveCart(cart)
}

export function updateCartQuantity(productId, quantity) {
  let cart = getCart()
  if (quantity <= 0) {
    cart = cart.filter(i => i.product_id !== productId)
  } else {
    const item = cart.find(i => i.product_id === productId)
    if (item) item.quantity = quantity
  }
  saveCart(cart)
}

export function removeFromCart(productId) {
  const cart = getCart().filter(i => i.product_id !== productId)
  saveCart(cart)
}

export function clearCart() {
  saveCart([])
}

async function request(endpoint, options = {}) {
  const url = `${API_BASE}${endpoint}`
  const headers = {
    'Accept': 'application/json',
    ...(options.headers || {})
  }

  if (!(options.body instanceof FormData)) {
    headers['Content-Type'] = 'application/json'
  }

  const token = getToken()
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const response = await fetch(url, {
    ...options,
    headers
  })

  let data
  try {
    data = await response.json()
  } catch (err) {
    data = { error: `HTTP ${response.status}: ${response.statusText}` }
  }

  if (!response.ok) {
    const errorMsg = data?.error || data?.message || `Хатогии сервер (${response.status})`
    throw new Error(errorMsg)
  }

  return data
}

export const api = {
 
  auth: {
    login: (creds) => request('/auth/login', {
      method: 'POST',
      body: JSON.stringify(creds)
    }),
    register: (userData) => request('/auth/register', {
      method: 'POST',
      body: JSON.stringify(userData)
    }),
    getMe: () => request('/auth/me', { method: 'GET' }),
    logout: () => {
      setToken('')
      setStoredUser(null)
    }
  },


  products: {
    list: (params = {}) => {
      const q = new URLSearchParams()
      for (const [k, v] of Object.entries(params)) {
        if (v !== undefined && v !== null && v !== '') {
          q.append(k, v)
        }
      }
      const qs = q.toString() ? `?${q.toString()}` : ''
      return request(`/products${qs}`, { method: 'GET' })
    },
    get: (idOrSlug) => request(`/products/${idOrSlug}`, { method: 'GET' })
  },


  categories: {
    list: () => request('/categories', { method: 'GET' })
  },

 
  shops: {
    list: () => request('/shops', { method: 'GET' }),
    get: (idOrSlug) => request(`/shops/${idOrSlug}`, { method: 'GET' })
  },


  orders: {
    create: (orderData) => request('/orders', {
      method: 'POST',
      body: JSON.stringify(orderData)
    }),
    getMyOrders: () => request('/orders/my', { method: 'GET' })
  },

  reviews: {
    getByProduct: (productId) => request(`/reviews/${productId}`, { method: 'GET' }),
    create: (data) => request('/reviews', {
      method: 'POST',
      body: JSON.stringify(data)
    })
  },

  seller: {

    getProducts: () => request('/seller/products', { method: 'GET' }),
    createProduct: (data) => request('/seller/products', {
      method: 'POST',
      body: JSON.stringify(data)
    }),
    updateProduct: (id, data) => request(`/seller/products/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data)
    }),
    deleteProduct: (id) => request(`/seller/products/${id}`, { method: 'DELETE' }),

    getOrders: (status = 'all') => request(`/seller/orders?status=${status}`, { method: 'GET' }),
    updateOrderStatus: (subOrderId, data) => request(`/seller/orders/${subOrderId}/status`, {
      method: 'PUT',
      body: JSON.stringify(data)
    }),

    getMonthlyStats: () => request('/seller/stats/monthly', { method: 'GET' }),

    getMyShop: () => request('/seller/my-shop', { method: 'GET' }),
    updateMyShop: (data) => request('/seller/my-shop', {
      method: 'PUT',
      body: JSON.stringify(data)
    })
  },

  moderator: {
    getPendingProducts: () => request('/moderator/products/pending', { method: 'GET' }),
    getAllProducts: (status = 'all') => request(`/moderator/products?status=${status}`, { method: 'GET' }),
    reviewProduct: (id, decision) => request(`/moderator/products/${id}/review`, {
      method: 'PUT',
      body: JSON.stringify(decision)
    }),
    createCategory: (data) => request('/moderator/categories', {
      method: 'POST',
      body: JSON.stringify(data)
    }),
    updateCategory: (id, data) => request(`/moderator/categories/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data)
    }),
    deleteCategory: (id) => request(`/moderator/categories/${id}`, { method: 'DELETE' }),
    getStats: () => request('/moderator/stats', { method: 'GET' })
  },

  upload: (file) => {
    const fd = new FormData()
    fd.append('file', file)
    return request('/upload', {
      method: 'POST',
      body: fd
    })
  }
}
