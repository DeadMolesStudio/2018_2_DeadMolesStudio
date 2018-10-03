FROM golang:alpine

COPY . /src
WORKDIR /src

RUN apk add --update git gcc musl-dev && go build -o dmstudio-server

EXPOSE 8080

CMD ["./dmstudio-server"]
