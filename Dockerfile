FROM golang:latest AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine

COPY --from=builder /app/main /app/
COPY --from=builder /app/migrations /app/migrations
#COPY --from=builder /app/app.env /app/

WORKDIR /app

CMD ["./main"]
