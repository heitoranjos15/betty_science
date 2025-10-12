FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app/main .

EXPOSE 8020

COPY ./entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

