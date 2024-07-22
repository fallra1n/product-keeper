## Запуск

Перед запуском нужно установить переменные окружения в файле .env
Нужно указать CONFIG_PATH и JWT_SECRET

Например:
```shell
CONFIG_PATH=config/local.yaml
JWT_SECRET=qwerty
```

```shell
docker-compose up --build 
```

## Описание
* API сервиса полностью описан в файле [api.yaml](api/api.yaml) с помощью [OpenAPI](https://www.openapis.org/). 
* БД: в качесве базы данных используется postgres, для работы с бд используется популярный фреймворк [sqlx](https://github.com/jmoiron/sqlx).
* HTTP: для реализации http сервера используется фреймворк [gin](https://github.com/gin-gonic/gin).
* Для авторизации используются JWT токены и библиотека [jwt](https://github.com/golang-jwt/jwt).
* БД и сам сервис запускаются в докер контейнере(скрипт [wait-for-postgres.sh](wait-for-postgres.sh) нужен для того, чтобы дождатся запуска бд перед стартом приложения).
* Используется стандартная структура проекта(cmd, internal, pkg, ...)
* Код разбит на слои: сервисный слой, слой для работы с бд, http handler'ы - папка internal.
* В папке pkg переиспользуемый код(logging, gracefull shutdown).


## Примеры взаимодействия с API

* Регистрация пользователя:
    ```shell
    curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"username":"${USERNAME?}","password":"${PASSWORD?}"}' \
    'http://localhost:8080/user/register'
    ```

* Авторизация:
    ```shell
    curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"username":"${USERNAME?}","password":"${PASSWORD?}"}' \
    'http://localhost:8080/user/login'
    ```

* Создание продукта:
    ```shell
    curl -X 'POST' \
    -H 'Content-Type: application/json' \
    -H 'Authorization: Bearer ${TOKEN?}' \
    -d '{
    "name": "gopher1",
    "price": 42,
    "quantity": 42
    }' \
    'http://localhost:8080/product/add'
    ```

* Получить продукт по ID:
    ```shell
    curl -X 'GET' \
    -H 'Content-Type: application/json' \
    -H 'Authorization: Bearer ${TOKEN?}' \
    'http://localhost:8080/product/${ID?}
    ```

* Редактировать продукт по ID:
    ```shell
    curl -X 'PUT' \
    -H 'Content-Type: application/json' \
    -H 'Authorization: Bearer ${TOKEN?}' \
    -d '{
      "name": "gopher1",
      "price": 43,
      "quantity": 43
    }' \
    'http://localhost:8080/product/${ID?}'
    ```
  
* Удалить продукт по ID:
    ```shell
    curl -X 'DELETE' \
    -H 'Content-Type: application/json' \
    -H 'Authorization: Bearer ${TOKEN?}' \
    'http://localhost:8080/product/${ID?}'
    ```
  
* Получить продукты:
    ```shell
    curl -X 'GET' \
    -H 'Content-Type: application/json' \
    -H 'Authorization: Bearer ${TOKEN?}' \
    'http://localhost:8080/products?name=aaaaaaaaaaaa&sort_by=last_create'
    ```
  