FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download -x

COPY . .

ENV CGO_ENABLED 0

RUN go build -o /app/bin/sso ./cmd/sso/main.go


FROM alpine:3.20

RUN apk add --no-cache \
  ca-certificates \
  tzdata 

WORKDIR /app
COPY sso.env /app/sso.env
COPY --from=builder /app/go.mod /app/go.mod
COPY --from=builder /app/bin/sso /app/bin/sso
COPY --from=builder /app/settings /app/settings/
COPY --from=builder /app/secrets/private.pem /app/secrets/private.pem


EXPOSE 44044 8090

CMD ["/app/bin/sso"]