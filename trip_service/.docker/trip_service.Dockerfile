FROM alpine:latest

COPY --from=builder /trip_service/main .
COPY --from=builder /trip_service/.config .
COPY --from=builder /trip_service/test .

CMD ["./main", "-c", "./trip_service.yaml"]
