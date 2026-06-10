albums-api — REST API для управления альбомами
Production-ready Go-сервис с чистой архитектурой и отдельным сервисом авторизации.

Архитектура
Проект построен по принципам Clean Architecture (Layered / Hexagonal):
textcmd/app/
└── main.go                  # Точка входа

internal/
├── config/                  # Конфигурация (env + godotenv)
├── db/                      # Подключение к PostgreSQL (pgxpool)
├── repository/              # Репозиторий (абстракция над БД)
├── handlers/                # HTTP-обработчики (Gin)
├── middleware/              # JWT Auth + Logging
├── client/                  # Клиент для auth-service
├── server/                  # Настройка роутера и graceful shutdown
└── models/                  # Доменные модели
Микросервисная структура:

albums-api (Go + Gin) — основной сервис
auth-service (Flask + Flask-RESTX) — отдельный сервис авторизации
Две PostgreSQL БД + Redis

Технические решения:

Чистая архитектура,"Легко тестировать, масштабировать и поддерживать"
Отдельный auth-сервис,"Разделение ответственности, возможность масштабирования авторизации независимо"
JWT + Redis,Быстрая валидация + возможность logout (инвалидация токенов)
Repository pattern,"Абстракция от БД, легко менять реализацию и писать тесты"
Middleware для JWT,Защита только мутирующих операций (POST/PUT/DELETE)
Multi-stage Docker,Маленький финальный образ + удобство разработки
Интеграционные тесты,Проверка реальной интеграции с auth-service
Graceful Shutdown,Корректное завершение сервиса

Как запустить:

hmake up              # Запуск всех сервисов

make test-albums     # Юнит + интеграционные тесты

make logs-albums     # Логи основного сервиса

make test-integration # Только интеграционные тесты

Swagger: http://localhost:8080/swagger/index.html
