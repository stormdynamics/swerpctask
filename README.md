# SWERPCTASK

Приложение сервис SWERPCTASK разработано в рамках тестового задания.

## Директории

Краткое описание директорий

- `docker` - директория для Docker файлов используемых при развертке приложения
- `migrations` - SQL (MySQL) файлы для исполнения миграций
- `src` - исходный код приложения

## Сборка

Для сборки необходимо клонировать приложение на ПК или сервер:

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