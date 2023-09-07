# Этап 1: Создание бинарного файла приложения
FROM golang:1.20-alpine3.18 AS builder

COPY . .

WORKDIR /app

RUN go mod download

RUN go build -o ./bin/bot main.go

# Этап 2: Запуск приложения в минимальном образе
FROM alpine:latest

WORKDIR /root/
# Копирование бинарного файла из предыдущего этапа
COPY --from=0 /app/bin/bot .
COPY .env /root/.env
COPY updated_Azazel.txt /root/updated_Azazel.txt

EXPOSE 80

CMD ["./bot"]
