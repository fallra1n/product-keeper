import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import { useState, useEffect } from 'react'
import Login from './components/Auth/Login'
import Register from './components/Auth/Register'
import Dashboard from './components/Dashboard/Dashboard'
import ProductList from './components/Products/ProductList'
import AddProduct from './components/Products/AddProduct'
import EditProduct from './components/Products/EditProduct'
import Navbar from './components/Layout/Navbar'
import './App.css'

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    // Проверяем наличие токена в localStorage
    const token = localStorage.getItem('token')
    if (token) {
      setIsAuthenticated(true)
    }
    setLoading(false)
  }, [])

  const handleLogin = (token) => {
    localStorage.setItem('token', token)
    setIsAuthenticated(true)
  }

  const handleLogout = () => {
    localStorage.removeItem('token')
    setIsAuthenticated(false)
  }

  if (loading) {
    return (
      <div className="loading-container">
        <div className="loading-spinner"></div>
        <p>Загрузка...</p>
      </div>
    )
  }

  return (
    <Router>
      <div className="app">
        {isAuthenticated && <Navbar onLogout={handleLogout} />}
        
        <main className="main-content">
          <Routes>
            {/* Публичные маршруты */}
            <Route 
              path="/login" 
              element={
                !isAuthenticated ? 
                <Login onLogin={handleLogin} /> : 
                <Navigate to="/dashboard" replace />
              } 
            />
            <Route 
              path="/register" 
              element={
                !isAuthenticated ? 
                <Register /> : 
                <Navigate to="/dashboard" replace />
              } 
            />
            
            {/* Защищенные маршруты */}
            <Route 
              path="/dashboard" 
              element={
                isAuthenticated ? 
                <Dashboard /> : 
                <Navigate to="/login" replace />
              } 
            />
            <Route 
              path="/products" 
              element={
                isAuthenticated ? 
                <ProductList /> : 
                <Navigate to="/login" replace />
              } 
            />
            <Route 
              path="/products/add" 
              element={
                isAuthenticated ? 
                <AddProduct /> : 
                <Navigate to="/login" replace />
              } 
            />
            <Route 
              path="/products/edit/:id" 
              element={
                isAuthenticated ? 
                <EditProduct /> : 
                <Navigate to="/login" replace />
              } 
            />
            
            {/* Перенаправление по умолчанию */}
            <Route 
              path="/" 
              element={
                <Navigate to={isAuthenticated ? "/dashboard" : "/login"} replace />
              } 
            />
          </Routes>
        </main>
      </div>
    </Router>
  )
}

export default App
