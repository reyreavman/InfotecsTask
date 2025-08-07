# Infotecs TechTask
#### Система обработки транзакций платежной системы

Дает возможность управлять информацией о кошельках и транзакциях. 
Включает в себя возможности чтения данных о кошелькак и создания транзакций. 

## Содержание
- [Разработка](#разработка)
- [Технологии](#технологии)
- [Тестирование](#тестирование)
- [Требования к запуску](#требования-к-запуску)
- [Запуск](#запуск)
- [Документация к API](#документация-к-api-и-примеры-запросов)
- [Функциональность](#функциональность)

### Разработка
---
- Как основная БД была использована PostgreSQL.
- БД, как и основное приложение, крутятся в Docker контейнерах.
- Также реализован отдельный контейнер, который накатывает миграции на БД.
- Конфигурация приведена в docker-compose.yaml 
- Реализован CI/CD пайплайн.
- Написаны тесты (в том числе конкурентные) для репозиториев с использованием TestContainers
- Написаны тесты для handler с использованием testify
- Реализован Rate limiter 

### Технологии и пакеты
---
- [github.com/gin-gonic/gin](https://pkg.go.dev/github.com/gin-gonic/gin@v1.10.1) - как основной веб-фреймворк
- [github.com/jackc/pgx](https://pkg.go.dev/github.com/jackc/pgx@v3.6.2+incompatible) - обертка над нативной реализацией пакета sql
- [github.com/lib/pq](https://pkg.go.dev/github.com/lib/pq@v1.10.9) - выбран в качестве Postgres-драйвера
- [github.com/google/uuid](https://pkg.go.dev/github.com/google/uuid@v1.6.0) - для генерации UUID
- [Docker](https://www.docker.com/) - платформа контейнеризации
- [Github Actions] - для прогона тестов

### Тестирование
---
Для тестирования в проекте были использованы:
- [github.com/testcontainers/testcontainers-go/modules/postgres](https://www.red-gate.com/products/flyway/community/) - Postgres TestContainer
- [github.com/stretchr/testify](https://assertj.github.io/doc/) - пакет для тестирования

Проект покрыт юнит и интеграционными тестами. Написаны юнит тесты для всех основных сценариев использования handlerа
Также написаны интеграционные тесты с помощью TestContainers для тестирования слоя репозиториев.

### Требования к запуску
---
Для запуска проекта необходим Docker.

### Запуск
---
Для запуска выполните команду:
```sh
$ docker compose up
```

По умолчанию проект запускается на порту 8080.

### Документация к API и примеры запросов
---

#### 1. **Отправка средств**  
**`POST /api/send`**  
Отправляет средства между кошельками.  

**Тело запроса (JSON)**:
```json
{
  "from": "e240d825d255af751f5f55af8d9671beabdf2236c0a3b4e2639b3e182d994c88e",
  "to": "9b3e182d994c88ee240d825d255af751f5f55af8d9671beabdf2236c0a3b4e26",
  "amount": 3.50
}
```

**Параметры**:
| Поле    | Тип     | Обязательно | Описание                     |
|---------|---------|-------------|------------------------------|
| from    | string  | Да          | Адрес кошелька-отправителя   |
| to      | string  | Да          | Адрес кошелька-получателя    |
| amount  | float   | Да          | Сумма перевода (>0)          |

**Пример запроса**:
```bash
curl -X POST http://localhost:8080/api/send \
  -H "Content-Type: application/json" \
  -d    '{
            "from": "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a10",
            "to": "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14",
            "amount": 3.50
        }'
```

**Успешный ответ** (`200 OK`):
```json
{
  "id": "h7hgyh54-3m7q-9h8t-xp6t-1cp8qe487b10",
  "from": "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a10",
  "to": "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14",
  "amount": 3.50,
  "timestamp": "2025-08-08T14:30:00Z"
}
```

**Ошибки**:
- `400 Bad Request` - Невалидные данные:
```json
{
  "error": "Validation failed",
  "details": [
      {
          "field": "FromAddress",
          "message": "Field is required"
      }
  ]
}
```

- `400 Bad Request` - Кошелек отправителя не найден:
```json
{
  "error": "Sender wallet not found"
}
```

- `400 Bad Request` - Кошелек получателя не найден:
```json
{
  "error": "Recipient wallet not found"
}
```

---

#### 2. **Получение последних транзакций**  
**`GET /api/transactions`**  
Возвращает N последних транзакций.  

**Query-параметры**:
| Параметр | Тип  | Обязательно | Описание                     |
|----------|------|-------------|------------------------------|
| count    | int  | Нет         | Количество транзакций        |

Если количество транзакций не указано, то возвращается список всех транзакций

**Пример запроса**:
```bash
curl "http://localhost:8080/api/transactions?count=2"
```

**Успешный ответ** (`200 OK`):
```json
[
    {
        "id": "h7hgyh54-3m7q-9h8t-xp6t-1cp8qe487b10",
        "from": "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a10",
        "to": "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14",
        "amount": 3.50,
        "status": "completed",
        "message": "Transaction completed",
        "timestamp": "2025-08-08T14:30:00Z"
    },
    {
        "id": "h7hgyh54-3m7q-9h8t-xp6t-1cp8qe487b11",
        "from": "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a10",
        "to": "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14",
        "amount": 3.50,
        "status": "failed",
        "message": "Sender does not have enough balance",
        "timestamp": "2025-08-08T14:30:00Z"
    }
]
```

**Ошибки**:
- `400 Bad Request` - Некорректный параметр count:
```json
{
  "error": "Invalid query parameters"
}
```

---

#### 3. **Получение баланса кошелька**  
**`GET /api/wallet/{address}/balance`**  
Возвращает баланс указанного кошелька.  

**Path-параметры**:
| Параметр | Тип    | Обязательно | Описание          |
|----------|---------|-------------|-------------------|
| address  | uuid  | Да          | Адрес кошелька    |

**Пример запроса**:
```bash
curl http://localhost:8080/api/wallet/b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a10/balance
```

**Успешный ответ** (`200 OK`):
```json
{
    "id": "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a10",
    "balance": 100      
}
```

**Ошибки**:
- `404 Not Found` - Кошелек не существует:
```json
{
  "error": "Wallet not found"
}
```

---

### Примеры сценариев

#### 📤 Успешный перевод средств
**Запрос**:
```json
POST /api/send
{
  "from": "wallet_A",
  "to": "wallet_B",
  "amount": 10.0
}
```
**Ответ**:
```json
200 OK
{
  "id": "txn_12345",
  "from": "wallet_A",
  "to": "wallet_B",
  "amount": 10.0,
  "status": "completed",
  "message": "Transaction completed",
  "timestamp": "2025-08-08T15:00:00Z"
}
```

#### ❌ Недостаточно средств
**Запрос**:
```json
POST /api/send
{
  "from": "wallet_A",
  "to": "wallet_B",
  "amount": 1000.0
}
```
**Ответ**:
```json
200 OK
{
  "id": "txn_12345",
  "from": "wallet_A",
  "to": "wallet_B",
  "amount": 1000.0,
  "status": "failed",
  "message": "Sender does not have enough balance",
  "timestamp": "2025-08-08T15:00:00Z"
}
```

#### 📜 Получение последних транзакций
**Запрос**:
```
GET /api/transactions?count=3
```
**Ответ**:
```json
200 OK
[
  {"id": "txn3", "amount": 5.0, ...},
  {"id": "txn2", "amount": 2.5, ...},
  {"id": "txn1", "amount": 10.0, ...}
]
```

#### 💰 Запрос баланса
**Запрос**:
```
GET /api/wallet/wallet_A/balance
```
**Ответ**:
```json
200 OK
{
  "id": "wallet_A",
  "balance": 85.30
}
```

---

### Примечания
1. **Формат адресов**:  
   Строка из 64 hex-символов (`[a-f0-9]{64}`)
   
2. **Формат сумм**:  
   - Все суммы с двумя знаками после точки
   - Минимальная сумма: 0.01

3. **Формат времени**:  
   ISO 8601 в UTC: `YYYY-MM-DDThh:mm:ssZ`

4. **Лимиты**:  
   - Максимальное количество запросов в течение 1 минуты с одного IP: 100

5. **Заголовки**:  
   Все запросы требуют:  
   ```http
   Content-Type: application/json
   ```


### Функциональность
Реализованный API имеет следующие методы:
- Получение баланса кошелька
- Получение ограниченного количества транзакций
- Получение всех транзакций
- Создание транзакциий

