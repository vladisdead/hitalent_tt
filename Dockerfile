FROM golang:1.23.12-alpine AS builder

WORKDIR /hitalent_tt

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o ./bin/hitalent_tt cmd/main.go

FROM alpine AS runner

COPY --from=builder /hitalent_tt/bin/hitalent_tt /
COPY config/config.yaml config/config.yaml

CMD ["/hitalent_tt"]