FROM golang:latest as builder

RUN apt-get update && apt-get install -y gcc-aarch64-linux-gnu

COPY . /app/pass-keeper-bot/
WORKDIR /app/pass-keeper-bot/

RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o ./bin/main main.go

FROM gcr.io/distroless/base:latest

COPY --from=0 /app/pass-keeper-bot/bin/main /main

EXPOSE 80

CMD ["./main"]