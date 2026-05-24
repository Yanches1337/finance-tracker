# Этап 1: Сборка статики React
FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
COPY vite.config.js ./

RUN npm install

# Копируем исходники
COPY frontend/ ./frontend/

# Запускаем сборку (Vite создаст папку внутри frontend/finance-app/dist)
RUN npm run build

# Этап 2: Раздача через Nginx
FROM nginx:1.25-alpine

# Копируем скомпилированный фронтенд из нового места
COPY --from=builder /app/dist /usr/share/nginx/html
# Копируем конфиг Nginx
COPY frontend/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]