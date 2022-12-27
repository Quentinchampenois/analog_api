FROM golang:alpine as builder

LABEL maintainer="Quentin Champenois"

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o analog_api .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/.env .
COPY --from=builder /app/public ./public
COPY --from=builder /app/analog_api .

EXPOSE 8080 8080
CMD ["./analog_api"]