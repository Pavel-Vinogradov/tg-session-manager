# Telegram Session Manager

Сервис для управления несколькими Telegram сессиями через gRPC.

## Возможности

- Создание нескольких независимых Telegram сессий
- QR-код аутентификация
- Отправка сообщений
- Получение сообщений в реальном времени

## Установка и запуск

1. Клонировать репозиторий:
```bash
git clone git@github.com:Pavel-Vinogradov/tg-session-manager.git
cd tg-session-manager
```

2. Настроить переменные окружения:
```bash
cp .env.example .env
```

3. Отредактировать `.env` файл с вашими Telegram API данными

4. Запустить:
```bash
make run
```

## Использование API

### Создать сессию
```bash
grpcurl -plaintext -d '{}' localhost:50051 pact.telegram.TelegramService/CreateSession
```

### Отправить сообщение
```bash
grpcurl -plaintext -d '{"session_id":"id","peer":"@username","text":"Привет"}' localhost:50051 pact.telegram.TelegramService/SendMessage
```

### Подписаться на сообщения
```bash
grpcurl -plaintext -d '{"session_id":"id"}' localhost:50051 pact.telegram.TelegramService/SubscribeMessages
```
