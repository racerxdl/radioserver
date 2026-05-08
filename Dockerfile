FROM golang:1.25-alpine AS build

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o radioserver ./cmd/server/

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /build/radioserver .

ENV LISTEN_ADDR=:4050
ENV CENTER_FREQUENCY=106300000
ENV GAIN=0

EXPOSE 4050

CMD ["./radioserver"]
