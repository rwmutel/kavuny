FROM hazelcast/hazelcast:latest AS build-stage

COPY --from=golang:1.22-alpine /usr/local/go/ /usr/local/go/
ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /app

USER root

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN go build -o /auth_service

FROM hazelcast/hazelcast:latest AS build-release-stage

COPY --from=build-stage /auth_service /auth_service

EXPOSE 8080

CMD /auth_service
