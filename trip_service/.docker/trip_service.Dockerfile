FROM golang:1.19-alpine as builder

RUN apk add --no-cache git

WORKDIR /trip_service

COPY /trip_service .


RUN --mount=type=cache,mode=0755,target=/go/pkg/mod go mod vendor
RUN go build ./cmd/main.go

FROM alpine:latest

COPY --from=builder /trip_service/main .
COPY --from=builder /trip_service/.config .
COPY --from=builder /trip_service/test .

CMD ["./main", "-c", "./trip_service.yaml"]
