FROM golang:alpine as builder

WORKDIR /src
COPY . .

RUN apk add --no-cache --virtual .build-deps \
		git
RUN CGO_ENABLED=0 go build -o dmstudio-server

FROM alpine:latest

WORKDIR /app
COPY --from=builder /src/dmstudio-server .

EXPOSE 8080
CMD ["./dmstudio-server"]
