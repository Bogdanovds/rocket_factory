## Rocket Factory

![Coverage](https://img.shields.io/badge/coverage-98%25-brightgreen)

### Описание проекта
Имитируем сборку ракет

### Структура проекта

Проект использует чистую архитектуру с разделением на слои:

```
.
├── inventory/          # Сервис управления деталями
│   └── internal/
│       ├── api/        # gRPC хендлеры
│       ├── service/    # Бизнес-логика
│       ├── repository/ # Хранилище данных
│       ├── model/      # Модели сервисного слоя
│       └── converter/  # Конвертеры между слоями
│
├── order/              # Сервис заказов
│   └── internal/
│       ├── api/        # REST/HTTP хендлеры
│       ├── service/    # Бизнес-логика
│       ├── repository/ # Хранилище данных
│       ├── client/     # Клиенты внешних сервисов
│       ├── model/      # Модели сервисного слоя
│       └── converter/  # Конвертеры между слоями
│
├── payment/            # Сервис оплаты
│   └── internal/
│       ├── api/        # gRPC хендлеры
│       ├── service/    # Бизнес-логика
│       └── model/      # Модели
│
└── shared/             # Общие компоненты
    ├── api/            # OpenAPI спецификации
    ├── proto/          # Protobuf определения
    └── pkg/            # Сгенерированный код
```

### Запуск проекта

Для того чтобы вызывать команды из Taskfile, необходимо установить Taskfile CLI:

```bash
brew install go-task
```

### Доступные команды

```bash
# Запуск юнит-тестов
task test

# Запуск тестов с покрытием
task test-coverage

# Генерация HTML-отчета о покрытии
task test-coverage-report

# Запуск API тестов
task test-api

# Линтинг
task lint

# Форматирование кода
task format

# Генерация кода (proto + OpenAPI)
task gen
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
