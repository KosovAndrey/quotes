# Цитатник (Quotes Service)

Простой REST API сервис для управления цитатами, написанный на Go.

Сервис использует хранение данных в памяти

## Описание

Цитатник - это REST API сервис, который позволяет:
- Добавлять новые цитаты
- Получать список всех цитат
- Получать случайную цитату
- Фильтровать цитаты по автору
- Удалять цитаты по ID

## Требования

- Go 1.24.3 или выше
- gorilla/mux (для маршрутизации)

## Установка

1. Клонируйте репозиторий:
```bash
git clone https://github.com/KosovAndrey/quotes.git
cd quotes
```

2. Установите зависимости:
```bash
go mod download
```

3. Запустите сервер:
```bash
go run main.go
```

Сервер будет запущен на `http://localhost:8080`

## API Endpoints

### Добавление новой цитаты
```bash
curl -X POST http://localhost:8080/quotes \
-H "Content-Type: application/json" \
-d "{\"author\":\"Confucius\", \"quote\":\"Life is simple, but we insist on making it complicated.\"}"
```

### Получение всех цитат
```bash
curl http://localhost:8080/quotes
```

### Получение случайной цитаты
```bash
curl http://localhost:8080/quotes/random
```

### Фильтрация по автору
```bash
curl http://localhost:8080/quotes?author=Confucius
```

### Удаление цитаты по ID
```bash
curl -X DELETE http://localhost:8080/quotes/1
```

## Структура проекта

```
quotes/
├── cmd/            # Точка входа приложения
├── internal/       # Внутренние пакеты приложения
│   ├── domain/     # Доменная логика и модели
│   ├── handlers/   # HTTP обработчики
│   ├── services/   # Бизнес-логика
│   └── storage/    # Хранилище данных
├── tests/          # Тесты
├── go.mod          # Файл зависимостей
├── go.sum          # Файл с хешами зависимостей
└── README.md       # Документация
```

## Тестирование

Для запуска тестов выполните:
```bash
go test ./tests/... 
```
