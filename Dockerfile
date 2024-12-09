# Используем базовый образ с Go
FROM golang:1.22.5-alpine3.19

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Устанавливаем git, чтобы Go мог загрузить зависимости
RUN apk update && apk add git

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod download

# Загружаем и очищаем зависимости Go
RUN go mod tidy

# Копируем остальные файлы
COPY . .

# Собираем исполняемый файл
RUN go build -o bot .

# Указываем команду запуска
CMD ["./bot"]
