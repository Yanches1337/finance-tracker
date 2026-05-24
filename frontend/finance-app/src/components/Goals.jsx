import React, { useState, useEffect } from 'react';
import { Target, Plus, Calendar, TrendingUp, CheckCircle, AlertCircle, RefreshCw } from 'lucide-react';
import { api } from '../api/axios'; // Импортируем Axios
import './Goals.css';

export const Goals = () => {
    const [goals, setGoals] = useState([]);
    const [loading, setLoading] = useState(false);

    // Состояние формы (соответствует CreateGoalRequest на бэкенде)
    const [name, setName] = useState('');
    const [targetAmount, setTargetAmount] = useState('');
    const [targetDate, setTargetDate] = useState('');
    const [description, setDescription] = useState('');
    const [error, setError] = useState('');

    // Загрузка целей при старте страницы
    const fetchGoals = async () => {
        setLoading(true);
        try {
            const response = await api.get('/protected/goals');
            setGoals(response.data || []);
        } catch (err) {
            console.error("Не удалось загрузить цели:", err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchGoals();
    }, []);

    const handleCreateGoal = async (e) => {
        e.preventDefault();
        setError('');

        if (!name.trim()) {
            setError('Название цели обязательно');
            return;
        }
        const amount = parseInt(targetAmount, 10);
        if (isNaN(amount) || amount <= 0) {
            setError('Целевая сумма должна быть больше нуля');
            return;
        }

        // Решаем баг бэкенда: если дата выбрана, приводим её к RFC3339 ("YYYY-MM-DDT00:00:00Z")
        // Если дата не выбрана — отправляем пустую строку (или null), чтобы Go проставил дефолт
        const formattedDate = targetDate ? `${targetDate}T00:00:00Z` : "";

        try {
            const response = await api.post('/protected/goals', {
                name: name.trim(),
                target_amount: amount,
                target_date: formattedDate, // <--- Теперь бэкенд распарсит это без ошибок
                description: description.trim()
            });

            if (response.status === 201 || response.status === 200) {
                fetchGoals();
                setName('');
                setTargetAmount('');
                setTargetDate('');
                setDescription('');
            }
        } catch (err) {
            const msg = err.response?.data?.error || 'Ошибка при создании финансовой цели';
            setError(msg);
        }
    };

    // Быстрое пополнение цели через PUT-запрос
    const handleQuickDeposit = async (goal) => {
        const depositStr = prompt('Введите сумму пополнения (₽):');
        const deposit = parseInt(depositStr, 10);

        if (isNaN(deposit) || deposit <= 0) return;

        // Рассчитываем новую сумму накоплений
        const nextAmount = goal.current_amount + deposit;

        try {
            // Твой роут на бэкенде: PUT /goals/:id
            await api.put(`/protected/goals/${goal.id}`, {
                name: goal.name,
                target_amount: goal.target_amount,
                current_amount: nextAmount,
                target_date: goal.target_date,
                description: goal.description
            });

            // Перечитываем актуальные данные из БД
            fetchGoals();
        } catch (err) {
            alert(err.response?.data?.error || 'Не удалось обновить прогресс цели');
        }
    };

    return (
        <div className="goals-container">
            <div className="goals-header">
                <h2>Финансовые цели {loading && <RefreshCw size={16} className="animate-spin" style={{ display: 'inline', marginLeft: '8px', color: '#6366f1' }} />}</h2>
                <p style={{ color: '#64748b', fontSize: '14px', marginTop: '4px' }}>
                    Планируйте крупные покупки и отслеживайте прогресс накоплений
                </p>
            </div>

            {/* ФОРМА СОЗДАНИЯ ЦЕЛИ */}
            <div className="goal-form-card">
                <h3 style={{ fontSize: '16px', marginBottom: '16px', color: '#f8fafc' }}>Создать новую цель</h3>
                {error && (
                    <div style={{ color: '#ef4444', fontSize: '13px', marginBottom: '12px', display: 'flex', alignItems: 'center', gap: '6px' }}>
                        <AlertCircle size={16} /> {error}
                    </div>
                )}
                <form onSubmit={handleCreateGoal}>
                    <div className="form-inline-group">
                        <div className="form-field">
                            <label>Что хотим купить / накопить?</label>
                            <input
                                type="text"
                                placeholder="Например: Отпуск на море"
                                value={name}
                                onChange={(e) => setName(e.target.value)}
                                required
                            />
                        </div>
                        <div className="form-field">
                            <label>Целевая сумма (₽)</label>
                            <input
                                type="number"
                                placeholder="50000"
                                value={targetAmount}
                                onChange={(e) => setTargetAmount(e.target.value)}
                                required
                            />
                        </div>
                        <div className="form-field">
                            <label>Срок (Дедлайн)</label>
                            <input
                                type="date"
                                value={targetDate}
                                onChange={(e) => setTargetDate(e.target.value)}
                            />
                        </div>
                    </div>
                    <div className="form-field" style={{ marginTop: '12px' }}>
                        <label>Описание или заметки (необязательно)</label>
                        <input
                            type="text"
                            placeholder="Дополнительные детали, марка, ссылки или мотивация..."
                            value={description}
                            onChange={(e) => setDescription(e.target.value)}
                        />
                    </div>
                    <button type="submit" className="submit-btn" style={{ width: 'auto', padding: '10px 20px', display: 'flex', alignItems: 'center', gap: '8px', marginTop: '12px' }}>
                        <Plus size={16} />
                        Поставить цель
                    </button>
                </form>
            </div>

            {/* СЕТКА С КАРТОЧКАМИ ЦЕЛЕЙ */}
            <div className="goals-grid">
                {goals.length === 0 && !loading ? (
                    <div style={{ color: '#64748b', fontSize: '14px', gridColumn: '1 / -1', textAlign: 'center', padding: '40px' }}>
                        Список целей пуст. Сформулируйте вашу первую цель выше!
                    </div>
                ) : (
                    goals.map((goal) => {
                        const progressPercent = Math.min(Math.round((goal.current_amount / goal.target_amount) * 100), 100);
                        // Форматируем дату ответа Go для аккуратного вывода (обрезаем таймштамп, если он прилетает)
                        const cleanDate = goal.target_date ? goal.target_date.split('T')[0] : 'Не указан';

                        return (
                            <div key={goal.id} className={`goal-card ${goal.is_completed ? 'completed' : ''}`}>
                                <div>
                                    <div className="goal-card-header">
                                        <h3>{goal.name}</h3>
                                        {goal.is_completed ? (
                                            <CheckCircle size={20} style={{ color: '#10b981' }} />
                                        ) : (
                                            <Target size={20} style={{ color: '#6366f1' }} />
                                        )}
                                    </div>

                                    {goal.description && <p className="goal-desc">{goal.description}</p>}
                                </div>

                                <div>
                                    <div className="goal-progress-info">
                                        <span style={{ color: '#cbd5e1' }}>
                                            {(goal.current_amount || 0).toLocaleString('ru-RU')} / {(goal.target_amount || 0).toLocaleString('ru-RU')} ₽
                                        </span>
                                        <span className="goal-percentage">{progressPercent}%</span>
                                    </div>

                                    <div className="goal-progress-bar-bg">
                                        <div
                                            className="goal-progress-bar-fill"
                                            style={{ width: `${progressPercent}%` }}
                                        />
                                    </div>

                                    <div className="goal-footer-meta">
                                        <div style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                                            <Calendar size={14} />
                                            <span>До: {cleanDate}</span>
                                        </div>

                                        {!goal.is_completed && (
                                            <button
                                                type="button"
                                                className="deposit-btn"
                                                onClick={() => handleQuickDeposit(goal)}
                                            >
                                                <TrendingUp size={14} />
                                                Пополнить
                                            </button>
                                        )}
                                    </div>
                                </div>
                            </div>
                        );
                    })
                )}
            </div>
        </div>
    );
};