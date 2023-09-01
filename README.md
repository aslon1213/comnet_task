task - тестовое задание_2023.crash

# To run application server

- firstly set SIGNING_SECRET variable in .env.example file to any value and rename the file to .env

```shell
    go mod tidy
    go run cmd/comnet_task/main.go
```

# To run tests

```shell
    go test -v -cover cmd/comnet_task/main_test.go
```
