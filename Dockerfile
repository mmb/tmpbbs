ARG UID="1000"

FROM golang:1.24 AS build
ARG TARGETARCH
ARG TARGETOS
ARG UID
RUN useradd --uid ${UID} tmpbbs
WORKDIR /app
COPY . .
ENV CGO_ENABLED=0
ENV GOARCH=${TARGETARCH}
ENV GOOS=${TARGETOS}
RUN go build -ldflags "-s -w -X main.version=$(git tag --points-at HEAD)-$(git rev-parse HEAD)" -v

FROM scratch
ARG UID
COPY --from=build /app/tmpbbs /tmpbbs
EXPOSE 8080/tcp
USER ${UID}
ENTRYPOINT ["/tmpbbs"]
