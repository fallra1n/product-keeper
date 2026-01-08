import axios from 'axios'

// Базовый URL для API (можно изменить в зависимости от окружения)
const API_BASE_URL = 'http://158.160.217.83:8080'

// Создаем экземпляр axios с базовой конфигурацией
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Интерцептор для добавления токена авторизации к запросам
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Интерцептор для обработки ответов и ошибок
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      // Если получили 401, удаляем токен и перенаправляем на логин
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// API методы для аутентификации
export const authAPI = {
  // Регистрация пользователя
  register: async (userData) => {
    try {
      console.log('Отправка запроса регистрации:', userData)
      console.log('URL:', `${API_BASE_URL}/user/register`)
      const response = await api.post('/user/register', userData)
      console.log('Ответ регистрации:', response.data)
      return { success: true, data: response.data }
    } catch (error) {
      console.error('Ошибка регистрации:', error)
      console.error('Статус ошибки:', error.response?.status)
      console.error('Данные ошибки:', error.response?.data)
      return {
        success: false,
        error: error.response?.data?.message || 'Ошибка регистрации'
      }
    }
  },

  // Авторизация пользователя
  login: async (credentials) => {
    try {
      console.log('Отправка запроса авторизации:', credentials)
      console.log('URL:', `${API_BASE_URL}/user/login`)
      const response = await api.post('/user/login', credentials)
      console.log('Ответ авторизации:', response.data)
      return { success: true, data: response.data }
    } catch (error) {
      console.error('Ошибка авторизации:', error)
      console.error('Статус ошибки:', error.response?.status)
      console.error('Данные ошибки:', error.response?.data)
      return {
        success: false,
        error: error.response?.data?.message || 'Ошибка авторизации'
      }
    }
  },
}

// API методы для работы с продуктами
export const productsAPI = {
  // Получить все продукты
  getProducts: async (params = {}) => {
    try {
      const response = await api.get('/products', { params })
      return { success: true, data: response.data }
    } catch (error) {
      return { 
        success: false, 
        error: error.response?.data?.message || 'Ошибка получения продуктов' 
      }
    }
  },

  // Получить продукт по ID
  getProduct: async (id) => {
    try {
      const response = await api.get(`/product/${id}`)
      return { success: true, data: response.data }
    } catch (error) {
      return { 
        success: false, 
        error: error.response?.data?.message || 'Ошибка получения продукта' 
      }
    }
  },

  // Создать новый продукт
  createProduct: async (productData) => {
    try {
      const response = await api.post('/product/add', productData)
      return { success: true, data: response.data }
    } catch (error) {
      return { 
        success: false, 
        error: error.response?.data?.message || 'Ошибка создания продукта' 
      }
    }
  },

  // Обновить продукт
  updateProduct: async (id, productData) => {
    try {
      const response = await api.put(`/product/${id}`, productData)
      return { success: true, data: response.data }
    } catch (error) {
      return { 
        success: false, 
        error: error.response?.data?.message || 'Ошибка обновления продукта' 
      }
    }
  },

  // Удалить продукт
  deleteProduct: async (id) => {
    try {
      const response = await api.delete(`/product/${id}`)
      return { success: true, data: response.data }
    } catch (error) {
      return { 
        success: false, 
        error: error.response?.data?.message || 'Ошибка удаления продукта' 
      }
    }
  },
}

// Утилитарные функции
export const apiUtils = {
  // Проверка валидности токена
  isTokenValid: () => {
    const token = localStorage.getItem('token')
    if (!token) return false
    
    try {
      // Простая проверка структуры JWT токена
      const payload = JSON.parse(atob(token.split('.')[1]))
      const currentTime = Date.now() / 1000
      return payload.exp > currentTime
    } catch (error) {
      return false
    }
  },

  // Получение данных пользователя из токена
  getUserFromToken: () => {
    const token = localStorage.getItem('token')
    if (!token) return null
    
    try {
      const payload = JSON.parse(atob(token.split('.')[1]))
      return {
        username: payload.username,
        exp: payload.exp
      }
    } catch (error) {
      return null
    }
  },

  // Форматирование ошибок для отображения
  formatError: (error) => {
    if (typeof error === 'string') return error
    if (error?.message) return error.message
    return 'Произошла неизвестная ошибка'
  }
}

export default api