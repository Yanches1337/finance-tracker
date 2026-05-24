import React, { useState, useEffect } from 'react';
import { PlusCircle, ArrowUpRight, ArrowDownLeft, RefreshCw, FolderPlus } from 'lucide-react';
import { api } from '../api/axios';
import './Operations.css';

// Список стандартных категорий, которые мы автоматически зальем в базу, если там пусто
const STANDARD_CATEGORIES = [
    { name: 'Продукты', type: 'expense' },
    { name: 'Кафе и рестораны', type: 'expense' },
    { name: 'Транспорт', type: 'expense' },
    { name: 'Развлечения', type: 'expense' },
    { name: 'Коммунальные платежи', type: 'expense' },
    { name: 'Зарплата', type: 'income' },
    { name: 'Фриланс', type: 'income' }
];

export const Operations = () => {
    const [categories, setCategories] = useState([]);
    const [transactions, setTransactions] = useState([]);
    const [loading, setLoading] = useState(false);

    // Состояния основной формы операции
    const [type, setType] = useState('expense');
    const [amount, setAmount] = useState('');
    const [categoryId, setCategoryId] = useState('');
    const [date, setDate] = useState(new Date().toISOString().split('T')[0]);
    const [comment, setComment] = useState(''); // Это состояние инпута
    const [error, setError] = useState('');

    // Состояние для создания новой категории
    const [newCatName, setNewCatName] = useState('');
    const [showCatForm, setShowCatForm] = useState(false);
    const [catError, setCatError] = useState('');

    // Функция проверки и инициализации категорий + загрузка транзакций
    const fetchData = async () => {
        setLoading(true);
        try {
            // 1. Запрашиваем категории с бэкенда
            const categoriesRes = await api.get('/protected/categories');
            let dbCategories = categoriesRes.data || [];

            // 2. ЕСЛИ В БАЗЕ ПУСТО: автоматически регистрируем стандартные категории на бэкенде
            if (dbCategories.length === 0) {
                console.log("База категорий пуста. Инициализируем стандартные категории...");

                // Отправляем все стандартные категории на бэкенд параллельно
                await Promise.all(
                    STANDARD_CATEGORIES.map(cat =>
                        api.post('/protected/categories', { name: cat.name, type: cat.type })
                    )
                );

                // Перезапрашиваем категории, теперь они придут с реальными ID из базы данных
                const updatedCategoriesRes = await api.get('/protected/categories');
                dbCategories = updatedCategoriesRes.data || [];
            }

            setCategories(dbCategories);

            // 3. Загружаем транзакции
            const transactionsRes = await api.get('/protected/transactions');
            setTransactions(transactionsRes.data || []);

        } catch (err) {
            console.error("Ошибка синхронизации с бэкендом:", err);
            setError("Ошибка связи с сервером при обновлении данных");
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchData();
    }, []);

    // Фильтруем категории (теперь все они со 100% валидными ID из БД)
    const filteredCategories = categories.filter(cat => cat.type === type);

    // Хендлер создания кастомной категории пользователем
    const handleCreateCategory = async (e) => {
        e.preventDefault();
        setCatError('');
        if (!newCatName.trim()) return;

        try {
            await api.post('/protected/categories', {
                name: newCatName.trim(),
                type: type
            });

            setNewCatName('');
            setShowCatForm(false);

            // Сразу обновляем список категорий из БД
            const categoriesRes = await api.get('/protected/categories');
            setCategories(categoriesRes.data || []);
        } catch (err) {
            setCatError(err.response?.data?.error || 'Не удалось создать категорию');
        }
    };

    // Хендлер добавления транзакции
    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');

        const parsedAmount = parseFloat(amount);
        if (isNaN(parsedAmount) || parsedAmount <= 0) {
            setError('Сумма должна быть положительным числом');
            return;
        }

        if (!categoryId) {
            setError('Пожалуйста, выберите категорию');
            return;
        }

        try {
            // Исправляем баг с комментарием: отправляем и comment, и description
            // чтобы точно попасть в ожидаемое бэкендом поле структуры.
            await api.post('/protected/transactions', {
                amount: parsedAmount,
                type: type,
                category_id: parseInt(categoryId, 10),
                date: `${date}T00:00:00Z`,
                description: comment.trim()   // На случай если в Go тег `json:"description"`
            });

            // Сбрасываем форму и обновляем списки
            setAmount('');
            setCategoryId('');
            setComment('');
            fetchData();
        } catch (err) {
            setError(err.response?.data?.error || 'Не удалось сохранить операцию');
        }
    };

    return (
        <div className="operations-container">
            {/* ЛЕВАЯ КОЛОНКА: ФОРМА ОПЕРАЦИИ */}
            <div className="op-card">
                <h3>Новая операция</h3>

                <div className="type-selector">
                    <button
                        type="button"
                        className={`type-btn ${type === 'expense' ? 'active expense' : ''}`}
                        onClick={() => { setType('expense'); setCategoryId(''); setError(''); }}
                    >
                        Расход
                    </button>
                    <button
                        type="button"
                        className={`type-btn ${type === 'income' ? 'active income' : ''}`}
                        onClick={() => { setType('income'); setCategoryId(''); setError(''); }}
                    >
                        Доход
                    </button>
                </div>

                {error && (
                    <div style={{ color: '#ef4444', fontSize: '13px', marginBottom: '12px', fontWeight: '500' }}>
                        {error}
                    </div>
                )}

                <form onSubmit={handleSubmit}>
                    <div className="form-field">
                        <label>Сумма операции</label>
                        <input
                            type="number"
                            step="0.01"
                            placeholder="0.00"
                            value={amount}
                            onChange={(e) => setAmount(e.target.value)}
                            required
                        />
                    </div>

                    <div className="form-field">
                        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '4px' }}>
                            <label style={{ margin: 0 }}>Категория</label>
                            <button
                                type="button"
                                onClick={() => setShowCatForm(!showCatForm)}
                                style={{ background: 'none', border: 'none', color: '#6366f1', fontSize: '12px', cursor: 'pointer', display: 'flex', alignItems: 'center', gap: '4px' }}
                            >
                                <FolderPlus size={14} />
                                {showCatForm ? 'Скрыть панель' : 'Добавить свою'}
                            </button>
                        </div>

                        {/* Форма добавления новой категории */}
                        {showCatForm && (
                            <div style={{ background: '#0f172a', padding: '10px', borderRadius: '6px', marginBottom: '10px', border: '1px solid #334155' }}>
                                <div style={{ display: 'flex', gap: '8px' }}>
                                    <input
                                        type="text"
                                        placeholder="Название категории"
                                        value={newCatName}
                                        onChange={(e) => setNewCatName(e.target.value)}
                                        style={{ flex: 1, padding: '6px', fontSize: '13px' }}
                                    />
                                    <button
                                        type="button"
                                        onClick={handleCreateCategory}
                                        style={{ padding: '6px 12px', backgroundColor: '#6366f1', color: '#fff', border: 'none', borderRadius: '4px', fontSize: '13px', cursor: 'pointer' }}
                                    >
                                        Создать
                                    </button>
                                </div>
                                {catError && <p style={{ color: '#ef4444', fontSize: '11px', margin: '4px 0 0 0' }}>{catError}</p>}
                            </div>
                        )}

                        <select
                            value={categoryId}
                            onChange={(e) => setCategoryId(e.target.value)}
                            required
                        >
                            <option value="">-- Выберите категорию --</option>
                            {filteredCategories.map(cat => (
                                <option key={cat.id} value={cat.id}>
                                    {cat.name}
                                </option>
                            ))}
                        </select>
                    </div>

                    <div className="form-field">
                        <label>Дата</label>
                        <input
                            type="date"
                            value={date}
                            onChange={(e) => setDate(e.target.value)}
                            required
                        />
                    </div>

                    <div className="form-field">
                        <label>Комментарий / Описание</label>
                        <textarea
                            rows="2"
                            placeholder="Например: Купил продукты на неделю"
                            value={comment}
                            onChange={(e) => setComment(e.target.value)}
                        />
                    </div>

                    <button type="submit" className="submit-btn" style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: '8px' }}>
                        <PlusCircle size={18} />
                        Добавить операцию
                    </button>
                </form>
            </div>

            {/* ПРАВАЯ КОЛОНКА: ИСТОРИЯ */}
            <div className="op-card">
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '12px' }}>
                    <h3>История операций</h3>
                    {loading && <RefreshCw size={16} className="animate-spin" style={{ color: '#6366f1' }} />}
                </div>

                <div className="transactions-list">
                    {transactions.length === 0 ? (
                        <div className="no-data">Операции в базе данных отсутствуют.</div>
                    ) : (
                        transactions.map((tx) => {
                            const catObj = categories.find(c => c.id === tx.category_id);
                            const categoryName = catObj ? catObj.name : `Категория #${tx.category_id}`;
                            const cleanDate = tx.date ? tx.date.split('T')[0] : '';

                            // Выводим то текстовое поле, которое придет от бэкенда
                            const displayComment = tx.description || tx.comment || '';

                            return (
                                <div key={tx.id} className={`tx-item ${tx.type}`}>
                                    <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
                                        <div style={{
                                            padding: '8px',
                                            borderRadius: '50%',
                                            backgroundColor: tx.type === 'income' ? 'rgba(16, 185, 129, 0.1)' : 'rgba(239, 68, 68, 0.1)',
                                            color: tx.type === 'income' ? '#10b981' : '#ef4444',
                                            display: 'flex'
                                        }}>
                                            {tx.type === 'income' ? <ArrowUpRight size={18} /> : <ArrowDownLeft size={18} />}
                                        </div>

                                        <div className="tx-info">
                                            <h4>{categoryName}</h4>
                                            {displayComment && <p style={{ color: '#94a3b8', marginBottom: '2px' }}>{displayComment}</p>}
                                            <p>{cleanDate}</p>
                                        </div>
                                    </div>

                                    <div className={`tx-amount ${tx.type}`}>
                                        {tx.type === 'income' ? '+' : '-'} {(tx.amount || 0).toLocaleString('ru-RU')} ₽
                                    </div>
                                </div>
                            );
                        })
                    )}
                </div>
            </div>
        </div>
    );
};