import { useState, useEffect } from 'react'
import { useNavigate, useParams, Link } from 'react-router-dom'
import { Save, ArrowLeft, AlertCircle, CheckCircle, Loader } from 'lucide-react'
import { productsAPI } from '../../services/api'

const EditProduct = () => {
  const navigate = useNavigate()
  const { id } = useParams()
  const [formData, setFormData] = useState({
    name: '',
    price: '',
    quantity: ''
  })
  const [originalData, setOriginalData] = useState(null)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')

  useEffect(() => {
    loadProduct()
  }, [id])

  const loadProduct = async () => {
    setLoading(true)
    setError('')

    try {
      const result = await productsAPI.getProduct(id)
      
      if (result.success) {
        const product = result.data
        const productData = {
          name: product.name,
          price: product.price.toString(),
          quantity: product.quantity.toString()
        }
        setFormData(productData)
        setOriginalData(product)
      } else {
        setError(result.error)
      }
    } catch (error) {
      setError('Ошибка загрузки продукта')
    } finally {
      setLoading(false)
    }
  }

  const handleChange = (e) => {
    const { name, value } = e.target
    setFormData(prev => ({
      ...prev,
      [name]: value
    }))
    // Очищаем сообщения при изменении полей
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

  const hasChanges = () => {
    if (!originalData) return false
    
    return (
      formData.name !== originalData.name ||
      Number(formData.price) !== originalData.price ||
      Number(formData.quantity) !== originalData.quantity
    )
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setSaving(true)
    setError('')
    setSuccess('')

    if (!validateForm()) {
      setSaving(false)
      return
    }

    if (!hasChanges()) {
      setError('Нет изменений для сохранения')
      setSaving(false)
      return
    }

    try {
      const productData = {
        name: formData.name.trim(),
        price: Number(formData.price),
        quantity: Number(formData.quantity)
      }

      const result = await productsAPI.updateProduct(id, productData)
      
      if (result.success) {
        setSuccess('Продукт успешно обновлен! Перенаправляем...')
        setTimeout(() => {
          navigate('/products')
        }, 1500)
      } else {
        setError(result.error)
      }
    } catch (error) {
      setError('Произошла ошибка при обновлении продукта')
    } finally {
      setSaving(false)
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

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('ru-RU', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  if (loading) {
    return (
      <div className="edit-product-loading">
        <div className="loading-spinner"></div>
        <p>Загрузка продукта...</p>
      </div>
    )
  }

  if (error && !originalData) {
    return (
      <div className="edit-product-error">
        <div className="container-sm">
          <div className="error-state card">
            <AlertCircle size={48} />
            <h2>Ошибка загрузки</h2>
            <p>{error}</p>
            <Link to="/products" className="btn btn-primary">
              <ArrowLeft size={20} />
              Вернуться к списку
            </Link>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="edit-product-page fade-in">
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
            <h1>Редактировать продукт</h1>
            <p className="page-subtitle">
              Измените информацию о продукте
            </p>
          </div>
        </div>

        {/* Информация о продукте */}
        {originalData && (
          <div className="product-info card mb-6">
            <h3>Информация о продукте</h3>
            <div className="info-grid">
              <div className="info-item">
                <span className="info-label">ID:</span>
                <span className="info-value">#{originalData.id}</span>
              </div>
              <div className="info-item">
                <span className="info-label">Владелец:</span>
                <span className="info-value">{originalData.owner_name}</span>
              </div>
              <div className="info-item">
                <span className="info-label">Создано:</span>
                <span className="info-value">{formatDate(originalData.created_at)}</span>
              </div>
            </div>
          </div>
        )}

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
                disabled={saving}
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
                  disabled={saving}
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
                  disabled={saving}
                  min="1"
                  step="1"
                />
                <small className="form-hint">
                  Целое число, больше 0
                </small>
              </div>
            </div>

            {/* Сравнение изменений */}
            {originalData && hasChanges() && (
              <div className="changes-preview">
                <h3>Изменения</h3>
                <div className="changes-grid">
                  <div className="changes-column">
                    <h4>Было</h4>
                    <div className="preview-card original">
                      <div className="preview-header">
                        <h5>{originalData.name}</h5>
                      </div>
                      <div className="preview-details">
                        <div className="preview-item">
                          <span>Цена:</span>
                          <span>{formatPrice(originalData.price)}</span>
                        </div>
                        <div className="preview-item">
                          <span>Количество:</span>
                          <span>{originalData.quantity} шт.</span>
                        </div>
                        <div className="preview-item">
                          <span>Общая стоимость:</span>
                          <span>{formatPrice(originalData.price * originalData.quantity)}</span>
                        </div>
                      </div>
                    </div>
                  </div>

                  <div className="changes-column">
                    <h4>Будет</h4>
                    <div className="preview-card updated">
                      <div className="preview-header">
                        <h5>{formData.name || 'Название продукта'}</h5>
                      </div>
                      <div className="preview-details">
                        <div className="preview-item">
                          <span>Цена:</span>
                          <span>{formatPrice(formData.price)}</span>
                        </div>
                        <div className="preview-item">
                          <span>Количество:</span>
                          <span>{formData.quantity || 0} шт.</span>
                        </div>
                        <div className="preview-item">
                          <span>Общая стоимость:</span>
                          <span>{formatPrice(calculateTotal())}</span>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            )}

            <div className="form-actions">
              <Link 
                to="/products" 
                className="btn btn-secondary"
                disabled={saving}
              >
                Отмена
              </Link>
              <button
                type="submit"
                className="btn btn-primary"
                disabled={saving || !hasChanges()}
              >
                {saving ? (
                  <>
                    <div className="loading-spinner small"></div>
                    Сохранение...
                  </>
                ) : (
                  <>
                    <Save size={20} />
                    Сохранить изменения
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

export default EditProduct