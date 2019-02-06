
FROM golang:alpine
WORKDIR /app
COPY . .
RUN apk add --update -t build-deps curl go git libc-dev gcc libgcc

RUN go build -o main .
ENTRYPOINT ["./main"]