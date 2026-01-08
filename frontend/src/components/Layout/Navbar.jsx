import { Link, useLocation } from 'react-router-dom'
import { LogOut, Package, Plus, Home, User } from 'lucide-react'
import { apiUtils } from '../../services/api'

const Navbar = ({ onLogout }) => {
  const location = useLocation()
  const user = apiUtils.getUserFromToken()

  const isActive = (path) => {
    return location.pathname === path ? 'nav-link active' : 'nav-link'
  }

  return (
    <nav className="navbar">
      <div className="navbar-container">
        {/* Логотип */}
        <div className="navbar-brand">
          <Package className="brand-icon" />
          <span className="brand-text">ProductKeeper</span>
        </div>

        {/* Навигационные ссылки */}
        <div className="navbar-nav">
          <Link to="/dashboard" className={isActive('/dashboard')}>
            <Home size={20} />
            <span>Главная</span>
          </Link>
          
          <Link to="/products" className={isActive('/products')}>
            <Package size={20} />
            <span>Продукты</span>
          </Link>
          
          <Link to="/products/add" className={isActive('/products/add')}>
            <Plus size={20} />
            <span>Добавить</span>
          </Link>
        </div>

        {/* Пользователь и выход */}
        <div className="navbar-user">
          <div className="user-info">
            <User size={20} />
            <span className="username">{user?.username || 'Пользователь'}</span>
          </div>
          
          <button onClick={onLogout} className="logout-btn" title="Выйти">
            <LogOut size={16} />
            <span className="logout-text">Выйти</span>
          </button>
        </div>
      </div>
    </nav>
  )
}

export default Navbar