# Toolbox

Набор полезных инструментов и утилит для разработки на Go.

## Описание

Toolbox - это библиотека, содержащая набор готовых компонентов для быстрой разработки приложений на Go. Библиотека предоставляет готовые решения для работы с базами данных, аутентификацией, HTTP-запросами, логированием, форматированием, RabbitMQ и другими часто используемыми функциями.

## Структура проекта

```
pkg/
├── auth        - Компоненты для аутентификации и авторизации
├── compressor  - Утилиты для сжатия данных
├── db          - Клиент для работы с базами данных
├── formater    - Инструменты для форматирования данных
├── http        - HTTP-клиент и утилиты
├── logger      - Компоненты для логирования
├── luna        - Специализированные инструменты
├── migrations  - Инструменты для миграций базы данных
├── paramsparser - Парсеры параметров
└── rabbitmq    - Клиент для работы с RabbitMQ
```

## Требования

- Go 1.24.1 или выше

## Установка

```bash
go get github.com/SmirnovND/toolbox
```

## Использование

### Подключение к базе данных

```go
import "github.com/SmirnovND/toolbox/pkg/db"

func main() {
    dsn := "postgres://user:password@localhost:5432/dbname?sslmode=disable"
    database := db.NewDB(dsn)
    // Используйте database для выполнения запросов
}
```

### Логирование HTTP-запросов

```go
import (
    "github.com/SmirnovND/toolbox/pkg/logger"
    "net/http"
)

func main() {
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Обработка запроса
    })

    // Оборачиваем обработчик в middleware для логирования
    loggedHandler := logger.WithLogging(handler)

    http.ListenAndServe(":8080", loggedHandler)
}
```

## Зависимости

- github.com/golang-jwt/jwt/v4 - Работа с JWT токенами
- github.com/golang-migrate/migrate/v4 - Миграции базы данных
- github.com/jmoiron/sqlx - Расширенный SQL клиент
- github.com/rs/zerolog - Логирование
- github.com/streadway/amqp - Работа с RabbitMQ

## Тестирование

Проект содержит тесты для всех пакетов. Для запуска тестов используйте команды из Makefile:

```bash
# Запуск всех тестов
make test

# Запуск тестов с подробным выводом
make test-verbose

# Запуск тестов с покрытием кода
make test-coverage
```

## Makefile

Проект содержит Makefile со следующими командами:

```bash
make                 # Показать справку (выполняется по умолчанию)
make help            # Показать справку
make test            # Запуск тестов
make test-verbose    # Запуск тестов с подробным выводом
make test-coverage   # Запуск тестов с покрытием кода
make build           # Сборка проекта
make clean           # Очистка временных файлов
```

## Лицензия

MIT

## Автор

SmirnovND
