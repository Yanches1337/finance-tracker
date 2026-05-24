import React, { useState } from 'react';
import { useAuth } from './context/AuthContext';
import { AuthForm } from './components/AuthForm';
import { MainLayout } from './components/MainLayout';
import { Operations } from './components/Operations';
import { Dashboard } from './components/Dashboard';
import { Goals } from './components/Goals';
import { Analytics } from './components/Analytics';
import { Reports } from './components/Reports';

function App() {
    const { user, loading } = useAuth();
    const [activeTab, setActiveTab] = useState('dashboard');

    // Если идет проверка сессии через Redis при обновлении страницы
    if (loading) {
        return (
            <div style={{
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                height: '100vh',
                backgroundColor: '#0f172a',
                color: '#ffffff'
            }}>
                Загрузка конфигурации...
            </div>
        );
    }

    // Если пользователь не авторизован — показываем красивый экран логина/регистрации
    if (!user) {
        return <AuthForm />;
    }

    const renderContent = () => {
        switch (activeTab) {
            case 'dashboard':
                return <Dashboard />;
            case 'operations':
                return <Operations />;
            case 'analytics':
                return <Analytics />;
            case 'goals':
                return <Goals />;
            case 'reports':
                return <Reports />;
            default:
                return <div>Экран не найден</div>;
        }
    };

    // Если авторизован — рендерим основную разметку с переключателем экранов
    return (
        <MainLayout activeTab={activeTab} setActiveTab={setActiveTab}>
            {renderContent()}
        </MainLayout>
    );
}

export default App;