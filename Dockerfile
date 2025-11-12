FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o finscheduler ./cmd/finscheduler

FROM alpine:latest AS run

WORKDIR /app
COPY --from=build /app/finscheduler .
COPY configs ./configs
EXPOSE 8080
CMD ["./finscheduler"]