import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { Plus, ArrowLeft, AlertCircle, CheckCircle } from 'lucide-react'
import { productsAPI } from '../../services/api'

const AddProduct = () => {
  const navigate = useNavigate()
  const [formData, setFormData] = useState({
    name: '',
    price: '',
    quantity: ''
  })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')

  const handleChange = (e) => {
    const { name, value } = e.target
    setFormData(prev => ({
      ...prev,
      [name]: value
    }))
    // Очищаем ошибки при изменении полей
    if (error) setError('')
    if (success) setSuccess('')
  }

  const validateForm = () => {
    if (!formData.name.trim()) {
      setError('Введите название продукта')
      return false
    }

    if (formData.name.length < 2) {
      setError('Название должно содержать минимум 2 символа')
      return false
    }

    if (!formData.price || isNaN(formData.price) || Number(formData.price) <= 0) {
      setError('Введите корректную цену (больше 0)')
      return false
    }

    if (!formData.quantity || isNaN(formData.quantity) || Number(formData.quantity) <= 0) {
      setError('Введите корректное количество (больше 0)')
      return false
    }

    if (!Number.isInteger(Number(formData.quantity))) {
      setError('Количество должно быть целым числом')
      return false
    }

    return true
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setLoading(true)
    setError('')
    setSuccess('')

    if (!validateForm()) {
      setLoading(false)
      return
    }

    try {
      const productData = {
        name: formData.name.trim(),
        price: Number(formData.price),
        quantity: Number(formData.quantity)
      }

      const result = await productsAPI.createProduct(productData)
      
      if (result.success) {
        setSuccess('Продукт успешно создан! Перенаправляем...')
        setTimeout(() => {
          navigate('/products')
        }, 1500)
      } else {
        setError(result.error)
      }
    } catch (error) {
      setError('Произошла ошибка при создании продукта')
    } finally {
      setLoading(false)
    }
  }

  const formatPrice = (price) => {
    if (!price || isNaN(price)) return '0 ₽'
    return new Intl.NumberFormat('ru-RU', {
      style: 'currency',
      currency: 'RUB'
    }).format(Number(price))
  }

  const calculateTotal = () => {
    const price = Number(formData.price) || 0
    const quantity = Number(formData.quantity) || 0
    return price * quantity
  }

  return (
    <div className="add-product-page fade-in">
      <div className="container-sm">
        {/* Заголовок */}
        <div className="page-header">
          <div className="header-nav">
            <Link to="/products" className="btn btn-secondary">
              <ArrowLeft size={20} />
              Назад к списку
            </Link>
          </div>
          <div className="header-content">
            <h1>Добавить новый продукт</h1>
            <p className="page-subtitle">
              Заполните информацию о продукте
            </p>
          </div>
        </div>

        {/* Форма */}
        <div className="product-form-container">
          <form onSubmit={handleSubmit} className="product-form card">
            {error && (
              <div className="error-message mb-6">
                <AlertCircle size={16} />
                {error}
              </div>
            )}

            {success && (
              <div className="success-message mb-6">
                <CheckCircle size={16} />
                {success}
              </div>
            )}

            <div className="form-group">
              <label htmlFor="name" className="form-label">
                Название продукта *
              </label>
              <input
                type="text"
                id="name"
                name="name"
                value={formData.name}
                onChange={handleChange}
                className={`form-input ${error && !formData.name ? 'error' : ''}`}
                placeholder="Введите название продукта"
                disabled={loading}
                maxLength={255}
              />
              <small className="form-hint">
                Минимум 2 символа, максимум 255
              </small>
            </div>

            <div className="form-row">
              <div className="form-group">
                <label htmlFor="price" className="form-label">
                  Цена за единицу *
                </label>
                <input
                  type="number"
                  id="price"
                  name="price"
                  value={formData.price}
                  onChange={handleChange}
                  className={`form-input ${error && (!formData.price || Number(formData.price) <= 0) ? 'error' : ''}`}
                  placeholder="0"
                  disabled={loading}
                  min="0.01"
                  step="0.01"
                />
                <small className="form-hint">
                  В рублях, больше 0
                </small>
              </div>

              <div className="form-group">
                <label htmlFor="quantity" className="form-label">
                  Количество *
                </label>
                <input
                  type="number"
                  id="quantity"
                  name="quantity"
                  value={formData.quantity}
                  onChange={handleChange}
                  className={`form-input ${error && (!formData.quantity || Number(formData.quantity) <= 0) ? 'error' : ''}`}
                  placeholder="0"
                  disabled={loading}
                  min="1"
                  step="1"
                />
                <small className="form-hint">
                  Целое число, больше 0
                </small>
              </div>
            </div>

            {/* Предварительный просмотр */}
            {(formData.name || formData.price || formData.quantity) && (
              <div className="product-preview">
                <h3>Предварительный просмотр</h3>
                <div className="preview-card">
                  <div className="preview-header">
                    <h4>{formData.name || 'Название продукта'}</h4>
                  </div>
                  <div className="preview-details">
                    <div className="preview-item">
                      <span>Цена за единицу:</span>
                      <span className="preview-value">{formatPrice(formData.price)}</span>
                    </div>
                    <div className="preview-item">
                      <span>Количество:</span>
                      <span className="preview-value">{formData.quantity || 0} шт.</span>
                    </div>
                    <div className="preview-item total">
                      <span>Общая стоимость:</span>
                      <span className="preview-value">{formatPrice(calculateTotal())}</span>
                    </div>
                  </div>
                </div>
              </div>
            )}

            <div className="form-actions">
              <Link 
                to="/products" 
                className="btn btn-secondary"
                disabled={loading}
              >
                Отмена
              </Link>
              <button
                type="submit"
                className="btn btn-primary"
                disabled={loading}
              >
                {loading ? (
                  <>
                    <div className="loading-spinner small"></div>
                    Создание...
                  </>
                ) : (
                  <>
                    <Plus size={20} />
                    Создать продукт
                  </>
                )}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  )
}

export default AddProduct