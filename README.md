## Запуск

Перед запуском нужно установить переменные окружения в файле .env
Нужно указать CONFIG_PATH и JWT_SECRET

Например:
```shell
CONFIG_PATH=config/local.yaml
JWT_SECRET=qwerty
```

```shell
make run
```

После примените все миграции из папки [migrations](migrations/).

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
  