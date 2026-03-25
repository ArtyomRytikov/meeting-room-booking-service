[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-22041afd0340ce965d47ae6ef1cefeee28c7c493a6346c4f15d667ab976d596c.svg)](https://classroom.github.com/a/xR-tWBKa)
# Room Booking Backend

Backend-сервис для бронирования переговорных комнат.

## Возможности

* Создание комнат (admin)
* Создание расписания для комнат (admin)
* Генерация слотов на основе расписания
* Бронирование слотов (user)
* Отмена бронирования
* Просмотр своих бронирований
* Просмотр всех бронирований с пагинацией (admin)
* JWT-аутентификация (dummy login)

---

## Технологии

* Go
* PostgreSQL
* Docker / Docker Compose
* JWT (аутентификация)

---

## Запуск проекта

### 1. Клонировать репозиторий

```bash
git clone <your-repo-url>
cd test-backend-1-ArtyomRytikov
```

### 2. Запустить через Docker

```bash
docker compose up --build
```

Сервис будет доступен:

```
http://localhost:8080
```

---

## Аутентификация

Используется dummy login.

### Получить токен admin

```bash
curl -X POST http://localhost:8080/dummyLogin -H "Content-Type: application/json" -d "{\"role\":\"admin\"}"
```

### Получить токен user

```bash
curl -X POST http://localhost:8080/dummyLogin -H "Content-Type: application/json" -d "{\"role\":\"user\"}"
```

---

## Основные эндпоинты

### Комнаты

* `POST /rooms/create` — создать комнату (admin)
* `GET /rooms/list` — список комнат

---

### Расписание

* `POST /rooms/{roomId}/schedule/create` — создать расписание (admin)
* `GET /rooms/{roomId}/slots/list?date=YYYY-MM-DD` — получить слоты

---

### Бронирования

* `POST /bookings/create` — создать бронирование (user)
* `POST /bookings/{id}/cancel` — отменить бронирование
* `GET /bookings/my` — мои бронирования
* `GET /bookings/list?page=1&pageSize=20` — все бронирования (admin)

---

## Тестирование

Запуск тестов:

```bash
go test ./...
```

Покрытие:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

### Покрытие

* Общее покрытие: **>40%**
* Unit тесты:

  * service слой
* Integration-like тесты:

  * handler слой (через httptest)
* Дополнительно:

  * auth
  * middleware
  * config

---

## Проверка сценариев (E2E)

Основные сценарии протестированы вручную через curl:

1. Создание комнаты
2. Создание расписания
3. Получение слотов
4. Бронирование
5. Повторное бронирование (ошибка)
6. Отмена бронирования
7. Повторная отмена
8. Проверка прав (admin/user)

---

##  Обработка ошибок

API возвращает структурированные ошибки:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "description"
  }
}
```

Примеры:

* `INVALID_REQUEST`
* `FORBIDDEN`
* `ROOM_NOT_FOUND`
* `SCHEDULE_EXISTS`
* `SLOT_ALREADY_BOOKED`

---

## Архитектура

Проект разделён на слои:

```
handler → service → repository → database
```

* **handler** — HTTP слой
* **service** — бизнес-логика
* **repository** — работа с БД

---

##  Особенности

* Проверка ролей (admin / user)
* Защита от двойного бронирования
* Генерация слотов по расписанию
* Пагинация для админских запросов
* Валидация входных данных

---
