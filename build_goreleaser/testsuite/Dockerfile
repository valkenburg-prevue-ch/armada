ARG BASE_IMAGE=alpine:3.17.0
FROM ${BASE_IMAGE}

RUN addgroup -S -g 2000 armada && adduser -S -u 1000 armada -G armada
USER armada

COPY testsuite /app/

WORKDIR /app

ENTRYPOINT ["./testsuite"]
