FROM alpine:3.10
WORKDIR /app
COPY tezos-indexer-api /app/

ENTRYPOINT ["/app/tezos-indexer-api"]
CMD ["-log-http", "-db", "postgres://inexer:indexer@db:5432/mainnet"]
