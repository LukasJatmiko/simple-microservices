FROM golang:1.21.3-alpine3.18

WORKDIR /app

COPY . ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /data-processor

EXPOSE 8081

CMD ["/data-processor"]