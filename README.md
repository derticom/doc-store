# doc-store

**doc-store** — это веб-сервис на для хранения, управления и раздачи электронных документов с авторизацией, кешированием и разграничением доступа.

Проект реализует REST API и включает:

- регистрацию и аутентификацию пользователей;
- загрузку, получение и удаление документов;
- разграничение доступа к документам;
- кэширование и хранение данных с использованием сторонних сервисов.

---

## Технологии

- **PostgreSQL** — хранение пользователей и метаинформации о документах;
- **Redis** — кеширование запросов;
- **MinIO** — файловое хранилище для документов.

---
## API

### 1. Регистрация пользователя

**POST** `/api/register`

**Тело запроса:**
```json
{
  "token": "admin_secret_token",
  "login": "user1234",
  "pswd": "P@ssw0rd1"
}
```

**Ответ:**
```json
{
  "response": {
    "login": "newuser"
  }
}
```

---

### 2. Аутентификация

**POST** `/api/auth`

**Тело запроса:**
```json
{
  "login": "user1234",
  "pswd": "P@ssw0rd1"
}
```

**Ответ:**
```json
{
  "response": {
    "token": "user_auth_token"
  }
}
```

---

### 3. Загрузка документа

**POST** `/api/docs`

**Форма запроса (multipart/form-data):**
- `meta`: JSON-объект, содержащий:
```json
{
  "name": "photo.jpg",
  "file": true,
  "public": false,
  "token": "user_auth_token",
  "mime": "image/jpg",
  "grant": ["login1", "login2"]
}
```
- `json`: JSON-данные (опционально)
- `file`: файл документа

**Ответ:**
```json
{
  "data": {
    "json": { ... },
    "file": "photo.jpg"
  }
}
```

---

### 4. Получение списка документов

**GET / HEAD** `/api/docs/?token=...`

**Параметры запроса (query):**
- `token` (обязательно)
- `login` (опционально) — получить документы другого пользователя
- `key` и `value` — фильтрация по полям
- `limit` — ограничение по количеству

**Ответ:**
```json
{
  "data": {
    "docs": [
      {
        "id": "uuid",
        "name": "photo.jpg",
        "mime": "image/jpg",
        "file": true,
        "public": false,
        "created": "2024-01-01 12:00:00",
        "grant": ["login1", "login2"]
      }
    ]
  }
}
```

---

### 5. Получение одного документа

**GET / HEAD** `/api/docs/<id>?token=...`

- Если файл — отдаётся файл
- Если JSON — возвращается содержимое

**Пример ответа:**
```json
{
  "data": {
    "some": "metadata"
  }
}
```

---

### 6. Удаление документа

**DELETE** `/api/docs/<id>?token=...`

**Ответ:**
```json
{
  "response": {
    "<id>": true
  }
}
```

---

### 7. Завершение сессии

**DELETE** `/api/auth/<token>`

**Ответ:**
```json
{
  "response": {
    "<token>": true
  }
}
```

---

## Запуск сервиса

Перед запуском необходимо:

- создать конфигурационный файл (или использовать демо _config/config.yml_) с параметрами подключения к:
   - PostgreSQL
   - Redis
   - MinIO
- задать административный токен для регистрации новых пользователей

--------

**Пример `config.yaml`:**
```yaml
log_level: "debug"
admin_token: "admin_secret_token"
postgres_url:  "postgres://docstore:docstore@localhost:5432/docstore?sslmode=disable"
redis_url: "redis://@localhost:6379"

minio_server:
  address: "localhost:9000"
  access_key: "minioadmin"
  secret_key: "minioadmin"
  bucket: "docstore"
  use_ssl: false

http_server:
  address: "localhost:8085"
  timeout: 5s
```
--------
Для локального запуска в корне проекта присутствует файл _docker-compose.yml_, который разворачивает следующие сервисы:

PostgreSQL — база данных docstore с пользователем docstore;
Redis — кеш;
MinIO — файловое хранилище с web-интерфейсом на порту 9001.

