FROM alpine

WORKDIR /app

COPY --from=build:develop /app/cmd/offering_service/app ./app

CMD ["/app/app", "-c", "config.yaml"]
