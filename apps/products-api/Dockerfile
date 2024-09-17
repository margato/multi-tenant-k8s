FROM golang:1.22.5-alpine

WORKDIR /app

COPY . .

RUN go get -d -v ./...
RUN go build -ldflags "-s -w"

ENV PORT=8080
ENV GIN_MODE=release
EXPOSE 8080

ENTRYPOINT ["./products-api"]