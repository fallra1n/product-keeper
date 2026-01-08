import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { Package, Plus, Search, Edit, Trash2, Eye, Filter, SortAsc, SortDesc } from 'lucide-react'
import { productsAPI } from '../../services/api'

const ProductList = () => {
  const [products, setProducts] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [searchTerm, setSearchTerm] = useState('')
  const [sortBy, setSortBy] = useState('last_create')
  const [showFilters, setShowFilters] = useState(false)

  useEffect(() => {
    loadProducts()
  }, [sortBy])

  useEffect(() => {
    const delayedSearch = setTimeout(() => {
      if (searchTerm !== '') {
        loadProducts()
      } else if (searchTerm === '') {
        loadProducts()
      }
    }, 500)

    return () => clearTimeout(delayedSearch)
  }, [searchTerm])

  const loadProducts = async () => {
    setLoading(true)
    setError('')

    try {
      const params = {}
      if (searchTerm) params.name = searchTerm
      if (sortBy) params.sort_by = sortBy

      const result = await productsAPI.getProducts(params)
      
      if (result.success) {
        setProducts(result.data || [])
      } else {
        setError(result.error)
      }
    } catch (error) {
      setError('Ошибка загрузки продуктов')
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (id, name) => {
    if (!window.confirm(`Вы уверены, что хотите удалить продукт "${name}"?`)) {
      return
    }

    try {
      const result = await productsAPI.deleteProduct(id)
      
      if (result.success) {
        setProducts(products.filter(product => product.id !== id))
      } else {
        setError(result.error)
      }
    } catch (error) {
      setError('Ошибка удаления продукта')
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

  const getSortIcon = () => {
    return sortBy === 'name' ? <SortAsc size={16} /> : <SortDesc size={16} />
  }

  if (loading) {
    return (
      <div className="products-loading">
        <div className="loading-spinner"></div>
        <p>Загрузка продуктов...</p>
      </div>
    )
  }

  return (
    <div className="products-page fade-in">
      <div className="container">
        {/* Заголовок */}
        <div className="page-header">
          <div>
            <h1>Управление продуктами</h1>
            <p className="page-subtitle">
              Просматривайте и управляйте своими продуктами
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

        {/* Фильтры и поиск */}
        <div className="products-controls">
          <div className="search-container">
            <Search size={20} className="search-icon" />
            <input
              type="text"
              placeholder="Поиск по названию..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="search-input"
            />
          </div>

          <div className="controls-right">
            <button
              onClick={() => setShowFilters(!showFilters)}
              className={`btn btn-secondary ${showFilters ? 'active' : ''}`}
            >
              <Filter size={20} />
              Фильтры
            </button>

            <select
              value={sortBy}
              onChange={(e) => setSortBy(e.target.value)}
              className="sort-select"
            >
              <option value="last_create">По дате создания</option>
              <option value="name">По названию</option>
            </select>
          </div>
        </div>

        {/* Расширенные фильтры */}
        {showFilters && (
          <div className="filters-panel slide-up">
            <div className="filters-content">
              <h3>Дополнительные фильтры</h3>
              <p className="text-muted">Скоро будет доступно больше фильтров</p>
            </div>
          </div>
        )}

        {/* Список продуктов */}
        {products.length > 0 ? (
          <div className="products-grid">
            {products.map((product) => (
              <div key={product.id} className="product-item">
                <div className="product-header">
                  <h3 className="product-title">{product.name}</h3>
                  <div className="product-actions">
                    <Link
                      to={`/products/edit/${product.id}`}
                      className="action-btn edit"
                      title="Редактировать"
                    >
                      <Edit size={14} />
                      <span>Изменить</span>
                    </Link>
                    <button
                      onClick={() => handleDelete(product.id, product.name)}
                      className="action-btn delete"
                      title="Удалить"
                    >
                      <Trash2 size={14} />
                      <span>Удалить</span>
                    </button>
                  </div>
                </div>

                <div className="product-details">
                  <div className="detail-item">
                    <span className="detail-label">Цена:</span>
                    <span className="detail-value price">{formatPrice(product.price)}</span>
                  </div>
                  
                  <div className="detail-item">
                    <span className="detail-label">Количество:</span>
                    <span className="detail-value">{product.quantity} шт.</span>
                  </div>
                  
                  <div className="detail-item">
                    <span className="detail-label">Общая стоимость:</span>
                    <span className="detail-value total">{formatPrice(product.price * product.quantity)}</span>
                  </div>
                </div>

                <div className="product-footer">
                  <span className="product-date">
                    Создано: {formatDate(product.created_at)}
                  </span>
                  <span className="product-owner">
                    Владелец: {product.owner_name}
                  </span>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="empty-state">
            <Package size={48} />
            <h3>
              {searchTerm ? 'Продукты не найдены' : 'Пока нет продуктов'}
            </h3>
            <p>
              {searchTerm 
                ? `По запросу "${searchTerm}" ничего не найдено. Попробуйте изменить поисковый запрос.`
                : 'Создайте свой первый продукт, чтобы начать работу'
              }
            </p>
            {!searchTerm && (
              <Link to="/products/add" className="btn btn-primary">
                <Plus size={20} />
                Добавить продукт
              </Link>
            )}
          </div>
        )}

        {/* Статистика */}
        {products.length > 0 && (
          <div className="products-stats">
            <div className="stats-item">
              <span className="stats-label">Всего продуктов:</span>
              <span className="stats-value">{products.length}</span>
            </div>
            <div className="stats-item">
              <span className="stats-label">Общее количество:</span>
              <span className="stats-value">
                {products.reduce((sum, product) => sum + product.quantity, 0)} шт.
              </span>
            </div>
            <div className="stats-item">
              <span className="stats-label">Общая стоимость:</span>
              <span className="stats-value">
                {formatPrice(products.reduce((sum, product) => sum + (product.price * product.quantity), 0))}
              </span>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

export default ProductList