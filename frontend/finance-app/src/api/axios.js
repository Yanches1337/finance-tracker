import axios from 'axios';

// Базовый инстанс для обычных запросов
export const api = axios.create({
    baseURL: '/api/v1',
    headers: {
        'Content-Type': 'application/json',
    },
});

// Хранилище токенов в памяти (in-memory) для безопасности,
// а refresh_token временно сохраним в localStorage для персистентности сессии.
let accessToken = null;

export const setAccessToken = (token) => {
    accessToken = token;
};

// Интерцептор ЗАПРОСА: добавляет Access Token в заголовки protected-роутов
api.interceptors.request.use(
    (config) => {
        // Если запрос идет к защищенному роуту и у нас есть токен — добавляем его
        if (config.url.includes('/protected') && accessToken) {
            config.headers['Authorization'] = `Bearer ${accessToken}`;
        }
        return config;
    },
    (error) => Promise.reject(error)
);

// Интерцептор ОТВЕТА: перехватывает 401 и обновляет токены
api.interceptors.response.use(
    (response) => response,
    async (error) => {
        const originalRequest = error.config;

        // Если ошибка 401 и мы еще не пробовали обновить токен (избегаем бесконечного цикла)
        if (error.response?.status === 401 && !originalRequest._retry) {
            originalRequest._retry = true;
            const refreshToken = localStorage.getItem('refresh_token');

            if (!refreshToken) {
                return Promise.reject(error);
            }

            try {
                // Запрос на обновление токенов.
                // Твой бэкенд принимает refresh-токен в заголовке X-Refresh-Token
                const response = await axios.post(
                    '/api/v1/auth/refresh',
                    {},
                    {
                        headers: { 'X-Refresh-Token': refreshToken },
                    }
                );

                const { access_token, refresh_token } = response.data;

                // Обновляем данные в приложении
                setAccessToken(access_token);
                localStorage.setItem('refresh_token', refresh_token);

                // Повторяем изначальный запрос с новым токеном
                originalRequest.headers['Authorization'] = `Bearer ${access_token}`;
                return axios(originalRequest);
            } catch (refreshError) {
                // Если рефреш тоже протух — разлогиниваем пользователя
                localStorage.removeItem('refresh_token');
                setAccessToken(null);
                window.location.href = '/login';
                return Promise.reject(refreshError);
            }
        }

        return Promise.reject(error);
    }
);