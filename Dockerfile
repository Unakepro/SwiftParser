FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY data.csv /app/data.csv  

WORKDIR /app/app
ENV GOARCH=amd64
RUN go build -v -o /app/main .


FROM golang:1.20 AS tester
WORKDIR /app
COPY --from=builder /app /app

FROM alpine:latest AS production

COPY --from=builder /app/main /main
COPY --from=builder /app/data.csv /data.csv  

CMD ["./main"]
