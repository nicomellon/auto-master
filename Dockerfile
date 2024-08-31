FROM golang AS builder
WORKDIR /app
COPY ./api .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine
COPY --from=builder /app/main .
CMD ["./main"]
