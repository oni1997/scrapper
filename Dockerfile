FROM golang:alpine

RUN apk add --no-cache chromium

WORKDIR /app
COPY . .

RUN go mod tidy && go build -o app

CMD ["./app"]
