import React, { useState, useEffect } from 'react';
import { FileText, Download, FileSpreadsheet, AlertTriangle, CheckCircle2, RefreshCw } from 'lucide-react';
import { api } from '../api/axios'; // Твой инстанс Axios
import './Reports.css';

export const Reports = () => {
    // Параметры GenerateReportRequest
    const [format, setFormat] = useState('csv');
    const [fromDate, setFromDate] = useState('2026-05-01');
    const [toDate, setToDate] = useState('2026-05-24');

    const [loading, setLoading] = useState(false);
    const [historyLoading, setHistoryLoading] = useState(false);
    const [statusMessage, setStatusMessage] = useState(null);

    // Сюда сохраняем реальный массив []domain.Report с бэкенда
    const [reportsHistory, setReportsHistory] = useState([]);

    // 1. Метод выгрузки истории отчетов из БД
    const fetchReportsHistory = async () => {
        setHistoryLoading(true);
        try {
            const response = await api.get('/protected/reports');
            setReportsHistory(response.data || []);
        } catch (err) {
            console.error("Не удалось загрузить историю отчетов:", err);
        } finally {
            setHistoryLoading(false);
        }
    };

    // Загружаем историю один раз при входе на вкладку
    useEffect(() => {
        fetchReportsHistory();
    }, []);

    // 2. Вызов POST /protected/reports для генерации новой записи
    const handleGenerateReport = async (e) => {
        e.preventDefault();
        if (format !== 'csv') return; // Защита фронтенда

        setLoading(true);
        setStatusMessage(null);

        // Подготовка JSON структуры под domain.GenerateReportRequest бэкенда
        const requestData = {
            format: format,
            // Передаем даты в чистом YYYY-MM-DD или ISO, в зависимости от парсера Go.
            // Если Go ждет чисто строку "2026-05-01", оставляем fromDate.
            // Если time.Time — то .toISOString()
            from_date: new Date(fromDate).toISOString(),
            to_date: new Date(toDate).toISOString()
        };

        try {
            const response = await api.post('/protected/reports', requestData);

            // Сервер вернул созданный объект отчета, добавляем его в начало списка
            if (response.data) {
                setReportsHistory(prev => [response.data, ...prev]);
            } else {
                // На всякий случай перечитываем всю историю, если бэкенд отдал пустой body
                fetchReportsHistory();
            }

            setStatusMessage({ type: 'success', text: 'Отчет успешно сгенерирован и добавлен в историю!' });
        } catch (err) {
            console.error(err);
            setStatusMessage({
                type: 'error',
                text: err.response?.data?.error || 'Не удалось отправить запрос на генерацию отчета'
            });
        } finally {
            setLoading(false);
        }
    };

    // 3. Вызов GET /protected/reports/{id}/download с обработкой Blob
    const handleDownloadFile = async (id, fileName) => {
        try {
            setStatusMessage({ type: 'success', text: `Началось скачивание файла ${fileName}...` });

            // Запрашиваем файл как бинарный поток (blob)
            const response = await api.get(`/protected/reports/${id}/download`, {
                responseType: 'blob'
            });

            // Создаем безопасную ссылку в памяти браузера
            const blob = new Blob([response.data], { type: 'text/csv;charset=utf-8;' });
            const downloadUrl = window.URL.createObjectURL(blob);

            // Программный клик для скачивания
            const link = document.createElement('a');
            link.href = downloadUrl;
            link.setAttribute('download', fileName || `report_${id}.csv`);
            document.body.appendChild(link);
            link.click();

            // Чистим за собой DOM и память
            document.body.removeChild(link);
            window.URL.revokeObjectURL(downloadUrl);
        } catch (err) {
            console.error("Ошибка скачивания файла:", err);
            setStatusMessage({ type: 'error', text: 'Ошибка при чтении файла с сервера.' });
        }
    };

    return (
        <div className="reports-container">
            <div className="analytics-header" style={{ marginBottom: '24px' }}>
                <h2>Выгрузка отчетов и документов</h2>
                <p style={{ color: '#64748b', fontSize: '14px', marginTop: '4px' }}>
                    Экспортируйте ваши финансовые транзакции за любые периоды
                </p>
            </div>

            <div className="reports-grid">

                {/* ЛЕВАЯ ПАНЕЛЬ: ЗАПРОС НА ГЕНЕРАЦИЮ */}
                <div className="report-panel-card">
                    <h3 style={{ fontSize: '16px', color: '#f8fafc', marginBottom: '16px' }}>Параметры отчета</h3>

                    <form onSubmit={handleGenerateReport}>

                        <label style={{ display: 'block', fontSize: '13px', color: '#94a3b8', marginBottom: '8px' }}>
                            Формат документа
                        </label>
                        <div className="format-selector-grid">
                            <div
                                className={`format-option ${format === 'csv' ? 'selected' : ''}`}
                                onClick={() => setFormat('csv')}
                            >
                                <span className="format-badge">CSV</span>
                                <span className="format-status" style={{ color: '#34d399' }}>Доступен</span>
                            </div>

                            <div className="format-option disabled" title="Временно недоступно на бэкенде">
                                <span className="format-badge" style={{ color: '#64748b' }}>PDF</span>
                                <span className="format-status">В разработке</span>
                            </div>
                        </div>

                        <div className="form-field" style={{ marginBottom: '12px' }}>
                            <label>С даты (From)</label>
                            <input
                                type="date"
                                value={fromDate}
                                onChange={(e) => setFromDate(e.target.value)}
                                required
                            />
                        </div>

                        <div className="form-field" style={{ marginBottom: '20px' }}>
                            <label>По дату (To)</label>
                            <input
                                type="date"
                                value={toDate}
                                onChange={(e) => setToDate(e.target.value)}
                                required
                            />
                        </div>

                        <button
                            type="submit"
                            className="submit-btn"
                            disabled={loading}
                            style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: '8px' }}
                        >
                            {loading ? (
                                <RefreshCw size={16} className="animate-spin" />
                            ) : (
                                <FileSpreadsheet size={16} />
                            )}
                            {loading ? 'Формирование...' : 'Сгенерировать отчет'}
                        </button>
                    </form>

                    {statusMessage && (
                        <div style={{
                            marginTop: '16px',
                            padding: '12px',
                            borderRadius: '8px',
                            fontSize: '13px',
                            backgroundColor: statusMessage.type === 'success' ? 'rgba(16, 185, 129, 0.1)' : 'rgba(239, 68, 68, 0.1)',
                            color: statusMessage.type === 'success' ? '#34d399' : '#f87171',
                            display: 'flex',
                            alignItems: 'center',
                            gap: '8px'
                        }}>
                            {statusMessage.type === 'success' ? <CheckCircle2 size={16} /> : <AlertTriangle size={16} />}
                            <span>{statusMessage.text}</span>
                        </div>
                    )}
                </div>

                {/* ПРАВАЯ ПАНЕЛЬ: ТАБЛИЦА ИСТОРИИ И ВЫГРУЗКИ */}
                <div className="history-table-container">
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
                        <h3 style={{ fontSize: '16px', color: '#f8fafc', margin: 0 }}>Готовые файлы</h3>
                        {historyLoading && <RefreshCw size={16} className="animate-spin" style={{ color: '#6366f1' }} />}
                    </div>

                    {reportsHistory.length === 0 && !historyLoading ? (
                        <p style={{ color: '#64748b', fontSize: '14px' }}>Вы еще не генерировали отчеты.</p>
                    ) : (
                        <table className="reports-table">
                            <thead>
                            <tr>
                                <th>Дата создания</th>
                                <th>Период отчета</th>
                                <th>Формат</th>
                                <th style={{ textAlign: 'right' }}>Действие</th>
                            </tr>
                            </thead>
                            <tbody>
                            {reportsHistory.map((report) => {
                                // Форматируем дату создания (убираем T...Z из таймштампа бэкенда Go)
                                const createdDate = report.created_at
                                    ? new Date(report.created_at).toLocaleDateString('ru-RU')
                                    : '—';

                                const cleanFrom = report.from_date ? report.from_date.split('T')[0] : '—';
                                const cleanTo = report.to_date ? report.to_date.split('T')[0] : '—';

                                return (
                                    <tr key={report.id}>
                                        <td>
                                            <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                                                <FileText size={16} style={{ color: '#a5b4fc' }} />
                                                <span>{createdDate}</span>
                                            </div>
                                        </td>
                                        <td style={{ fontSize: '13px', color: '#94a3b8' }}>
                                            {cleanFrom} — {cleanTo}
                                        </td>
                                        <td>
                                            <span style={{
                                                fontSize: '11px',
                                                padding: '2px 6px',
                                                backgroundColor: '#0f172a',
                                                borderRadius: '4px',
                                                textTransform: 'uppercase',
                                                fontWeight: '700',
                                                color: '#34d399'
                                            }}>
                                                {report.format}
                                            </span>
                                        </td>
                                        <td style={{ textAlign: 'right' }}>
                                            <button
                                                className="download-link-btn"
                                                onClick={() => handleDownloadFile(report.id, report.file_name)}
                                            >
                                                <Download size={14} />
                                                Скачать
                                            </button>
                                        </td>
                                    </tr>
                                );
                            })}
                            </tbody>
                        </table>
                    )}
                </div>

            </div>
        </div>
    );
};