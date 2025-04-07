## Run

Set env variables CONFIG_PATH and JWT_SECRET before run:
```shell
CONFIG_PATH=config/local.yaml
JWT_SECRET=qwerty
```
* Generate certificates:

```shell
openssl req -x509 -newkey rsa:4096 -keyout .cert/key.pem -out .cert/cert.pem -days 365 -nodes -config config/localhost.cnf
```


```shell
make run
```

After apply all migrations from folder [migrations](migrations/).

## Examples

* User registration:
    ```shell
    curl --cacert .cert/cert.pem -X POST \
    -H "Content-Type: application/json" \
    -d '{"username":"${USERNAME?}","password":"${PASSWORD?}"}' \
    'https://localhost:8080/user/register'
    ```

* Login:
    ```shell
    curl --cacert .cert/cert.pem -X POST \
    -H "Content-Type: application/json" \
    -d '{"username":"${USERNAME?}","password":"${PASSWORD?}"}' \
    'https://localhost:8080/user/login'
    ```

* Create a product:
    ```shell
    curl --cacert .cert/cert.pem -X 'POST' \
    -H 'Content-Type: application/json' \
    -H 'Authorization: Bearer ${TOKEN?}' \
    -d '{
    "name": "gopher1",
    "price": 42,
    "quantity": 42
    }' \
    'https://localhost:8080/product/add'
    ```

* Get product by id:
    ```shell
    curl --cacert .cert/cert.pem -X 'GET' \
    -H 'Content-Type: application/json' \
    -H 'Authorization: Bearer ${TOKEN?}' \
    'https://localhost:8080/product/${ID?}
    ```

* Update product by id:
    ```shell
    curl --cacert .cert/cert.pem -X 'PUT' \
    -H 'Content-Type: application/json' \
    -H 'Authorization: Bearer ${TOKEN?}' \
    -d '{
      "name": "gopher1",
      "price": 43,
      "quantity": 43
    }' \
    'https://localhost:8080/product/${ID?}'
    ```
  
* Delete product by id:
    ```shell
    curl --cacert .cert/cert.pem -X 'DELETE' \
    -H 'Content-Type: application/json' \
    -H 'Authorization: Bearer ${TOKEN?}' \
    'https://localhost:8080/product/${ID?}'
    ```
  
* Get all products:
    ```shell
    curl --cacert .cert/cert.pem -X 'GET' \
    -H 'Content-Type: application/json' \
    -H 'Authorization: Bearer ${TOKEN?}' \
    'https://localhost:8080/products?name=aaaaaaaaaaaa&sort_by=last_create'
    ```
  
