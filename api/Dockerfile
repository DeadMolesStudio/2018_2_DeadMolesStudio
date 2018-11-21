FROM golang:alpine as builder

WORKDIR /src
COPY . .

RUN apk add --no-cache --virtual .build-deps \
		git
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -a -installsuffix cgo -ldflags="-w -s" -o dmstudio-server

FROM scratch

WORKDIR /app
COPY --from=builder /src/dmstudio-server .
COPY logger/logger-config.json logger/logger-config.json
COPY migrations migrations

VOLUME ["/var/log/dmstudio", "/app/static"]

EXPOSE 8080
CMD ["./dmstudio-server"]
