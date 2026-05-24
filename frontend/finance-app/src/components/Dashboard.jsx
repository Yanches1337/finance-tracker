import React, { useState, useEffect } from 'react';
import { Wallet, ArrowUpCircle, ArrowDownCircle, Calendar, RefreshCw } from 'lucide-react';
import { api } from '../api/axios'; // Импортируем настроенный Axios инстанс
import './Dashboard.css';

export const Dashboard = () => {
    const getInitialDates = () => {
        const today = new Date();
        const yyyy = today.getFullYear();
        const mm = String(today.getMonth() + 1).padStart(2, '0');
        const dd = String(today.getDate()).padStart(2, '0');

        return {
            from: `${yyyy}-${mm}-01`,
            to: `${yyyy}-${mm}-${dd}`
        };
    };

    const [dateRange, setDateRange] = useState(getInitialDates());
    const [loading, setLoading] = useState(false);

    // Локальное состояние под структуру Go-домена domain.Dashboard
    const [dashboardData, setDashboardData] = useState({
        balance: 0,
        total_income: 0,
        total_expense: 0,
        expenses_by_category: [],
        income_by_category: []
    });

    const getMaxCategoryTotal = (categories) => {
        if (!categories || categories.length === 0) return 1;
        return Math.max(...categories.map(c => c.total));
    };

    const handleDateChange = (e) => {
        const { name, value } = e.target;
        setDateRange(prev => ({ ...prev, [name]: value }));
    };

    // Реальный асинхронный запрос к твоему DashboardHandler в Go
    const fetchDashboardData = async () => {
        setLoading(true);
        try {
            const response = await api.get(`/protected/dashboards?from=${dateRange.from}&to=${dateRange.to}`);

            // Записываем данные ответа, защищаясь от null-ответов бэкенда по массивам
            setDashboardData({
                balance: response.data.balance || 0,
                total_income: response.data.total_income || 0,
                total_expense: response.data.total_expense || 0,
                expenses_by_category: response.data.expenses_by_category || [],
                income_by_category: response.data.income_by_category || []
            });
        } catch (err) {
            console.error("Не удалось обновить сводную панель:", err);
        } finally {
            setLoading(false);
        }
    };

    // Обновляем данные каждый раз, когда пользователь меняет фильтр дат
    useEffect(() => {
        fetchDashboardData();
    }, [dateRange.from, dateRange.to]);

    const maxExpense = getMaxCategoryTotal(dashboardData.expenses_by_category);
    const maxIncome = getMaxCategoryTotal(dashboardData.income_by_category);

    return (
        <div className="dashboard-container">
            <div className="dashboard-title-zone">
                <div>
                    <h2>Сводная панель</h2>
                    <p style={{ color: '#64748b', fontSize: '14px', marginTop: '4px' }}>
                        Обзор финансового состояния за выбранный период
                    </p>
                </div>

                <div className="date-filter-panel">
                    <Calendar size={16} style={{ color: '#94a3b8' }} />
                    <input
                        type="date"
                        name="from"
                        value={dateRange.from}
                        onChange={handleDateChange}
                    />
                    <span style={{ color: '#64748b' }}>—</span>
                    <input
                        type="date"
                        name="to"
                        value={dateRange.to}
                        onChange={handleDateChange}
                    />
                    {loading && <RefreshCw size={14} className="animate-spin" style={{ color: '#6366f1' }} />}
                </div>
            </div>

            <div className="metrics-grid">
                <div className="metric-card balance">
                    <div className="metric-icon-box">
                        <Wallet size={24} />
                    </div>
                    <div className="metric-info">
                        <p>Текущий Баланс</p>
                        <div className="metric-value">
                            {(dashboardData.balance).toLocaleString('ru-RU')} ₽
                        </div>
                    </div>
                </div>

                <div className="metric-card income">
                    <div className="metric-icon-box">
                        <ArrowUpCircle size={24} />
                    </div>
                    <div className="metric-info">
                        <p>Всего доходов</p>
                        <div className="metric-value" style={{ color: '#10b981' }}>
                            + {(dashboardData.total_income).toLocaleString('ru-RU')} ₽
                        </div>
                    </div>
                </div>

                <div className="metric-card expense">
                    <div className="metric-icon-box">
                        <ArrowDownCircle size={24} />
                    </div>
                    <div className="metric-info">
                        <p>Всего расходов</p>
                        <div className="metric-value" style={{ color: '#fca5a5' }}>
                            - {(dashboardData.total_expense).toLocaleString('ru-RU')} ₽
                        </div>
                    </div>
                </div>
            </div>

            <div className="dashboard-details-grid">
                {/* Расходы по категориям */}
                <div className="category-stats-block">
                    <div className="category-stats-title">Расходы по категориям</div>
                    <div className="category-rows">
                        {dashboardData.expenses_by_category.length === 0 ? (
                            <p style={{ color: '#64748b', fontSize: '14px' }}>Нет расходов за этот период</p>
                        ) : (
                            dashboardData.expenses_by_category.map((cat) => {
                                const percentage = maxExpense > 0 ? (cat.total / maxExpense) * 100 : 0;
                                return (
                                    <div key={cat.category_id} className="category-stat-row">
                                        <div className="category-stat-info">
                                            <span className="cat-name">{cat.category_name}</span>
                                            <span className="cat-total" style={{ color: '#fca5a5' }}>
                                                {cat.total.toLocaleString('ru-RU')} ₽
                                            </span>
                                        </div>
                                        <div className="progress-bar-bg">
                                            <div
                                                className="progress-bar-fill expense"
                                                style={{ width: `${percentage}%` }}
                                            />
                                        </div>
                                    </div>
                                );
                            })
                        )}
                    </div>
                </div>

                {/* Доходы по категориям */}
                <div className="category-stats-block">
                    <div className="category-stats-title">Источники доходов</div>
                    <div className="category-rows">
                        {dashboardData.income_by_category.length === 0 ? (
                            <p style={{ color: '#64748b', fontSize: '14px' }}>Нет доходов за этот период</p>
                        ) : (
                            dashboardData.income_by_category.map((cat) => {
                                const percentage = maxIncome > 0 ? (cat.total / maxIncome) * 100 : 0;
                                return (
                                    <div key={cat.category_id} className="category-stat-row">
                                        <div className="category-stat-info">
                                            <span className="cat-name">{cat.category_name}</span>
                                            <span className="cat-total" style={{ color: '#a7f3d0' }}>
                                                {cat.total.toLocaleString('ru-RU')} ₽
                                            </span>
                                        </div>
                                        <div className="progress-bar-bg">
                                            <div
                                                className="progress-bar-fill income"
                                                style={{ width: `${percentage}%` }}
                                            />
                                        </div>
                                    </div>
                                );
                            })
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
};