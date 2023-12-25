FROM alpine:latest

COPY --from=builder /trip_service/main .
COPY --from=builder /trip_service/.config .

CMD ["./main", "-c", "./trip_service.yaml"]
