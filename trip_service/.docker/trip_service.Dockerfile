FROM alpine

WORKDIR /app

COPY --from=build:develop /app/cmd/trip_service/app ./app

CMD ["/app/app", "-c", "config.yaml"]
