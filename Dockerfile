ARG UID="1000"

FROM golang:1.24 AS build
ARG TARGETARCH
ARG TARGETOS
ARG VERSION
ARG COMMIT
ARG UID
RUN useradd --uid ${UID} tmpbbs
WORKDIR /app
COPY . .
ENV CGO_ENABLED=0
ENV GOARCH=${TARGETARCH}
ENV GOOS=${TARGETOS}
RUN go build -ldflags "-s -w -X github.com/mmb/tmpbbs/internal/tmpbbs.Version=${VERSION:-$(git tag --points-at HEAD)} -X github.com/mmb/tmpbbs/internal/tmpbbs.Commit=${COMMIT:-$(git rev-parse HEAD)}" -v

FROM scratch
ARG UID
COPY --from=build /app/tmpbbs /tmpbbs
EXPOSE 8080/tcp
EXPOSE 8081/tcp
USER ${UID}
ENV TMPBBS_JSON_LOG=true
ENTRYPOINT ["/tmpbbs"]
