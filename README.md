## Run

Set env variables CONFIG_PATH and JWT_SECRET before run:
```shell
CONFIG_PATH=config/local.yaml
JWT_SECRET=qwerty
```

```shell
make run
```

After apply all migrations from folder [migrations](migrations/).

## Examples

* User registration:
    ```shell
    curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"username":"${USERNAME?}","password":"${PASSWORD?}"}' \
    'http://localhost:8080/user/register'
    ```

* Login:
    ```shell
    curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"username":"${USERNAME?}","password":"${PASSWORD?}"}' \
    'http://localhost:8080/user/login'
    ```

* Create a product:
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

* Get product by id:
    ```shell
    curl -X 'GET' \
    -H 'Content-Type: application/json' \
    -H 'Authorization: Bearer ${TOKEN?}' \
    'http://localhost:8080/product/${ID?}
    ```

* Update product by id:
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
  
* Delete product by id:
    ```shell
    curl -X 'DELETE' \
    -H 'Content-Type: application/json' \
    -H 'Authorization: Bearer ${TOKEN?}' \
    'http://localhost:8080/product/${ID?}'
    ```
  
* Get all products:
    ```shell
    curl -X 'GET' \
    -H 'Content-Type: application/json' \
    -H 'Authorization: Bearer ${TOKEN?}' \
    'http://localhost:8080/products?name=aaaaaaaaaaaa&sort_by=last_create'
    ```
  
