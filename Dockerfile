ARG UID="1000"

FROM golang:1.25 AS build
ARG TARGETARCH
ARG TARGETOS
ARG VERSION
ARG COMMIT
ARG UID
RUN useradd --uid ${UID} tmpbbs
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
ENV GOARCH=${TARGETARCH}
ENV GOOS=${TARGETOS}
RUN go build -ldflags "-s -w -X main.version=${VERSION:-$(git tag --points-at HEAD)} -X main.commit=${COMMIT:-$(git rev-parse HEAD)}" -v

FROM scratch
ARG UID
COPY --from=build /app/tmpbbs /tmpbbs
EXPOSE 8080/tcp
EXPOSE 8081/tcp
USER ${UID}
ENV TMPBBS_JSON_LOG=true
ENTRYPOINT ["/tmpbbs"]
