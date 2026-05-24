import React, { createContext, useContext, useState, useEffect } from 'react';
import { api, setAccessToken } from '../api/axios';

const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);

    // Проверка сессии при первой загрузке вкладки
    useEffect(() => {
        const checkAuth = async () => {
            const refreshToken = localStorage.getItem('refresh_token');
            if (refreshToken) {
                try {
                    // Пробуем получить новый access_token с помощью имеющегося refresh
                    const response = await api.post('/auth/refresh', {}, {
                        headers: { 'X-Refresh-Token': refreshToken }
                    });

                    setAccessToken(response.data.access_token);
                    localStorage.setItem('refresh_token', response.data.refresh_token);

                    // Запрашиваем данные текущего пользователя из профиля
                    const userResponse = await api.get('/protected/me');
                    setUser(userResponse.data);
                } catch (err) {
                    console.error("Session restoration failed", err);
                    localStorage.removeItem('refresh_token');
                }
            }
            setLoading(false);
        };

        checkAuth();
    }, []);

    const login = async (email, password) => {
        const response = await api.post('/auth/login', { email, password });
        const { access_token, refresh_token, user: userData } = response.data;

        setAccessToken(access_token);
        localStorage.setItem('refresh_token', refresh_token);
        setUser(userData);
        return userData;
    };

    const register = async (name, email, password) => {
        await api.post('/auth/register', { name, email, password });
    };

    const logout = async () => {
        const refreshToken = localStorage.getItem('refresh_token');
        try {
            if (refreshToken) {
                await api.post('/auth/logout', {}, {
                    headers: { 'X-Refresh-Token': refreshToken }
                });
            }
        } catch (err) {
            console.error("Logout on backend failed", err);
        } finally {
            setAccessToken(null);
            localStorage.removeItem('refresh_token');
            setUser(null);
        }
    };

    return (
        <AuthContext.Provider value={{ user, loading, login, register, logout }}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => useContext(AuthContext);