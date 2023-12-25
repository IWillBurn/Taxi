FROM alpine

WORKDIR /app

COPY --from=build:develop /app/cmd/app ./app

CMD ["/app/app", "-c", "config.yaml"]
