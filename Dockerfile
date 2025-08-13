# Используем официальный образ Go 
FROM golang:1.24-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы зависимостей и загружаем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код 
COPY . .

# Собираем статичный бинарный файл для Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o /server main.go

# Используем минимальный базовый образ
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем скомпилированный бинарный файл из этапа сборки
COPY --from=builder /server .

# Копируем папку с файлами фронтенда из этапа сборки
COPY --from=builder /app/web ./web

# Указываем, что контейнер будет слушать этот порт
EXPOSE 7540

# Команда для запуска сервера при старте контейнера
CMD ["./server"]