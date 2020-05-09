FROM golang:1.14.2-alpine3.11 AS build

WORKDIR /src

RUN apk add --no-cache \
	upx

ARG LDFLAGS

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o producer -a -ldflags="$LDFLAGS" cmd/main.go && \
    upx --ultra-brute producer

FROM scratch

COPY --from=build /src/producer /producer

ENV GIN_MODE release
ENTRYPOINT [ "/producer"  ]
