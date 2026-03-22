# Telegram Session Manager

Сервис на Go для управления множественными соединениями с Telegram через библиотеку gotd/td с gRPC API.

## Возможности

- 🔄 Динамическое создание и удаление соединений с Telegram
- 📤 Отправка текстовых сообщений
- 📥 Получение входящих сообщений в реальном времени
- 🔌 Изолированные соединения (проблемы с одним соединением не влияют на другие)
- 🌐 gRPC API для взаимодействия
- 💾 Хранение состояния в памяти

## Требования

- Go 1.26+
- Telegram API ID и API Hash

## Установка и запуск

### 1. Клонирование и установка зависимостей

```bash
git clone <repository-url>
cd tg-session-manager
go mod download
```

### 2. Получение Telegram API credentials

1. Перейдите на https://my.telegram.org/apps
2. Войдите в свой аккаунт Telegram
3. Создайте новое приложение
4. Сохраните `api_id` и `api_hash`

### 3. Конфигурация

Скопируйте файл переменных окружения:

```bash
cp .env.example .env
```

Отредактируйте `.env` файл:

```env
# gRPC Server Configuration
TGSM_SERVER_HOST=localhost
TGSM_SERVER_GRPC_PORT=50051

# Telegram Configuration
TGSM_TELEGRAM_API_ID=ваш_api_id
TGSM_TELEGRAM_API_HASH=ваш_api_hash
TGSM_SESSION_DIR=./sessions

# Logging
TGSM_LOG_LEVEL=info
```

### 4. Запуск сервера

```bash
go run cmd/server/main.go
```

Или соберите и запустите:

```bash
go build -o bin/server cmd/server/main.go
./bin/server
```

## gRPC API

### Создание соединения

```protobuf
rpc CreateConnection(CreateConnectionRequest) returns (CreateConnectionResponse);
```

### Удаление соединения

```protobuf
rpc DeleteConnection(DeleteConnectionRequest) returns (DeleteConnectionResponse);
```

### Отправка сообщения

```protobuf
rpc SendMessage(SendMessageRequest) returns (SendMessageResponse);
```

### Подписка на сообщения

```protobuf
rpc ReceiveMessages(ReceiveMessagesRequest) returns (stream ReceiveMessagesResponse);
```

### Получение статуса

```protobuf
rpc GetConnectionStatus(GetConnectionStatusRequest) returns (GetConnectionStatusResponse);
```

### Список соединений

```protobuf
rpc ListConnections(ListConnectionsRequest) returns (ListConnectionsResponse);
```

## Пример использования с grpcurl

### Создание соединения

```bash
grpcurl -plaintext -d '{
  "session_id": "session1",
  "phone_number": "+1234567890",
  "api_id": "123456",
  "api_hash": "abcdef1234567890"
}' localhost:50051 telegram.TelegramService/CreateConnection
```

### Отправка сообщения

```bash
grpcurl -plaintext -d '{
  "session_id": "session1",
  "chat_id": "123456789",
  "text": "Hello, World!"
}' localhost:50051 telegram.TelegramService/SendMessage
```

### Получение входящих сообщений

```bash
grpcurl -plaintext -d '{
  "session_id": "session1"
}' localhost:50051 telegram.TelegramService/ReceiveMessages
```

## Структура проекта

```
tg-session-manager/
├── api/proto/              # Proto файлы и сгенерированный код
├── cmd/server/            # Основной файл сервера
├── internal/
│   ├── config/            # Конфигурация приложения
│   ├── grpc/              # gRPC сервер и обработчики
│   └── telegram/          # Telegram клиент и менеджер соединений
├── configs/               # Файлы конфигурации
└── .env.example          # Пример переменных окружения
```

## Архитектура

1. **Manager** - управляет множественными Telegram клиентами
2. **Client** - обертка над gotd/td клиентом
3. **gRPC Server** - предоставляет API для внешнего взаимодействия
4. **Config** - управление конфигурацией через файлы и переменные окружения

## Важные замечания

- Сессии сохраняются в директории `./sessions` (настраивается через `TGSM_SESSION_DIR`)
- Каждое соединение изолировано и работает в отдельной goroutine
- При ошибке соединения автоматического переподключения нет (можно добавить)
- Состояние хранится только в памяти (при перезапуске все соединения теряются)

## Лицензия

MIT