FROM alpine

WORKDIR /app

COPY --from=build:develop /app/cmd/client_service/app ./app

CMD ["/app/app", "-c", "config.yaml"]
