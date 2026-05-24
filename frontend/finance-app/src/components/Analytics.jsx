import React, { useState, useEffect } from 'react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts';
import { Calendar, BarChart3, PieChart as PieIcon, TrendingUp, RefreshCw } from 'lucide-react';
import { api } from '../api/axios'; // Импортируем наш Axios инстанс
import './Analytics.css';

export const Analytics = () => {
    // Даты для фильтрации (соответствуют Query-параметрам бэкенда)
    const [fromDate, setFromDate] = useState('2026-05-01');
    const [toDate, setToDate] = useState('2026-05-24');
    const [loading, setLoading] = useState(false);

    // Сюда будем сохранять реальный ответ от Go-бэкенда
    const [serverData, setServerData] = useState({
        total_income: 0,
        total_expense: 0,
        expenses_by_category: [],
        income_by_category: []
    });

    // Загрузка данных аналитики с сервера
    useEffect(() => {
        const fetchAnalyticsData = async () => {
            setLoading(true);
            try {
                const response = await api.get(`/protected/dashboards?from=${fromDate}&to=${toDate}`);

                setServerData({
                    total_income: response.data.total_income || 0,
                    total_expense: response.data.total_expense || 0,
                    expenses_by_category: response.data.expenses_by_category || [],
                    income_by_category: response.data.income_by_category || []
                });
            } catch (err) {
                console.error("Ошибка при получении данных аналитики:", err);
            } finally {
                setLoading(false);
            }
        };

        fetchAnalyticsData();
    }, [fromDate, toDate]);

    // Разделяем доходы и расходы на два отдельных столбца по оси X, чтобы график выглядел воздушно и гармонично
    const compareData = [
        {
            name: 'Доходы',
            'Сумма': serverData.total_income,
            fillColor: '#10b981'
        },
        {
            name: 'Расходы',
            'Сумма': serverData.total_expense,
            fillColor: '#f87171'
        }
    ];

    // Палитра для круговой диаграммы
    const COLORS = ['#6366f1', '#10b981', '#f59e0b', '#ef4444', '#ec4899', '#8b5cf6'];

    // Исправленный, аккуратный CustomTooltip в темных тонах для обоих графиков
    const CustomTooltip = ({ active, payload }) => {
        if (active && payload && payload.length) {
            return (
                <div style={{
                    backgroundColor: '#1e293b',
                    border: '1px solid #334155',
                    padding: '10px 14px',
                    borderRadius: '8px',
                    boxShadow: '0 10px 15px -3px rgba(0, 0, 0, 0.3)'
                }}>
                    <p style={{ margin: '0 0 4px 0', color: '#94a3b8', fontSize: '12px', fontWeight: '500' }}>
                        {payload[0].name || payload[0].payload.category_name}
                    </p>
                    <p style={{ margin: 0, color: payload[0].payload.fillColor || payload[0].color || '#f8fafc', fontSize: '15px', fontWeight: '600' }}>
                        {payload[0].value.toLocaleString('ru-RU')} ₽
                    </p>
                </div>
            );
        }
        return null;
    };

    const savingsRate = serverData.total_income > 0
        ? Math.round(((serverData.total_income - serverData.total_expense) / serverData.total_income) * 100)
        : 0;

    return (
        <div className="analytics-container">
            <div className="analytics-title-zone" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px' }}>
                <div className="analytics-header">
                    <h2>Аналитика и графики</h2>
                    <p style={{ color: '#64748b', fontSize: '14px', marginTop: '4px' }}>
                        Наглядный визуальный анализ ваших денежных потоков
                    </p>
                </div>

                <div className="date-filter-panel">
                    <Calendar size={16} style={{ color: '#94a3b8' }} />
                    <input type="date" value={fromDate} onChange={(e) => setFromDate(e.target.value)} />
                    <span style={{ color: '#64748b' }}>—</span>
                    <input type="date" value={toDate} onChange={(e) => setToDate(e.target.value)} />
                    {loading && <RefreshCw size={14} className="animate-spin" style={{ color: '#6366f1', marginLeft: '8px' }} />}
                </div>
            </div>

            <div className="charts-grid">
                {/* Левый график: Сравнение общих объемов (ИСПРАВЛЕННЫЙ) */}
                <div className="chart-card">
                    <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '16px' }}>
                        <BarChart3 size={18} style={{ color: '#6366f1' }} />
                        <h3>Соотношение доходов и расходов</h3>
                    </div>
                    <div className="responsive-chart-container">
                        <ResponsiveContainer width="100%" height="100%">
                            <BarChart data={compareData} margin={{ top: 10, right: 10, left: 10, bottom: 5 }}>
                                <CartesianGrid strokeDasharray="3 3" stroke="#334155" vertical={false} />
                                <XAxis dataKey="name" stroke="#94a3b8" tickLine={false} />
                                <YAxis stroke="#94a3b8" tickFormatter={(v) => `${v / 1000}k`} tickLine={false} axisLine={false} />

                                {/* cursor={false} полностью убирает уродливый скачущий задний фон при наведении */}
                                <Tooltip content={<CustomTooltip />} cursor={false} />

                                {/* dataKey="Сумма" отрисует аккуратные раздельные столбики с динамическим цветом fill */}
                                <Bar
                                    dataKey="Сумма"
                                    radius={[6, 6, 0, 0]}
                                    maxBarSize={50}
                                    animationDuration={300}
                                >
                                    {compareData.map((entry, index) => (
                                        <Cell
                                            key={`cell-${index}`}
                                            fill={entry.fillColor}
                                            style={{ transition: 'opacity 0.2s ease', cursor: 'pointer' }}
                                            // Добавляем микро-эффект затухания при наведении вместо серого окна
                                            onMouseEnter={(e) => e.target.style.opacity = 0.85}
                                            onMouseLeave={(e) => e.target.style.opacity = 1}
                                        />
                                    ))}
                                </Bar>
                            </BarChart>
                        </ResponsiveContainer>
                    </div>
                </div>

                {/* Правый график: Круговая диаграмма расходов */}
                <div className="chart-card">
                    <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '16px' }}>
                        <PieIcon size={18} style={{ color: '#f59e0b' }} />
                        <h3>Доли расходов</h3>
                    </div>
                    <div className="responsive-chart-container">
                        <ResponsiveContainer width="100%" height="100%">
                            {serverData.expenses_by_category.length === 0 ? (
                                <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100%', color: '#64748b' }}>
                                    Нет данных для отображения диаграммы
                                </div>
                            ) : (
                                <PieChart>
                                    <Pie
                                        data={serverData.expenses_by_category}
                                        dataKey="total"
                                        nameKey="category_name"
                                        cx="50%"
                                        cy="50%"
                                        innerRadius={65}
                                        outerRadius={90}
                                        paddingAngle={3}
                                        animationDuration={300}
                                    >
                                        {serverData.expenses_by_category.map((entry, index) => (
                                            <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                                        ))}
                                    </Pie>
                                    <Tooltip content={<CustomTooltip />} />
                                </PieChart>
                            )}
                        </ResponsiveContainer>
                    </div>

                    <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '8px', marginTop: '12px', fontSize: '12px', maxHeight: '70px', overflowY: 'auto' }}>
                        {serverData.expenses_by_category.map((cat, index) => (
                            <div key={cat.category_id} style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
                                <div style={{ width: '8px', height: '8px', borderRadius: '50%', backgroundColor: COLORS[index % COLORS.length], flexShrink: 0 }} />
                                <span style={{ color: '#94a3b8', textOverflow: 'ellipsis', overflow: 'hidden', whiteSpace: 'nowrap' }}>{cat.category_name}</span>
                            </div>
                        ))}
                    </div>
                </div>
            </div>

            {/* Карточка эффективности */}
            <div className="chart-card" style={{ marginTop: '24px' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '10px' }}>
                    <TrendingUp size={18} style={{ color: '#10b981' }} />
                    <h3>Финансовая эффективность периода</h3>
                </div>
                <p style={{ color: '#94a3b8', fontSize: '14px', margin: '0' }}>
                    {serverData.total_income > 0 ? (
                        <>
                            Вы сохраняете <strong style={{ color: '#10b981' }}>{savingsRate}%</strong> от общего объема входящих средств.
                            {savingsRate > 0 ? ' Это отличный показатель сбережений!' : ' Стоит пересмотреть структуру расходов.'}
                        </>
                    ) : (
                        'Внесите ваши доходы за этот период, чтобы рассчитать процент финансовой эффективности.'
                    )}
                </p>
            </div>
        </div>
    );
};