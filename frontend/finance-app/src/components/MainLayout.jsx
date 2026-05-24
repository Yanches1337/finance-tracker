import React from 'react';
import { useAuth } from '../context/AuthContext';
import {
    LayoutDashboard,
    Wallet,
    PieChart,
    Target,
    FileText,
    LogOut,
    User as UserIcon
} from 'lucide-react';
import './MainLayout.css';

export const MainLayout = ({ children, activeTab, setActiveTab }) => {
    const { user, logout } = useAuth();

    const navigationItems = [
        { id: 'dashboard', label: 'Обзор', icon: LayoutDashboard },
        { id: 'operations', label: 'Операции', icon: Wallet },
        { id: 'analytics', label: 'Аналитика', icon: PieChart },
        { id: 'goals', label: 'Цели', icon: Target },
        { id: 'reports', label: 'Отчеты', icon: FileText },
    ];

    const handleLogout = async () => {
        if (window.confirm('Вы уверены, что хотите выйти?')) {
            await logout();
        }
    };

    return (
        <div className="layout-container">
            {/* --- ДЕСКТОПНЫЙ САЙДБАР --- */}
            <aside className="sidebar">
                <div>
                    <div className="sidebar-brand">
                        <Wallet size={24} />
                        <span>FinanceTracker</span>
                    </div>

                    <nav className="sidebar-menu">
                        {navigationItems.map((item) => {
                            const Icon = item.icon;
                            return (
                                <button
                                    key={item.id}
                                    className={`menu-item ${activeTab === item.id ? 'active' : ''}`}
                                    onClick={() => setActiveTab(item.id)}
                                >
                                    <Icon size={20} />
                                    <span>{item.label}</span>
                                </button>
                            );
                        })}
                    </nav>
                </div>

                <div className="sidebar-footer">
                    <div className="user-info" style={{ display: 'flex', flexDirection: 'column', alignItems: 'flex-start', gap: '2px' }}>
                        <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                            <UserIcon size={18} className="text-slate-400" />
                            <span className="user-name">{user?.name || 'Пользователь'}</span>
                        </div>
                        {user?.email && <span style={{ fontSize: '11px', color: '#64748b', paddingLeft: '26px' }}>{user.email}</span>}
                    </div>
                    <button className="menu-item" onClick={handleLogout} style={{ color: '#f87171' }}>
                        <LogOut size={20} />
                        <span>Выйти</span>
                    </button>
                </div>
            </aside>

            {/* --- МОБИЛЬНАЯ ШАПКА --- */}
            <header className="mobile-header">
                <div className="mobile-brand">FinanceTracker</div>
                <button className="logout-btn-mobile" onClick={handleLogout}>
                    <LogOut size={22} style={{ color: '#f87171' }} />
                </button>
            </header>

            {/* --- ОСНОВНАЯ ЗОНА КОНТЕНТА --- */}
            <main className="main-content">
                {children}
            </main>

            {/* --- МОБИЛЬНОЕ МЕНЮ --- */}
            <nav className="mobile-nav">
                {navigationItems.map((item) => {
                    const Icon = item.icon;
                    return (
                        <button
                            key={item.id}
                            className={`mobile-menu-item ${activeTab === item.id ? 'active' : ''}`}
                            onClick={() => setActiveTab(item.id)}
                        >
                            <Icon size={22} />
                            <span>{item.label}</span>
                        </button>
                    );
                })}
            </nav>
        </div>
    );
};