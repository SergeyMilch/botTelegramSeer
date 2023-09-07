# Этап 1: Создание бинарного файла приложения
FROM golang:1.20-alpine3.18 AS builder

WORKDIR /app

# Копирование всех файлов в контекст сборки
COPY . .

RUN go mod download

# Сборка бинарного файла
RUN go build -o ./bin/bot .

# Этап 2: Запуск приложения в минимальном образе
FROM alpine:latest

WORKDIR /root/

# Копирование бинарного файла и других файлов
COPY --from=builder /app/bin/bot .

COPY updated_Azazel.txt /root/updated_Azazel.txt

EXPOSE 80

CMD ["./bot"]
