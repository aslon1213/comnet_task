go mod tidy
go test -v -cover cmd/comnet_task/main_test.go
go run cmd/comnet_task/main.go