# SWERPCTASK

Приложение JSON RPC 2.0 API сервис разработан в рамках тестового задания.

## API функции

Вызовы API функций приложения происходят с использованием протокола HTTP
и JSON RPC, пример:

```json
--> {"jsonrpc": "2.0", "method": "Api.GetNumber", "params": [], "id": 1}
<-- {"jsonrpc": "2.0", "result": 0, "id": 1}
```
### Список API функций

- `Api.GetNumber` - Текущее значение счетчика
- `Api.IncrementNumber` - Увеличение счетчика
- `Api.SetSettings` - Установка максимального значения счетчика и значения приращения

## База данных

В качестве БД рекомендуется к использованию современная СУБД MySQL - MySQL/MariaDB

## Директории

Краткое описание директорий

- `docker` - директория для Docker файлов используемых при развертке приложения
- `migrations` - SQL (MySQL) файлы для исполнения миграций
- `src` - исходный код приложения

## Сборка

Для сборки необходимо клонировать данный репозиторий на ПК или сервер:

```git clone https://github.com/stormdynamics/swerpctask.git```

### Сборка из директории на локальном ПК

- Установить [GoLang](https://golang.org/dl/) не ниже версии 1.14
- Выполнить команду - `go mod download`
- Выполнить команду - `go build`

### Сбока Docker-Compose

Выполнить команду:

```docker-compose -f docker-compose.yml up -d --build```

## Зависимости

Для основного приложения используются следующие зависимости:

- драйвер MySQL - `github.com/go-sql-driver/mysql`
- менеджер миграций - `github.com/pressly/goose`

Для тестов используются зависимости:

- макеты(mockup) БД - `github.com/DATA-DOG/go-sqlmock`