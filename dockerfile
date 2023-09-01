FROM golang 
COPY . /app
WORKDIR /app
RUN go mod tidy
RUN go build -o main cmd/comnet_task/main.go
ENTRYPOINT [ "./main" ]
