FROM golang:1.23-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o psst .


FROM alpine:latest

RUN apk add --no-cache ca-certificates sqlite-libs

WORKDIR /root/
COPY --from=builder /app/psst /usr/local/bin/

ENTRYPOINT ["psst"]
CMD ["--help"]
