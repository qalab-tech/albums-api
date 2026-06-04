-- =============================================
-- Автоматическая инициализация authdb
-- =============================================

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL,
    email VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Удаляем старого тестового пользователя
DELETE FROM users WHERE username = 'test';

-- Создаём тестового пользователя (password = "test")
INSERT INTO users (username, hashed_password, email)
VALUES (
    'test',
    '$2b$12$Y8z9K7vL3mN5pQ7rT9vW2uX4yZ6aB8cD0eF2gH4iJ6kL8mN0pQ2r',
    'test@example.com'
);

SELECT '✅ Auth DB initialized successfully with test user' as message;
SELECT * FROM users;