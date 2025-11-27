## Rocket Factory

![Coverage](https://img.shields.io/badge/coverage-98%25-brightgreen)

### Описание проекта
Имитируем сборку ракет

### Структура проекта

Проект использует чистую архитектуру с разделением на слои:

```
.
├── deploy/
│   ├── compose/
│   │   ├── core/           # Общая сеть микросервисов
│   │   ├── inventory/      # MongoDB для InventoryService
│   │   └── order/          # PostgreSQL для OrderService
│   └── env/                # Шаблоны переменных окружения
│       ├── .env.template   # Главный шаблон
│       ├── generate-env.sh # Скрипт генерации
│       ├── order.env.template
│       ├── inventory.env.template
│       └── payment.env.template
│
├── inventory/              # Сервис управления деталями (MongoDB)
│   └── internal/
│       ├── api/            # gRPC хендлеры
│       ├── service/        # Бизнес-логика
│       ├── repository/     # Хранилище данных (MongoDB)
│       ├── model/          # Модели сервисного слоя
│       └── converter/      # Конвертеры между слоями
│
├── order/                  # Сервис заказов (PostgreSQL)
│   ├── migrations/         # SQL миграции (goose)
│   └── internal/
│       ├── api/            # REST/HTTP хендлеры
│       ├── service/        # Бизнес-логика
│       ├── repository/     # Хранилище данных (PostgreSQL)
│       ├── migrator/       # Миграции при старте
│       ├── client/         # Клиенты внешних сервисов
│       ├── model/          # Модели сервисного слоя
│       └── converter/      # Конвертеры между слоями
│
├── payment/                # Сервис оплаты
│   └── internal/
│       ├── api/            # gRPC хендлеры
│       ├── service/        # Бизнес-логика
│       └── model/          # Модели
│
├── platform/               # Платформенные библиотеки
│   └── pkg/
│       ├── closer/         # Graceful shutdown
│       ├── logger/         # Структурированное логирование (zap)
│       └── testcontainers/ # Testcontainers для интеграционных тестов
│           ├── mongo/      # MongoDB контейнер
│           ├── app/        # Контейнер приложения
│           ├── network/    # Docker сеть
│           └── path/       # Утилиты путей
│
└── shared/                 # Общие компоненты
    ├── api/                # OpenAPI спецификации
    ├── proto/              # Protobuf определения
    └── pkg/                # Сгенерированный код
```

### Запуск проекта

Для того чтобы вызывать команды из Taskfile, необходимо установить Taskfile CLI:

```bash
brew install go-task
```

### Переменные окружения

Генерация `.env` файлов из шаблонов:

```bash
# Генерирует .env файлы для всех сервисов
task env:generate
```

При первом запуске создаётся файл `deploy/env/.env` из шаблона `.env.template`. 
Отредактируйте его при необходимости и запустите команду снова.

### Docker Compose

Запуск инфраструктуры для локальной разработки:

```bash
# Сначала сгенерируйте .env файлы
task env:generate

# Поднять общую сеть
task up-core

# Поднять PostgreSQL для OrderService
task up-order

# Поднять MongoDB для InventoryService
task up-inventory

# Поднять всё вместе
task up-all

# Остановить всё
task down-all
```

#### Конфигурация баз данных

**PostgreSQL (OrderService):**
- Host: `localhost:5433`
- Database: `order-service`
- User: `order-service-user`
- Password: `order-service-password`

**MongoDB (InventoryService):**
- URI: `mongodb://inventory-service-user:inventory-service-password@localhost:27017`
- Database: `inventory-service`

### Миграции

Миграции для OrderService выполняются автоматически при запуске сервиса с использованием [goose](https://github.com/pressly/goose).

Миграции расположены в `order/migrations/`:
- `20250404191615_create_uuid_ossp_extension.sql` - расширение для UUID
- `20250404191624_create_orders_table.sql` - таблица заказов

### Доступные команды

```bash
# Docker Compose
task up-core           # Поднять общую сеть
task up-order          # Поднять PostgreSQL
task up-inventory      # Поднять MongoDB
task up-all            # Поднять всё
task down-all          # Остановить всё

# Переменные окружения
task env:generate      # Генерация .env файлов из шаблонов

# Тесты
task test              # Запуск юнит-тестов
task test-coverage     # Тесты с покрытием
task coverage:html     # HTML-отчет о покрытии
task test-integration  # Запуск интеграционных тестов
task test-api          # Запуск API тестов

# Разработка
task lint              # Линтинг
task format            # Форматирование кода
task gen               # Генерация кода (proto + OpenAPI)
task mockery:gen       # Генерация моков
task deps:update       # Обновление зависимостей
```

### Покрытие тестами

| Сервис    | Покрытие |
|-----------|----------|
| inventory | 97.6%    |
| order     | 97.8%    |
| payment   | 100%     |

Для получения актуального покрытия запустите:

```bash
task test-coverage
```

### CI/CD

Проект использует GitHub Actions для непрерывной интеграции и доставки. Основные workflow:

- **CI** (`.github/workflows/ci.yml`) - проверяет код при каждом push и pull request
    - Извлечение версий из Taskfile
    - Линтинг кода
    - Запуск юнит-тестов
    - Запуск интеграционных тестов
    - Проверка покрытия тестами

### Platform библиотека

Модуль `platform` содержит общие утилиты:

#### Closer (Graceful Shutdown)

```go
import "github.com/bogdanovds/rocket_factory/platform/pkg/closer"

// Добавление функции для закрытия
closer.AddNamed("database", func(ctx context.Context) error {
    return db.Close()
})

// Настройка обработки сигналов
closer.Configure(syscall.SIGINT, syscall.SIGTERM)

// Закрытие всех ресурсов
closer.CloseAll(ctx)
```

#### Logger (Zap)

```go
import "github.com/bogdanovds/rocket_factory/platform/pkg/logger"

// Инициализация
logger.Init("debug", false)

// Использование
logger.Info(ctx, "message", zap.String("key", "value"))
logger.Error(ctx, "error", zap.Error(err))
```

#### Testcontainers

```go
import (
    "github.com/bogdanovds/rocket_factory/platform/pkg/testcontainers/mongo"
    "github.com/bogdanovds/rocket_factory/platform/pkg/testcontainers/network"
)

// Создание сети
net, _ := network.NewNetwork(ctx, "test")

// Создание MongoDB контейнера
container, _ := mongo.NewContainer(ctx,
    mongo.WithNetworkName(net.Name()),
    mongo.WithDatabase("test"),
)
defer container.Terminate(ctx)

// Получение клиента
client := container.Client()
```

### Интеграционные тесты

Интеграционные тесты используют тег `integration`:

```go
//go:build integration

package mypackage_test

func TestIntegration(t *testing.T) {
    // Тест с реальной базой данных
}
```

Запуск интеграционных тестов:

```bash
task test-integration
```

### Переменные окружения

**OrderService:**
- `POSTGRES_HOST` - хост PostgreSQL (default: `localhost`)
- `POSTGRES_PORT` - порт PostgreSQL (default: `5432`)
- `EXTERNAL_POSTGRES_PORT` - внешний порт (default: `5433`)
- `POSTGRES_USER` - пользователь (default: `order-service-user`)
- `POSTGRES_PASSWORD` - пароль (default: `order-service-password`)
- `POSTGRES_DB` - база данных (default: `order-service`)

**InventoryService:**
- `MONGO_HOST` - хост MongoDB (default: `mongo-inventory`)
- `MONGO_PORT` - порт MongoDB (default: `27017`)
- `EXTERNAL_MONGO_PORT` - внешний порт (default: `27017`)
- `MONGO_DATABASE` - база данных (default: `inventory-service`)
- `MONGO_INITDB_ROOT_USERNAME` - пользователь (default: `inventory-service-user`)
- `MONGO_INITDB_ROOT_PASSWORD` - пароль (default: `inventory-service-password`)

### Разработка

#### Добавление новых тестов

Тесты располагаются рядом с тестируемым кодом в файлах `*_test.go`. Для организации тестов используется `testify/suite`:

```go
type MyServiceTestSuite struct {
    suite.Suite
    mockRepo *mocks.MockRepository
    service  *Service
}

func (s *MyServiceTestSuite) SetupTest() {
    s.mockRepo = mocks.NewMockRepository()
    s.service = NewService(s.mockRepo)
}

func TestMyServiceTestSuite(t *testing.T) {
    suite.Run(t, new(MyServiceTestSuite))
}
```

#### Генерация моков

Моки хранятся в директориях `mocks/` рядом с интерфейсами:
- `internal/repository/mocks/` - моки репозиториев
- `internal/service/mocks/` - моки сервисов
- `internal/client/grpc/mocks/` - моки gRPC клиентов
