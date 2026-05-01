## Microservices

Проект комбинирует синхронные вызовы (HTTP/gRPC) и асинхронные события (Kafka), использует outbox-паттерн, централизованный сбор логов и OTEL-метрики.

---

## Что внутри

Система состоит из 6 доменных сервисов и 2 общих модулей:

- `order` — HTTP API заказов, orchestration бизнес-процесса.
- `inventory` — gRPC API деталей/комплектующих.
- `payment` — gRPC API оплаты (симуляция платежного шлюза).
- `assembly` — Kafka consumer/producer, «сборка» после оплаты.
- `iam` — gRPC сервис аутентификации и сессий.
- `notification` — Kafka consumers + Telegram нотификации.
- `platform` — общая библиотека (middleware, kafka, logger, metrics, tx manager, redis client и т.д.).
- `shared` — контракты (`proto`, OpenAPI), генерируемые типы.

---

## Карта сервисов

| Сервис | Протокол (вход) | Порт по умолчанию | Основные зависимости | Хранилище |
|---|---|---|---|---|
| `order` | HTTP (OpenAPI/ogen) | `8080` | IAM gRPC, Inventory gRPC, Payment gRPC, Kafka | PostgreSQL |
| `inventory` | gRPC | `50051` | IAM gRPC (auth interceptor) | MongoDB |
| `payment` | gRPC | `50052` | — | нет (stateless) |
| `assembly` | Kafka consumer/producer | — | Kafka | нет (stateless) |
| `iam` | gRPC | `50053` | PostgreSQL, Redis | PostgreSQL + Redis |
| `notification` | Kafka consumer | — | Kafka, Telegram API | нет (stateless) |

---

## Аутентификация и передача сессии

### Вход в систему

1. Клиент вызывает IAM:
   - `Register`
   - `Login` → получает `session_uuid`
2. Для HTTP вызовов в `order` клиент передает:
   - `X-Session-Uuid: <uuid>`

### Как проверяется доступ

- `order` использует HTTP middleware `platform/pkg/middleware/http`:
  - берет `X-Session-Uuid`,
  - вызывает `AuthService.Whoami`,
  - кладет `user` и `session_uuid` в `context`.
- `order` при вызове `inventory` прокидывает `session_uuid` в gRPC metadata через `ForwardSessionUUIDToGRPC`.
- `inventory` на входе защищен gRPC interceptor `platform/pkg/middleware/grpc`:
  - читает metadata `session-uuid`,
  - валидирует через `IAM.Whoami`.

---

## Полный жизненный цикл запроса

Ниже описан основной пользовательский путь и то, как запрос проходит через систему.

### 1. Регистрация и логин

1. Клиент вызывает `IAM.Register` и создает пользователя.
2. Клиент вызывает `IAM.Login` и получает `session_uuid`.
3. Клиент сохраняет `session_uuid` и передает его в `Order` через заголовок `X-Session-Uuid`.

Результат: у клиента есть действующая сессия, которую `order` и `inventory` могут валидировать через `IAM.Whoami`.

### 2. Создание заказа (`POST /api/v1/orders`)

1. HTTP-запрос приходит в `order`.
2. `order` middleware читает `X-Session-Uuid`, вызывает `IAM.Whoami` и кладет пользователя в `context`.
3. `order` вызывает `Inventory.ListParts` по gRPC, чтобы получить детали и цены.
4. Перед gRPC вызовом `order` прокидывает `session_uuid` в metadata (`session-uuid`).
5. `inventory` interceptor проверяет `session-uuid` через `IAM.Whoami`.
6. `order` считает итоговую цену и сохраняет заказ в PostgreSQL со статусом `PENDING_PAYMENT`.
7. `order` увеличивает метрику `orders_total`.
8. Клиент получает `order_uuid` и `total_price`.

Результат: заказ создан, но еще не оплачен.

### 3. Оплата заказа (`POST /api/v1/orders/{order_uuid}/pay`)

1. HTTP-запрос приходит в `order` и проходит ту же auth-проверку middleware.
2. Service слой `order` читает заказ из PostgreSQL и проверяет, что статус допускает оплату.
3. `order` вызывает `Payment.Pay` по gRPC.
4. `payment` возвращает `transaction_uuid`.
5. `order` открывает транзакцию в PostgreSQL.
6. В этой транзакции `order` обновляет заказ: статус `PAID`, `transaction_uuid`, `payment_method`.
7. В этой же транзакции `order` пишет событие `OrderPaid` в outbox-таблицу `events` со статусом `PENDING`.
8. Транзакция коммитится.
9. `order` увеличивает метрику `orders_revenue_total` на сумму заказа.
10. Клиент получает `transaction_uuid`.

Результат: деньги подтверждены, событие о платеже надежно сохранено в outbox.

### 4. Публикация outbox и сборка корабля

1. Фоновый воркер `order` (`outbox sender`) батчем читает `PENDING` события из `events`.
2. Для каждого события воркер пытается отправить сообщение в Kafka.
3. Если отправка успешна, событие помечается как `PUBLISHED`.
4. Если отправка неуспешна, событие остается `PENDING` и будет ретраиться.
5. `assembly` consumer читает `OrderPaid` из Kafka.
6. `assembly` валидирует payload события.
7. `assembly` имитирует сборку (задержка), формирует событие `ShipAssembled`.
8. `assembly` публикует `ShipAssembled` в Kafka.
9. `assembly` пишет метрику `assembly_duration_seconds`.

Результат: после оплаты заказ проходит этап сборки и генерирует событие завершения сборки.

### 5. Завершение заказа и нотификации

1. `order` consumer читает `ShipAssembled` из Kafka.
2. `order` переводит заказ в статус `COMPLETED`.
3. `notification` consumer читает `OrderPaid` и отправляет уведомление об оплате в Telegram.
4. `notification` consumer читает `ShipAssembled` и отправляет уведомление о сборке в Telegram.

Результат: жизненный цикл заказа завершается статусом `COMPLETED`, пользователь получает уведомления.

### 6. Чтение и отмена заказа

1. `GET /api/v1/orders/{order_uuid}` возвращает заказ по UUID после auth-проверки.
2. `POST /api/v1/orders/{order_uuid}/cancel` отменяет заказ только из допустимых статусов.
3. Если статус уже `PAID`, `CANCELED` или `COMPLETED`, сервис возвращает `conflict`.

Результат: клиент может читать и отменять только корректные по состоянию заказы.

---

## Хранилища и данные

### Order (PostgreSQL)

- Таблица заказов (`orders`)
- Таблица outbox-событий (`events`)
- Миграции: `order/migrations`

### IAM

- PostgreSQL: пользователи и профильные данные (`iam/migrations`)
- Redis: сессии (`session:<session_uuid> -> user_uuid`, TTL)

### Inventory (MongoDB)

- Коллекция деталей
- Индексы и seed при инициализации

---

## Контракты и кодогенерация

### gRPC (`shared/proto`)

- `auth/v1/auth.proto`
- `inventory/v1/inventory.proto`
- `payment/v1/payment.proto`
- `events/v1/events.proto` (Kafka payload модели)


### HTTP OpenAPI (`shared/api/order/v1`)

- `order.openapi.yaml`
- пути:
  - `POST /api/v1/orders`
  - `GET /api/v1/orders/{order_uuid}`
  - `POST /api/v1/orders/{order_uuid}/pay`
  - `POST /api/v1/orders/{order_uuid}/cancel`


---

## Observability

### Метрики (OTEL + Prometheus + Grafana)

Pipeline:

`Service OTLP export -> OTEL Collector -> Prometheus scrape -> Grafana`

Конфиги:
- OTEL Collector: `deploy/compose/core/otel-collector/config.yaml`
- Prometheus: `deploy/compose/core/prometheus/prometheus.yml`

Кастомные метрики:
- `order`
  - `orders_total`
  - `orders_revenue_total`
- `assembly`
  - `assembly_duration_seconds` (histogram)

### Логи (Filebeat + Logstash + Elasticsearch + Kibana)

Pipeline:

`local *.log files -> Filebeat -> Logstash -> Elasticsearch -> Kibana`

---

## Локальный запуск

### 1. Поднять инфраструктуру

```bash
task up-all
```

### 2. Запустить сервисы из корня репозитория

```bash
go run ./iam/cmd/main.go
go run ./inventory/cmd/main.go
go run ./payment/cmd/main.go
go run ./assembly/cmd/main.go
go run ./order/cmd/main.go
# опционально:
go run ./notification/cmd/main.go
```

### 3. Полезные команды

```bash
task lint
task test
task test-iam
task test-api
```

---

## Порты по умолчанию

- `order` HTTP: `localhost:8080`
- `inventory` gRPC: `localhost:50051`
- `payment` gRPC: `localhost:50052`
- `iam` gRPC: `localhost:50053`
- `postgres-order`: `localhost:5435`
- `postgres-iam`: `localhost:5444`
- `redis-iam`: `localhost:6333`
- `mongo-inventory`: `localhost:27018`
- `kafka`: `localhost:9092`
- `kafka-ui`: `localhost:8090`
- `prometheus`: `localhost:9090`
- `grafana`: `localhost:3000`
- `kibana`: `localhost:5601`
- `elasticsearch`: `localhost:9200`

---

## Структура репозитория

```text
.
├── order/         # HTTP API + orchestration + outbox + consumers
├── inventory/     # gRPC API деталей + Mongo
├── payment/       # gRPC API оплаты
├── assembly/      # Kafka consumer/producer сборки
├── iam/           # gRPC auth + users (Postgres + Redis sessions)
├── notification/  # Kafka consumers + Telegram notifications
├── platform/      # shared runtime libs (kafka/logger/middleware/metrics/tx/cache)
├── shared/        # OpenAPI + proto + generated contracts
└── deploy/        # compose/env/infra configs
```

---

## ограничения текущей реализации

- Нет единого API Gateway/BFF для внешнего клиента (entrypoint сейчас — `order` по HTTP и отдельный `iam` по gRPC).
- Не все сервисы упакованы как отдельные docker-compose приложения (часть запускается напрямую `go run`).

---
