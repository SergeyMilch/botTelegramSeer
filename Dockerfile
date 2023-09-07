# Этап 1: Создание бинарного файла приложения
FROM golang:1.20-alpine3.18 AS builder

COPY . /github.com/SergeyMilch/botTelegramSeer/

WORKDIR /github.com/SergeyMilch/botTelegramSeer/

RUN go mod download

RUN go build -o ./bin/bot main.go

# Этап 2: Запуск приложения в минимальном образе
FROM alpine:latest

WORKDIR /root/
# Копирование бинарного файла из предыдущего этапа
COPY --from=0 /github.com/SergeyMilch/botTelegramSeer/bin/bot .
COPY .env /root/.env
COPY updated_Azazel.txt /root/updated_Azazel.txt

EXPOSE 80

CMD ["./bot"]
