FROM golang:1.21.1-alpine AS build

WORKDIR /github-login
COPY . .

RUN apk add --no-cache git
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o github-login

FROM alpine:latest

WORKDIR /github-login
COPY --from=build /github-login/github-login /github-login/github-login
COPY .env /github-login/.env

EXPOSE 8086

CMD ["./github-login"]
