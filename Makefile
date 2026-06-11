.PHONY: up down build logs test test-positive test-negative shell db-shell redis-shell fresh help

# ================== Основные команды ==================

up:
	docker-compose up --build -d

down:
	docker-compose down -v

build:
	docker-compose build --no-cache

logs-auth:
	docker-compose logs -f auth-service
logs-albums:
	docker-compose logs -f app

fresh:
	docker-compose down -v && docker-compose up --build -d

# ================== Тесты ==================

test-auth:
	docker-compose exec auth-service pytest tests/ -v --tb=short

test-auth-positive:
	docker-compose exec auth-service pytest tests/test_users_positive_cases.py -v --tb=short

test-auth-negative:
	docker-compose exec auth-service pytest tests/test_users_negative_cases.py -v --tb=short

test-albums:
	docker-compose exec app go test ./internal/handlers/tests -v

test-albums-pytest-allure:
	docker-compose exec -T auth-service pytest tests/tests_albums_api/ -v --alluredir=allure-results
	allure serve allure-results

test-albums-integration:
	docker-compose exec app go test ./internal/handlers/tests -run TestIntegration -v

# ================== Доступ к контейнерам ==================

shell-auth:
	docker-compose exec auth-service bash

# ================== Работа с базами ==================

db-auth-shell:
	docker-compose exec auth-postgres psql -U postgres -d authdb

db-albums-shell:
	docker-compose exec postgres psql -U postgres -d albumsdb

redis-shell:
	docker-compose exec redis redis-cli

# ================== Полезные команды ==================

restart-auth:
	docker-compose restart auth-service

status:
	docker-compose ps

# Полная очистка + запуск
reset:
	docker-compose down -v && docker-compose up --build -d

help:
	@echo "=== Основные команды ==="
	@echo "  make up              - Запустить все сервисы"
	@echo "  make down            - Остановить и удалить контейнеры"
	@echo "  make fresh / reset   - Полная пересборка с нуля"
	@echo ""
	@echo "=== Тесты auth-service ==="
	@echo "  make test-auth            - Запустить все тесты auth-service"
	@echo "  make test-auth-positive   - Только позитивные тесты auth-service"
	@echo "  make test-auth-negative   - Только негативные тесты auth-service"
	@echo ""
	@echo "=== Тесты albums ==="
	@echo "  make test-albums            - Запустить все тесты albums"
	@echo "  make test-albums-integration - Запустить только интеграционные тесты albums"
	@echo "  make test-albums-pytest-allure - Запустить pytest с Allure для albums API"
	@echo "=== Доступ к сервисам ==="
	@echo "  make shell           - Зайти в auth-service (bash)"
	@echo "  make db-shell        - Зайти в psql authdb"
	@echo "  make db-albums       - Зайти в psql albumsdb"
	@echo "  make redis-shell     - Зайти в redis-cli"
	@echo "  make logs-auth            - Посмотреть логи auth-service"
	@echo "  make logs-albums          - Посмотреть логи albums"
	@echo ""
	@echo "  make help            - Показать эту справку"