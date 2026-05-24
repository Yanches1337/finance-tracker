# Finance Tracker — Трекер Личных Финансов

Современное веб-приложение для учёта доходов, расходов, финансового планирования и аналитики.

![Go](https://img.shields.io/badge/Go-1.23-00ADD8)
![React](https://img.shields.io/badge/React-19-61DAFB)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1)
![Redis](https://img.shields.io/badge/Redis-DC382D)
![Docker](https://img.shields.io/badge/Docker-2496ED)

## Основные возможности

- Учёт доходов и расходов по категориям
- Гибкое управление категориями (с цветами и иконками)
- Постановка и отслеживание финансовых целей
- Подробная аналитика и визуализация данных
- Интерактивный дашборд
- Генерация и экспорт отчётов (CSV, JSON, PDF)
- Адаптивный дизайн (работает на мобильных устройствах)
- Безопасная авторизация с JWT + Refresh Token

## Технологический стек

**Backend:**
- Go 1.23
- Clean Architecture + DDD
- Gin Framework
- PostgreSQL + pgxpool
- Redis (хранение refresh-токенов и сессий)
- Docker + Docker Compose

**Frontend:**
- React 19 + Vite
- Tailwind CSS
- Axios с автоматическим обновлением токенов

---

## Быстрый запуск

### 1. Клонируйте репозиторий

```bash
git clone https://github.com/Yanches1337/finance-tracker.git
cd finance-tracker
```
### 2. Запустите проект

```bash
docker compose up -d --build
```

После запуска приложение будет доступно по адресам:

- Фронтенд: http://localhost:3000
- API: http://localhost:8000
- Swagger-документация: http://localhost:8000/swagger/index.html

## Структура проекта
### Backend (Go)
Проект построен по принципам **Clean Architecture** (Hexagonal Architecture)

```text
backend/
├── cmd/
│   └── server/
│       └── main.go                  # Точка входа приложения
├── configs/
│   └── config.yaml                  # Конфигурация
├── internal/
│   ├── adapters/
│   │   ├── interfaces/              # Интерфейсы репозиториев
│   │   ├── postgres/                # Реализация PostgreSQL
│   │   │   ├── db.go
│   │   │   └── repository_impls/
│   │   └── redis/                   # Работа с Redis
│   ├── api/
│   │   └── http/
│   │       ├── handlers/            # HTTP-обработчики
│   │       ├── middleware/          # Middleware (auth, cors и т.д.)
│   │       └── router.go
│   ├── domain/                      # Основные бизнес-сущности
│   │   ├── user.go
│   │   ├── transaction.go
│   │   ├── category.go
│   │   ├── goal.go
│   │   └── types.go
│   ├── services/                    # Бизнес-логика
│   │   ├── auth_service.go
│   │   ├── transaction_service.go
│   │   ├── category_service.go
│   │   ├── goal_service.go
│   │   ├── analytics_service.go
│   │   └── report_service.go
│   └── utils/                       # Вспомогательные утилиты
├── migrations/                      # SQL-миграции
├── test/                            # Тесты
└── Dockerfile
```

### Frontend (React)
```text
frontend/
└── finance-app/
├── src/
│   ├── api/                     # Настройка axios
│   ├── components/              # Переиспользуемые компоненты
│   ├── context/                 # React Context (Auth и др.)
│   ├── App.jsx
│   └── main.jsx
├── index.html
├── Dockerfile
└── vite.config.js
```

## Основные возможности
- Авторизация — регистрация, вход, обновление токенов
- Операции — добавление доходов и расходов
- Категории — создание, редактирование, удаление
- Цели — постановка финансовых целей с отслеживанием прогресса
- Дашборд — текущий баланс, статистика, диаграммы
- Аналитика — графики по категориям и периодам
- Отчёты — генерация и скачивание в разных форматах

## Тестирование API
Перейдите по ссылке:
http://localhost:8000/swagger/index.html

## Планы развития
- Мультивалютность
- Импорт транзакций из CSV и банковских выписок
- Уведомления и напоминания
- Бюджетирование по категориям
- Искусственный интеллект для более глубокой и полной аналитики 
- Полноценное мобильное приложение