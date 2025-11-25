## Rocket Factory

![Coverage](https://img.shields.io/badge/coverage-98%25-brightgreen)

### Описание проекта
Имитируем сборку ракет

### Структура проекта

Проект использует чистую архитектуру с разделением на слои:

```
.
├── deploy/
│   └── compose/
│       ├── core/           # Общая сеть микросервисов
│       ├── inventory/      # MongoDB для InventoryService
│       └── order/          # PostgreSQL для OrderService
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

### Docker Compose

Запуск инфраструктуры для локальной разработки:

```bash
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
- Host: `localhost:5432`
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

# Тесты
task test              # Запуск юнит-тестов
task test-coverage     # Тесты с покрытием
task test-coverage-report  # HTML-отчет о покрытии
task test-api          # Запуск API тестов

# Разработка
task lint              # Линтинг
task format            # Форматирование кода
task gen               # Генерация кода (proto + OpenAPI)
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
    - Линтинг кода
    - Запуск юнит-тестов
    - Проверка покрытия тестами
    - Обновление бейджа покрытия

### Переменные окружения

**OrderService:**
- `POSTGRES_HOST` - хост PostgreSQL (default: `localhost`)
- `POSTGRES_PORT` - порт PostgreSQL (default: `5432`)
- `POSTGRES_USER` - пользователь (default: `order-service-user`)
- `POSTGRES_PASSWORD` - пароль (default: `order-service-password`)
- `POSTGRES_DB` - база данных (default: `order-service`)

**InventoryService:**
- `MONGO_URI` - URI для подключения к MongoDB (default: `mongodb://inventory-service-user:inventory-service-password@localhost:27017`)
- `MONGO_DB` - база данных (default: `inventory-service`)

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
