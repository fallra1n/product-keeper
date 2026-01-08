import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { Package, Plus, TrendingUp, Users, ShoppingCart, BarChart3 } from 'lucide-react'
import { productsAPI, apiUtils } from '../../services/api'

const Dashboard = () => {
  const [stats, setStats] = useState({
    totalProducts: 0,
    recentProducts: []
  })
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  
  const user = apiUtils.getUserFromToken()

  useEffect(() => {
    loadDashboardData()
  }, [])

  const loadDashboardData = async () => {
    setLoading(true)
    try {
      const result = await productsAPI.getProducts({ sort_by: 'last_create' })
      
      if (result.success) {
        const products = result.data || []
        setStats({
          totalProducts: products.length,
          recentProducts: products.slice(0, 5) // Последние 5 продуктов
        })
      } else {
        setError(result.error)
      }
    } catch (error) {
      setError('Ошибка загрузки данных')
    } finally {
      setLoading(false)
    }
  }

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('ru-RU', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  const formatPrice = (price) => {
    return new Intl.NumberFormat('ru-RU', {
      style: 'currency',
      currency: 'RUB'
    }).format(price)
  }

  if (loading) {
    return (
      <div className="dashboard-loading">
        <div className="loading-spinner"></div>
        <p>Загрузка данных...</p>
      </div>
    )
  }

  return (
    <div className="dashboard fade-in">
      <div className="container">
        {/* Заголовок */}
        <div className="dashboard-header">
          <div>
            <h1>Добро пожаловать, {user?.username}!</h1>
            <p className="dashboard-subtitle">
              Управляйте своими продуктами эффективно
            </p>
          </div>
          <Link to="/products/add" className="btn btn-primary">
            <Plus size={20} />
            Добавить продукт
          </Link>
        </div>

        {error && (
          <div className="error-message mb-6">
            {error}
          </div>
        )}

        {/* Статистические карточки */}
        <div className="stats-grid">
          <div className="stat-card">
            <div className="stat-icon">
              <Package size={24} />
            </div>
            <div className="stat-content">
              <h3>{stats.totalProducts}</h3>
              <p>Всего продуктов</p>
            </div>
          </div>

          <div className="stat-card">
            <div className="stat-icon success">
              <TrendingUp size={24} />
            </div>
            <div className="stat-content">
              <h3>
                {stats.recentProducts.reduce((sum, product) => sum + (product.quantity || 0), 0)}
              </h3>
              <p>Общее количество</p>
            </div>
          </div>

          <div className="stat-card">
            <div className="stat-icon warning">
              <ShoppingCart size={24} />
            </div>
            <div className="stat-content">
              <h3>
                {formatPrice(stats.recentProducts.reduce((sum, product) => sum + (product.price * product.quantity || 0), 0))}
              </h3>
              <p>Общая стоимость</p>
            </div>
          </div>

          <div className="stat-card">
            <div className="stat-icon info">
              <BarChart3 size={24} />
            </div>
            <div className="stat-content">
              <h3>
                {stats.totalProducts > 0 ? 
                  formatPrice(stats.recentProducts.reduce((sum, product) => sum + (product.price || 0), 0) / stats.totalProducts) : 
                  formatPrice(0)
                }
              </h3>
              <p>Средняя цена</p>
            </div>
          </div>
        </div>

        {/* Последние продукты */}
        <div className="dashboard-section">
          <div className="section-header">
            <h2>Последние продукты</h2>
            <Link to="/products" className="btn btn-secondary">
              Посмотреть все
            </Link>
          </div>

          {stats.recentProducts.length > 0 ? (
            <div className="recent-products">
              {stats.recentProducts.map((product) => (
                <div key={product.id} className="product-card">
                  <div className="product-info">
                    <h3 className="product-name">{product.name}</h3>
                    <p className="product-details">
                      Цена: {formatPrice(product.price)} • 
                      Количество: {product.quantity} шт.
                    </p>
                    <p className="product-date">
                      Создано: {formatDate(product.created_at)}
                    </p>
                  </div>
                  <div className="product-actions">
                    <Link 
                      to={`/products/edit/${product.id}`}
                      className="btn btn-secondary btn-sm"
                    >
                      Редактировать
                    </Link>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="empty-state">
              <Package size={48} />
              <h3>Пока нет продуктов</h3>
              <p>Создайте свой первый продукт, чтобы начать работу</p>
              <Link to="/products/add" className="btn btn-primary">
                <Plus size={20} />
                Добавить продукт
              </Link>
            </div>
          )}
        </div>

        {/* Быстрые действия */}
        <div className="dashboard-section">
          <h2>Быстрые действия</h2>
          <div className="quick-actions">
            <Link to="/products/add" className="quick-action-card">
              <Plus size={32} />
              <h3>Добавить продукт</h3>
              <p>Создать новый продукт в каталоге</p>
            </Link>

            <Link to="/products" className="quick-action-card">
              <Package size={32} />
              <h3>Управление продуктами</h3>
              <p>Просмотр и редактирование продуктов</p>
            </Link>

            <div className="quick-action-card disabled">
              <BarChart3 size={32} />
              <h3>Аналитика</h3>
              <p>Скоро будет доступно</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default Dashboard