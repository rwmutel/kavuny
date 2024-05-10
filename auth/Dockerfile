FROM hazelcast/hazelcast:latest AS build-stage

COPY --from=golang:1.22-alpine /usr/local/go/ /usr/local/go/
ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /app

COPY go.mod go.sum *.go ./
RUN go mod download

USER root

RUN go build -o /auth_service

FROM hazelcast/hazelcast:latest AS build-release-stage

COPY --from=build-stage /auth_service /auth_service

EXPOSE 8080

CMD ["/bin/bash", "-c", "hz start & sleep 20 && /auth_service"]
