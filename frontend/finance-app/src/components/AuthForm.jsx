import React, { useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { Mail, Lock, User, Eye, EyeOff } from 'lucide-react';
import './AuthForm.css';

export const AuthForm = ({ onAuthSuccess }) => {
    const { login, register } = useAuth();
    const [isLogin, setIsLogin] = useState(true);
    const [showPassword, setShowPassword] = useState(false);
    const [isSubmitting, setIsSubmitting] = useState(false);

    const [formData, setFormData] = useState({ name: '', email: '', password: '' });
    const [errors, setErrors] = useState({});
    const [apiError, setApiError] = useState('');

    const handleInputChange = (e) => {
        const { name, value } = e.target;
        setFormData((prev) => ({ ...prev, [name]: value }));
        if (errors[name]) setErrors((prev) => ({ ...prev, [name]: '' }));
    };

    const validateForm = () => {
        const newErrors = {};
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

        if (!isLogin && (!formData.name || formData.name.trim().length < 2)) {
            newErrors.name = 'Имя должно содержать не менее 2 символов';
        }

        if (!formData.email) {
            newErrors.email = 'Email обязателен для заполнения';
        } else if (!emailRegex.test(formData.email)) {
            newErrors.email = 'Введите корректный email адрес';
        }

        if (!formData.password) {
            newErrors.password = 'Пароль обязателен для заполнения';
        } else if (formData.password.length < 8) {
            newErrors.password = 'Пароль должен быть не менее 8 символов';
        }

        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setApiError('');

        if (!validateForm()) return;

        setIsSubmitting(true);
        try {
            if (isLogin) {
                await login(formData.email, formData.password);
                onAuthSuccess?.();
            } else {
                // Регистрируем на бэкенде через /api/v1/auth/register
                await register(formData.name, formData.email, formData.password);
                // Сразу автоматически авторизуем пользователя
                await login(formData.email, formData.password);
                onAuthSuccess?.();
            }
        } catch (err) {
            // Перехватываем то точечное сообщение об ошибке, которое отдал c.JSON() в Go
            const message = err.response?.data?.error || 'Произошла непредвиденная ошибка. Попробуйте позже.';
            setApiError(message);
        } finally {
            setIsSubmitting(false);
        }
    };

    const toggleMode = () => {
        setIsLogin(!isLogin);
        setFormData({ name: '', email: '', password: '' });
        setErrors({});
        setApiError('');
    };

    return (
        <div className="auth-container">
            <div className="auth-card">
                <div className="auth-header">
                    <h2>{isLogin ? 'Вход в систему' : 'Регистрация'}</h2>
                    <p>{isLogin ? 'Добро пожаловать обратно!' : 'Создайте аккаунт для учета финансов'}</p>
                </div>

                {apiError && <div className="api-error-box">{apiError}</div>}

                <form onSubmit={handleSubmit} noValidate>
                    {!isLogin && (
                        <div className="form-group">
                            <label htmlFor="name">Ваше имя</label>
                            <div className="input-wrapper">
                                <User className="input-icon" size={18} />
                                <input
                                    type="text"
                                    id="name"
                                    name="name"
                                    value={formData.name}
                                    onChange={handleInputChange}
                                    className={errors.name ? 'input-error' : ''}
                                    placeholder="Константин"
                                />
                            </div>
                            {errors.name && <span className="error-text">{errors.name}</span>}
                        </div>
                    )}

                    <div className="form-group">
                        <label htmlFor="email">Электронная почта</label>
                        <div className="input-wrapper">
                            <Mail className="input-icon" size={18} />
                            <input
                                type="email"
                                id="email"
                                name="email"
                                value={formData.email}
                                onChange={handleInputChange}
                                className={errors.email ? 'input-error' : ''}
                                placeholder="example@mail.com"
                            />
                        </div>
                        {errors.email && <span className="error-text">{errors.email}</span>}
                    </div>

                    <div className="form-group">
                        <label htmlFor="password">Пароль</label>
                        <div className="input-wrapper">
                            <Lock className="input-icon" size={18} />
                            <input
                                type={showPassword ? 'text' : 'password'}
                                id="password"
                                name="password"
                                value={formData.password}
                                onChange={handleInputChange}
                                className={errors.password ? 'input-error' : ''}
                                placeholder="••••••••"
                            />
                            <button
                                type="button"
                                style={{
                                    position: 'absolute',
                                    right: '12px',
                                    background: 'none',
                                    border: 'none',
                                    color: '#64748b',
                                    cursor: 'pointer',
                                    display: 'flex',
                                    alignItems: 'center'
                                }}
                                onClick={() => setShowPassword(!showPassword)}
                            >
                                {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                            </button>
                        </div>
                        {errors.password && <span className="error-text">{errors.password}</span>}
                    </div>

                    <button type="submit" className="submit-btn" disabled={isSubmitting}>
                        {isSubmitting ? 'Загрузка...' : isLogin ? 'Войти' : 'Зарегистрироваться'}
                    </button>
                </form>

                <div className="auth-toggle">
                    {isLogin ? 'Впервые у нас?' : 'Уже есть аккаунт?'}
                    <button type="button" className="toggle-link" onClick={toggleMode}>
                        {isLogin ? 'Создать аккаунт' : 'Войти в профиль'}
                    </button>
                </div>
            </div>
        </div>
    );
};